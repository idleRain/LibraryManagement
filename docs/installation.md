# 安装指南

本文档将指导您完成图书管理系统的安装。

## 📋 系统要求

### 必需软件

| 软件 | 版本要求 | 说明 |
|------|----------|------|
| Docker | 20.0+ | 容器运行环境 |
| Docker Compose | 2.0+ | 容器编排工具 |

### 开发环境额外要求

| 软件 | 版本要求 | 说明 |
|------|----------|------|
| Go | 1.21+ | 后端开发语言 |
| Node.js | 20.0+ | 前端运行环境 |
| pnpm | 8.0+ | 包管理器 |
| Make | - | 构建工具 |

## 🚀 方式一：Docker 部署（推荐）

### 1. 克隆项目

```bash
git clone https://github.com/idleRain/LibraryManagement.git
cd LibraryManagement
```

### 2. 启动服务

```bash
# 使用一键脚本
bash scripts/start.sh docker

# 或使用 Makefile
make docker-up

# 或使用 pnpm
pnpm docker:up
```

### 3. 初始化数据库

```bash
bash scripts/db-init.sh
```

### 4. 访问系统

- 前端：http://localhost:3000
- 后端：http://localhost:8080
- 默认账户：admin / admin123

## 💻 方式二：本地开发

### 1. 安装依赖

```bash
# 安装前端依赖
cd apps/frontend
pnpm install

# 安装后端依赖
cd ../backend
go mod download
```

### 2. 启动数据库

```bash
# 启动 Docker 数据库服务
make docker-up

# 初始化数据库
make db-init
```

### 3. 配置环境变量

创建 `apps/backend/.env` 文件：

```env
# 服务配置
SERVER_PORT=8080

# MySQL 配置
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=root123456
MYSQL_DATABASE=library_system

# MongoDB 配置
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=library_borrow

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_ENABLED=true

# JWT 配置
JWT_SECRET=your-secret-key
JWT_EXPIRE=86400
```

### 4. 启动开发服务

```bash
# 使用一键脚本
bash scripts/start.sh dev

# 或分别启动
make dev-backend  # 终端1
make dev-frontend # 终端2
```

## 🔧 验证安装

### 检查服务状态

```bash
# 检查 Docker 服务
docker-compose -f docker/docker-compose.yml ps

# 检查后端健康
curl http://localhost:8080/health

# 检查前端
curl http://localhost:3000
```

### 检查数据库连接

```bash
# MySQL
docker exec -it library-mysql mysql -u root -p

# MongoDB
docker exec -it library-mongodb mongosh

# Redis
docker exec -it library-redis redis-cli ping
```

## 📦 目录权限

确保以下目录有正确的写入权限：

```bash
chmod -R 755 uploads/
chmod -R 755 backups/
```

## 🔄 更新系统

```bash
# 拉取最新代码
git pull origin master

# 更新依赖
make update

# 重启服务
make docker-restart
```

## ❓ 常见问题

### 端口被占用

```bash
# 查看端口占用
lsof -i :3000
lsof -i :8080

# 修改 docker-compose.yml 中的端口映射
```

### Docker 服务启动失败

```bash
# 查看日志
docker-compose -f docker/docker-compose.yml logs

# 重置 Docker
docker-compose -f docker/docker-compose.yml down -v
docker-compose -f docker/docker-compose.yml up -d
```

### 数据库连接失败

1. 确认 Docker 服务已启动
2. 检查环境变量配置
3. 等待数据库完全启动（约 30 秒）

## 📚 下一步

- [快速开始](./quick-start.md) - 了解基本功能
- [配置说明](./configuration.md) - 详细配置选项
- [API 文档](./api/README.md) - 接口文档
