# Получение VK App Secret для валидации sign

## Зачем нужен VK App Secret?

VK App Secret используется для валидации подписи (sign parameter) запросов от VK Mini Apps. Это критически важно для безопасности - без проверки подписи любой злоумышленник может подделать `vk_user_id` и получить доступ к чужому профилю.

## Как получить VK App Secret

### Шаг 1: Открыть настройки приложения

1. Перейдите на https://vk.com/apps?act=manage
2. Войдите в аккаунт VK если требуется
3. Найдите ваше приложение в списке (ID: 54517931)
4. Кликните на приложение чтобы открыть его

### Шаг 2: Открыть раздел "Настройки"

1. В левом меню найдите раздел **"Настройки"**
2. Кликните на него

### Шаг 3: Скопировать Защищённый ключ

1. Найдите поле **"Защищённый ключ"** (или "Secret Key")
2. Скопируйте значение ключа

**Важно:**
- Это **секретный** ключ - никогда не публикуйте его в публичных репозиториях
- Если ключ скомпрометирован - его можно пересоздать в настройках VK

### Шаг 4: Добавить в .env файл

Откройте файл `.env` в корне проекта и обновите строку:

```bash
VK_APP_SECRET=ваш_секретный_ключ_из_vk
```

Замените `ЗАМЕНИТЕ_НА_РЕАЛЬНЫЙ_SECRET_ИЗ_VK` на реальный ключ.

### Шаг 5: Перезапустить backend

```bash
docker compose -f docker/docker-compose.dev.yml restart backend
```

## Проверка что валидация работает

### Development режим

В dev режиме (ENV=development) валидация sign **отключена** для удобства разработки. Все запросы будут пропускаться.

### Production режим

В production режиме (ENV=production) каждый запрос с VK параметрами проверяется:

1. Собираются все `vk_*` параметры из query string
2. Параметры сортируются по ключу
3. Формируется строка `key1=value1&key2=value2...`
4. Вычисляется `HMAC-SHA256(строка, VK_APP_SECRET)`
5. Результат кодируется в base64 URL-safe формат
6. Сравнивается с параметром `sign`

Если подписи не совпадают - запрос отклоняется с кодом **401 Unauthorized**.

## Логирование

Middleware логирует все попытки валидации:

```
# Успешная валидация:
[VK Auth] Valid VK signature for user: 12345678

# Невалидная подпись:
[VK Auth] Invalid VK signature for user: 12345678

# Отсутствует sign:
[VK Auth] Missing sign parameter

# Dev режим:
[VK Auth] Development mode: skipping sign validation
```

## Тестирование

Запустите тесты чтобы убедиться что валидация работает корректно:

```bash
cd api
go test -v ./internal/api/middleware/vk_auth_test.go ./internal/api/middleware/vk_auth.go
```

Все 5 тестов должны пройти:
- ✅ TestVKAuthMiddleware_ValidSign
- ✅ TestVKAuthMiddleware_InvalidSign
- ✅ TestVKAuthMiddleware_MissingSign
- ✅ TestVKAuthMiddleware_DevelopmentMode
- ✅ TestVKAuthMiddleware_NoVKParams

## Troubleshooting

### Ошибка "VK_APP_SECRET not configured"

Проверьте что переменная VK_APP_SECRET установлена в `.env` файле и backend перезапущен.

### Все запросы отклоняются с 401

1. Проверьте что вы используете правильный VK_APP_SECRET
2. Проверьте что в dev режиме (ENV=development) валидация отключена
3. Проверьте логи backend: `docker logs child_bot_backend_dev --tail 50`

### Запросы проходят без валидации

Проверьте переменную ENV:
- `ENV=development` - валидация отключена (норма для разработки)
- `ENV=production` - валидация включена

## Безопасность

**Никогда не:**
- ❌ Не коммитьте VK_APP_SECRET в git
- ❌ Не передавайте VK_APP_SECRET в frontend
- ❌ Не логируйте VK_APP_SECRET
- ❌ Не отправляйте VK_APP_SECRET в API ответах

**Всегда:**
- ✅ Храните VK_APP_SECRET в `.env` (который в .gitignore)
- ✅ Используйте разные secrets для dev и production
- ✅ Ротируйте secret если он скомпрометирован
- ✅ Включайте валидацию в production (ENV=production)
