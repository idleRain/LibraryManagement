#!/bin/bash

# =============================================================================
# 数据库初始化脚本
# =============================================================================

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DOCKER_DIR="$PROJECT_ROOT/docker"

# 默认配置
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-root123456}"
MYSQL_DATABASE="${MYSQL_DATABASE:-library_system}"

MONGODB_URI="${MONGODB_URI:-mongodb://localhost:27017}"
MONGODB_DATABASE="${MONGODB_DATABASE:-library_borrow}"

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}     图书管理系统 - 数据库初始化${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 检查 Docker 是否运行
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}✗ Docker 未运行，请先启动 Docker${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ Docker 运行正常${NC}"
}

# 检查 MySQL 连接
check_mysql() {
    echo -e "${YELLOW}检查 MySQL 连接...${NC}"
    if docker exec library-mysql mysqladmin ping -h localhost -u root -p"$MYSQL_PASSWORD" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ MySQL 连接成功${NC}"
        return 0
    else
        echo -e "${RED}✗ MySQL 连接失败${NC}"
        return 1
    fi
}

# 检查 MongoDB 连接
check_mongodb() {
    echo -e "${YELLOW}检查 MongoDB 连接...${NC}"
    if docker exec library-mongodb mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ MongoDB 连接成功${NC}"
        return 0
    else
        echo -e "${RED}✗ MongoDB 连接失败${NC}"
        return 1
    fi
}

# 初始化 MySQL
init_mysql() {
    echo -e "${YELLOW}初始化 MySQL 数据库...${NC}"
    
    # 执行初始化 SQL
    if [ -f "$DOCKER_DIR/init/mysql/01_schema.sql" ]; then
        docker exec -i library-mysql mysql -u root -p"$MYSQL_PASSWORD" < "$DOCKER_DIR/init/mysql/01_schema.sql"
        echo -e "${GREEN}✓ MySQL 表结构初始化完成${NC}"
    fi
    
    # 创建默认管理员账户
    echo -e "${YELLOW}创建默认管理员账户...${NC}"
    docker exec -i library-mysql mysql -u root -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" <<EOF
-- 创建管理员账户（密码: admin123）
INSERT IGNORE INTO users (username, password, email, real_name, status, created_at, updated_at)
VALUES ('admin', '\$2a\$10\$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.iW8jY9D6aV3xV9xKCe', 'admin@library.com', '管理员', 1, NOW(), NOW());

-- 创建默认角色
INSERT IGNORE INTO roles (name, code, description, created_at, updated_at)
VALUES 
    ('超级管理员', 'super_admin', '拥有所有权限', NOW(), NOW()),
    ('管理员', 'admin', '系统管理员', NOW(), NOW()),
    ('图书管理员', 'librarian', '图书管理', NOW(), NOW()),
    ('普通用户', 'user', '普通用户', NOW(), NOW());

-- 关联管理员角色
INSERT IGNORE INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r 
WHERE u.username = 'admin' AND r.code = 'super_admin';

-- 创建默认权限
INSERT IGNORE INTO permissions (name, code, resource, action, description, created_at, updated_at)
VALUES 
    ('用户管理', 'user:manage', 'user', 'manage', '用户管理权限', NOW(), NOW()),
    ('角色管理', 'role:manage', 'role', 'manage', '角色管理权限', NOW(), NOW()),
    ('图书管理', 'book:manage', 'book', 'manage', '图书管理权限', NOW(), NOW()),
    ('库存管理', 'stock:manage', 'stock', 'manage', '库存管理权限', NOW(), NOW()),
    ('借阅管理', 'borrow:manage', 'borrow', 'manage', '借阅管理权限', NOW(), NOW()),
    ('销售管理', 'sale:manage', 'sale', 'manage', '销售管理权限', NOW(), NOW()),
    ('采购管理', 'purchase:manage', 'purchase', 'manage', '采购管理权限', NOW(), NOW()),
    ('系统配置', 'config:manage', 'config', 'manage', '系统配置权限', NOW(), NOW());

-- 给超级管理员分配所有权限
INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p WHERE r.code = 'super_admin';
EOF
    
    echo -e "${GREEN}✓ MySQL 初始化完成${NC}"
}

# 初始化 MongoDB
init_mongodb() {
    echo -e "${YELLOW}初始化 MongoDB 数据库...${NC}"
    
    # 执行初始化 JS
    if [ -f "$DOCKER_DIR/init/mongo/01_schema.js" ]; then
        docker exec -i library-mongodb mongosh "$MONGODB_DATABASE" < "$DOCKER_DIR/init/mongo/01_schema.js"
        echo -e "${GREEN}✓ MongoDB 集合初始化完成${NC}"
    fi
    
    # 创建借阅规则
    echo -e "${YELLOW}创建借阅规则...${NC}"
    docker exec -i library-mongodb mongosh "$MONGODB_DATABASE" <<EOF
db.borrow_rules.insertMany([
    {
        name: '普通用户',
        max_books: 5,
        max_days: 30,
        max_renew_times: 2,
        fine_per_day: 0.5,
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        name: 'VIP用户',
        max_books: 10,
        max_days: 60,
        max_renew_times: 3,
        fine_per_day: 0.3,
        created_at: new Date(),
        updated_at: new Date()
    }
]);
EOF
    
    echo -e "${GREEN}✓ MongoDB 初始化完成${NC}"
}

# 主函数
main() {
    echo -e "${YELLOW}开始初始化数据库...${NC}"
    echo ""
    
    check_docker
    
    # 启动 Docker 服务（如果未运行）
    if ! docker ps | grep -q library-mysql; then
        echo -e "${YELLOW}启动 Docker 服务...${NC}"
        cd "$DOCKER_DIR" && docker-compose up -d
        sleep 10
    fi
    
    # 检查连接
    check_mysql || exit 1
    check_mongodb || exit 1
    
    echo ""
    
    # 初始化数据库
    init_mysql
    echo ""
    init_mongodb
    
    echo ""
    echo -e "${GREEN}======================================${NC}"
    echo -e "${GREEN}     数据库初始化完成！${NC}"
    echo -e "${GREEN}======================================${NC}"
    echo ""
    echo -e "默认管理员账户:"
    echo -e "  用户名: ${YELLOW}admin${NC}"
    echo -e "  密码:   ${YELLOW}admin123${NC}"
    echo ""
}

main "$@"
