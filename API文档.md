# r0Website-server API 文档

## 概述

本文档描述了 r0Website-server 的所有 API 接口。API 分为两个主要分组：
- **Base API**：公共接口，无需认证
- **Admin API**：管理接口，需要 JWT 认证

所有 API 接口都以 `/api` 为前缀。

---

## 认证说明

### JWT 认证
Admin API 大部分接口需要 JWT 认证。在请求头中添加：
```
Authorization: Bearer <your_jwt_token>
```

### 获取 JWT Token
通过登录接口获取 JWT token：
- Base 用户登录：`POST /api/base/login`
- Admin 用户登录：`POST /api/admin/login`

---

## Base API（公共接口）

### 用户管理

#### 用户注册
```http
POST /api/base/register
```

**请求体：**
```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

**响应：**
```json
{
  "code": 200,
  "msg": "注册成功",
  "data": {
    "token": "jwt_token_string",
    "user": {
      "id": "user_id",
      "username": "username",
      "email": "email"
    }
  }
}
```

#### 用户登录
```http
POST /api/base/login
```

**请求体：**
```json
{
  "username": "string",
  "password": "string"
}
```

**响应：**
```json
{
  "code": 200,
  "msg": "登录成功",
  "data": {
    "token": "jwt_token_string",
    "user": {
      "id": "user_id",
      "username": "username",
      "email": "email"
    }
  }
}
```

---

### 文章管理

#### 搜索文章
```http
GET /api/base/article?keyword=搜索关键词&page=1&size=10
```

**查询参数：**
- `keyword`（可选）：搜索关键词，支持模糊搜索
- `page`（可选）：页码，默认 1
- `size`（可选）：每页大小，默认 10

**响应：**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "articles": [
      {
        "id": "article_id",
        "title": "文章标题",
        "summary": "文章摘要",
        "cover": "封面图片URL",
        "category": "分类名称",
        "tags": ["标签1", "标签2"],
        "pv": 100,
        "praise": 10,
        "create_time": "2023-01-01T00:00:00Z",
        "update_time": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "size": 10
  }
}
```

#### 获取文章详情
```http
GET /api/base/article/{id}
```

**路径参数：**
- `id`：文章ID

**响应：**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "id": "article_id",
    "title": "文章标题",
    "content": "文章内容（Markdown格式）",
    "html_content": "HTML格式的内容",
    "summary": "文章摘要",
    "cover": "封面图片URL",
    "category": "分类名称",
    "tags": ["标签1", "标签2"],
    "pv": 100,
    "praise": 10,
    "md_words": ["分词1", "分词2"],
    "title_words": ["标题分词1", "标题分词2"],
    "create_time": "2023-01-01T00:00:00Z",
    "update_time": "2023-01-01T00:00:00Z"
  }
}
```

#### 增加文章浏览量
```http
PUT /api/base/article/{id}/pv
```

**路径参数：**
- `id`：文章ID

**响应：**
```json
{
  "code": 200,
  "msg": "更新成功",
  "data": {
    "pv": 101
  }
}
```

#### 增加文章点赞数
```http
PUT /api/base/article/{id}/praise
```

**路径参数：**
- `id`：文章ID

**响应：**
```json
{
  "code": 200,
  "msg": "点赞成功",
  "data": {
    "praise": 11
  }
}
```

#### 获取分类下的文章
```http
GET /api/base/article/category/{name}?page=1&size=10
```

**路径参数：**
- `name`：分类名称

**查询参数：**
- `page`（可选）：页码，默认 1
- `size`（可选）：每页大小，默认 10

**响应：** 同搜索文章接口

---

### 分类管理

#### 获取所有分类
```http
GET /api/base/category/all
```

**响应：**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "categories": [
      {
        "id": "category_id",
        "name": "分类名称",
        "description": "分类描述",
        "article_count": 10,
        "create_time": "2023-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### 图床管理

#### 图集管理

##### 创建图集
```http
POST /api/base/picbed/album
```

**请求体：**
```json
{
  "name": "图集名称",
  "description": "图集描述",
  "cover": "封面图片ID",
  "tags": ["标签1", "标签2"],
  "author": "作者名称"
}
```

##### 获取图集详情
```http
GET /api/base/picbed/album/{id}
```

**路径参数：**
- `id`：图集ID

##### 获取图集列表
```http
GET /api/base/picbed/album?page=1&size=10&sort=create_time&order=desc
```

**查询参数：**
- `page`（可选）：页码，默认 1
- `size`（可选）：每页大小，默认 10
- `sort`（可选）：排序字段
- `order`（可选）：排序方式（asc/desc）

##### 按标签查询图集
```http
GET /api/base/picbed/album/tag/{tag}?page=1&size=10
```

##### 按作者查询图集
```http
GET /api/base/picbed/album/author/{author}?page=1&size=10
```

##### 搜索图集
```http
GET /api/base/picbed/album/search/{keyword}?page=1&size=10
```

##### 更新图集信息
```http
PUT /api/base/picbed/album/{id}
```

**请求体：**
```json
{
  "name": "新图集名称",
  "description": "新图集描述",
  "cover": "新封面图片ID",
  "tags": ["新标签1", "新标签2"]
}
```

##### 删除图集
```http
DELETE /api/base/picbed/album/{id}
```

#### 图片管理

##### 上传图片
```http
POST /api/base/picbed/image
```

**请求体：**
```json
{
  "name": "图片名称",
  "description": "图片描述",
  "url": "图片URL",
  "width": 1920,
  "height": 1080,
  "size": 2048000,
  "format": "jpg",
  "tags": ["标签1", "标签2"],
  "author": "作者名称"
}
```

##### 获取图片详情
```http
GET /api/base/picbed/image/{id}
```

##### 获取图片列表
```http
GET /api/base/picbed/image?page=1&size=10&sort=create_time&order=desc
```

##### 按标签查询图片
```http
GET /api/base/picbed/image/tag/{tag}?page=1&size=10
```

##### 搜索图片
```http
GET /api/base/picbed/image/search/{keyword}?page=1&size=10
```

##### 查询图片所在图集
```http
GET /api/base/picbed/image/{id}/albums
```

##### 删除图片
```http
DELETE /api/base/picbed/image/{id}
```

#### 图集图片关联管理

##### 添加图片到图集
```http
PUT /api/base/picbed/album/{albumId}/image
```

**请求体：**
```json
{
  "image_id": "图片ID",
  "layout": "布局信息",
  "description": "在图集中的描述"
}
```

##### 更新图片布局
```http
PUT /api/base/picbed/album/{albumId}/image/{imageId}/layout
```

**请求体：**
```json
{
  "layout": "新布局信息"
}
```

##### 从图集移除图片
```http
DELETE /api/base/picbed/album/{albumId}/image/{imageId}
```

##### 移动图片到另一个图集
```http
PUT /api/base/picbed/image/move
```

**请求体：**
```json
{
  "image_id": "图片ID",
  "from_album_id": "原图集ID",
  "to_album_id": "目标图集ID"
}
```

---

## Admin API（管理接口）

### 管理员登录
```http
POST /api/admin/login
```

**请求体：**
```json
{
  "username": "admin",
  "password": "password"
}
```

**响应：**
```json
{
  "code": 200,
  "msg": "登录成功",
  "data": {
    "token": "jwt_token_string",
    "user": {
      "id": "admin_id",
      "username": "admin",
      "role": "admin"
    }
  }
}
```

### 文章管理（需要JWT认证）

#### 创建文章
```http
POST /api/admin/article
```

**请求体：**
```json
{
  "title": "文章标题",
  "content": "文章内容（Markdown格式）",
  "summary": "文章摘要",
  "cover": "封面图片URL",
  "category": "分类名称",
  "tags": ["标签1", "标签2"]
}
```

#### 更新文章
```http
POST /api/admin/article/{id}
```

**路径参数：**
- `id`：文章ID

**请求体：** 同创建文章

#### 删除文章
```http
DELETE /api/admin/article/{id}
```

#### 通过文件上传创建/更新文章
```http
POST /api/admin/article/upload
POST /api/admin/article/upload/{id}
```

**请求体：** multipart/form-data
- `file`：Markdown文件
- 其他字段同创建文章接口

### 分类管理（需要JWT认证）

#### 归档文章分类
```http
POST /api/admin/category/archive
```

**请求体：**
```json
{
  "category_name": "要归档的分类名称"
}
```

---

## 响应格式说明

### 统一响应格式
所有 API 接口都使用以下统一响应格式：

```json
{
  "code": 200,        // 状态码
  "msg": "操作成功",   // 消息
  "data": {}          // 数据（可选）
}
```

### 状态码说明
- `200`：操作成功
- `400`：请求参数错误
- `401`：未认证或认证失败
- `403`：权限不足
- `404`：资源不存在
- `500`：服务器内部错误

---

## 分页参数说明

支持分页的接口使用以下参数：
- `page`：页码，从 1 开始
- `size`：每页大小，默认 10
- `sort`：排序字段
- `order`：排序方式（asc/desc）

分页响应格式：
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "list": [],        // 数据列表
    "total": 100,      // 总记录数
    "page": 1,         // 当前页码
    "size": 10,        // 每页大小
    "pages": 10        // 总页数
  }
}
```

---

## 错误处理

当发生错误时，API 会返回相应的错误信息：

```json
{
  "code": 400,
  "msg": "错误描述信息",
  "data": null
}
```

---

## 更新记录

- 2024.01：新增图床管理功能，包括图集和图片管理
- 2023.12：增加文章点赞功能
- 2023.11：优化中文分词和搜索功能
- 2023.10：初始版本发布

---

## 注意事项

1. **JWT Token 有效期**：JWT token 默认有效期为 3000 秒（50 分钟），需要在过期前重新获取
2. **中文搜索**：文章搜索支持中文分词，会自动对标题和内容进行分词处理
3. **图片上传**：图床功能需要先上传图片获取 URL，然后再创建图片记录
4. **分类归档**：归档操作会将分类下的所有文章设置为未分类状态
5. **权限控制**：Admin 接口需要管理员权限，普通用户无法访问