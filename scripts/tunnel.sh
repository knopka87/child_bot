#!/bin/bash

# Скрипт для запуска локального туннеля для разработки

set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Локальный туннель для разработки ===${NC}\n"

# Проверяем наличие ngrok
if ! command -v ngrok &> /dev/null; then
    echo -e "${RED}❌ ngrok не установлен${NC}"
    echo -e "Установите ngrok одним из способов:"
    echo -e "  • brew install ngrok (macOS)"
    echo -e "  • https://ngrok.com/download (все платформы)"
    exit 1
fi

# Проверяем авторизацию ngrok
if ! ngrok config check &> /dev/null; then
    echo -e "${YELLOW}⚠️  ngrok не настроен${NC}"
    echo -e "1. Зарегистрируйтесь на https://dashboard.ngrok.com/signup"
    echo -e "2. Получите токен: https://dashboard.ngrok.com/get-started/your-authtoken"
    echo -e "3. Выполните: ngrok config add-authtoken <ваш-токен>"
    exit 1
fi

# Определяем что туннелировать
MODE=${1:-frontend}

case $MODE in
    frontend)
        PORT=5173
        echo -e "${GREEN}🌐 Запуск туннеля для Frontend (Vite dev server)${NC}"
        echo -e "   Локальный адрес: ${YELLOW}http://localhost:$PORT${NC}\n"
        ;;
    backend)
        PORT=8080
        echo -e "${GREEN}🔧 Запуск туннеля для Backend API${NC}"
        echo -e "   Локальный адрес: ${YELLOW}http://localhost:$PORT${NC}\n"
        ;;
    prod)
        PORT=80
        echo -e "${GREEN}🚀 Запуск туннеля для Production Frontend${NC}"
        echo -e "   Локальный адрес: ${YELLOW}http://localhost:$PORT${NC}\n"
        ;;
    *)
        echo -e "${RED}❌ Неверный режим: $MODE${NC}"
        echo -e "Использование: $0 [frontend|backend|prod]"
        echo -e "  frontend - туннель для Vite dev server (порт 5173)"
        echo -e "  backend  - туннель для Backend API (порт 8080)"
        echo -e "  prod     - туннель для Production Frontend (порт 80)"
        exit 1
        ;;
esac

echo -e "${YELLOW}📋 Следующие шаги после запуска:${NC}"
echo -e "1. Скопируйте HTTPS URL (https://xxxx.ngrok.io)"
echo -e "2. Добавьте URL в настройки VK Mini App"
echo -e "3. Обновите ALLOWED_ORIGINS в .env файле"
echo -e "4. Для frontend: обновите VITE_API_BASE_URL если нужно\n"

echo -e "${GREEN}🚀 Запуск ngrok...${NC}\n"
ngrok http $PORT --log=stdout