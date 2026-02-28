#!/bin/bash

# =============================================================================
# 数据库重置脚本
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${RED}警告: 此操作将删除所有数据！${NC}"
read -p "确定要重置数据库吗？(yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo -e "${YELLOW}操作已取消${NC}"
    exit 0
fi

echo -e "${YELLOW}重置数据库...${NC}"

# 停止并删除容器
cd "$PROJECT_ROOT/docker"
docker-compose down -v

# 重新启动
docker-compose up -d

# 等待服务启动
sleep 15

# 重新初始化
cd "$PROJECT_ROOT"
bash scripts/db-init.sh

echo -e "${GREEN}✓ 数据库重置完成${NC}"
