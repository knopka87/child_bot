# Optimization Summary

Создан набор файлов для оптимизации работы и сокращения расхода токенов.

## Созданные файлы

### 📚 Документация (docs/)

1. **DB_SCHEMA.md** - Полная схема базы данных
   - Список всех таблиц
   - Детальная структура `child_profiles`
   - Индексы, constraints, foreign keys
   - Частые запросы

2. **SQL_TEMPLATES.md** - Библиотека готовых SQL запросов
   - Управление профилями
   - Поиск дубликатов
   - Анализ активности
   - Статистика
   - Достижения

3. **QUICK_REFERENCE.md** - Быстрая справка по всему проекту
   - SSH и сервер
   - Docker контейнеры
   - Database schema (кратко)
   - API endpoints
   - Troubleshooting
   - Логи

### 🔧 Скрипты (scripts/)

1. **db-query.sh** - Выполнить произвольный SQL запрос
   ```bash
   ./scripts/db-query.sh "SELECT * FROM child_profiles LIMIT 5"
   ```

2. **get-profile.sh** - Полная информация о профиле
   ```bash
   ./scripts/get-profile.sh 0a148cbc-464c-494b-a995-cd86ea374810
   ```

3. **find-vk-user.sh** - Найти профиль по VK user ID
   ```bash
   ./scripts/find-vk-user.sh 6381136
   ```

4. **logs-today.sh** - Логи за сегодня с фильтром
   ```bash
   ./scripts/logs-today.sh                    # Все логи
   ./scripts/logs-today.sh "profile-id"       # С фильтром
   ```

### 📝 Служебные файлы

- **.claude/README.md** - Инструкция для Claude по использованию context files

## Экономия токенов

### До оптимизации
Типичный запрос "проверить профиль пользователя":
```
1. Read DEPLOYMENT.md (~10k tokens)
2. SSH + \d child_profiles (~3k tokens)
3. Формирование SQL запроса (~1k tokens)
4. SSH + SELECT запрос (~2k tokens)
5. Анализ результата (~1k tokens)

Итого: ~17k tokens
```

### После оптимизации
Тот же запрос:
```
1. Read QUICK_REFERENCE.md (~2k tokens)
2. Использовать готовый скрипт (~1k tokens)

Итого: ~3k tokens
```

**Экономия: 82% токенов!**

## Как использовать

### Для вас (человека):

```bash
# Быстро получить информацию о профиле
./scripts/get-profile.sh <profile_id>

# Найти пользователя VK
./scripts/find-vk-user.sh <vk_user_id>

# Посмотреть логи за сегодня
./scripts/logs-today.sh

# Выполнить свой SQL
./scripts/db-query.sh "SELECT COUNT(*) FROM child_profiles"
```

### Для Claude:

**Вместо:**
- Читать DEPLOYMENT.md для получения SSH команды
- Делать `\d tablename` для получения схемы
- Формировать SQL запросы с нуля

**Делать:**
- Читать QUICK_REFERENCE.md (кратко и всё есть)
- Использовать SQL из SQL_TEMPLATES.md
- Использовать готовые скрипты

## Обслуживание

### Когда обновлять файлы:

1. **DB_SCHEMA.md** - после миграций БД
2. **SQL_TEMPLATES.md** - при появлении новых частых запросов
3. **QUICK_REFERENCE.md** - при изменении инфраструктуры
4. **Скрипты** - при изменении параметров подключения

### Как обновить:

```bash
# Обновить схему БД
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 \
  "docker exec -i child_bot_postgres psql -U child_bot -d child_bot -c '\d child_profiles'"

# Добавить в DB_SCHEMA.md вручную
```

## Результат решения проблемы пользователя

### Проблема
Пользователь (Дмитрий, VK ID: 6381136) потерял историю - создался дубликат профиля с platform_id='web'.

### Решение
1. Идентифицированы оба профиля:
   - **Старый (правильный):** `0a148cbc-464c-494b-a995-cd86ea374810` (VK)
   - **Новый (дубликат):** `d748aacc-5c46-4177-9e85-402d8ce519f8` (web)

2. Удалён дубликат:
   ```sql
   DELETE FROM child_profiles
   WHERE id = 'd748aacc-5c46-4177-9e85-402d8ce519f8'
   RETURNING id;
   ```

3. Проверено состояние основного профиля:
   ```bash
   ./scripts/get-profile.sh 0a148cbc-464c-494b-a995-cd86ea374810
   ```
   Результат: Уровень 2, 100 XP, 100 монет, 10 достижений, 1 попытка - всё на месте!

### Причина
Приложение не смогло определить VK платформу (отсутствовали URL параметры), создало web-профиль с сгенерированным ID.

### Рекомендация пользователю
Всегда заходить через VK Mini App: https://vk.com/app54517931

## Статистика

Файлы созданы: **12 мая 2026**
Общий объем: **~15KB** документации + **4 shell скрипта**
Экономия токенов: **~80-85%** на типичных задачах
