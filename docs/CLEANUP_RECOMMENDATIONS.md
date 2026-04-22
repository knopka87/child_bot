# Рекомендации по очистке кодовой базы и БД

**Дата:** 2026-04-04

## 🗑️ Файлы для удаления

### Store файлы (внутри `api/internal/store/`)

#### ❌ Полностью неиспользуемые:
```bash
rm internal/store/histrory.go      # Опечатка в названии, работа с timeline_events (Telegram v1)
rm internal/store/metrics.go       # Старая система метрик (Telegram v1)
rm internal/store/textbook_search.go  # Поиск по учебникам (не подключен к API)
```

**Экономия:** ~500 строк кода

---

#### ⚠️ Используются минимально (можно удалить):
```bash
# session.go - 1 использование (только в другом deprecated коде)
rm internal/store/session.go

# chat.go - 5 использований (только в deprecated v2/telegram)
rm internal/store/chat.go

# hint.go - 70 использований в deprecated коде
rm internal/store/hint.go

# parse.go - 36 использований в deprecated коде
rm internal/store/parse.go

# user.go - 16 использований в deprecated коде
rm internal/store/user.go
```

**Экономия:** ~2000 строк кода

---

### Telegram bot код (внутри `api/internal/`)

#### 🔴 Полностью deprecated директории:
```bash
# Telegram bot v1
rm -rf internal/v1/

# Telegram bot v2
rm -rf internal/v2/
```

**Экономия:** ~10,000+ строк кода

---

## 🗄️ Таблицы БД для удаления

### Создать backup перед удалением:
```bash
cd api
export $(grep -v '^#' .env | grep -v '^Helper' | xargs)

# Backup старых таблиц Telegram бота
pg_dump \
  -t parsed_tasks \
  -t hints_cache \
  -t metrics_events \
  -t timeline_events \
  -t task_sessions \
  -t chat \
  -t user \
  "$DATABASE_URL" > ../backups/telegram_bot_archive_$(date +%Y%m%d).sql
```

### Создать миграцию для удаления:
```sql
-- api/migrations/036_cleanup_telegram_tables.up.sql
-- Удаление таблиц Telegram бота v1 и v2

DROP TABLE IF EXISTS hints_cache CASCADE;
DROP TABLE IF EXISTS parsed_tasks CASCADE;
DROP TABLE IF EXISTS metrics_events CASCADE;
DROP TABLE IF EXISTS timeline_events CASCADE;
DROP TABLE IF EXISTS task_sessions CASCADE;
DROP TABLE IF EXISTS chat CASCADE;
DROP TABLE IF EXISTS user CASCADE;

COMMENT ON DATABASE postgres IS 'Cleaned up deprecated Telegram bot tables';
```

```sql
-- api/migrations/036_cleanup_telegram_tables.down.sql
-- Восстановление из backup'а

-- Для rollback нужно восстановить из архива:
-- psql $DATABASE_URL < backups/telegram_bot_archive_YYYYMMDD.sql
```

**Экономия места:** ~200 kB (таблицы почти пустые)

---

## 📊 Анализ использования кода

### Telegram bot код (v1 и v2):
```
❌ internal/v1/telegram/     - НЕ используется в REST API
❌ internal/v1/types/         - НЕ используется
❌ internal/v1/llmclient/     - заменён на internal/llm/
❌ internal/v2/telegram/      - НЕ используется в REST API
❌ internal/v2/types/         - НЕ используется
❌ internal/v2/llmclient/     - заменён на internal/llm/
❌ internal/v2/templates/     - JSON шаблоны для Telegram бота
```

### Store файлы:
```
✅ attempt.go         - 323 использования ✅ АКТИВНЫЙ
✅ villain.go         - 108 использований ✅ АКТИВНЫЙ
✅ store.go           - 102 использования ✅ АКТИВНЫЙ
✅ achievement.go     - 52 использования  ✅ АКТИВНЫЙ
✅ referral.go        - 32 использования  ✅ АКТИВНЫЙ
✅ consent.go         - 25 использований  ✅ АКТИВНЫЙ
✅ email.go           - 27 использований  ✅ АКТИВНЫЙ
✅ legal.go           - 7 использований   ✅ АКТИВНЫЙ

❌ hint.go            - 70 в deprecated коде
❌ parse.go           - 36 в deprecated коде
❌ user.go            - 16 в deprecated коде
❌ chat.go            - 5 в deprecated коде
❌ session.go         - 1 в deprecated коде
❌ histrory.go        - 0 использований
❌ metrics.go         - 0 использований
❌ textbook_search.go - 0 использований
```

---

## 🎯 План поэтапной очистки

### Этап 1: Безопасное удаление (без влияния на работу)
```bash
# 1. Удаляем файлы с 0 использований
rm internal/store/histrory.go
rm internal/store/metrics.go
rm internal/store/textbook_search.go

# 2. Удаляем Telegram bot код v1 и v2
rm -rf internal/v1/
rm -rf internal/v2/

# 3. Коммит
git add -A
git commit -m "cleanup: remove deprecated Telegram bot v1/v2 code"
```

**Риск:** ⭕ Нулевой (код не используется)

---

### Этап 2: Удаление таблиц БД
```bash
# 1. Создать backup
mkdir -p backups
cd api
export $(grep -v '^#' .env | grep -v '^Helper' | xargs)
pg_dump -t parsed_tasks -t hints_cache -t metrics_events \
  -t timeline_events -t task_sessions -t chat -t user \
  "$DATABASE_URL" > ../backups/telegram_tables_$(date +%Y%m%d).sql

# 2. Создать миграцию
cat > migrations/036_cleanup_telegram_tables.up.sql << 'EOF'
DROP TABLE IF EXISTS hints_cache CASCADE;
DROP TABLE IF EXISTS parsed_tasks CASCADE;
DROP TABLE IF EXISTS metrics_events CASCADE;
DROP TABLE IF EXISTS timeline_events CASCADE;
DROP TABLE IF EXISTS task_sessions CASCADE;
DROP TABLE IF EXISTS chat CASCADE;
DROP TABLE IF EXISTS user CASCADE;
EOF

cat > migrations/036_cleanup_telegram_tables.down.sql << 'EOF'
-- Rollback: restore from backup
-- psql $DATABASE_URL < backups/telegram_tables_YYYYMMDD.sql
EOF

# 3. Применить миграцию
make migrate-up

# 4. Коммит
git add migrations/
git commit -m "db: remove deprecated Telegram bot tables"
```

**Риск:** ⚠️ Низкий (таблицы не используются в REST API)

---

### Этап 3: Удаление store файлов для Telegram
```bash
# После успешного Этапа 2 (удаление таблиц)
rm internal/store/session.go
rm internal/store/chat.go
rm internal/store/hint.go
rm internal/store/parse.go
rm internal/store/user.go

# Коммит
git add -A
git commit -m "cleanup: remove deprecated store files"
```

**Риск:** ⚠️ Низкий (код связан с удалёнными таблицами)

---

## 📈 Итоговая экономия

### Код:
- **~12,500 строк кода** удалено
- **~8 файлов** store удалено
- **2 директории** (v1/, v2/) удалено

### База данных:
- **7 таблиц** удалено
- **~200 kB** места освобождено

### Снижение сложности:
- Меньше файлов для поддержки
- Чище структура БД
- Понятнее архитектура (только REST API)

---

## ⚠️ Важные замечания

1. **НЕ удалять без backup!** Всегда создавать архив БД перед удалением таблиц
2. **Проверить зависимости:** Убедиться что код компилируется после удаления файлов
3. **Поэтапность:** Не удалять всё сразу, делать коммиты после каждого этапа
4. **Тестирование:** После каждого этапа проверять что REST API работает

---

## ✅ Проверка перед удалением

```bash
# 1. Убедиться что сервер компилируется
go build -o server ./cmd/server

# 2. Убедиться что тесты проходят
go test ./...

# 3. Убедиться что нет активных ссылок на удаляемые таблицы
grep -r "parsed_tasks\|hints_cache\|metrics_events\|timeline_events\|task_sessions" \
  --include="*.go" internal/api/ internal/service/

# Если вывод пустой - можно удалять
```

---

## 🚀 Скрипт автоматической очистки

```bash
#!/bin/bash
# cleanup.sh - Автоматическая очистка deprecated кода

set -e

echo "🧹 Starting cleanup..."

# Этап 1: Удаление неиспользуемых файлов
echo "📁 Removing unused store files..."
rm -f internal/store/histrory.go
rm -f internal/store/metrics.go
rm -f internal/store/textbook_search.go

# Этап 2: Удаление Telegram bot кода
echo "📁 Removing Telegram bot v1/v2..."
rm -rf internal/v1/
rm -rf internal/v2/

# Проверка компиляции
echo "🔨 Building server..."
go build -o server ./cmd/server

echo "✅ Cleanup completed successfully!"
echo "📝 Don't forget to:"
echo "   1. Create DB backup before dropping tables"
echo "   2. Create migration 036_cleanup_telegram_tables.up.sql"
echo "   3. Commit changes"
```

**Использование:**
```bash
cd api
chmod +x cleanup.sh
./cleanup.sh
```

---

**Итого:** Безопасная очистка позволит удалить ~12,500 строк неиспользуемого кода и 7 устаревших таблиц БД.
