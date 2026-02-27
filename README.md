# 📚 图书管理系统 (Library Management System)

一个功能完整、现代化的图书管理系统，采用前后端分离架构，支持图书管理、库存管理、采购管理、销售管理、借阅管理等核心功能。

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)
![Svelte](https://img.shields.io/badge/SvelteKit-latest-FF3E00.svg)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1.svg)
![MongoDB](https://img.shields.io/badge/MongoDB-7.0-47A248.svg)

## ✨ 功能特性

### 📖 图书管理
- 图书信息 CRUD 操作
- ISBN 自动识别
- 分类管理
- 图书搜索与筛选
- 批量导入导出

### 📦 库存管理
- 实时库存监控
- 入库/出库操作
- 库存变动记录
- 低库存预警
- 库存盘点

### 🛒 采购管理
- 供应商管理
- 采购订单创建
- 订单审核流程
- 采购入库
- 采购统计

### 💰 销售管理
- 销售订单管理
- 购物车功能
- 多种支付方式
- 销售统计分析
- 日/月/年报表

### 📕 借阅管理
- 图书借阅/归还
- 借阅续期
- 逾期管理
- 罚款计算
- 借阅统计

### 👥 用户权限
- 用户管理
- 角色管理
- 权限控制 (RBAC)
- 操作日志

## 🏗️ 技术架构

### 前端技术栈
| 技术 | 版本 | 说明 |
|------|------|------|
| SvelteKit | latest | 前端框架 |
| TypeScript | 5.0+ | 类型支持 |
| Tailwind CSS | 3.4+ | 样式框架 |
| shadcn-svelte | latest | UI 组件库 |
| Lucide Icons | latest | 图标库 |
| svelte-sonner | latest | Toast 通知 |

### 后端技术栈
| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.21+ | 编程语言 |
| Gin | latest | Web 框架 |
| GORM | latest | ORM 框架 |
| MySQL | 8.0+ | 主数据库 |
| MongoDB | 7.0+ | 文档数据库 |
| JWT | - | 认证方案 |

## 📁 项目结构

```
.
├── apps/
│   ├── frontend/                # SvelteKit 前端应用
│   │   ├── src/
│   │   │   ├── lib/
│   │   │   │   ├── api/         # API 客户端
│   │   │   │   ├── components/  # UI 组件
│   │   │   │   │   └── ui/      # shadcn-svelte 组件
│   │   │   │   ├── stores/      # 状态管理
│   │   │   │   ├── types/       # TypeScript 类型
│   │   │   │   └── utils/       # 工具函数
│   │   │   └── routes/          # 页面路由
│   │   ├── static/              # 静态资源
│   │   ├── package.json
│   │   └── ...
│   └── backend/                 # Go 后端应用
│       ├── cmd/
│       │   └── server/          # 入口文件
│       ├── internal/
│       │   ├── config/          # 配置管理
│       │   ├── controllers/     # 控制器
│       │   ├── middleware/      # 中间件
│       │   ├── models/          # 数据模型
│       │   ├── routes/          # 路由配置
│       │   └── utils/           # 工具函数
│       ├── pkg/
│       │   └── database/        # 数据库连接
│       ├── go.mod
│       └── ...
├── docker/
│   ├── docker-compose.yml       # Docker 编排
│   └── init/
│       ├── mysql/               # MySQL 初始化
│       └── mongo/               # MongoDB 初始化
└── README.md
```

## 🚀 快速开始

### 方式一：Docker Compose（推荐）

```bash
# 进入 docker 目录
cd docker

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

服务启动后访问：
- 🌐 前端：http://localhost:3000
- 🔌 后端 API：http://localhost:8080
- 🗄️ MySQL：localhost:3306
- 🍃 MongoDB：localhost:27017
- 🔴 Redis：localhost:6379

### 方式二：本地开发

#### 环境要求
- Go 1.21+
- Node.js 20+
- pnpm 8+
- MySQL 8.0+
- MongoDB 7.0+

#### 后端启动

```bash
cd apps/backend

# 安装依赖
go mod download

# 设置环境变量
export MYSQL_HOST=localhost
export MYSQL_PORT=3306
export MYSQL_USER=root
export MYSQL_PASSWORD=your_password
export MYSQL_DATABASE=library_system
export MONGODB_URI=mongodb://localhost:27017
export MONGODB_DATABASE=library_borrow
export JWT_SECRET=your-secret-key

# 运行
go run cmd/server/main.go
```

#### 前端启动

```bash
cd apps/frontend

# 安装依赖
pnpm install

# 开发模式
pnpm dev

# 构建
pnpm build
```

## 📝 API 文档

### API 风格

所有 API 均使用 **POST** 方法，请求和响应均为 JSON 格式。

### 请求示例

```bash
# 登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 获取图书列表
curl -X POST http://localhost:8080/api/books/list \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"page":1,"page_size":10}'
```

### 响应格式

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

### 主要接口

| 接口 | 说明 |
|------|------|
| POST /api/login | 用户登录 |
| POST /api/register | 用户注册 |
| POST /api/books/list | 图书列表 |
| POST /api/books/create | 创建图书 |
| POST /api/stocks/list | 库存列表 |
| POST /api/stocks/in | 入库操作 |
| POST /api/purchases/list | 采购订单 |
| POST /api/sales/list | 销售订单 |
| POST /api/borrows/list | 借阅记录 |
| POST /api/users/list | 用户列表 |

## 🔐 默认账户

| 角色 | 用户名 | 密码 |
|------|--------|------|
| 管理员 | admin | admin123 |

## ⚙️ 环境变量

### 后端环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 服务端口 | 8080 |
| MYSQL_HOST | MySQL 主机 | localhost |
| MYSQL_PORT | MySQL 端口 | 3306 |
| MYSQL_USER | MySQL 用户 | root |
| MYSQL_PASSWORD | MySQL 密码 | - |
| MYSQL_DATABASE | MySQL 数据库 | library_system |
| MONGODB_URI | MongoDB 连接串 | mongodb://localhost:27017 |
| MONGODB_DATABASE | MongoDB 数据库 | library_borrow |
| JWT_SECRET | JWT 密钥 | - |
| JWT_EXPIRE | JWT 过期时间(秒) | 86400 |

## 📊 数据库设计

### MySQL 表结构

- `users` - 用户表
- `roles` - 角色表
- `permissions` - 权限表
- `user_roles` - 用户角色关联表
- `role_permissions` - 角色权限关联表
- `books` - 图书表
- `stocks` - 库存表
- `stock_records` - 库存变动记录表
- `suppliers` - 供应商表
- `purchase_orders` - 采购订单表
- `purchase_order_items` - 采购订单明细表
- `sale_orders` - 销售订单表
- `sale_order_items` - 销售订单明细表
- `carts` - 购物车表
- `operation_logs` - 操作日志表
- `system_configs` - 系统配置表

### MongoDB 集合

- `borrow_records` - 借阅记录
- `borrow_rules` - 借阅规则
- `reservations` - 预约记录
- `borrow_statistics` - 借阅统计

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 🙏 致谢

- [SvelteKit](https://kit.svelte.dev/)
- [shadcn-svelte](https://www.shadcn-svelte.com/)
- [Gin](https://gin-gonic.com/)
- [GORM](https://gorm.io/)

---

<p align="center">
  Made with ❤️ by Library System Team
</p>
