# Monetization Strategy

## Overview

Приложение использует гибридную модель монетизации:
- **Подписка (Subscription)** - основной источник дохода
- **Внутриигровые покупки (IAP)** - дополнительные функции
- **Реклама (Ads)** - для бесплатных пользователей

## Subscription Plans

### Trial (Пробный период)
- **Длительность**: 7 дней
- **Цена**: Бесплатно
- **Особенности**:
  - Полный доступ ко всем функциям
  - Без рекламы
  - Автоматически активируется при регистрации

### Monthly (Месячная подписка)
- **Цена**: 299 ₽/месяц
- **Особенности**:
  - Неограниченные задания
  - Приоритетная поддержка
  - Без рекламы
  - Доступ к эксклюзивным стикерам
  - Еженедельный отчёт родителю

### Yearly (Годовая подписка)
- **Цена**: 2499 ₽/год (40% скидка)
- **Особенности**:
  - Всё из месячной подписки
  - 40% экономия
  - Приоритет в обработке заданий

## Payment Methods

### VK Pay
- **Платформы**: VK, VK Max
- **Комиссия**: 5%
- **Поддержка**: Карты, СБП, Кошелёк VK

### Telegram Stars
- **Платформы**: Telegram
- **Комиссия**: 7%
- **Поддержка**: Telegram Stars

### Web Payments
- **Платформы**: Web
- **Провайдер**: CloudPayments / YooMoney
- **Комиссия**: 3.5%
- **Поддержка**: Карты, СБП, Apple Pay, Google Pay

## In-App Purchases

### Consumables (Расходуемые)
- **Монеты (Coins)**: 100, 500, 1000, 5000
- **Подсказки (Hints)**: Пакеты по 5, 10, 20 штук
- **Попытки (Attempts)**: Дополнительные попытки на задания

### Non-Consumables (Постоянные)
- **Эксклюзивные аватары**: 49-99 ₽
- **Стикеры**: 29-79 ₽
- **Темы оформления**: 99 ₽
- **Отключение рекламы навсегда**: 399 ₽

## Ads Integration

### VK Ads SDK
```typescript
const adConfig: AdConfig = {
  enabled: !hasActiveSubscription,
  provider: 'vk_ads',
  placementId: 'YOUR_PLACEMENT_ID',
  frequency: 5, // Показывать каждые 5 заданий
};
```

### Ad Types
- **Interstitial**: После завершения задания (бесплатные пользователи)
- **Rewarded**: За просмотр - дополнительная попытка или монеты
- **Banner**: Снизу экрана на главной странице

### Ad Rules
- ❌ Не показывать во время решения задания
- ❌ Не показывать пользователям с активной подпиской
- ✅ Показывать только после полезного действия (завершение задания)
- ✅ Опция "Пропустить после 5 секунд"
- ✅ Rewarded ads всегда опциональны

## Paywall Strategy

### Триггеры показа Paywall:
1. **Окончание пробного периода** (день 7)
2. **Лимит заданий** (5 заданий/день для бесплатных)
3. **Доступ к премиум-функциям**:
   - Отчёты родителю
   - Эксклюзивные аватары
   - Отключение рекламы
4. **Подсказки закончились** (предложение купить)

### UI/UX Paywall:
- Чёткое отображение преимуществ
- Сравнение планов (Monthly vs Yearly)
- Отзывы довольных пользователей
- "Попробуй 7 дней бесплатно"
- Кнопка "Восстановить покупки" (iOS)

## Analytics Events

### Subscription Events
```typescript
trackEvent('paywall_opened', {
  trigger: 'trial_ended' | 'task_limit' | 'premium_feature',
  plan_shown: 'monthly' | 'yearly'
});

trackEvent('payment_started', {
  plan_id: string,
  method: string
});

trackEvent('payment_success', {
  plan_id: string,
  amount: number,
  currency: string
});

trackEvent('payment_failed', {
  plan_id: string,
  error: string
});
```

### IAP Events
```typescript
trackEvent('iap_product_viewed', {
  product_id: string,
  product_type: 'consumable' | 'non_consumable'
});

trackEvent('iap_purchase_completed', {
  product_id: string,
  price: number
});
```

## Implementation Checklist

### API Integration
- [x] Типы для подписок и платежей
- [x] API методы для монетизации
- [ ] Интеграция с VK Pay SDK
- [ ] Интеграция с Telegram Stars
- [ ] Интеграция с CloudPayments/YooMoney

### UI Components
- [ ] Paywall screen с планами подписки
- [ ] Payment method selector
- [ ] Subscription management screen
- [ ] IAP shop screen
- [ ] Success/Error payment modals

### Business Logic
- [ ] Проверка статуса подписки
- [ ] Лимиты для бесплатных пользователей
- [ ] Автоматическое продление
- [ ] Обработка отмены подписки
- [ ] Восстановление покупок

### Ads
- [ ] VK Ads SDK integration
- [ ] Interstitial ads после заданий
- [ ] Rewarded ads за монеты/попытки
- [ ] Banner ads на главной

### Testing
- [ ] Тестовые платежи
- [ ] Sandbox окружение
- [ ] Восстановление покупок
- [ ] Отмена подписки
- [ ] Обработка ошибок

## Compliance & Legal

### App Store Guidelines
- Прозрачное отображение цен
- Возможность отмены подписки
- Чёткие условия возврата
- Родительский контроль для детей

### VK Модерация
- Соответствие правилам VK Pay
- Корректные описания товаров
- Отсутствие обмана пользователей

### GDPR & Privacy
- Согласие на обработку платёжных данных
- Безопасное хранение токенов
- Шифрование sensitive данных
