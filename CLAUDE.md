# CLAUDE.md

必须遵循的原则：Claude Code 的回复、生成的注释、文档，全部使用中文。

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

r0Website-server 是一个基于 Go 语言的个人网站后端服务器，使用 Gin 框架和 MongoDB 数据库。提供文章/博客管理、用户认证、分类管理、图床服务、中文分词和全文搜索功能。

## 开发命令

### 构建和运行
```bash
# 开发服务器（创建日志目录，构建并运行）
./run-server.sh

# 停止服务器
./stop-server.sh

# 重启服务器
./restart-server.sh

# Docker 构建和部署
./docker-build.sh
```

### 手动构建
```bash
# 创建日志目录
mkdir logs
touch logs/r0Website_server.log

# 构建二进制文件
go build -o ./r0Website-server ./main/main.go

# 运行服务器
./r0Website-server
```

### 依赖管理
```bash
# 更新依赖
go mod tidy

# 设置 Go 代理（适用于中国开发者）
go env -w GOPROXY=https://goproxy.cn,direct
```

## 架构设计

### 自定义依赖注入（r0Ioc）
项目使用基于反射的自定义 DI 容器，位于 `r0Ioc/` 目录。组件自动注册和注入。关键点：
- 组件在 `initialize/utils_init.go` 中自我注册
- DI 容器使用退出函数管理组件生命周期
- 所有主要组件（API、服务、DAO）都由 DI 管理

### 分层架构
- **API 层**（`api/`）：HTTP 处理器，分为 `admin/`（受保护）和 `base/`（公共）
- **服务层**（`service/`）：业务逻辑实现
- **DAO 层**（`dao/`）：使用 MongoDB 驱动的数据库操作
- **模型层**（`models/`）：数据结构 - PO（持久化对象）、VO（值对象）、BO（业务对象）

### 关键组件
- **认证系统**：基于 JWT，中间件位于 `middleware/jwt_midw.go`
- **中文文本处理**：使用 go-ego/gse 对文章进行分词
- **MongoDB 集合**：users、articles、categories、albums/images
- **配置管理**：基于 YAML，位于 `config/config.yml`

### 数据库要求
- MongoDB 5.0.9+，需要 dbOwner 权限
- `md_words` 和 `title_words` 字段需要全文搜索索引
- 中文分词结果存储在分词字段中

### API 结构
- 基础路由（公共）：文章浏览、分类、图片、点赞
- 管理路由（JWT 保护）：文章 CRUD、分类管理、用户管理
- 服务器运行在 8202 端口（可在 config.yml 中配置）

## 重要说明

- 没有正式的测试结构 - 需要手动验证更改
- 主要文档为中文
- Docker 网络使用自定义子网 172.97.0.0/16
- 日志写入 `./logs/r0Website_server.log`
- 可通过 `export GIN_MODE=release` 设置环境

---

## 项目结构详解

### 分层架构
```
api/          # API层 - HTTP请求处理
├── admin/   # 管理端API（需要JWT认证）
└── base/    # 公共API（无需认证）

service/     # 服务层 - 业务逻辑实现
├── albums_service.go    # 相册服务
├── article_service.go   # 文章服务
├── category_service.go  # 分类服务
├── image_service.go     # 图片服务
└── user_service.go      # 用户服务

dao/         # 数据访问层 - MongoDB数据库操作
├── basic_dao.go         # 基础MongoDB连接
├── article_dao.go       # 文章数据访问
├── category_dao.go      # 分类数据访问
├── images_dao.go        # 图片数据访问
└── user_dao.go          # 用户数据访问

models/      # 数据模型层
├── po/      # 持久化对象（数据库结构）
├── vo/      # 值对象（API数据传输）
└── bo/      # 业务对象（业务逻辑封装）
```

### 开发架构规范

#### 1. API层开发规范
- **控制器结构**：使用依赖注入，通过 `R0Ioc:"true"` 标签注入服务
- **参数接收**：使用VO结构体绑定请求参数
- **响应格式**：统一使用 `msg.NewMsg().Success()` 和 `msg.NewMsg().Failed()`
- **错误处理**：参数验证错误返回400，业务逻辑错误在Service层处理

```go
// 标准API控制器示例
type UserController struct {
    UserService *service.UserService `R0Ioc:"true"`
}

func (c *UserController) Login(ctx *gin.Context) {
    var params vo.LoginVo
    if err := ctx.ShouldBind(&params); err != nil {
        ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed("参数错误"))
        return
    }

    result, err := c.UserService.UserLogin(params)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, msg.NewMsg().Failed(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, msg.NewMsg().Success(result))
}
```

#### 2. Service层开发规范
- **服务注入**：使用 `R0Ioc:"true"` 标签注入需要的DAO组件
- **业务逻辑**：所有业务逻辑都在这里实现，不直接操作数据库
- **错误返回**：返回具体的错误信息，让API层决定HTTP状态码
- **VO转换**：负责PO到VO的转换，API层只处理VO

```go
type UserService struct {
    UserDao *dao.UserDao `R0Ioc:"true"`
}

func (s *UserService) UserLogin(params vo.LoginVo) (*vo.LoginResultVo, error) {
    // 业务逻辑实现
    user, err := s.UserDao.FindObjByEmail(params.Email)
    if err != nil {
        return nil, errors.New("用户不存在")
    }

    // 返回VO，不是PO
    return &vo.LoginResultVo{
        Token: token,
        Username: user.Username,
        Email: user.Email,
    }, nil
}
```

#### 3. DAO层开发规范
- **基础结构**：所有DAO都嵌入 `BasicDaoMongo` 获得基础MongoDB操作
- **集合名称**：使用常量定义集合名，如 `const UserColl = "users"`
- **错误处理**：返回具体的错误，上层决定如何处理
- **索引管理**：在初始化时确保必要的索引存在

#### 4. 模型定义规范
- **PO模型**：直接对应数据库集合结构，包含所有数据库字段
- **VO模型**：用于API数据传输，只包含需要暴露的字段
- **命名规范**：PO模型名与集合名对应，VO模型名明确表明用途

```go
// PO - 数据库模型
type User struct {
    ID       primitive.ObjectID `bson:"_id"`
    Username string             `bson:"username"`
    Password string             `bson:"password"` // 敏感信息
    Email    string             `bson:"email"`
}

// VO - API传输模型
type UserInfoVo struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    // 不包含密码等敏感信息
}
```

#### 5. 依赖注入使用规范
- **组件注册**：在 `initialize/utils_init.go` 中注册所有组件
- **自动注入**：通过结构体标签 `R0Ioc:"true"` 实现自动注入
- **生命周期**：支持退出函数管理组件生命周期
- **单例模式**：所有组件默认都是单例

#### 6. 新增功能开发流程
1. **定义数据模型**：在 `models/po/` 中定义PO，在 `models/vo/` 中定义VO
2. **创建DAO层**：在 `dao/` 中创建对应的数据访问对象
3. **实现Service层**：在 `service/` 中创建业务逻辑服务
4. **编写API层**：在 `api/base/` 或 `api/admin/` 中创建API控制器
5. **配置路由**：在 `router/` 中添加对应的路由配置
6. **注册组件**：在 `initialize/utils_init.go` 中注册新组件

#### 7. API设计原则
- **RESTful风格**：使用标准的HTTP方法（GET/POST/PUT/DELETE）
- **统一响应**：所有API返回统一格式的JSON响应
- **状态码规范**：[详细状态码说明]
  - 200：操作成功
  - 400：请求参数错误
  - 401：认证失败
  - 403：权限不足
  - 404：资源不存在
  - 500：服务器内部错误

#### 8. 中文分词集成
- **分词库**：使用 `github.com/go-ego/gse` 实现中文分词
- **分词字段**：文章内容分词存储在 `md_words`，标题分词存储在 `title_words`
- **搜索功能**：基于分词结果实现全文搜索
- **索引要求**：确保相关字段建立文本索引

#### 9. 图片管理规范
- **存储方式**：使用腾讯云COS对象存储
- **文件路径**：统一存储在 `/somnium/primitive/` 目录下
- **格式支持**：JPG、PNG、GIF、WebP格式
- **文件大小**：单文件最大支持10MB
- **分类管理**：支持图片在多个分类中的位置管理
- **标签同步**：上传图片时自动同步标签信息

#### 10. 错误处理最佳实践
- **分层处理**：DAO层返回原始错误，Service层包装业务错误，API层决定HTTP响应
- **错误类型**：定义具体的错误类型，如 `UniqueError` 表示唯一性约束错误
- **日志记录**：使用 `global.Logger` 记录错误日志，包含足够的上下文信息
- **用户友好**：返回给用户的错误信息要清晰易懂，不包含技术细节

### 开发注意事项

1. **代码风格**：所有注释和文档必须使用中文
2. **依赖管理**：使用 `go mod tidy` 管理依赖，设置国内代理加速下载
3. **测试验证**：没有自动化测试，需要手动验证所有功能
4. **日志管理**：日志写入 `logs/` 目录，按日期分文件存储
5. **安全配置**：JWT密钥和数据库密码通过配置文件管理，不要硬编码
6. **性能优化**：MongoDB查询需要合理设计索引，特别是搜索功能
7. **Docker部署**：使用自定义网络 `172.97.0.0/16`，确保端口映射正确