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