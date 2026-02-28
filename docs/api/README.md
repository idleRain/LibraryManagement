# API 文档

## 📋 概述

- **基础 URL**: `http://localhost:8080/api`
- **请求方式**: 所有接口使用 `POST` 方法
- **数据格式**: JSON
- **认证方式**: JWT Bearer Token

## 🔐 认证

除了登录和注册接口外，其他接口都需要在请求头中携带 Token：

```
Authorization: Bearer <token>
```

## 📝 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "参数错误",
  "data": null
}
```

### 分页响应

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 🔑 认证接口

### 登录

```
POST /api/login
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名 |
| password | string | 是 | 密码 |

**响应示例**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@library.com"
    },
    "permissions": ["user:manage", "book:manage"]
  }
}
```

### 注册

```
POST /api/register
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| username | string | 是 | 用户名（3-50字符） |
| password | string | 是 | 密码（至少6位） |
| email | string | 是 | 邮箱地址 |
| real_name | string | 否 | 真实姓名 |

### 登出

```
POST /api/logout
```

需要认证。

### 获取当前用户信息

```
POST /api/profile
```

需要认证。

---

## 📊 仪表盘接口

### 获取概览数据

```
POST /api/dashboard/overview
```

**响应示例**

```json
{
  "code": 200,
  "data": {
    "books": {
      "total": 1000,
      "active": 950,
      "inactive": 50
    },
    "users": {
      "total": 500,
      "active": 480,
      "new_today": 5
    },
    "stocks": {
      "total": 1000,
      "low_stock": 20,
      "out_of_stock": 5
    },
    "borrows": {
      "borrowed": 150,
      "overdue": 10
    },
    "sales": {
      "today_orders": 25,
      "today_amount": 3500.00
    }
  }
}
```

### 获取借阅趋势

```
POST /api/dashboard/borrow-trend
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | int | 否 | 天数，默认7天 |

### 获取销售趋势

```
POST /api/dashboard/sales-trend
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| days | int | 否 | 天数，默认7天 |

### 获取预警信息

```
POST /api/dashboard/alerts
```

### 全局搜索

```
POST /api/search
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| keyword | string | 是 | 搜索关键词 |
| limit | int | 否 | 每类结果数量，默认5 |

---

## 📚 图书接口

### 获取图书列表

```
POST /api/books/list
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认1 |
| page_size | int | 否 | 每页数量，默认10 |
| title | string | 否 | 书名（模糊搜索） |
| author | string | 否 | 作者（模糊搜索） |
| category | string | 否 | 分类 |
| status | int | 否 | 状态：1-上架，0-下架 |

### 获取图书详情

```
POST /api/books/detail
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 图书ID |

### 创建图书

```
POST /api/books/create
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| isbn | string | 是 | ISBN |
| title | string | 是 | 书名 |
| author | string | 是 | 作者 |
| publisher | string | 否 | 出版社 |
| category | string | 否 | 分类 |
| price | float | 否 | 价格 |
| cover_image | string | 否 | 封面图片URL |
| description | string | 否 | 描述 |

### 更新图书

```
POST /api/books/update
```

**请求参数**

同创建图书，额外需要 `id` 参数。

### 删除图书

```
POST /api/books/delete
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int | 是 | 图书ID |

### 获取分类列表

```
POST /api/books/categories
```

---

## 📦 库存接口

### 获取库存列表

```
POST /api/stocks/list
```

### 入库操作

```
POST /api/stocks/in
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| book_id | int | 是 | 图书ID |
| quantity | int | 是 | 数量 |
| reason | string | 否 | 原因说明 |

### 出库操作

```
POST /api/stocks/out
```

### 获取库存变动记录

```
POST /api/stocks/records
```

### 获取低库存预警

```
POST /api/stocks/low
```

---

## 📕 借阅接口

### 获取借阅记录

```
POST /api/borrows/list
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| status | string | 否 | 状态：borrowed/returned/overdue |
| user_id | int | 否 | 用户ID |

### 借阅图书

```
POST /api/borrows/borrow
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| book_id | int | 是 | 图书ID |
| days | int | 否 | 借阅天数，默认30天 |

### 归还图书

```
POST /api/borrows/return
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | string | 是 | 借阅记录ID |

### 续借

```
POST /api/borrows/renew
```

### 获取逾期记录

```
POST /api/borrows/overdue
```

### 支付罚款

```
POST /api/borrows/pay-fine
```

---

## 💰 销售接口

### 获取销售订单

```
POST /api/sales/list
```

### 创建销售订单

```
POST /api/sales/create
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| customer_name | string | 否 | 客户姓名 |
| customer_phone | string | 否 | 客户电话 |
| items | array | 是 | 订单项 |
| items[].book_id | int | 是 | 图书ID |
| items[].quantity | int | 是 | 数量 |
| items[].unit_price | float | 是 | 单价 |

### 取消订单

```
POST /api/sales/cancel
```

### 获取销售统计

```
POST /api/sales/stats
```

---

## 🛒 购物车接口

### 获取购物车

```
POST /api/cart/list
```

### 添加到购物车

```
POST /api/cart/add
```

### 从购物车移除

```
POST /api/cart/remove
```

### 清空购物车

```
POST /api/cart/clear
```

---

## 📤 文件上传接口

### 上传图片

```
POST /api/upload/image
```

Content-Type: `multipart/form-data`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 图片文件（jpg/png/gif/webp，最大5MB） |

**响应示例**

```json
{
  "code": 200,
  "data": {
    "url": "/uploads/2024/02/28/abc123.jpg",
    "filename": "cover.jpg",
    "size": 102400
  }
}
```

---

## 📥 数据导出接口

### 导出图书

```
POST /api/export/books
```

返回 Excel 文件下载。

### 导出借阅记录

```
POST /api/export/borrows
```

### 导出销售记录

```
POST /api/export/sales
```

### 导出用户

```
POST /api/export/users
```

---

## 📦 批量操作接口

### 批量导入图书

```
POST /api/batch/import-books
```

Content-Type: `multipart/form-data`

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | CSV文件 |

### 批量更新图书

```
POST /api/batch/update-books
```

**请求参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| ids | array | 是 | 图书ID数组 |
| update | object | 是 | 更新内容 |

### 批量删除图书

```
POST /api/batch/delete-books
```

### 下载导入模板

```
GET /api/batch/template?type=books
GET /api/batch/template?type=users
```

---

## ⚙️ 系统配置接口

### 获取配置

```
POST /api/config/list
```

### 更新配置

```
POST /api/config/update
```

### 获取借阅规则

```
POST /api/config/borrow-rules
```

### 更新借阅规则

```
POST /api/config/update-borrow-rules
```

---

## ❌ 错误码说明

| 错误码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证或认证过期 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |
