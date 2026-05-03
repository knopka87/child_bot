#!/bin/bash
# Скрипт для проверки подключения к PostgreSQL

set -e

echo "🔍 Проверка подключения к PostgreSQL..."
echo ""

# Загружаем переменные из .env
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "❌ Файл .env не найден"
    exit 1
fi

# Проверяем что контейнер запущен
echo "1️⃣ Проверка статуса контейнера..."
if docker ps | grep -q child_bot_postgres; then
    echo "✅ Контейнер PostgreSQL запущен"
else
    echo "❌ Контейнер PostgreSQL не запущен"
    echo "   Запустите: docker-compose up -d postgres"
    exit 1
fi

echo ""
echo "2️⃣ Проверка подключения к БД..."

# Пробуем подключиться
if docker exec child_bot_postgres psql -U "${POSTGRES_USER:-child_bot}" -d "${POSTGRES_DB:-child_bot}" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ Подключение успешно"
else
    echo "❌ Не удалось подключиться к БД"
    echo "   Проверьте логи: docker logs child_bot_postgres"
    exit 1
fi

echo ""
echo "3️⃣ Информация о базе данных..."

# Версия PostgreSQL
PG_VERSION=$(docker exec child_bot_postgres psql -U "${POSTGRES_USER:-child_bot}" -d "${POSTGRES_DB:-child_bot}" -t -c "SELECT version();" | head -n 1 | xargs)
echo "   PostgreSQL: $PG_VERSION"

# Размер базы данных
DB_SIZE=$(docker exec child_bot_postgres psql -U "${POSTGRES_USER:-child_bot}" -d "${POSTGRES_DB:-child_bot}" -t -c "SELECT pg_size_pretty(pg_database_size('${POSTGRES_DB:-child_bot}'));" | xargs)
echo "   Размер БД: $DB_SIZE"

# Количество таблиц
TABLE_COUNT=$(docker exec child_bot_postgres psql -U "${POSTGRES_USER:-child_bot}" -d "${POSTGRES_DB:-child_bot}" -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" | xargs)
echo "   Количество таблиц: $TABLE_COUNT"

echo ""
echo "4️⃣ Параметры подключения для IDE:"
echo ""
echo "   Host:     localhost (или 127.0.0.1)"
echo "   Port:     ${POSTGRES_PORT:-5432}"
echo "   Database: ${POSTGRES_DB:-child_bot}"
echo "   User:     ${POSTGRES_USER:-child_bot}"
echo "   Password: ${POSTGRES_PASSWORD}"
echo ""

echo "5️⃣ Статистика основных таблиц:"
echo ""

docker exec child_bot_postgres psql -U "${POSTGRES_USER:-child_bot}" -d "${POSTGRES_DB:-child_bot}" -c "
SELECT
    'child_profiles' as table_name,
    COUNT(*) as records
FROM child_profiles
UNION ALL
SELECT 'parent_profiles', COUNT(*) FROM parent_profiles
UNION ALL
SELECT 'attempts', COUNT(*) FROM attempts
UNION ALL
SELECT 'achievements', COUNT(*) FROM achievements
UNION ALL
SELECT 'child_achievements', COUNT(*) FROM child_achievements
ORDER BY table_name;
"

echo ""
echo "✅ Всё работает!"
echo ""
echo "📖 Подробная документация: docs/DATABASE_CONNECTION.md"
echo "🚀 Быстрый старт: docs/QUICK_DB_SETUP.md"
echo ""
echo "Для прямого подключения к БД выполните:"
echo "   docker exec -it child_bot_postgres psql -U ${POSTGRES_USER:-child_bot} -d ${POSTGRES_DB:-child_bot}"
echo ""
