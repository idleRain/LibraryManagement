#!/bin/bash

# =============================================================================
# 数据库备份脚本
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
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# 配置
MYSQL_HOST="${MYSQL_HOST:-localhost}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD:-root123456}"
MYSQL_DATABASE="${MYSQL_DATABASE:-library_system}"

MONGODB_DATABASE="${MONGODB_DATABASE:-library_borrow}"

# 创建备份目录
mkdir -p "$BACKUP_DIR"

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}     数据库备份${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

# 备份 MySQL
backup_mysql() {
    echo -e "${YELLOW}备份 MySQL 数据库...${NC}"
    
    local backup_file="$BACKUP_DIR/mysql_$TIMESTAMP.sql"
    
    docker exec library-mysql mysqldump \
        -h "$MYSQL_HOST" \
        -P "$MYSQL_PORT" \
        -u "$MYSQL_USER" \
        -p"$MYSQL_PASSWORD" \
        --single-transaction \
        --routines \
        --triggers \
        "$MYSQL_DATABASE" > "$backup_file"
    
    # 压缩
    gzip "$backup_file"
    
    echo -e "${GREEN}✓ MySQL 备份完成: $backup_file.gz${NC}"
}

# 备份 MongoDB
backup_mongodb() {
    echo -e "${YELLOW}备份 MongoDB 数据库...${NC}"
    
    local backup_path="$BACKUP_DIR/mongodb_$TIMESTAMP"
    
    docker exec library-mongodb mongodump \
        --db "$MONGODB_DATABASE" \
        --out "/tmp/backup"
    
    docker cp library-mongodb:/tmp/backup "$backup_path"
    
    docker exec library-mongodb rm -rf /tmp/backup
    
    # 压缩
    tar -czf "$backup_path.tar.gz" -C "$BACKUP_DIR" "$(basename "$backup_path")"
    rm -rf "$backup_path"
    
    echo -e "${GREEN}✓ MongoDB 备份完成: $backup_path.tar.gz${NC}"
}

# 清理旧备份（保留最近7天）
cleanup_old_backups() {
    echo -e "${YELLOW}清理旧备份文件...${NC}"
    find "$BACKUP_DIR" -type f -mtime +7 -delete
    echo -e "${GREEN}✓ 清理完成${NC}"
}

# 主函数
main() {
    backup_mysql
    echo ""
    backup_mongodb
    echo ""
    cleanup_old_backups
    
    echo ""
    echo -e "${GREEN}======================================${NC}"
    echo -e "${GREEN}     备份完成！${NC}"
    echo -e "${GREEN}======================================${NC}"
    echo ""
    echo -e "备份文件位置: ${YELLOW}$BACKUP_DIR${NC}"
    ls -lh "$BACKUP_DIR" | tail -5
}

main "$@"
