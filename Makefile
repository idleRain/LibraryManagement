.PHONY: all dev build start clean test lint help docker-up docker-down db-init db-reset db-backup

# 变量定义
GO := go
GOFLAGS := -v
BACKEND_DIR := apps/backend
FRONTEND_DIR := apps/frontend
DOCKER_DIR := docker
BINARY_NAME := library-backend

# Go 相关变量
GOCMD := $(GO)
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# 颜色输出
GREEN  := \033[0;32m
YELLOW := \033[0;33m
BLUE   := \033[0;34m
NC     := \033[0m

# 默认目标
all: help

## ============================================================================
## 开发相关
## ============================================================================

# 开发模式
dev: dev-backend dev-frontend
	@echo "$(GREEN)✓ 开发环境已启动$(NC)"

# 后端开发模式
dev-backend:
	@echo "$(BLUE)启动后端开发服务器...$(NC)"
	cd $(BACKEND_DIR) && $(GO) run cmd/server/main.go

# 前端开发模式
dev-frontend:
	@echo "$(BLUE)启动前端开发服务器...$(NC)"
	cd $(FRONTEND_DIR) && pnpm dev

## ============================================================================
## 构建相关
## ============================================================================

# 构建所有
build: build-backend build-frontend
	@echo "$(GREEN)✓ 构建完成$(NC)"

# 构建后端
build-backend:
	@echo "$(BLUE)构建后端...$(NC)"
	cd $(BACKEND_DIR) && $(GOBUILD) $(GOFLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go

# 构建前端
build-frontend:
	@echo "$(BLUE)构建前端...$(NC)"
	cd $(FRONTEND_DIR) && pnpm build

# 构建生产版本
build-prod: build-backend-prod build-frontend-prod
	@echo "$(GREEN)✓ 生产版本构建完成$(NC)"

build-backend-prod:
	@echo "$(BLUE)构建后端生产版本...$(NC)"
	cd $(BACKEND_DIR) && CGO_ENABLED=0 GOOS=linux $(GOBUILD) -ldflags="-s -w" -o bin/$(BINARY_NAME) cmd/server/main.go

build-frontend-prod:
	@echo "$(BLUE)构建前端生产版本...$(NC)"
	cd $(FRONTEND_DIR) && pnpm build

## ============================================================================
## 运行相关
## ============================================================================

# 启动服务
start: start-backend start-frontend
	@echo "$(GREEN)✓ 服务已启动$(NC)"

# 启动后端
start-backend:
	@echo "$(BLUE)启动后端服务...$(NC)"
	cd $(BACKEND_DIR) && ./bin/$(BINARY_NAME)

# 启动前端
start-frontend:
	@echo "$(BLUE)启动前端服务...$(NC)"
	cd $(FRONTEND_DIR) && pnpm preview

## ============================================================================
## 测试相关
## ============================================================================

# 运行所有测试
test: test-backend test-frontend
	@echo "$(GREEN)✓ 测试完成$(NC)"

# 后端测试
test-backend:
	@echo "$(BLUE)运行后端测试...$(NC)"
	cd $(BACKEND_DIR) && $(GOTEST) -v ./...

# 前端测试
test-frontend:
	@echo "$(BLUE)运行前端测试...$(NC)"
	cd $(FRONTEND_DIR) && pnpm test

## ============================================================================
## 代码质量
## ============================================================================

# 代码检查
lint: lint-backend lint-frontend
	@echo "$(GREEN)✓ 代码检查完成$(NC)"

# 后端代码检查
lint-backend:
	@echo "$(BLUE)检查后端代码...$(NC)"
	@which golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	cd $(BACKEND_DIR) && golangci-lint run ./...

# 前端代码检查
lint-frontend:
	@echo "$(BLUE)检查前端代码...$(NC)"
	cd $(FRONTEND_DIR) && pnpm lint

# 格式化代码
fmt: fmt-backend fmt-frontend
	@echo "$(GREEN)✓ 代码格式化完成$(NC)"

fmt-backend:
	@echo "$(BLUE)格式化后端代码...$(NC)"
	cd $(BACKEND_DIR) && $(GO) fmt ./...

fmt-frontend:
	@echo "$(BLUE)格式化前端代码...$(NC)"
	cd $(FRONTEND_DIR) && pnpm format

## ============================================================================
## Docker 相关
## ============================================================================

# 启动 Docker 服务
docker-up:
	@echo "$(BLUE)启动 Docker 服务...$(NC)"
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml up -d

# 停止 Docker 服务
docker-down:
	@echo "$(YELLOW)停止 Docker 服务...$(NC)"
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml down

# 重启 Docker 服务
docker-restart:
	@echo "$(BLUE)重启 Docker 服务...$(NC)"
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml restart

# 查看 Docker 日志
docker-logs:
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml logs -f

# 构建 Docker 镜像
docker-build:
	@echo "$(BLUE)构建 Docker 镜像...$(NC)"
	docker-compose -f $(DOCKER_DIR)/docker-compose.yml build

# 完整 Docker 部署
docker-deploy: docker-build docker-up
	@echo "$(GREEN)✓ Docker 部署完成$(NC)"

## ============================================================================
## 数据库相关
## ============================================================================

# 初始化数据库
db-init:
	@echo "$(BLUE)初始化数据库...$(NC)"
	@bash scripts/db-init.sh

# 重置数据库
db-reset:
	@echo "$(YELLOW)重置数据库...$(NC)"
	@bash scripts/db-reset.sh

# 备份数据库
db-backup:
	@echo "$(BLUE)备份数据库...$(NC)"
	@bash scripts/db-backup.sh

# 恢复数据库
db-restore:
	@echo "$(BLUE)恢复数据库...$(NC)"
	@bash scripts/db-restore.sh

## ============================================================================
## 依赖管理
## ============================================================================

# 安装所有依赖
install: install-backend install-frontend
	@echo "$(GREEN)✓ 依赖安装完成$(NC)"

install-backend:
	@echo "$(BLUE)安装后端依赖...$(NC)"
	cd $(BACKEND_DIR) && $(GOMOD) download

install-frontend:
	@echo "$(BLUE)安装前端依赖...$(NC)"
	cd $(FRONTEND_DIR) && pnpm install

# 更新所有依赖
update: update-backend update-frontend
	@echo "$(GREEN)✓ 依赖更新完成$(NC)"

update-backend:
	@echo "$(BLUE)更新后端依赖...$(NC)"
	cd $(BACKEND_DIR) && $(GOMOD) tidy && $(GOMOD) download

update-frontend:
	@echo "$(BLUE)更新前端依赖...$(NC)"
	cd $(FRONTEND_DIR) && pnpm update

## ============================================================================
## 清理相关
## ============================================================================

# 清理所有
clean: clean-backend clean-frontend
	@echo "$(GREEN)✓ 清理完成$(NC)"

# 清理后端
clean-backend:
	@echo "$(YELLOW)清理后端...$(NC)"
	cd $(BACKEND_DIR) && rm -rf bin/ && $(GOCLEAN)

# 清理前端
clean-frontend:
	@echo "$(YELLOW)清理前端...$(NC)"
	cd $(FRONTEND_DIR) && rm -rf .svelte-kit/ build/ node_modules/

## ============================================================================
## 工具命令
## ============================================================================

# 生成 Swagger 文档
swagger:
	@echo "$(BLUE)生成 Swagger 文档...$(NC)"
	@which swag > /dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	cd $(BACKEND_DIR) && swag init -g cmd/server/main.go -o docs/swagger

# 生成 Mock 数据
mock:
	@echo "$(BLUE)生成 Mock 数据...$(NC)"
	@bash scripts/generate-mock.sh

# 检查依赖安全
security-check:
	@echo "$(BLUE)检查依赖安全...$(NC)"
	cd $(BACKEND_DIR) && $(GO) list -json -m all | nancy sleuth
	cd $(FRONTEND_DIR) && pnpm audit

## ============================================================================
## 帮助信息
## ============================================================================

help:
	@echo ""
	@echo "$(GREEN)图书管理系统 - Makefile 命令帮助$(NC)"
	@echo ""
	@echo "$(YELLOW)开发命令:$(NC)"
	@echo "  make dev              启动开发环境"
	@echo "  make dev-backend      启动后端开发服务器"
	@echo "  make dev-frontend     启动前端开发服务器"
	@echo ""
	@echo "$(YELLOW)构建命令:$(NC)"
	@echo "  make build            构建所有项目"
	@echo "  make build-backend    构建后端"
	@echo "  make build-frontend   构建前端"
	@echo "  make build-prod       构建生产版本"
	@echo ""
	@echo "$(YELLOW)运行命令:$(NC)"
	@echo "  make start            启动所有服务"
	@echo "  make start-backend    启动后端服务"
	@echo "  make start-frontend   启动前端服务"
	@echo ""
	@echo "$(YELLOW)测试命令:$(NC)"
	@echo "  make test             运行所有测试"
	@echo "  make test-backend     运行后端测试"
	@echo "  make test-frontend    运行前端测试"
	@echo ""
	@echo "$(YELLOW)代码质量:$(NC)"
	@echo "  make lint             代码检查"
	@echo "  make fmt              格式化代码"
	@echo ""
	@echo "$(YELLOW)Docker 命令:$(NC)"
	@echo "  make docker-up        启动 Docker 服务"
	@echo "  make docker-down      停止 Docker 服务"
	@echo "  make docker-build     构建 Docker 镜像"
	@echo "  make docker-deploy    完整 Docker 部署"
	@echo ""
	@echo "$(YELLOW)数据库命令:$(NC)"
	@echo "  make db-init          初始化数据库"
	@echo "  make db-reset         重置数据库"
	@echo "  make db-backup        备份数据库"
	@echo "  make db-restore       恢复数据库"
	@echo ""
	@echo "$(YELLOW)依赖管理:$(NC)"
	@echo "  make install          安装所有依赖"
	@echo "  make update           更新所有依赖"
	@echo ""
	@echo "$(YELLOW)清理命令:$(NC)"
	@echo "  make clean            清理所有构建产物"
	@echo ""
	@echo "$(YELLOW)其他命令:$(NC)"
	@echo "  make swagger          生成 Swagger 文档"
	@echo "  make security-check   检查依赖安全"
	@echo "  make help             显示帮助信息"
	@echo ""
