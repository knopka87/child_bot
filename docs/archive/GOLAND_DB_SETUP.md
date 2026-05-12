# Настройка БД в GoLand - Пошаговая инструкция

## 🎯 Ваши параметры

**SSH подключение:**
```
Host: 77.222.60.149
Port: 22
User: root
Private Key: /Users/a.yanover/Downloads/id_rsa_1/id_rsa
```

**PostgreSQL:**
```
Database: child_bot
User: child_bot
Password: [из .env или .env.production]
```

---

## 📋 Локальное подключение (Docker на Mac)

### Шаг 1: Убедитесь что PostgreSQL запущен

```bash
docker ps | grep postgres
```

Если не запущен:
```bash
docker-compose up -d postgres
```

### Шаг 2: Откройте Database panel в GoLand

**Горячая клавиша:** `Cmd+Shift+A` → наберите "Database" → Enter

Или через меню: `View → Tool Windows → Database`

### Шаг 3: Создайте подключение

1. В Database panel нажмите `+` (или правый клик → New)
2. Выберите `Data Source → PostgreSQL`

### Шаг 4: Заполните параметры (вкладка General)

```
Name: Child Bot (Local)

Host: localhost
Port: 5432
Database: child_bot

User: child_bot
Password: [скопируйте из .env файла]
```

### Шаг 5: Скачайте драйвер (если нужно)

- Внизу окна появится сообщение "Download missing driver files"
- Нажмите "Download"
- Дождитесь завершения

### Шаг 6: Тест подключения

- Нажмите кнопку "Test Connection"
- Должно появиться ✅ "Succeeded"
- Нажмите "OK"

### Готово! ✅

Теперь в Database panel вы увидите:
```
├─ Child Bot (Local)
   ├─ schemas
   │  └─ public
   │     ├─ tables
   │     │  ├─ achievements
   │     │  ├─ attempts
   │     │  ├─ child_profiles
   │     │  └─ ...
```

---

## 🚀 Production подключение (через SSH туннель)

### Шаг 1: Создайте новое подключение

1. Database panel → `+` → `Data Source → PostgreSQL`

### Шаг 2: Настройте SSH туннель

1. Перейдите на вкладку **SSH/SSL**
2. Поставьте галочку ✅ **Use SSH tunnel**
3. Нажмите `+` (рядом с выпадающим списком SSH конфигураций)

### Шаг 3: Заполните SSH параметры

**В открывшемся окне "SSH Configurations":**

```
Host: 77.222.60.149
Port: 22
User name: root

Auth type: Key pair (OpenSSH or PuTTY)

Private key file: /Users/a.yanover/Downloads/id_rsa_1/id_rsa
  → Нажмите на иконку папки справа
  → Найдите файл id_rsa
  → Выберите его

Passphrase: [оставьте пустым если ключ без пароля]
```

### Шаг 4: Проверьте SSH подключение

- Нажмите кнопку **"Test Connection"** (внизу окна SSH Configurations)
- Должно появиться сообщение о успешном подключении
- Нажмите **"OK"**

### Шаг 5: Настройте PostgreSQL

Вернитесь на вкладку **General** и заполните:

```
Name: Child Bot (Production)

⚠️ ВАЖНО: Host должен быть localhost, НЕ IP сервера!

Host: localhost
Port: 5432
Database: child_bot

User: child_bot
Password: [из .env.production на сервере]
```

**Почему localhost?**
- SSH туннель пробрасывает порт 5432 с сервера на вашу машину
- GoLand подключается к "localhost:5432" который на самом деле ведёт на сервер
- Это безопасно - данные шифруются через SSH

### Шаг 6: Тест подключения к БД

- Вернитесь на вкладку **General**
- Нажмите **"Test Connection"**
- Должно быть ✅ "Succeeded"
- Нажмите **"OK"**

### Готово! ✅

Теперь у вас есть два подключения:
- 🟢 Child Bot (Local) - локальная разработка
- 🔴 Child Bot (Production) - продакшн данные

---

## 💡 Полезные советы GoLand

### Выполнение SQL запросов

1. **Открыть SQL консоль:**
   - Правый клик на подключении → "New → Query Console"
   - Или: `Cmd+Shift+A` → "Query Console"

2. **Выполнить запрос:**
   - Напишите SQL запрос
   - `Cmd+Enter` - выполнить запрос под курсором
   - `Cmd+Shift+Enter` - выполнить все запросы

3. **Экспорт результатов:**
   - Правый клик на таблице с результатами
   - Export Data → CSV / JSON / SQL

### Быстрые действия

| Действие | Горячая клавиша |
|----------|----------------|
| Открыть Database panel | `Cmd+Shift+A` → Database |
| SQL Console | `Cmd+Shift+F10` |
| Выполнить запрос | `Cmd+Enter` |
| Автодополнение | `Ctrl+Space` |
| История запросов | Database panel → History |
| Поиск в таблице | `Cmd+F` |

### Просмотр данных

1. **Открыть таблицу:**
   - Двойной клик на таблице в Database panel
   - Или: правый клик → "Open Query Console"

2. **Фильтрация данных:**
   - Нажмите на иконку воронки в заголовке колонки
   - Введите условие фильтра

3. **Сортировка:**
   - Клик на заголовок колонки

4. **Редактирование:**
   - Двойной клик на ячейку
   - Отредактируйте значение
   - `Cmd+Enter` для сохранения

### Диаграммы БД

1. Правый клик на базе данных → "Diagrams → Show Visualization"
2. Drag & drop таблицы на диаграмму
3. GoLand автоматически покажет связи (foreign keys)

---

## 🔧 Troubleshooting

### "Can't connect: Connection refused"

**Для локального подключения:**
```bash
# Проверить что PostgreSQL запущен
docker ps | grep postgres

# Перезапустить если нужно
docker-compose restart postgres
```

**Для production:**
```bash
# Проверить SSH подключение
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149

# Проверить что PostgreSQL запущен на сервере
docker ps | grep postgres
```

### "Authentication failed"

- Проверьте пароль в `.env` (локально) или `.env.production` (на сервере)
- Убедитесь что нет лишних пробелов в начале/конце пароля

### "No suitable driver found"

- Нажмите "Download" в окне настройки подключения
- Или: Settings → Build, Execution, Deployment → Database → Drivers → PostgreSQL → Download

### "SSH: Auth fail"

- Проверьте путь к приватному ключу: `/Users/a.yanover/Downloads/id_rsa_1/id_rsa`
- Проверьте права на файл ключа: `ls -la /Users/a.yanover/Downloads/id_rsa_1/id_rsa`
- Должно быть: `-rw-------` (600)
- Если нет: `chmod 600 /Users/a.yanover/Downloads/id_rsa_1/id_rsa`

---

## 📚 Полезные запросы для начала

### Проверка подключения

```sql
-- Версия PostgreSQL
SELECT version();

-- Текущая база данных
SELECT current_database();

-- Текущий пользователь
SELECT current_user;
```

### Статистика профилей

```sql
-- Топ пользователей по XP
SELECT id, name, xp_total, level, coins_balance
FROM child_profiles
ORDER BY xp_total DESC
LIMIT 10;
```

### Проверка XP для конкретного пользователя

```sql
-- Замените 'user-id' на реальный ID
SELECT
    id,
    name,
    xp_total,
    level,
    coins_balance,
    tasks_solved_total,
    current_streak_days
FROM child_profiles
WHERE id = 'user-id';
```

### Последние попытки

```sql
SELECT
    id,
    status,
    is_correct,
    hints_requested,
    created_at
FROM attempts
ORDER BY created_at DESC
LIMIT 20;
```

---

## 📖 Дополнительная документация

- Полная документация: [DATABASE_CONNECTION.md](./DATABASE_CONNECTION.md)
- Быстрый старт: [QUICK_DB_SETUP.md](./QUICK_DB_SETUP.md)
- Скрипт проверки: `./scripts/check-db-connection.sh`
