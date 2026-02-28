#!/bin/bash

# =============================================================================
# 数据库恢复脚本
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="$PROJECT_ROOT/backups"

# 配置
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-root123456}"
MYSQL_DATABASE="${MYSQL_DATABASE:-library_system}"

MONGODB_DATABASE="${MONGODB_DATABASE:-library_borrow}"

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}     数据库恢复${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 列出可用备份
list_backups() {
    echo -e "${YELLOW}可用的备份文件:${NC}"
    echo ""
    echo "MySQL 备份:"
    ls -lht "$BACKUP_DIR"/mysql_*.sql.gz 2>/dev/null | head -5 || echo "  无备份文件"
    echo ""
    echo "MongoDB 备份:"
    ls -lht "$BACKUP_DIR"/mongodb_*.tar.gz 2>/dev/null | head -5 || echo "  无备份文件"
    echo ""
}

# 恢复 MySQL
restore_mysql() {
    local backup_file="$1"
    
    if [ ! -f "$backup_file" ]; then
        echo -e "${RED}✗ 备份文件不存在: $backup_file${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}恢复 MySQL 数据库...${NC}"
    
    # 解压并恢复
    zcat "$backup_file" | docker exec -i library-mysql mysql \
        -h "$MYSQL_HOST" \
        -P "$MYSQL_PORT" \
        -u "$MYSQL_USER" \
        -p"$MYSQL_PASSWORD" \
        "$MYSQL_DATABASE"
    
    echo -e "${GREEN}✓ MySQL 恢复完成${NC}"
}

# 恢复 MongoDB
restore_mongodb() {
    local backup_file="$1"
    
    if [ ! -f "$backup_file" ]; then
        echo -e "${RED}✗ 备份文件不存在: $backup_file${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}恢复 MongoDB 数据库...${NC}"
    
    # 解压
    local temp_dir=$(mktemp -d)
    tar -xzf "$backup_file" -C "$temp_dir"
    
    # 恢复
    docker cp "$temp_dir"/* library-mongodb:/tmp/restore
    docker exec library-mongodb mongorestore \
        --db "$MONGODB_DATABASE" \
        "/tmp/restore/$MONGODB_DATABASE"
    
    # 清理
    docker exec library-mongodb rm -rf /tmp/restore
    rm -rf "$temp_dir"
    
    echo -e "${GREEN}✓ MongoDB 恢复完成${NC}"
}

# 主函数
main() {
    list_backups
    
    # 选择 MySQL 备份
    echo -e "${YELLOW}请输入 MySQL 备份文件路径（或按 Enter 跳过）:${NC}"
    read -r mysql_backup
    
    if [ -n "$mysql_backup" ]; then
        restore_mysql "$mysql_backup"
    fi
    
    echo ""
    
    # 选择 MongoDB 备份
    echo -e "${YELLOW}请输入 MongoDB 备份文件路径（或按 Enter 跳过）:${NC}"
    read -r mongodb_backup
    
    if [ -n "$mongodb_backup" ]; then
        restore_mongodb "$mongodb_backup"
    fi
    
    echo ""
    echo -e "${GREEN}======================================${NC}"
    echo -e "${GREEN}     恢复完成！${NC}"
    echo -e "${GREEN}======================================${NC}"
}

main "$@"
