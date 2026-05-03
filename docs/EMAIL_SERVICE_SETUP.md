# Настройка Email Сервиса

## Зачем нужен Email сервис?

Email сервис используется для отправки кодов верификации родителям при подключении к профилю ребенка. Это важная часть системы безопасности и родительского контроля.

## Поддерживаемые провайдеры

### SendGrid (рекомендуется)
- ✅ Полностью реализован
- ✅ Простая настройка
- ✅ Бесплатный тариф: 100 писем/день
- ✅ HTML email поддержка
- ✅ Хорошая доставляемость

### AWS SES
- ⏳ Запланирован (TODO)
- Подходит если уже используете AWS

### Mailgun
- ⏳ Запланирован (TODO)
- Альтернатива SendGrid

## Настройка SendGrid

### Шаг 1: Создать аккаунт SendGrid

1. Перейдите на https://signup.sendgrid.com/
2. Зарегистрируйтесь (можно использовать бесплатный план)
3. Подтвердите email адрес

### Шаг 2: Создать API Key

1. Войдите в SendGrid Dashboard
2. Перейдите в **Settings** → **API Keys**
3. Нажмите **Create API Key**
4. Выберите тип: **Full Access** (или ограниченный с доступом к Mail Send)
5. Введите имя ключа: `obiasnyatel-dz-production` (или `dev`)
6. Скопируйте созданный API Key (он показывается только один раз!)

**Важно:**
- Сохраните ключ в безопасном месте
- Никогда не коммитьте ключ в git
- Используйте разные ключи для dev и production

### Шаг 3: Верифицировать отправителя (Sender)

SendGrid требует верификации email адреса отправителя:

1. Перейдите в **Settings** → **Sender Authentication**
2. Выберите **Single Sender Verification**
3. Нажмите **Create New Sender**
4. Заполните форму:
   - **From Name**: Объяснятель ДЗ
   - **From Email Address**: noreply@obiasnyatel-dz.ru (или ваш домен)
   - **Reply To**: support@obiasnyatel-dz.ru
   - **Company Address**: ваш адрес
5. Подтвердите email отправителя (SendGrid отправит письмо)

**Для production:**
Настройте Domain Authentication для лучшей доставляемости:
1. **Settings** → **Sender Authentication** → **Authenticate Your Domain**
2. Следуйте инструкциям для добавления DNS записей

### Шаг 4: Добавить в .env файл

Откройте файл `.env` в корне проекта:

```bash
# Email Service
EMAIL_PROVIDER=sendgrid
EMAIL_API_KEY=SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
EMAIL_FROM=noreply@obiasnyatel-dz.ru
```

Замените значения:
- `EMAIL_API_KEY` - ваш API ключ из шага 2
- `EMAIL_FROM` - верифицированный email из шага 3

### Шаг 5: Перезапустить backend

```bash
docker compose -f docker/docker-compose.dev.yml restart backend
```

## Тестирование

### Development режим

В dev режиме (`ENV=development`) email **не отправляется**. Вместо этого:
- Код логируется в консоль backend
- Код возвращается в API ответе в поле `devCode`

Пример ответа в dev:
```json
{
  "message": "Verification code sent to email",
  "expiresAt": "2024-04-19T15:30:00Z",
  "devCode": "123456"
}
```

### Production режим

В production (`ENV=production`):
- Email отправляется через SendGrid
- Код **НЕ** возвращается в ответе
- Пользователь получает красивое HTML письмо

Пример ответа в production:
```json
{
  "message": "Verification code sent to email",
  "expiresAt": "2024-04-19T15:30:00Z"
}
```

### Проверка отправки

1. Установите `ENV=production` в `.env`
2. Перезапустите backend
3. Отправьте POST запрос:

```bash
curl -X POST http://localhost:8080/api/v1/email/verify/send \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-test-email@example.com",
    "parentUserId": "test-parent-123"
  }'
```

4. Проверьте почту (включая спам)
5. Проверьте логи backend:

```bash
docker logs child_bot_backend_dev --tail 50 | grep EmailService
```

Успешная отправка:
```
[EmailService] Email sent successfully to your-test-email@example.com via SendGrid
```

## Шаблон Email

Отправляемое письмо включает:
- 📚 Логотип и брендинг "Объяснятель ДЗ"
- Большой 6-значный код верификации
- Время истечения кода (15 минут)
- Предупреждение о безопасности
- Контактную информацию поддержки
- Адаптивный дизайн для мобильных устройств

Шаблон находится в `api/internal/service/email.go` в функции `renderVerificationEmail()`.

## Troubleshooting

### Ошибка "EMAIL_API_KEY not configured"

Проверьте что переменная `EMAIL_API_KEY` установлена в `.env` и backend перезапущен.

```bash
grep EMAIL_API_KEY .env
docker compose -f docker/docker-compose.dev.yml restart backend
```

### Email не приходит

1. **Проверьте статус в SendGrid:**
   - Dashboard → Activity Feed
   - Смотрите события отправки

2. **Проверьте папку спам**

3. **Проверьте логи backend:**
   ```bash
   docker logs child_bot_backend_dev --tail 100 | grep Email
   ```

4. **Проверьте что sender верифицирован:**
   - Settings → Sender Authentication
   - Email должен быть Verified

### SendGrid возвращает 401

API Key невалидный или отозван:
- Создайте новый API Key в SendGrid
- Обновите `EMAIL_API_KEY` в `.env`
- Перезапустите backend

### SendGrid возвращает 403

Sender не верифицирован:
- Пройдите Single Sender Verification
- Используйте верифицированный email в `EMAIL_FROM`

## Мониторинг

### Логи отправки

Все отправки логируются:

```
# Успешная отправка:
[EmailService] Email sent successfully to parent@example.com via SendGrid

# Ошибка отправки:
[EmailHandler] Failed to send email: SendGrid returned status 403
```

### SendGrid Dashboard

Проверяйте статистику в SendGrid:
- **Dashboard** → **Stats** - общая статистика
- **Activity Feed** - детали каждой отправки
- **Suppressions** - список заблокированных адресов

## Безопасность

**Никогда не:**
- ❌ Не коммитьте `EMAIL_API_KEY` в git
- ❌ Не передавайте API key в frontend
- ❌ Не логируйте полные email адреса в production
- ❌ Не отправляйте реальные email в dev режиме

**Всегда:**
- ✅ Храните API key в `.env` (который в .gitignore)
- ✅ Используйте разные API keys для dev и production
- ✅ Ротируйте ключ если он скомпрометирован
- ✅ Включайте отправку только в production (`ENV=production`)
- ✅ Проверяйте SendGrid Activity Feed на подозрительную активность

## Лимиты

### Бесплатный план SendGrid
- 100 писем в день
- Все основные функции
- Подходит для начальной стадии

### Платные планы
- Essentials: $19.95/месяц - 50,000 писем
- Pro: от $89.95/месяц - 100,000+ писем
- Смотрите актуальные цены: https://sendgrid.com/pricing/

## Переход на другой провайдер

Если нужно переключиться на AWS SES или Mailgun:

1. Реализуйте соответствующую функцию в `api/internal/service/email.go`:
   - `sendViaSES()` - для AWS SES
   - `sendViaMailgun()` - для Mailgun

2. Обновите `.env`:
   ```bash
   EMAIL_PROVIDER=ses  # или mailgun
   ```

3. Добавьте необходимые credentials для нового провайдера

## Ссылки

- [SendGrid Документация](https://docs.sendgrid.com/)
- [SendGrid API Reference](https://docs.sendgrid.com/api-reference/mail-send/mail-send)
- [Sender Authentication Guide](https://docs.sendgrid.com/ui/account-and-settings/how-to-set-up-domain-authentication)
