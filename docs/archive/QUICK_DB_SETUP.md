# Быстрый старт: Подключение к БД

## 🚀 За 2 минуты

### 1. Убедитесь что PostgreSQL запущен

```bash
docker ps | grep postgres
```

Должен быть запущен контейнер `child_bot_postgres`.

Если не запущен:
```bash
docker-compose up -d postgres
```

### 2. Получите параметры подключения

```bash
# Посмотреть креды
cat .env | grep POSTGRES
```

Вы увидите:
```
POSTGRES_DB=child_bot
POSTGRES_USER=child_bot
POSTGRES_PASSWORD=ваш_пароль
POSTGRES_PORT=5432
```

### 3. Подключитесь из GoLand

#### Локальное подключение (Docker)

1. Откройте Database panel (View → Tool Windows → Database)
2. Нажмите `+` → Data Source → PostgreSQL
3. Заполните на вкладке General:
   ```
   Host: localhost
   Port: 5432
   Database: child_bot
   User: child_bot
   Password: [из .env]
   ```
4. Test Connection → Download driver (если нужно) → OK

#### Production подключение

1. Database panel → `+` → PostgreSQL
2. Вкладка SSH/SSL:
   - ✅ Use SSH tunnel → `+`
   - Host: `77.222.60.149`
   - Port: `22`
   - User name: `root`
   - Auth type: `Key pair`
   - Private key: `/Users/a.yanover/Downloads/id_rsa_1/id_rsa`
   - Test Connection (SSH)
3. Вкладка General:
   ```
   Host: localhost (не IP сервера!)
   Port: 5432
   Database: child_bot
   User: child_bot
   Password: [из .env.production на сервере]
   ```
4. Test Connection → OK

#### Альтернатива: psql через терминал

**Локально:**
```bash
docker exec -it child_bot_postgres psql -U child_bot -d child_bot
```

**Production (через SSH туннель):**
```bash
# Терминал 1: создать туннель
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa -L 5433:localhost:5432 root@77.222.60.149

# Терминал 2: подключиться к БД
psql -h localhost -p 5433 -U child_bot -d child_bot
```

### 4. Проверка подключения

Выполните тестовый запрос:
```sql
SELECT COUNT(*) FROM child_profiles;
```

## 📊 Полезные запросы для начала

### Профили пользователей
```sql
-- Топ по XP
SELECT id, name, xp_total, level FROM child_profiles ORDER BY xp_total DESC LIMIT 10;

-- Поиск пользователя
SELECT * FROM child_profiles WHERE name ILIKE '%имя%';
```

### Достижения
```sql
-- Статистика достижений
SELECT
    COUNT(*) FILTER (WHERE is_unlocked) as unlocked,
    COUNT(*) FILTER (WHERE NOT is_unlocked) as locked
FROM child_achievements
WHERE child_profile_id = 'your-id';
```

### Попытки
```sql
-- Последние попытки
SELECT id, status, is_correct, created_at
FROM attempts
ORDER BY created_at DESC
LIMIT 20;
```

## 🔥 Production (через SSH туннель)

### Способ 1: Встроенный SSH туннель в GoLand (проще)

См. инструкцию в разделе "Production подключение" выше ↑

### Способ 2: Ручной SSH туннель через терминал

```bash
# Создать SSH туннель
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa -L 5433:localhost:5432 root@77.222.60.149

# Держите этот терминал открытым!
# В другом терминале или GoLand подключайтесь:
```

В GoLand (пока туннель активен):
```
Host: localhost
Port: 5433  ← не 5432!
Database: child_bot
User: child_bot
Password: [из .env.production]
SSH: не настраивать (туннель уже работает)
```

## 📖 Полная документация

См. [DATABASE_CONNECTION.md](./DATABASE_CONNECTION.md)

## ⚠️ Безопасность

- ❌ Не коммитьте `.env` файлы
- ✅ Используйте SSH туннель для production
- ✅ Делайте регулярные бэкапы:
  ```bash
  docker exec child_bot_postgres pg_dump -U child_bot child_bot > backup.sql
  ```
