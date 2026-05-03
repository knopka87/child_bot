# 🚀 Быстрое подключение к БД - Шпаргалка

## ✅ Production сервер ГОТОВ к подключениям

**Статус:** PostgreSQL на production (77.222.60.149) настроен и принимает подключения через SSH туннель.

Порт `5432` пробрасывается только на `127.0.0.1` (безопасно, не виден из интернета).

---

## ✅ РЕКОМЕНДУЕМЫЙ способ: GoLand с новым ключом

### Параметры SSH (вкладка SSH/SSL в GoLand)
```
✅ Use SSH tunnel

Host: 77.222.60.149
Port: 22
User name: root
Private key: /Users/a.yanover/Downloads/id_rsa_1/id_rsa_goland
Passphrase: [пусто]
```

### Параметры PostgreSQL (вкладка General)
```
Host: localhost (ВАЖНО!)
Port: 5432
Database: child_bot
User: child_bot
Password: [получить с сервера, см. ниже]
```

### Получить пароль БД:
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 "grep POSTGRES_PASSWORD /root/child_bot/.env.production"
```

---

## 🔧 Альтернатива: Ручной SSH туннель

### Терминал 1 (не закрывать):
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa -L 5433:localhost:5432 root@77.222.60.149
```

### GoLand подключение:
```
SSH: выключен ❌
Host: localhost
Port: 5433 (не 5432!)
Database: child_bot
User: child_bot
Password: [см. выше]
```

---

## 📊 Полезные запросы

### Проверка XP пользователя
```sql
SELECT id, name, xp_total, level, coins_balance
FROM child_profiles
WHERE id = 'user-id-here';
```

### Топ по XP
```sql
SELECT id, name, xp_total, level
FROM child_profiles
ORDER BY xp_total DESC
LIMIT 10;
```

### Последние попытки
```sql
SELECT id, status, is_correct, created_at
FROM attempts
ORDER BY created_at DESC
LIMIT 20;
```

---

## 🆘 Troubleshooting

### Проверка SSH:
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 "echo OK"
```

### Проверка PostgreSQL на сервере:
```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 "docker ps | grep postgres"
```

### Полная документация:
- `docs/GOLAND_SSH_TROUBLESHOOTING.md` - решение проблем с SSH
- `docs/GOLAND_DB_SETUP.md` - пошаговая инструкция
- `docs/DATABASE_CONNECTION.md` - полная документация
