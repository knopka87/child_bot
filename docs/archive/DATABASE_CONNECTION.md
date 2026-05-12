# Подключение к базе данных

## Обзор

PostgreSQL уже настроен для внешних подключений через `docker-compose.yml`. Порт 5432 пробрасывается на локальный хост.

## Параметры подключения

### Локальная разработка (Docker)

```
Host:     localhost (или 127.0.0.1)
Port:     5432 (или значение из .env: POSTGRES_PORT)
Database: child_bot (или значение из .env: POSTGRES_DB)
User:     child_bot (или значение из .env: POSTGRES_USER)
Password: [из .env: POSTGRES_PASSWORD]
```

### Production (удалённый сервер)

⚠️ **Важно:** Для продакшена используйте SSH туннель (см. раздел "Безопасное подключение к production")

```
Host:     [адрес вашего сервера]
Port:     5432
Database: child_bot
User:     child_bot
Password: [из .env.production: POSTGRES_PASSWORD]
```

## Подключение из разных IDE

### 1. GoLand / DataGrip / IntelliJ IDEA

#### Локальное подключение (Docker на локальной машине)

**Шаг 1: Создать новое подключение**
1. Откройте Database panel (View → Tool Windows → Database)
2. Нажмите `+` → Data Source → PostgreSQL

**Шаг 2: Настроить параметры на вкладке General**
```
Host: localhost
Port: 5432
Database: child_bot
User: child_bot
Password: [из .env]
```

**Шаг 3: Тест подключения**
- Нажмите "Test Connection"
- Если нужно, скачайте драйвер PostgreSQL
- Если успешно → нажмите "OK"

#### Production подключение через SSH туннель

**Шаг 1: Создать новое подключение PostgreSQL**

**Шаг 2: Настроить SSH туннель на вкладке SSH/SSL**
1. Перейдите на вкладку "SSH/SSL"
2. Поставьте галочку "Use SSH tunnel"
3. Нажмите `+` для создания нового SSH конфига
4. Заполните параметры SSH:
   ```
   Host: 77.222.60.149
   Port: 22
   User name: root
   Auth type: Key pair (OpenSSH or PuTTY)
   Private key file: /Users/a.yanover/Downloads/id_rsa_1/id_rsa
   ```
5. Нажмите "Test Connection" для проверки SSH
6. Если успешно → "OK"

**Шаг 3: Настроить PostgreSQL на вкладке General**
```
Host: localhost (важно! не IP сервера)
Port: 5432
Database: child_bot
User: child_bot
Password: [из .env.production на сервере]
```

**Шаг 4: Тест подключения**
- Вернитесь на вкладку "General"
- Нажмите "Test Connection"
- Должно быть успешно → "OK"

**Объяснение:**
- SSH туннель создаёт безопасное соединение с сервером
- PostgreSQL подключается через этот туннель к `localhost:5432` (на сервере)
- Данные шифруются через SSH, порт 5432 не открыт в интернет

**Полезные фичи GoLand Database:**
- SQL Console (Ctrl+Shift+F10) - выполнение запросов
- Database Diagrams - визуализация схемы
- Data Editor - редактирование данных прямо в таблице
- Query History - история выполненных запросов
- Auto-completion - автодополнение SQL

### 2. VS Code (с расширением PostgreSQL)

**Установка расширения:**
```
Ctrl+P → ext install ckolkman.vscode-postgres
```

**Подключение:**
1. Нажмите на иконку PostgreSQL в боковой панели
2. Add Connection
3. Заполните параметры:
   ```
   Hostname: localhost
   User: child_bot
   Password: [из .env]
   Port: 5432
   Use SSL: No (для локальной разработки)
   Database: child_bot
   ```

**Альтернатива - Database Client:**
```
Ctrl+P → ext install cweijan.vscode-database-client2
```

### 3. TablePlus

**Создание подключения:**
1. File → New → PostgreSQL
2. Заполните параметры:
   ```
   Name: Child Bot (Local)
   Host: localhost
   Port: 5432
   User: child_bot
   Password: [из .env]
   Database: child_bot
   ```
3. Test → Connect

### 4. DBeaver

**Создание подключения:**
1. Database → New Database Connection
2. Выберите PostgreSQL
3. Заполните параметры:
   ```
   Host: localhost
   Port: 5432
   Database: child_bot
   Username: child_bot
   Password: [из .env]
   ```
4. Test Connection → Finish

### 5. psql (командная строка)

**Локальное подключение:**
```bash
# Через Docker
docker exec -it child_bot_postgres psql -U child_bot -d child_bot

# Если psql установлен локально
psql -h localhost -p 5432 -U child_bot -d child_bot
# Введите пароль из .env
```

**Полезные команды psql:**
```sql
-- Список таблиц
\dt

-- Описание таблицы
\d child_profiles

-- Список баз данных
\l

-- Список схем
\dn

-- Переключиться на другую БД
\c database_name

-- История команд
\s

-- Выход
\q
```

## Безопасное подключение к production

⚠️ **Никогда не открывайте порт PostgreSQL напрямую в интернет!**

### Вариант 1: SSH туннель через командную строку

**Создание SSH туннеля:**
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa -L 5433:localhost:5432 root@77.222.60.149
```

**Что происходит:**
- `-L 5433:localhost:5432` - пробрасывает локальный порт 5433 на удалённый порт 5432
- Теперь `localhost:5433` на вашей машине = `localhost:5432` на сервере

**Подключение к БД (пока туннель активен):**
```bash
psql -h localhost -p 5433 -U child_bot -d child_bot
# Или используйте эти параметры в GoLand
```

**В GoLand (если используете туннель вручную):**
```
Host: localhost
Port: 5433 (не 5432!)
Database: child_bot
User: child_bot
Password: [из .env.production]
SSH: не настраивать (туннель уже запущен в терминале)
```

### Вариант 2: SSH туннель через GoLand (рекомендуется)

**Преимущества:**
- Не нужно держать открытым терминал
- GoLand автоматически управляет туннелем
- Удобнее для ежедневной работы

**Настройка в GoLand:**
1. Database panel → `+` → PostgreSQL
2. Вкладка SSH/SSL → "Use SSH tunnel" → `+`
3. Параметры SSH:
   ```
   Host: 77.222.60.149
   Port: 22
   User name: root
   Auth type: Key pair (OpenSSH or PuTTY)
   Private key file: /Users/a.yanover/Downloads/id_rsa_1/id_rsa
   ```
4. Вкладка General:
   ```
   Host: localhost (не IP!)
   Port: 5432
   Database: child_bot
   User: child_bot
   Password: [из .env.production]
   ```
5. Test Connection → OK

### Вариант 2: VPN

Если у вас настроен VPN к серверу:
1. Подключитесь к VPN
2. Используйте внутренний IP адрес сервера
3. Подключайтесь как обычно

### Вариант 3: Bastion host

Если используется bastion host:
```bash
ssh -J bastion-user@bastion-host your-user@database-server -L 5433:localhost:5432
```

## Проверка подключения

### 1. Проверка что PostgreSQL запущен

```bash
# Проверить статус контейнера
docker ps | grep postgres

# Проверить порт
netstat -an | grep 5432  # Linux/Mac
Get-NetTCPConnection -LocalPort 5432  # Windows PowerShell
```

### 2. Тестовый запрос

После подключения выполните:
```sql
SELECT version();
SELECT current_database();
SELECT current_user;
```

### 3. Проверка таблиц

```sql
-- Список всех таблиц
SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;

-- Количество записей в основных таблицах
SELECT
    'child_profiles' as table_name, COUNT(*) as count FROM child_profiles
UNION ALL
SELECT 'parent_profiles', COUNT(*) FROM parent_profiles
UNION ALL
SELECT 'attempts', COUNT(*) FROM attempts
UNION ALL
SELECT 'achievements', COUNT(*) FROM achievements;
```

## Полезные запросы

### Проверка XP и уровней

```sql
-- Топ пользователей по XP
SELECT id, name, xp_total, level, coins_balance
FROM child_profiles
ORDER BY xp_total DESC
LIMIT 10;

-- Профиль конкретного пользователя
SELECT
    id,
    name,
    xp_total,
    level,
    coins_balance,
    tasks_solved_total,
    tasks_correct_total,
    hints_used_total,
    current_streak_days,
    last_active_at
FROM child_profiles
WHERE id = 'your-profile-id-here';
```

### Проверка достижений

```sql
-- Все достижения пользователя
SELECT
    a.title,
    a.description,
    ca.current_progress,
    a.requirement_value,
    ca.is_unlocked,
    ca.unlocked_at
FROM child_achievements ca
JOIN achievements a ON ca.achievement_id = a.id
WHERE ca.child_profile_id = 'your-profile-id-here'
ORDER BY ca.unlocked_at DESC NULLS LAST;

-- Статистика по достижениям
SELECT
    COUNT(*) FILTER (WHERE is_unlocked = TRUE) as unlocked,
    COUNT(*) FILTER (WHERE is_unlocked = FALSE) as locked,
    COUNT(*) as total
FROM child_achievements
WHERE child_profile_id = 'your-profile-id-here';
```

### Проверка попыток

```sql
-- Последние попытки пользователя
SELECT
    id,
    status,
    is_correct,
    hints_requested,
    created_at,
    updated_at
FROM attempts
WHERE child_profile_id = 'your-profile-id-here'
ORDER BY created_at DESC
LIMIT 20;

-- Статистика попыток
SELECT
    status,
    COUNT(*) as count,
    COUNT(*) FILTER (WHERE is_correct = TRUE) as correct,
    COUNT(*) FILTER (WHERE is_correct = FALSE) as incorrect
FROM attempts
WHERE child_profile_id = 'your-profile-id-here'
GROUP BY status;
```

### Проверка отчётов

```sql
-- Настройки отчётов
SELECT
    child_profile_id,
    parent_email,
    weekly_report_enabled,
    created_at,
    updated_at
FROM report_settings
WHERE child_profile_id = 'your-profile-id-here';

-- История отчётов
SELECT
    id,
    report_date,
    sent_at,
    created_at
FROM weekly_reports
WHERE user_id = 'your-profile-id-here'
ORDER BY report_date DESC;
```

## Troubleshooting

### Ошибка: "Connection refused"

**Причины:**
1. PostgreSQL контейнер не запущен
2. Неправильный порт
3. Firewall блокирует подключение

**Решение:**
```bash
# Проверить статус контейнера
docker ps | grep postgres

# Перезапустить если нужно
docker-compose restart postgres

# Проверить логи
docker logs child_bot_postgres

# Проверить порт в .env
cat .env | grep POSTGRES_PORT
```

### Ошибка: "Authentication failed"

**Причины:**
1. Неправильный пароль
2. Неправильное имя пользователя

**Решение:**
```bash
# Проверить креды в .env
cat .env | grep POSTGRES

# Сбросить пароль (если нужно)
docker exec -it child_bot_postgres psql -U postgres -c "ALTER USER child_bot WITH PASSWORD 'new_password';"
```

### Ошибка: "Database does not exist"

**Решение:**
```bash
# Создать базу данных
docker exec -it child_bot_postgres psql -U postgres -c "CREATE DATABASE child_bot;"

# Или пересоздать контейнер
docker-compose down -v
docker-compose up -d
```

### Медленные запросы

**Проверка активных запросов:**
```sql
SELECT
    pid,
    usename,
    application_name,
    state,
    query,
    query_start
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY query_start;
```

**Остановка долгого запроса:**
```sql
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE pid = [process_id];
```

## Безопасность

### ⚠️ Важные правила

1. **Никогда не коммитьте .env файлы с реальными паролями**
2. **Используйте сильные пароли для production**
3. **Для production всегда используйте SSH туннель**
4. **Регулярно делайте бэкапы:**
   ```bash
   docker exec child_bot_postgres pg_dump -U child_bot child_bot > backup_$(date +%Y%m%d).sql
   ```
5. **Ограничьте права пользователя БД** (только необходимые таблицы)

### Рекомендации для production

```sql
-- Создать read-only пользователя для аналитики
CREATE USER analytics_user WITH PASSWORD 'strong_password';
GRANT CONNECT ON DATABASE child_bot TO analytics_user;
GRANT USAGE ON SCHEMA public TO analytics_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO analytics_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO analytics_user;
```

## Docker команды для работы с БД

```bash
# Войти в контейнер PostgreSQL
docker exec -it child_bot_postgres bash

# Посмотреть логи PostgreSQL
docker logs child_bot_postgres -f

# Перезапустить только PostgreSQL
docker-compose restart postgres

# Остановить всё
docker-compose down

# Остановить и удалить все данные (осторожно!)
docker-compose down -v

# Бэкап базы данных
docker exec child_bot_postgres pg_dump -U child_bot child_bot > backup.sql

# Восстановление из бэкапа
docker exec -i child_bot_postgres psql -U child_bot child_bot < backup.sql

# Копировать бэкап из контейнера
docker cp child_bot_postgres:/backup.sql ./backup.sql
```

## Полезные ссылки

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [DataGrip Documentation](https://www.jetbrains.com/help/datagrip/)
- [psql Commands](https://www.postgresql.org/docs/current/app-psql.html)
- [PostgreSQL Performance](https://wiki.postgresql.org/wiki/Performance_Optimization)
