#!/bin/bash

# Localtunnel tunnels script for child_bot development
# Использует фиксированные subdomain'ы для стабильных URL

set -e

FRONTEND_PORT=5173
BACKEND_PORT=8080
FRONTEND_SUBDOMAIN=childbot-dz-fe
BACKEND_SUBDOMAIN=childbotbe

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Child Bot Localtunnel Setup ===${NC}\n"

# Проверка установки localtunnel
if ! command -v lt &> /dev/null; then
    echo -e "${RED}❌ localtunnel не установлен${NC}"
    echo "Установите: npm install -g localtunnel"
    exit 1
fi

# Остановка существующих процессов
echo -e "${YELLOW}Остановка существующих туннелей...${NC}"
pkill -f "lt --port" || true
sleep 2

# Запуск frontend туннеля
echo -e "\n${GREEN}Запуск frontend туннеля...${NC}"
nohup lt --port $FRONTEND_PORT --subdomain $FRONTEND_SUBDOMAIN > /tmp/lt-frontend.log 2>&1 &
FRONTEND_PID=$!
sleep 3

# Запуск backend туннеля
echo -e "${GREEN}Запуск backend туннеля...${NC}"
nohup lt --port $BACKEND_PORT --subdomain $BACKEND_SUBDOMAIN > /tmp/lt-backend.log 2>&1 &
BACKEND_PID=$!
sleep 3

# Получение URL из логов
FRONTEND_URL=$(grep -o 'https://[^ ]*' /tmp/lt-frontend.log | head -1)
BACKEND_URL=$(grep -o 'https://[^ ]*' /tmp/lt-backend.log | head -1)

echo -e "\n${GREEN}✅ Туннели запущены успешно!${NC}\n"
echo -e "Frontend URL: ${GREEN}$FRONTEND_URL${NC}"
echo -e "Frontend PID: $FRONTEND_PID"
echo -e "\nBackend URL:  ${GREEN}$BACKEND_URL${NC}"
echo -e "Backend PID:  $BACKEND_PID"

echo -e "\n${YELLOW}📝 Настройки VK Mini App:${NC}"
echo -e "   Адрес приложения: $FRONTEND_URL"

echo -e "\n${YELLOW}ℹ️  Логи туннелей:${NC}"
echo -e "   Frontend: tail -f /tmp/lt-frontend.log"
echo -e "   Backend:  tail -f /tmp/lt-backend.log"

echo -e "\n${YELLOW}🛑 Остановить туннели:${NC}"
echo -e "   pkill -f \"lt --port\""

echo -e "\n${GREEN}Готово! Можете обновить страницу VK Mini App${NC}"