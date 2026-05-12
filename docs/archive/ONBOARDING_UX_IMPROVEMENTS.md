# Улучшения UX онбординга

## Дата: 2026-04-04

## Проблемы и решения

### 1. ❌ Письмо с кодом не приходит

**Проблема:**
Пользователи не получают письма с кодом подтверждения email.

**Причина:**
Реальная отправка email не реализована - это TODO в коде. Сейчас код только генерируется и логируется на сервере.

**Решение для dev режима:**
✅ Добавлено отображение кода прямо в интерфейсе для тестирования:

```typescript
// Backend возвращает devCode в ответе
{
  "message": "Verification code sent to email",
  "expiresAt": "...",
  "devCode": "123456"  // Для dev режима
}

// Frontend показывает код в зелёном блоке
{devCode && (
  <div className="bg-green-50 border-2 border-green-200">
    <p>🔧 Режим разработки</p>
    <p>Письма пока не отправляются. Используйте код:</p>
    <p className="text-3xl font-mono">{devCode}</p>
  </div>
)}
```

**Скриншот (концептуально):**
```
┌─────────────────────────────────────┐
│ Проверка email                      │
│ Мы отправили код на test@email.com  │
│                                     │
│ ┌─────────────────┐                │
│ │   0 0 0 0 0 0   │ [Ввод кода]   │
│ └─────────────────┘                │
│                                     │
│ ⏱️ Код действителен 15 минут       │
│                                     │
│ ┌───────────────────────────────┐  │
│ │ 🔧 Режим разработки           │  │
│ │ Письма пока не отправляются.  │  │
│ │ Используйте код:              │  │
│ │                               │  │
│ │      1 2 3 4 5 6              │  │
│ └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

**Для production:**
- TODO: Интеграция с email-сервисом (SendGrid, AWS SES, Mailgun)
- TODO: Удалить `devCode` из ответа API
- TODO: Убрать зелёный блок из UI

---

### 2. ❌ Непонятный прогресс "17%" на первом шаге

**Проблема:**
На странице выбора класса сразу показывается 17% прогресса, хотя пользователь ещё ничего не сделал.

**Причина:**
Математический расчёт прогресса: `((шаг_0 + 1) / 6_шагов) * 100 = 16.67%`

Хотя это правильно математически, визуально это сбивает с толку пользователя.

**Решение:**
✅ Изменено отображение с процентов на "Шаг X из Y":

**Было:**
```
┌──────────────────────────────────┐
│  [←]  ████░░░░░░░░  17%         │
└──────────────────────────────────┘
```

**Стало:**
```
┌──────────────────────────────────┐
│  [←]  ████░░░░░░░░  Шаг 1 из 6  │
└──────────────────────────────────┘
```

**Преимущества:**
- ✅ Понятно на каком шаге находимся
- ✅ Видно сколько всего шагов осталось
- ✅ Не вводит в заблуждение на первом шаге
- ✅ Более информативно для пользователя

---

## Технические изменения

### Frontend: OnboardingPageNew.tsx

#### 1. Добавлен state для devCode

```typescript
const [devCode, setDevCode] = useState<string>('');
```

#### 2. Функция sendVerificationCode сохраняет devCode

```typescript
const result = await onboardingAPI.sendEmailVerification({...});

if (result.devCode) {
  setDevCode(result.devCode);
  console.log('[Onboarding] Dev code received:', result.devCode);
}
```

#### 3. Изменена функция прогресса

**Было:**
```typescript
const progressPercent = () => {
  const steps = ['grade', 'avatar', ...];
  const currentIndex = steps.indexOf(currentStep);
  return ((currentIndex + 1) / steps.length) * 100;
};

// В UI:
<div>{Math.round(progressPercent())}%</div>
```

**Стало:**
```typescript
const getProgressInfo = () => {
  const steps = ['grade', 'avatar', ...];
  const currentIndex = steps.indexOf(currentStep);
  const currentStepNumber = currentIndex + 1;
  const totalSteps = steps.length;
  const percent = (currentStepNumber / totalSteps) * 100;

  return {
    currentStep: currentStepNumber,
    totalSteps,
    percent,
  };
};

// В UI:
<div>Шаг {getProgressInfo().currentStep} из {getProgressInfo().totalSteps}</div>
```

#### 4. Добавлен UI блок для devCode

```tsx
{devCode && (
  <div className="bg-green-50 border-2 border-green-200 rounded-2xl p-4 mt-4">
    <p className="text-sm text-green-800 font-semibold mb-2">
      🔧 Режим разработки
    </p>
    <p className="text-sm text-green-700 mb-2">
      Письма пока не отправляются. Используйте код:
    </p>
    <div className="bg-white rounded-xl p-3 text-center">
      <p className="text-3xl font-mono font-bold text-green-600 tracking-[0.5em]">
        {devCode}
      </p>
    </div>
  </div>
)}
```

---

## Инструкции по тестированию

### Тест 1: DevCode отображается

1. Откройте http://localhost:5173
2. Пройдите до шага "Email"
3. Введите любой email, нажмите "Далее"
4. ✅ **Ожидается:** Вы видите зелёный блок с кодом (6 цифр)
5. Скопируйте код, вставьте в поле ввода
6. ✅ **Ожидается:** Код принимается, переход к следующему шагу

### Тест 2: Прогресс "Шаг X из Y"

1. Откройте http://localhost:5173
2. На первом шаге (выбор класса) посмотрите в правый верхний угол
3. ✅ **Ожидается:** "Шаг 1 из 6" (не "17%")
4. Выберите класс, нажмите "Далее"
5. ✅ **Ожидается:** "Шаг 2 из 6"
6. Пройдите ещё несколько шагов
7. ✅ **Ожидается:** Номер шага увеличивается корректно

### Тест 3: Backend логи

```bash
docker compose -f docker/docker-compose.dev.yml logs -f backend | grep EmailHandler
```

✅ **Ожидается:**
```
[EmailHandler] Verification code for test@example.com: 123456 (expires at ...)
```

### Тест 4: API ответ

```bash
curl -X POST http://localhost:8080/email/verify/send \
  -H 'Content-Type: application/json' \
  -H 'X-Platform-ID: web' \
  -d '{"email":"test@example.com","parentUserId":"123"}'
```

✅ **Ожидается:**
```json
{
  "message": "Verification code sent to email",
  "expiresAt": "2026-04-04T10:00:00Z",
  "devCode": "123456"
}
```

---

## Roadmap: Email интеграция

### Phase 1: Выбор провайдера (production ready)

**Рекомендации:**

1. **SendGrid** (рекомендуется)
   - ✅ Бесплатно 100 писем/день
   - ✅ Простая интеграция
   - ✅ Go SDK: `github.com/sendgrid/sendgrid-go`

2. **AWS SES**
   - ✅ Дёшево ($0.10 за 1000 писем)
   - ✅ Высокая надёжность
   - ⚠️ Требует верификации домена

3. **Mailgun**
   - ✅ Бесплатно 5000 писем/месяц
   - ✅ Хорошая документация

### Phase 2: Реализация

**Файл:** `api/internal/email/sender.go` (новый)

```go
package email

import (
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Sender struct {
    apiKey string
    fromEmail string
    fromName string
}

func (s *Sender) SendVerificationCode(to, code string) error {
    from := mail.NewEmail(s.fromName, s.fromEmail)
    toEmail := mail.NewEmail("", to)
    subject := "Код подтверждения - Объяснятель ДЗ"

    htmlContent := fmt.Sprintf(`
        <h2>Код подтверждения</h2>
        <p>Ваш код: <strong>%s</strong></p>
        <p>Код действителен 15 минут.</p>
    `, code)

    message := mail.NewSingleEmail(from, subject, toEmail, "", htmlContent)
    client := sendgrid.NewSendClient(s.apiKey)

    _, err := client.Send(message)
    return err
}
```

**Файл:** `api/internal/api/handler/email.go` (обновить)

```go
// Отправка email через email service
err = h.emailSender.SendVerificationCode(req.Email, code)
if err != nil {
    log.Printf("[EmailHandler] Failed to send email: %v", err)
    // Не возвращаем ошибку - код уже в БД, пользователь может запросить повторно
}

// В production НЕ возвращаем devCode!
response.OK(w, map[string]interface{}{
    "message": "Verification code sent to email",
    "expiresAt": expiresAt,
    // devCode убран в production
})
```

### Phase 3: HTML шаблоны

Создать красивые email-шаблоны:
- Логотип приложения
- Брендированные цвета
- Адаптивная вёрстка для мобильных
- Кнопка для быстрого копирования кода

---

## Заключение

✅ **Завершено:**
- Dev режим: код показывается в UI
- Прогресс изменён на "Шаг X из Y"
- Улучшена понятность интерфейса

🔜 **TODO (production):**
- Интеграция с email-провайдером
- Создание HTML-шаблонов для писем
- Убрать devCode из API
- Настроить домен для отправки (SPF, DKIM)

---

**Статус:** ✅ Готово к тестированию в dev режиме
**Next Step:** Выбрать email-провайдер для production
