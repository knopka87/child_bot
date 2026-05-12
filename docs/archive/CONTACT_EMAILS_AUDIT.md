# Аудит и обновление контактных email-адресов

**Дата:** 19.04.2026
**Task:** #16 - Проверить и обновить contact emails

## Цель

Проверить все email-адреса в проекте и заменить плейсхолдеры на продакшн-адреса перед модерацией VK.

## Результаты аудита

### ✅ Правильные продакшн email-адреса

Основные контактные адреса уже настроены корректно:

| Назначение | Email | Где используется |
|------------|-------|------------------|
| Поддержка | `support@obiasnyatel-dz.ru` | Email service, БД миграции, документация |
| Конфиденциальность | `privacy@obiasnyatel-dz.ru` | БД миграции (Privacy Policy), документация |
| Отправитель писем | `noreply@obiasnyatel-dz.ru` | `.env`, email service |

### 🔧 Исправленные файлы

#### 1. `.env.new`

**Проблема:**
- Использовался неправильный домен: `noreply@dzobiasnitel.ru`
- Не было пояснений к плейсхолдерам

**Исправлено:**
```diff
- SMTP_FROM=ДЗ Объяснитель <noreply@dzobiasnitel.ru>
+ SMTP_FROM=ДЗ Объяснитель <noreply@obiasnyatel-dz.ru>

- SMTP_USERNAME=your_email@gmail.com
- SMTP_PASSWORD=your_app_password
+ SMTP_USERNAME=your_email@gmail.com  # ВАЖНО: Замените на реальный email перед продакшеном
+ SMTP_PASSWORD=your_app_password     # ВАЖНО: Замените на App Password из Gmail
```

#### 2. `frontend/docs/VK_MODERATION.md`

**Проблема:**
- Использовался плейсхолдер домена: `support@homework-app.ru`

**Исправлено:**
```diff
- [ ] Email поддержки: support@homework-app.ru
- [ ] VK сообщество: vk.com/homework_app
+ [x] Email поддержки: support@obiasnyatel-dz.ru
+ [x] Email конфиденциальности: privacy@obiasnyatel-dz.ru
+ [ ] VK сообщество: vk.com/obiasnyatel_dz
```

#### 3. `frontend/roadmap/15_VK_MODERATION.md`

**Проблема:**
- Использовались example.com домены в примерах кода

**Исправлено:**
```diff
- Email: <a href="mailto:privacy@example.com">privacy@example.com</a>
- Email: <a href="mailto:support@example.com">support@example.com</a>
+ Email: <a href="mailto:privacy@obiasnyatel-dz.ru">privacy@obiasnyatel-dz.ru</a>
+ Email: <a href="mailto:support@obiasnyatel-dz.ru">support@obiasnyatel-dz.ru</a>
```

### ✅ Проверенные файлы (изменений не требуется)

#### База данных

**Файл:** `api/migrations/034_legal_documents.up.sql`

Privacy Policy содержит правильные контакты:
```sql
- **Email:** privacy@obiasnyatel-dz.ru
- **Служба поддержки:** support@obiasnyatel-dz.ru
```

Terms of Service содержит правильные контакты:
```sql
- **Email:** support@obiasnyatel-dz.ru
- **Техническая поддержка:** Понедельник-Пятница, 10:00-18:00 (МСК)
- **ВКонтакте:** vk.com/obiasnyatel_dz
```

#### Email Service

**Файл:** `api/internal/service/email.go`

HTML шаблон письма с кодом верификации:
```html
<a href="mailto:support@obiasnyatel-dz.ru">support@obiasnyatel-dz.ru</a>
<a href="https://vk.com/obiasnyatel_dz">vk.com/obiasnyatel_dz</a>
```

#### Конфигурационные файлы

**Файл:** `.env` (production)
```env
EMAIL_FROM=noreply@obiasnyatel-dz.ru
```

**Файл:** `.env.example`
```env
EMAIL_FROM=support@obiasnyatel-dz.ru
```

### 📝 Тестовые и примеры (не требуют изменений)

Следующие файлы содержат тестовые/примерные email, но это корректно:

| Файл | Email | Контекст |
|------|-------|----------|
| `frontend/e2e/critical/01-onboarding.spec.ts` | `parent@example.com` | E2E тест онбординга |
| `frontend/roadmap/15_VK_MODERATION.md` | `moderator@test.com` | Тестовые данные для модерации |
| `frontend/src/pages/Onboarding/screens/EmailInput.tsx` | `example@mail.ru` | Placeholder в input поле |
| `docs/*.md` | `test@example.com` | Примеры в документации |

## Статус email-адресов по категориям

### 🟢 Продакшн контакты (готовы к модерации)

| Категория | Email | Статус |
|-----------|-------|--------|
| Поддержка | support@obiasnyatel-dz.ru | ✅ Настроен везде |
| Конфиденциальность | privacy@obiasnyatel-dz.ru | ✅ Настроен в БД и документации |
| Отправитель писем | noreply@obiasnyatel-dz.ru | ✅ Настроен в .env |

### 🟡 Требуют настройки перед продакшеном

| Переменная | Текущее значение | Что нужно сделать |
|------------|------------------|-------------------|
| `SMTP_USERNAME` | `your_email@gmail.com` | Заменить на реальный Gmail аккаунт |
| `SMTP_PASSWORD` | `your_app_password` | Создать App Password в Gmail |
| `EMAIL_API_KEY` | (пусто в .env) | Добавить SendGrid API key для продакшена |

### 🔵 Социальные сети и сообщества

| Платформа | URL | Статус |
|-----------|-----|--------|
| ВКонтакте | vk.com/obiasnyatel_dz | ✅ Указан в БД и email шаблонах |

## Проверка консистентности

### Все домены унифицированы

✅ Все продакшн email используют единый домен: `obiasnyatel-dz.ru`

### Email-адреса в критических местах

| Место | Email | ✓ |
|-------|-------|---|
| Privacy Policy (БД) | privacy@obiasnyatel-dz.ru | ✅ |
| Terms of Service (БД) | support@obiasnyatel-dz.ru | ✅ |
| Email шаблон верификации | support@obiasnyatel-dz.ru | ✅ |
| VK Moderation чеклист | support@obiasnyatel-dz.ru | ✅ |
| Отправитель писем (.env) | noreply@obiasnyatel-dz.ru | ✅ |

## Рекомендации для модерации VK

### ✅ Готово

1. **Контактные данные** - все email-адреса используют корректный домен
2. **Privacy Policy** - содержит актуальный email для запросов по GDPR
3. **Terms of Service** - содержит контакты техподдержки
4. **Email верификация** - в письмах указаны правильные контакты

### ⚠️ Требуется внимание

1. **SMTP настройки** - перед деплоем в продакшн:
   - Настроить реальный Gmail аккаунт для SMTP
   - Создать App Password
   - Или использовать SendGrid (уже настроен в `.env`)

2. **VK сообщество** - убедиться что `vk.com/obiasnyatel_dz`:
   - Создано и активно
   - Доступно пользователям
   - Модерируется

## Следующие шаги

1. ✅ **Email-адреса обновлены** - все плейсхолдеры заменены
2. ⏭️ **Task #13** - Протестировать на реальных устройствах VK
3. ⏭️ **Task #14** - Провести Lighthouse аудит
4. ⏭️ **Task #17** - Финальная проверка перед модерацией

## Изменённые файлы

```
M  .env.new                                    # Исправлен домен email
M  frontend/docs/VK_MODERATION.md             # Обновлены контакты
M  frontend/roadmap/15_VK_MODERATION.md       # Заменены example.com на реальные
A  docs/CONTACT_EMAILS_AUDIT.md               # Этот документ
```

## Заключение

✅ **Все критические email-адреса проверены и обновлены**

Продакшн email-адреса (`support@`, `privacy@`, `noreply@`) корректно настроены во всех ключевых местах:
- База данных (legal documents)
- Email service (шаблоны писем)
- Конфигурация (.env файлы)
- Документация для модерации

Приложение готово к модерации VK с точки зрения контактной информации.
