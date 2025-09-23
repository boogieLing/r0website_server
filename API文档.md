# r0Website-server API 文档

## 概述

本文档描述了 r0Website-server 的所有 API 接口。API 分为两个主要分组：
- **Base API**：公共接口，无需认证
- **Admin API**：管理接口，需要 JWT 认证

所有 API 接口都以 `/api` 为前缀。

### 主要功能模块
- **用户管理**：用户注册、登录、认证
- **文章管理**：文章CRUD、搜索、点赞、分类管理
- **图床管理**：图集管理、图片上传、分类管理、标签管理
- **分类系统**：支持图片在多个分类中的位置管理
- **标签系统**：标签同步、热门标签、标签搜索

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
  "username": "string",    // 用户名
  "password": "string",    // 密码
  "email": "string",       // 邮箱
  "phone": "string"        // 手机号（可选）
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
  "email": "string",       // 邮箱
  "password": "string"     // 密码
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
GET /api/base/article?search_text=搜索关键词&page_number=1&page_size=10&author=作者名
```

**查询参数：**
- `search_text`（可选）：搜索关键词，支持模糊搜索，允许空格
- `page_number`（可选）：页码，从1开始，默认 1
- `page_size`（可选）：每页大小，默认 10
- `author`（可选）：作者名
- `create_time_sort`（可选）：创建时间排序方向（1表示升序，-1表示降序）
- `update_time_sort`（可选）：更新时间排序方向（1表示升序，-1表示降序）
- `lazy`（可选）：是否懒惰加载，如果为true则不返回实体内容

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
GET /api/base/article/category/{name}?page_number=1&page_size=10&create_time_sort=-1&update_time_sort=-1
```

**路径参数：**
- `name`：分类名称

**查询参数：**
- `page_number`（可选）：页码，从1开始，默认 1
- `page_size`（可选）：每页大小，默认 10
- `create_time_sort`（可选）：创建时间排序方向（1表示升序，-1表示降序）
- `update_time_sort`（可选）：更新时间排序方向（1表示升序，-1表示降序）
- `lazy`（可选）：是否懒惰加载，如果为true则不返回实体内容

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
  "title": "图集标题",           // 必填
  "description": "图集描述",     // 必填
  "author": "作者名称"           // 必填
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

**请求体：** 使用 `models/po.Album` 结构体，包含以下主要字段：
```json
{
  "title": "图集标题",
  "description": "图集描述",
  "cover_image": "封面图片ID",
  "image_refs": [
    {
      "image_id": "图片ID",
      "position": {
        "x": 0.1,
        "y": 0.2,
        "width": 0.3,
        "height": 0.4,
        "unit": "%"
      },
      "caption": "图片标题",
      "description": "图片描述"
    }
  ],
  "tags": ["标签1", "标签2"],
  "author": "作者名称",
  "visibility": "public"
}
```

##### 删除图集
```http
DELETE /api/base/picbed/album/{id}
```

#### 图片管理

##### 上传图片（直接文件上传）
```http
POST /api/base/picbed/image
```

**请求体：** `multipart/form-data` 格式
```
file: 图片文件                    // 必填，支持的格式：jpg, jpeg, png, gif, webp
name: 图片名称                    // 可选，留空则使用文件名
tags: 标签1,标签2,标签3          // 可选，逗号分隔的标签
```

**响应：**
```json
{
  "code": 200,
  "msg": "图片上传成功",
  "data": {
    "id": "图片ID",
    "name": "图片名称",
    "cos_url": "腾讯云COS图片地址",
    "width": 1920,
    "height": 1080,
    "size": 2048000,
    "format": "image/jpeg",
    "tags": ["标签1", "标签2"],
    "positions": {
      "nexus": {
        "category_id": "nexus",
        "position": 1,
        "x": 0,
        "y": 0
      }
    },
    "uploaded_at": "2024-01-01T00:00:00Z"
  }
}
```

**注意：** 图片会自动上传到腾讯云COS，目录为 `/somnium/primitive/`，并自动添加到 "nexus" 分类

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

#### 图片分类管理

##### 创建分类
```http
POST /api/base/picbed/category
```

**请求体：**
```json
{
  "name": "分类名称",              // 必填
  "description": "分类描述"        // 可选
}
```

**响应：**
```json
{
  "code": 200,
  "msg": "分类创建成功",
  "data": {
    "id": "分类ID",
    "name": "分类名称",
    "description": "分类描述",
    "imageCount": 0,
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

##### 获取所有分类
```http
GET /api/base/picbed/category
```

**响应：**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "categories": [
      {
        "id": "分类ID",
        "name": "分类名称",
        "description": "分类描述",
        "imageCount": 10,
        "coverImage": "封面图片ID",
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 5
  }
}
```

##### 获取分类详情
```http
GET /api/base/picbed/category/{id}
```

##### 更新分类
```http
PUT /api/base/picbed/category/{id}
```

**请求体：**
```json
{
  "name": "新分类名称",
  "description": "新分类描述"
}
```

##### 删除分类
```http
DELETE /api/base/picbed/category/{id}
```

##### 获取分类中的图片
```http
GET /api/base/picbed/category/{id}/images?page=1&size=10&sort=position&order=asc
```

**查询参数：**
- `page`（可选）：页码，默认 1
- `size`（可选）：每页大小，默认 10
- `sort`（可选）：排序字段，支持 position, createdAt
- `order`（可选）：排序方式（asc/desc）

**响应：**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "images": [
      {
        "id": "图片ID",
        "name": "图片名称",
        "cos_url": "腾讯云COS地址",
        "width": 1920,
        "height": 1080,
        "position": {
          "category_id": "分类ID",
          "position": 1,
          "x": 0,
          "y": 0
        }
      }
    ],
    "total": 20,
    "page": 1,
    "size": 10
  }
}
```

##### 添加图片到分类
```http
POST /api/base/picbed/category/{id}/images
```

**请求体：**
```json
{
  "imageId": "图片ID",
  "position": 1,                    // 可选，位置序号
  "x": 0,                          // 可选，x坐标
  "y": 0                           // 可选，y坐标
}
```

##### 从分类移除图片
```http
DELETE /api/base/picbed/category/{id}/images?imageId=图片ID
```

##### 更新图片排序
```http
PUT /api/base/picbed/category/{id}/images/sort
```

**请求体：**
```json
{
  "imageIds": ["图片ID1", "图片ID2", "图片ID3"]
}
```

##### 设置分类封面
```http
PUT /api/base/picbed/category/{id}/cover
```

**请求体：**
```json
{
  "coverImageId": "封面图片ID"
}
```

#### 标签管理

##### 创建标签
```http
POST /api/base/picbed/tag
```

**请求体：**
```json
{
  "name": "标签名称",              // 必填
  "displayName": "显示名称",       // 可选
  "description": "标签描述",       // 可选
  "color": "#FF6B6B",             // 可选，颜色代码
  "category": "分类"               // 可选，标签分类
}
```

**响应：**
```json
{
  "code": 200,
  "msg": "标签创建成功",
  "data": {
    "tagId": "标签ID"
  }
}
```

##### 获取所有标签
```http
GET /api/base/picbed/tag?category=分类名称
```

**查询参数：**
- `category`（可选）：按分类筛选

**响应：**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "tags": [
      {
        "id": "标签ID",
        "name": "标签名称",
        "displayName": "显示名称",
        "description": "标签描述",
        "color": "#FF6B6B",
        "category": "分类",
        "imageCount": 10,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 15
  }
}
```

##### 获取热门标签
```http
GET /api/base/picbed/tag/popular?limit=20&category=分类名称
```

**查询参数：**
- `limit`（可选）：返回数量，默认 20
- `category`（可选）：按分类筛选

##### 搜索标签
```http
GET /api/base/picbed/tag/search?keyword=搜索关键词
```

**查询参数：**
- `keyword`（必填）：搜索关键词，支持名称和显示名称模糊搜索

##### 批量创建标签
```http
POST /api/base/picbed/tag/batch
```

**请求体：**
```json
{
  "tags": [
    {
      "name": "标签名称1",
      "displayName": "显示名称1",
      "description": "描述1",
      "color": "#FF6B6B",
      "category": "分类1"
    },
    {
      "name": "标签名称2",
      "displayName": "显示名称2",
      "description": "描述2",
      "color": "#4ECDC4",
      "category": "分类2"
    }
  ]
}
```

**响应：**
```json
{
  "code": 200,
  "msg": "批量创建标签完成",
  "data": {
    "createdTags": ["标签ID1", "标签ID2"],
    "count": 2
  }
}
```

##### 获取标签详情
```http
GET /api/base/picbed/tag/{id}
```

##### 更新标签
```http
PUT /api/base/picbed/tag/{id}
```

**请求体：**
```json
{
  "name": "新标签名称",
  "displayName": "新显示名称",
  "description": "新描述",
  "color": "#新颜色",
  "category": "新分类"
}
```

##### 删除标签
```http
DELETE /api/base/picbed/tag/{id}
```

##### 获取标签中的图片
```http
GET /api/base/picbed/tag/{id}/images?page=1&size=20
```

**响应：**
```json
{
  "code": 200,
  "msg": "获取成功",
  "data": {
    "images": [
      {
        "id": "图片ID",
        "name": "图片名称",
        "cos_url": "腾讯云COS地址",
        "width": 1920,
        "height": 1080
      }
    ],
    "total": 50,
    "page": 1,
    "pageSize": 20
  }
}
```

---

## Admin API（管理接口）

### 管理员登录
```http
POST /api/admin/login
```

**注意：** 当前管理员登录接口未实现，仅返回测试信息

**响应：**
```json
{
  "code": 200,
  "msg": "Hello",
  "data": "Hello"
}
```

**预计实现后的请求体：**
```json
{
  "email": "admin@example.com",    // 管理员邮箱
  "password": "admin_password"     // 管理员密码
}
```

### 文章管理（需要JWT认证）

#### 通过表单创建文章
```http
POST /api/admin/article
```

**请求体：** `form-data` 格式，包含以下字段：
```
title: 文章标题                    // 必填
author: 作者名称                   // 可选，留空则从上下文获取
synopsis: 文章备注/简介            // 可选
tags: 标签1,标签2,标签3           // 可选，逗号分隔
categories: 分类1,分类2           // 可选，逗号分隔
draft_flag: false                 // 可选，是否为草稿（true/false）
overhead: false                   // 可选，是否置顶（true/false）
pic_url: 封面图片链接             // 可选
markdown: 文章内容（Markdown格式） // 必填
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

**请求体：** `multipart/form-data` 格式，包含以下字段：
```
file: Markdown文件                    // 必填，文件上传
title: 文章标题                       // 必填
author: 作者名称                      // 可选，留空则从上下文获取
synopsis: 文章备注/简介               // 可选
tags: 标签1,标签2,标签3              // 可选，逗号分隔
categories: 分类1,分类2              // 可选，逗号分隔
draft_flag: false                    // 可选，是否为草稿（true/false）
overhead: false                      // 可选，是否置顶（true/false）
pic_url: 封面图片链接                // 可选
```

### 分类管理（需要JWT认证）

#### 归档文章分类
```http
POST /api/admin/category/archive
```

**注意：** 当前接口实现中未找到具体的请求结构体定义，需要根据实际情况确定参数格式

**预计的请求体：**
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

- 2024.09：重构图片管理功能，新增分类管理和标签管理，支持图片-分类倒排索引和标签同步
- 2024.09：新增直接图片上传接口，支持腾讯云COS自动上传，文件格式验证和标签同步
- 2024.01：新增图床管理功能，包括图集和图片管理
- 2023.12：增加文章点赞功能
- 2023.11：优化中文分词和搜索功能
- 2023.10：初始版本发布

---

## 注意事项

1. **JWT Token 有效期**：JWT token 默认有效期为 3000 秒（50 分钟），需要在过期前重新获取
2. **中文搜索**：文章搜索支持中文分词，会自动对标题和内容进行分词处理
3. **图片上传**：
   - 直接上传接口：支持 multipart/form-data 格式，自动上传到腾讯云COS
   - 手动创建接口：需要先将图片上传到腾讯云COS，获取URL后再创建图片记录
   - 所有新图片自动添加到 "nexus" 分类
4. **分类管理**：支持图片在多个分类中的位置管理，使用倒排索引实现高效查询
5. **标签同步**：上传图片时会自动同步标签，维护标签-图片倒排索引
6. **腾讯云COS**：图片存储在 `/somnium/primitive/` 目录下，支持自动文件格式验证
7. **分类归档**：归档操作会将分类下的所有文章设置为未分类状态
8. **权限控制**：Admin 接口需要管理员权限，普通用户无法访问
9. **管理员登录**：当前管理员登录接口未实现，仅返回测试数据
10. **请求格式**：Admin文章接口使用`form-data`格式，不是JSON格式
11. **图集更新**：更新图集时需要提供完整的`models/po.Album`结构体
12. **批量操作**：标签支持批量创建，提高操作效率