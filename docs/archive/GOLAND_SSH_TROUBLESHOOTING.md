# Решение проблемы "SSH tunnel creation failed: Connection refused" в GoLand

## Проблема

При попытке подключиться к production БД через SSH tunnel в GoLand появляется ошибка:
```
SSH tunnel creation failed: Connection refused.
```

## Решение

### Вариант 1: Используйте новый ключ без extended attributes (рекомендуется)

Создан чистый ключ: `/Users/a.yanover/Downloads/id_rsa_1/id_rsa_goland`

**Настройка в GoLand:**

1. Database panel → `+` → PostgreSQL
2. Вкладка **SSH/SSL**:
   - ✅ Use SSH tunnel
   - Нажмите `+` для создания нового SSH config
3. Заполните SSH параметры:
   ```
   Host: 77.222.60.149
   Port: 22
   User name: root

   Auth type: Key pair (OpenSSH or PuTTY)
   Private key file: /Users/a.yanover/Downloads/id_rsa_1/id_rsa_goland

   ⚠️ ВАЖНО: Оставьте Passphrase пустым
   ```
4. Нажмите **Test Connection** в окне SSH Configurations
   - Если запросит fingerprint → "Yes" (или "Accept")
   - Должно быть успешно
5. Нажмите **OK** чтобы сохранить SSH config
6. Вкладка **General**:
   ```
   Host: localhost (НЕ IP сервера!)
   Port: 5432
   Database: child_bot
   User: child_bot
   Password: [из .env.production на сервере]
   ```
7. Test Connection → OK

### Вариант 2: Ручной SSH туннель (если GoLand не работает)

**Шаг 1: Создайте SSH туннель в терминале**

```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa -L 5433:localhost:5432 root@77.222.60.149
```

**НЕ ЗАКРЫВАЙТЕ** этот терминал!

**Шаг 2: Подключитесь в GoLand БЕЗ SSH туннеля**

1. Database panel → `+` → PostgreSQL
2. Вкладка **SSH/SSL**:
   - ❌ Отключите "Use SSH tunnel"
3. Вкладка **General**:
   ```
   Host: localhost
   Port: 5433 (не 5432!)
   Database: child_bot
   User: child_bot
   Password: [из .env.production]
   ```
4. Test Connection → OK

**Минусы этого способа:**
- Нужно держать открытым терминал с SSH туннелем
- При перезапуске компьютера нужно заново запускать туннель

**Плюсы:**
- Работает 100%, если SSH подключение работает
- Не зависит от настроек GoLand

### Вариант 3: Проверка SSH конфигурации GoLand

Если оба варианта выше не работают:

1. GoLand → Settings (Cmd+,)
2. Tools → SSH Configurations
3. Нажмите `+` для создания нового
4. Заполните:
   ```
   Host: 77.222.60.149
   Port: 22
   User name: root
   Authentication type: Key pair
   Private key file: /Users/a.yanover/Downloads/id_rsa_1/id_rsa_goland
   ```
5. Test Connection
6. Если успешно → используйте этот SSH config в Database settings

## Проверка что SSH работает

### Из командной строки:

```bash
# Базовая проверка
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 "echo OK"

# Проверка с новым ключом
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa_goland root@77.222.60.149 "echo OK"

# Проверка что PostgreSQL работает на сервере
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 "docker ps | grep postgres"
```

Все команды должны работать без ошибок.

## Получение пароля БД с сервера

Если не помните пароль из .env.production:

```bash
ssh -i /Users/a.yanover/Downloads/id_rsa_1/id_rsa root@77.222.60.149 "cat ~/child_bot/.env.production | grep POSTGRES_PASSWORD"
```

## Типичные ошибки GoLand SSH

### "Auth fail"
- Проверьте путь к ключу
- Убедитесь что passphrase пустой (если ключ без пароля)
- Попробуйте использовать id_rsa_goland вместо id_rsa

### "Connection timed out"
- Проверьте firewall на сервере
- Проверьте что порт 22 открыт
- Проверьте интернет соединение

### "Permission denied"
- Проверьте права на ключ: `chmod 600 /path/to/key`
- Проверьте что используется правильный username (root)

### "Host key verification failed"
- GoLand предложит добавить fingerprint - согласитесь
- Или добавьте сервер в known_hosts:
  ```bash
  ssh-keyscan -H 77.222.60.149 >> ~/.ssh/known_hosts
  ```

## Альтернатива: Используйте TablePlus или DBeaver

Если GoLand продолжает давать проблемы, можно использовать другие инструменты:

### TablePlus (проще настроить)
1. File → New → PostgreSQL
2. SSH enabled: ✅
3. SSH Host: 77.222.60.149
4. SSH User: root
5. SSH Key: /Users/a.yanover/Downloads/id_rsa_1/id_rsa_goland
6. Host: localhost
7. Port: 5432
8. Database: child_bot
9. Test → Connect

### DBeaver (бесплатный)
1. Database → New Connection → PostgreSQL
2. Main tab:
   - Host: localhost
   - Port: 5432
   - Database: child_bot
3. SSH tab:
   - ✅ Use SSH Tunnel
   - Host: 77.222.60.149
   - Port: 22
   - User: root
   - Authentication: Public Key
   - Private key: /Users/a.yanover/Downloads/id_rsa_1/id_rsa_goland
4. Test Connection → OK

## Проверка подключения после настройки

После успешного подключения выполните:

```sql
-- Проверка версии PostgreSQL
SELECT version();

-- Проверка что это production БД
SELECT COUNT(*) as total_users FROM child_profiles;

-- Список таблиц
SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;
```

Если запросы выполняются - всё работает! ✅

## Полезные ссылки

- [Основная документация](./DATABASE_CONNECTION.md)
- [Быстрый старт](./QUICK_DB_SETUP.md)
- [Пошаговая инструкция GoLand](./GOLAND_DB_SETUP.md)
