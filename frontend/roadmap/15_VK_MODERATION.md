# Phase 14: VK Moderation Preparation

**Длительность:** 2-3 дня
**Приоритет:** Критический
**Зависимости:** Все предыдущие phases

---

## Цель

Подготовить приложение к подаче на модерацию VK, пройти все требования, избежать типичных ошибок отказа.

---

## Часть 1: VK Moderation Checklist

### 1.1. Технические требования

#### ✅ HTTPS обязателен
```bash
# Проверка:
curl -I https://your-miniapp.com

# Должен вернуть 200 OK с валидным SSL сертификатом
```

#### ✅ Sign Validation на Backend
```go
// КРИТИЧНО: backend должен проверять sign параметр
func ValidateVKSign(params map[string]string, secretKey string) bool {
    // Реализация из 01_SETUP.md
    // ...
}
```

**Проверка:**
- Backend проверяет `sign` при каждом запросе аутентификации
- Невалидный `sign` возвращает 403 Forbidden
- Логи показывают успешную валидацию

#### ✅ Bundle Size < 10 MB
```bash
# Проверка размера bundle
npm run build
du -sh dist/

# Должно быть:
# dist/ < 10 MB (сжато)
```

**Оптимизация (если превышен лимит):**
```typescript
// vite.config.ts
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'react-vendor': ['react', 'react-dom'],
          'vk-vendor': ['@vkontakte/vkui', '@vkontakte/vk-bridge'],
        },
      },
    },
  },
  esbuild: {
    drop: ['console', 'debugger'], // Удаляем в prod
  },
});
```

#### ✅ Работает на всех платформах
- [ ] iOS (Safari)
- [ ] Android (Chrome)
- [ ] Desktop Web (Chrome/Firefox)
- [ ] VK Mobile App (iOS/Android)

**Тестирование:**
```typescript
// Используйте VK Bridge для определения платформы
const launchParams = await bridge.send('VKWebAppGetLaunchParams');
console.log('Platform:', launchParams.vk_platform);
```

#### ✅ Нет ошибок в консоли
```javascript
// В production НЕ должно быть:
// - Uncaught errors
// - Failed API requests (кроме ожидаемых 404)
// - CORS errors
// - Mixed content warnings

// Проверка:
window.addEventListener('error', (e) => {
  console.error('Global error:', e);
});
```

---

### 1.2. Функциональные требования

#### ✅ Приложение полностью работает
- Все основные flows завершаются успешно
- Нет "coming soon" заглушек
- Нет битых ссылок или пустых страниц
- Loading states показываются корректно
- Error states обрабатываются gracefully

#### ✅ UI соответствует VK Guidelines
```typescript
// Используем VKUI компоненты
import { Button, Card, Group, Panel } from '@vkontakte/vkui';
import '@vkontakte/vkui/dist/vkui.css';

// Адаптивный layout для разных платформ
<ConfigProvider platform={platform}>
  <AdaptivityProvider>
    <AppRoot>
      {/* Ваше приложение */}
    </AppRoot>
  </AdaptivityProvider>
</ConfigProvider>
```

#### ✅ Обработка всех edge cases
```typescript
// Примеры edge cases:
// 1. Нет интернета
if (!navigator.onLine) {
  showOfflineMessage();
}

// 2. API недоступен
try {
  await api.getData();
} catch (error) {
  if (error.code === 'NETWORK_ERROR') {
    showRetryButton();
  }
}

// 3. Пользователь закрыл приложение во время загрузки
useEffect(() => {
  const handleVisibilityChange = () => {
    if (document.visibilityState === 'visible') {
      refetchData();
    }
  };

  document.addEventListener('visibilitychange', handleVisibilityChange);
  return () => document.removeEventListener('visibilitychange', handleVisibilityChange);
}, []);

// 4. Backend вернул неожиданный ответ
try {
  const data = await api.getData();
  if (!data || !data.profile) {
    throw new Error('Invalid response format');
  }
} catch (error) {
  showGenericError();
}
```

---

### 1.3. Контент и юридические требования

#### ✅ Privacy Policy (Политика конфиденциальности)

**Файл:** `public/privacy-policy.html`

```html
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Политика конфиденциальности</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      line-height: 1.6;
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
    }
    h1, h2 { color: #333; }
  </style>
</head>
<body>
  <h1>Политика конфиденциальности</h1>
  <p><strong>Дата вступления в силу:</strong> 01.01.2026</p>

  <h2>1. Общие положения</h2>
  <p>
    Настоящая Политика конфиденциальности регулирует порядок обработки и защиты
    персональных данных пользователей мини-приложения "Объяснятель ДЗ" (далее — "Приложение").
  </p>

  <h2>2. Какие данные мы собираем</h2>
  <ul>
    <li>ID пользователя ВКонтакте</li>
    <li>Имя и фамилия (из профиля ВКонтакте)</li>
    <li>Фотография профиля (из профиля ВКонтакте)</li>
    <li>Email адрес (если предоставлен пользователем)</li>
    <li>Загруженные изображения домашних заданий</li>
    <li>История попыток и результатов</li>
    <li>Аналитические данные (события, время использования)</li>
  </ul>

  <h2>3. Как мы используем данные</h2>
  <ul>
    <li>Для предоставления функциональности Приложения</li>
    <li>Для улучшения качества сервиса</li>
    <li>Для персонализации опыта пользователя</li>
    <li>Для аналитики и статистики</li>
  </ul>

  <h2>4. Защита данных</h2>
  <p>
    Мы применяем современные методы защиты данных, включая шифрование при передаче
    и хранении, ограничение доступа к данным, регулярные аудиты безопасности.
  </p>

  <h2>5. Передача данных третьим лицам</h2>
  <p>
    Мы не передаем ваши персональные данные третьим лицам, за исключением случаев,
    предусмотренных законодательством РФ.
  </p>

  <h2>6. Права пользователей</h2>
  <p>Вы имеете право:</p>
  <ul>
    <li>Запросить копию ваших данных</li>
    <li>Запросить удаление ваших данных</li>
    <li>Отозвать согласие на обработку данных</li>
  </ul>

  <h2>7. Контакты</h2>
  <p>
    По вопросам конфиденциальности обращайтесь:
    <br>Email: <a href="mailto:privacy@example.com">privacy@example.com</a>
  </p>
</body>
</html>
```

#### ✅ Terms of Service (Пользовательское соглашение)

**Файл:** `public/terms-of-service.html`

```html
<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Пользовательское соглашение</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      line-height: 1.6;
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
    }
    h1, h2 { color: #333; }
  </style>
</head>
<body>
  <h1>Пользовательское соглашение</h1>
  <p><strong>Дата вступления в силу:</strong> 01.01.2026</p>

  <h2>1. Общие положения</h2>
  <p>
    Настоящее Соглашение регулирует отношения между пользователем и администрацией
    мини-приложения "Объяснятель ДЗ".
  </p>

  <h2>2. Предмет соглашения</h2>
  <p>
    Приложение предоставляет образовательные услуги по помощи с домашними заданиями
    для школьников через искусственный интеллект.
  </p>

  <h2>3. Права и обязанности пользователя</h2>
  <p>Пользователь обязуется:</p>
  <ul>
    <li>Не использовать Приложение в противоправных целях</li>
    <li>Не загружать материалы, защищенные авторским правом</li>
    <li>Не пытаться обойти технические ограничения Приложения</li>
  </ul>

  <h2>4. Права и обязанности администрации</h2>
  <p>Администрация имеет право:</p>
  <ul>
    <li>Изменять функциональность Приложения</li>
    <li>Приостанавливать доступ пользователя при нарушении правил</li>
    <li>Удалять контент, нарушающий законодательство</li>
  </ul>

  <h2>5. Ограничение ответственности</h2>
  <p>
    Приложение предоставляется "как есть". Администрация не несет ответственности
    за точность ответов ИИ и рекомендует проверять результаты.
  </p>

  <h2>6. Контакты</h2>
  <p>
    Email: <a href="mailto:support@example.com">support@example.com</a>
  </p>
</body>
</html>
```

#### ✅ Ссылки на документы в приложении

```typescript
// src/pages/Profile/SettingsPage.tsx
import { Group, Cell } from '@vkontakte/vkui';
import { Icon24ExternalLinkOutline } from '@vkontakte/icons';
import bridge from '@vkontakte/vk-bridge';

function SettingsPage() {
  const openPrivacyPolicy = () => {
    bridge.send('VKWebAppOpenLink', {
      url: 'https://your-domain.com/privacy-policy.html',
    });
  };

  const openTerms = () => {
    bridge.send('VKWebAppOpenLink', {
      url: 'https://your-domain.com/terms-of-service.html',
    });
  };

  return (
    <Group>
      <Cell
        after={<Icon24ExternalLinkOutline />}
        onClick={openPrivacyPolicy}
      >
        Политика конфиденциальности
      </Cell>
      <Cell
        after={<Icon24ExternalLinkOutline />}
        onClick={openTerms}
      >
        Пользовательское соглашение
      </Cell>
    </Group>
  );
}
```

---

### 1.4. Монетизация (если используется)

#### ✅ VK Pay интегрирован корректно
- Используется официальный VK Pay API
- Webhook обрабатывает уведомления
- Подписки активируются автоматически
- Есть возможность отменить подписку

#### ✅ Реклама соответствует правилам
- Используется VK Ads SDK
- Реклама не мешает основному функционалу
- Rewarded ads дают реальную награду
- Интервал между interstitial ads разумный (не чаще 1 раз в 5 минут)

---

## Часть 2: Типичные причины отказа

### ❌ Причина 1: "Приложение не работает"

**Проблема:** Модератор не может войти или использовать основной функционал.

**Решение:**
1. Создайте тестовый аккаунт для модератора
2. Убедитесь, что онбординг понятен
3. Добавьте подсказки для первого использования
4. Протестируйте на чистом профиле

```typescript
// Добавьте debug режим для модераторов
const DEBUG_USER_IDS = [123456789]; // VK ID модератора

if (DEBUG_USER_IDS.includes(currentUserId)) {
  console.log('[DEBUG] Модератор вошел в приложение');
  // Логируйте все важные события
}
```

---

### ❌ Причина 2: "Нарушение privacy"

**Проблема:** Нет политики конфиденциальности или она неполная.

**Решение:**
- Добавьте ссылку на privacy policy в настройках приложения
- Опишите все собираемые данные
- Укажите цели обработки данных
- Добавьте контакты для запросов

---

### ❌ Причина 3: "UI не соответствует VK Guidelines"

**Проблема:** Кастомный дизайн не похож на VK.

**Решение:**
```typescript
// Используйте VKUI везде
import { Button, Card, Group, Panel } from '@vkontakte/vkui';
import '@vkontakte/vkui/dist/vkui.css';

// НЕ создавайте кастомные компоненты без необходимости
```

---

### ❌ Причина 4: "Ошибки в работе"

**Проблема:** Модератор столкнулся с ошибкой.

**Решение:**
```typescript
// Глобальный error boundary
<ErrorBoundary
  fallback={
    <Placeholder
      icon={<Icon56ErrorOutline />}
      header="Что-то пошло не так"
      action={<Button onClick={() => window.location.reload()}>Обновить</Button>}
    >
      Мы уже работаем над исправлением
    </Placeholder>
  }
>
  <App />
</ErrorBoundary>

// Логируйте все ошибки
window.addEventListener('error', (e) => {
  // Отправьте в Sentry/LogRocket
  console.error('Error:', e);
});
```

---

### ❌ Причина 5: "Монетизация не соответствует правилам"

**Проблема:** Подписка или реклама работают неправильно.

**Решение:**
- VK Pay: проверьте webhook, sign validation
- Ads: используйте только VK Ads SDK
- Не блокируйте весь функционал для free users
- Дайте возможность попробовать перед покупкой

---

## Часть 3: Timing и процесс модерации

### 3.1. Когда подавать на модерацию

**Лучшее время:**
- **Вторник-Среда** (10:00-16:00 МСК)
- Избегайте пятницы и понедельника
- Не подавайте в праздники

**Релизы:**
- VK выпускает обновления по **четвергам**
- Если прошли модерацию в среду, релиз будет в четверг
- Если прошли в четверг, релиз будет через неделю

### 3.2. Сроки модерации

- **Быстрая:** 1-2 дня (если все ОК)
- **Средняя:** 3-5 дней (если есть вопросы)
- **Долгая:** 7+ дней (если много замечаний)

### 3.3. Что писать в описании для модератора

```
Тестовые данные:
- Email: moderator@test.com
- Класс: 5
- Для тестирования используйте любое фото с математическим примером

Основной функционал:
1. Онбординг: выбор класса, имени, аватара
2. Помощь с ДЗ: загрузка фото, получение подсказок
3. Проверка работы: загрузка задания и ответа, получение обратной связи
4. Достижения: система наград за прогресс
5. Профиль: история попыток, настройки

Монетизация:
- Подписка через VK Pay (тестовый режим)
- Rewarded ads для дополнительных подсказок
- Interstitial ads каждые 5 попыток

Контакты:
- Email: support@example.com
- Telegram: @support_bot
```

---

## Часть 4: После одобрения

### 4.1. Первый релиз

После одобрения:
1. Проверьте, что приложение доступно в каталоге VK
2. Протестируйте на production окружении
3. Мониторьте ошибки в Sentry/LogRocket
4. Отслеживайте метрики: DAU, retention, конверсия

### 4.2. Обновления

Для обновлений:
- Мелкие фиксы (bugfix) проходят быстрее
- Новые features требуют полной модерации
- Изменения в монетизации проверяются строго
- UI изменения проверяются на соответствие guidelines

---

## Чеклист перед подачей на модерацию

### Технические требования
- [ ] HTTPS настроен и работает
- [ ] Sign validation на backend реализована
- [ ] Bundle size < 10 MB
- [ ] Приложение работает на iOS
- [ ] Приложение работает на Android
- [ ] Приложение работает на Desktop
- [ ] Нет критических ошибок в консоли
- [ ] Все API endpoints отвечают корректно

### Функциональные требования
- [ ] Все основные flows завершаются
- [ ] Нет заглушек "coming soon"
- [ ] Loading states показываются
- [ ] Error states обрабатываются
- [ ] Onboarding понятен и завершаем
- [ ] UI использует VKUI компоненты
- [ ] Адаптивность под разные платформы

### Контент и юридические
- [ ] Privacy Policy создана и доступна
- [ ] Terms of Service созданы и доступны
- [ ] Ссылки на документы в приложении
- [ ] Нет нарушений авторских прав
- [ ] Контент соответствует правилам VK

### Монетизация (если есть)
- [ ] VK Pay интегрирован корректно
- [ ] Webhook обрабатывает платежи
- [ ] Подписки активируются автоматически
- [ ] VK Ads SDK интегрирован
- [ ] Rewarded ads работают
- [ ] Interstitial ads не раздражают
- [ ] Есть бесплатный функционал

### Подготовка к подаче
- [ ] Создан тестовый аккаунт для модератора
- [ ] Написано описание для модератора
- [ ] Указаны тестовые данные
- [ ] Загружены скриншоты (минимум 3)
- [ ] Заполнены все обязательные поля
- [ ] Выбран правильный timing (вторник-среда)

### После одобрения
- [ ] Проверить доступность в каталоге
- [ ] Протестировать на production
- [ ] Настроить мониторинг ошибок
- [ ] Отслеживать метрики
- [ ] Собрать обратную связь от пользователей

---

## Escalation Process

Если модерация затягивается или отказывают без понятной причины:

1. **Повторная подача** (если исправили замечания)
   - Опишите, что конкретно исправили
   - Прикрепите скриншоты до/после

2. **Обращение в техподдержку VK**
   - vk.com/support
   - Приложите:
     - App ID
     - Дату подачи
     - Причину отказа
     - Скриншоты исправлений

3. **Community VK Разработчиков**
   - vk.com/apiclub
   - Задайте вопрос в обсуждениях
   - Другие разработчики могут помочь

---

## Полезные ссылки

- [VK Mini Apps Guidelines](https://dev.vk.com/mini-apps/development/creating-app)
- [VK Moderation Rules](https://dev.vk.com/mini-apps/review)
- [VK Bridge Documentation](https://dev.vk.com/bridge/overview)
- [VKUI Components](https://vkcom.github.io/VKUI/)

---

## Следующие шаги

После прохождения модерации:
1. Мониторинг метрик (DAU, retention, конверсия)
2. Сбор обратной связи
3. Планирование следующих features
4. A/B тестирование улучшений
