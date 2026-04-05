#!/bin/bash

# Альтернативный скрипт для туннеля через Cloudflare (бесплатный, без регистрации)

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Cloudflare Tunnel для разработки ===${NC}\n"

# Проверяем наличие cloudflared
if ! command -v cloudflared &> /dev/null; then
    echo -e "${RED}❌ cloudflared не установлен${NC}"
    echo -e "Установите cloudflared:"
    echo -e "  • brew install cloudflared (macOS)"
    echo -e "  • https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation"
    exit 1
fi

MODE=${1:-frontend}

case $MODE in
    frontend)
        PORT=5173
        echo -e "${GREEN}🌐 Запуск туннеля для Frontend${NC}"
        ;;
    backend)
        PORT=8080
        echo -e "${GREEN}🔧 Запуск туннеля для Backend API${NC}"
        ;;
    prod)
        PORT=80
        echo -e "${GREEN}🚀 Запуск туннеля для Production Frontend${NC}"
        ;;
    *)
        echo -e "${RED}❌ Неверный режим: $MODE${NC}"
        echo -e "Использование: $0 [frontend|backend|prod]"
        exit 1
        ;;
esac

echo -e "   Локальный адрес: ${YELLOW}http://localhost:$PORT${NC}\n"
echo -e "${GREEN}🚀 Запуск cloudflared...${NC}\n"

cloudflared tunnel --url http://localhost:$PORT