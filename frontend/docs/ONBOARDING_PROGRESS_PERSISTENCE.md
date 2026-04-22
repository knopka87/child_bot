# Сохранение прогресса онбординга

## Обзор

Реализована функциональность автоматического сохранения и восстановления прогресса онбординга, позволяющая пользователю продолжить с того места, где остановился после перезагрузки страницы или закрытия приложения.

## Функциональность

### Автоматическое сохранение

Прогресс онбординга автоматически сохраняется в VK Storage (или localStorage для веб-платформы) при каждом изменении следующих данных:

- **Текущий шаг** (`onboarding_step`): grade, avatar, email, email_verification, consent, completed
- **Выбранный класс** (`onboarding_grade`): 1-4
- **Выбранный аватар** (`onboarding_avatar`): ID аватара
- **Email родителя** (`onboarding_email`): введённый email
- **Статус верификации** (`onboarding_email_verified`): true/false
- **Имя ребёнка** (`onboarding_display_name`): отображаемое имя
- **Согласия** (`onboarding_consents`): JSON с флагами adultConsent, privacyAccepted, termsAccepted

### Восстановление при загрузке

При открытии страницы онбординга:

1. Загружаются данные пользователя из VK Bridge (имя)
2. Проверяется наличие сохранённого прогресса в storage
3. Если найден сохранённый прогресс, состояние восстанавливается:
   - Устанавливается текущий шаг
   - Восстанавливаются все введённые данные
   - Восстанавливается статус верификации email
   - Восстанавливаются отмеченные согласия

### Очистка после завершения

После успешного завершения онбординга все временные данные автоматически удаляются из storage для предотвращения конфликтов при повторном использовании.

## Технические детали

### Storage Keys

Добавлены новые ключи в `vk-storage.ts`:

```typescript
export const storageKeys = {
  // ... существующие ключи

  // Onboarding progress
  ONBOARDING_STEP: 'onboarding_step',
  ONBOARDING_GRADE: 'onboarding_grade',
  ONBOARDING_AVATAR: 'onboarding_avatar',
  ONBOARDING_EMAIL: 'onboarding_email',
  ONBOARDING_EMAIL_VERIFIED: 'onboarding_email_verified',
  ONBOARDING_DISPLAY_NAME: 'onboarding_display_name',
  ONBOARDING_CONSENTS: 'onboarding_consents',
} as const;
```

### Логика сохранения

Используется `useEffect` с зависимостями на все поля онбординга:

```typescript
useEffect(() => {
  const saveProgress = async () => {
    // Сохраняем все данные в storage
    await vkStorage.setItem(storageKeys.ONBOARDING_STEP, currentStep);
    // ... остальные поля
  };

  if (currentStep !== 'completed') {
    saveProgress();
  }
}, [currentStep, grade, avatarId, email, ...]);
```

### Логика восстановления

При монтировании компонента:

```typescript
useEffect(() => {
  const initOnboarding = async () => {
    // 1. Загрузка данных VK
    // 2. Восстановление сохранённого прогресса
    const savedStep = await vkStorage.getItem(storageKeys.ONBOARDING_STEP);
    if (savedStep) {
      setCurrentStep(savedStep as OnboardingStep);
      // ... восстановление остальных полей
    }
  };

  initOnboarding();
}, []);
```

## Поддерживаемые сценарии

### ✅ Поддерживается

1. **Перезагрузка страницы** на любом шаге онбординга
2. **Закрытие и повторное открытие** приложения
3. **Переключение между вкладками** (для VK Mini Apps)
4. **Частичное заполнение данных** с последующим продолжением
5. **Email verification** - код остаётся активным 15 минут

### ⚠️ Ограничения

1. **Срок хранения**: данные хранятся до завершения онбординга или очистки storage
2. **Email код**: имеет срок действия 15 минут (серверная проверка)
3. **Один пользователь**: прогресс привязан к устройству/браузеру, а не к аккаунту VK

## Пользовательский сценарий

### Пример 1: Прерывание на этапе email verification

1. Пользователь вводит класс → **сохранено**
2. Выбирает аватар → **сохранено**
3. Вводит email → **сохранено**
4. Получает код, но **закрывает приложение**
5. Открывает приложение заново → **восстанавливается на шаге email_verification**
6. Вводит код из письма → продолжает онбординг

### Пример 2: Возврат после перерыва

1. Пользователь проходит несколько шагов
2. **Перезагружает страницу**
3. Видит тот же шаг, на котором остановился
4. Все введённые данные сохранены
5. Продолжает с того же места

## Логирование

Для отладки добавлены логи:

```typescript
console.log('[Onboarding] Progress saved:', { step, grade, avatarId, ... });
console.log('[Onboarding] Restoring progress from storage:', { ... });
console.log('[Onboarding] Temporary onboarding data cleared');
```

## Аналитика

События аналитики учитывают восстановление прогресса:
- `onboarding_opened` - отправляется при каждом открытии
- Остальные события отправляются при выполнении действий (не при восстановлении)

## Безопасность

- Email хранится локально только до завершения онбординга
- Код верификации **не хранится** на клиенте (только на сервере)
- После завершения все временные данные удаляются

## Совместимость

Работает на всех платформах:
- ✅ VK Mini Apps (VK Storage)
- ✅ Web (localStorage)
- ✅ Telegram Mini Apps (localStorage)
- ✅ MAX Mini Apps (localStorage)

---

**Дата реализации:** 2026-04-04
**Версия:** 1.0
