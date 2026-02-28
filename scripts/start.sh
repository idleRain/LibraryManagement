#!/bin/bash

# =============================================================================
# 一键启动脚本
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 显示 Banner
show_banner() {
    echo ""
    echo -e "${GREEN}"
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║                                                               ║"
    echo "║     📚 图书管理系统 (Library Management System)               ║"
    echo "║                                                               ║"
    echo "║     Version: 1.2.0                                            ║"
    echo "║     Author:  idleRain                                         ║"
    echo "║                                                               ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo ""
}

# 检查依赖
check_dependencies() {
    echo -e "${YELLOW}检查系统依赖...${NC}"
    
    local missing=()
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        missing+=("docker")
    fi
    
    # 检查 Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        missing+=("docker-compose")
    fi
    
    # 检查 Go
    if ! command -v go &> /dev/null; then
        missing+=("go")
    fi
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        missing+=("node")
    fi
    
    # 检查 pnpm
    if ! command -v pnpm &> /dev/null; then
        missing+=("pnpm")
    fi
    
    if [ ${#missing[@]} -gt 0 ]; then
        echo -e "${RED}✗ 缺少以下依赖:${NC}"
        for dep in "${missing[@]}"; do
            echo -e "  - $dep"
        done
        echo ""
        echo -e "${YELLOW}请先安装缺少的依赖${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✓ 所有依赖已安装${NC}"
    echo ""
}

# 检查端口
check_ports() {
    echo -e "${YELLOW}检查端口占用...${NC}"
    
    local ports=(3000 8080 3306 27017 6379)
    local occupied=()
    
    for port in "${ports[@]}"; do
        if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1; then
            occupied+=($port)
        fi
    done
    
    if [ ${#occupied[@]} -gt 0 ]; then
        echo -e "${YELLOW}以下端口已被占用:${NC}"
        for port in "${occupied[@]}"; do
            local service=""
            case $port in
                3000) service="前端" ;;
                8080) service="后端" ;;
                3306) service="MySQL" ;;
                27017) service="MongoDB" ;;
                6379) service="Redis" ;;
            esac
            echo -e "  - $port ($service)"
        done
        echo ""
        read -p "是否继续？(y/n): " continue
        if [ "$continue" != "y" ]; then
            exit 0
        fi
    else
        echo -e "${GREEN}✓ 端口检查通过${NC}"
    fi
    echo ""
}

# 启动 Docker 服务
start_docker() {
    echo -e "${YELLOW}启动 Docker 服务...${NC}"
    
    cd "$PROJECT_ROOT/docker"
    
    # 检查是否已运行
    if docker-compose ps | grep -q "Up"; then
        echo -e "${GREEN}✓ Docker 服务已在运行${NC}"
    else
        docker-compose up -d
        echo -e "${GREEN}✓ Docker 服务已启动${NC}"
    fi
    
    echo ""
}

# 等待数据库就绪
wait_for_db() {
    echo -e "${YELLOW}等待数据库就绪...${NC}"
    
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if docker exec library-mysql mysqladmin ping -h localhost -u root -proot123456 >/dev/null 2>&1; then
            echo -e "${GREEN}✓ MySQL 就绪${NC}"
            break
        fi
        attempt=$((attempt + 1))
        sleep 1
    done
    
    if [ $attempt -eq $max_attempts ]; then
        echo -e "${RED}✗ MySQL 启动超时${NC}"
        exit 1
    fi
    
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if docker exec library-mongodb mongosh --eval "db.adminCommand('ping')" >/dev/null 2>&1; then
            echo -e "${GREEN}✓ MongoDB 就绪${NC}"
            break
        fi
        attempt=$((attempt + 1))
        sleep 1
    done
    
    if [ $attempt -eq $max_attempts ]; then
        echo -e "${RED}✗ MongoDB 启动超时${NC}"
        exit 1
    fi
    
    echo ""
}

# 初始化数据库（首次运行）
init_database() {
    echo -e "${YELLOW}检查数据库初始化状态...${NC}"
    
    # 检查是否已初始化
    local tables=$(docker exec library-mysql mysql -u root -proot123456 -N -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'library_system'" 2>/dev/null || echo "0")
    
    if [ "$tables" -eq 0 ]; then
        echo -e "${YELLOW}首次运行，初始化数据库...${NC}"
        bash "$SCRIPT_DIR/db-init.sh"
    else
        echo -e "${GREEN}✓ 数据库已初始化${NC}"
    fi
    
    echo ""
}

# 安装依赖
install_dependencies() {
    echo -e "${YELLOW}安装依赖...${NC}"
    
    # 前端依赖
    if [ ! -d "$PROJECT_ROOT/apps/frontend/node_modules" ]; then
        echo -e "${BLUE}安装前端依赖...${NC}"
        cd "$PROJECT_ROOT/apps/frontend" && pnpm install
    else
        echo -e "${GREEN}✓ 前端依赖已安装${NC}"
    fi
    
    # 后端依赖
    echo -e "${BLUE}检查后端依赖...${NC}"
    cd "$PROJECT_ROOT/apps/backend" && go mod download
    echo -e "${GREEN}✓ 后端依赖已安装${NC}"
    
    echo ""
}

# 启动开发服务
start_dev() {
    echo -e "${YELLOW}启动开发服务...${NC}"
    echo ""
    echo -e "${GREEN}======================================${NC}"
    echo -e "${GREEN}  服务已启动！${NC}"
    echo -e "${GREEN}======================================${NC}"
    echo ""
    echo -e "  前端地址: ${BLUE}http://localhost:3000${NC}"
    echo -e "  后端地址: ${BLUE}http://localhost:8080${NC}"
    echo -e "  API 文档: ${BLUE}http://localhost:8080/swagger${NC}"
    echo ""
    echo -e "  默认账户: ${YELLOW}admin / admin123${NC}"
    echo ""
    echo -e "${YELLOW}按 Ctrl+C 停止服务${NC}"
    echo ""
    
    # 启动后端（后台）
    cd "$PROJECT_ROOT/apps/backend"
    go run cmd/server/main.go &
    BACKEND_PID=$!
    
    # 启动前端（前台）
    cd "$PROJECT_ROOT/apps/frontend"
    pnpm dev
    
    # 清理
    trap "kill $BACKEND_PID 2>/dev/null" EXIT
}

# 主函数
main() {
    show_banner
    
    # 解析参数
    case "${1:-dev}" in
        dev)
            check_dependencies
            check_ports
            start_docker
            wait_for_db
            init_database
            install_dependencies
            start_dev
            ;;
        docker)
            start_docker
            wait_for_db
            init_database
            echo -e "${GREEN}✓ Docker 服务已启动${NC}"
            ;;
        init)
            start_docker
            wait_for_db
            init_database
            echo -e "${GREEN}✓ 数据库初始化完成${NC}"
            ;;
        install)
            install_dependencies
            echo -e "${GREEN}✓ 依赖安装完成${NC}"
            ;;
        *)
            echo "用法: $0 {dev|docker|init|install}"
            exit 1
            ;;
    esac
}

main "$@"
