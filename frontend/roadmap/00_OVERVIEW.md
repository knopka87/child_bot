# Frontend MiniApp Roadmap — Общий План

**Дата создания:** 2026-03-29
**Проект:** Объяснятель ДЗ MiniApp
**Платформы:** VK (primary), Max, Telegram (future)

---

## Цель

Разработать кроссплатформенное миниапп-приложение для помощи школьникам с домашними заданиями с поддержкой VKontakte, Max и Telegram.

---

## Структура проекта

### Анализ завершен
- ✅ Проанализирован дизайн (21 экран)
- ✅ Извлечены динамические данные
- ✅ Построена карта навигации
- ✅ Определены требования к Backend API

### Динамические данные с Backend

| Переменная | Тип | Экраны | API Endpoint |
|-----------|------|--------|-------------|
| `level` | integer | Главная, Header | `GET /api/v1/child-profile/:id` |
| `level_progress_percent` | integer (0-100) | Header | `GET /api/v1/child-profile/:id` |
| `coins_balance` | integer | Главная, Header | `GET /api/v1/child-profile/:id` |
| `tasks_solved_correct_count` | integer | Главная, Header | `GET /api/v1/child-profile/:id` |
| `achievements_unlocked` | integer | Достижения | `GET /api/v1/achievements` |
| `achievements_total` | integer | Достижения | `GET /api/v1/achievements` |
| `invited_count` | integer | Друзья | `GET /api/v1/referrals` |
| `invited_target` | integer | Друзья | `GET /api/v1/referrals` |
| `villain_id` | string | Главная | `GET /api/v1/villain/active` |
| `villain_health_percent` | integer (0-100) | Главная | `GET /api/v1/villain/active` |
| `mascot_id` | string | Главная | `GET /api/v1/child-profile/:id` |
| `mascot_state` | enum | Главная | `GET /api/v1/child-profile/:id` |
| `unfinished_attempt` | object | Главная | `GET /api/v1/attempts/unfinished` |
| `recent_attempts` | array | Главная | `GET /api/v1/attempts/recent` |

---

## Этапы разработки

### Phase 0: Подготовка (1-2 дня)
**Файл:** `01_SETUP.md`
- Инициализация проекта (Vite + React + TypeScript)
- Настройка VK Bridge SDK
- Настройка роутинга и state management
- Конфигурация аналитики
- Настройка темизации для платформ

### Phase 1: Core Infrastructure (3-5 дней)
**Файл:** `02_CORE.md`
- Базовые компоненты UI Kit
- Layout и навигация
- API client и типизация
- Error boundary и обработка ошибок
- Интеграция с VK miniapp

### Phase 2: Онбординг и Регистрация (3-4 дня)
**Файл:** `03_ONBOARDING.md`
- Экраны онбординга
- Выбор класса и аватара
- Ввод имени
- Согласия и email verification
- Сохранение профиля

### Phase 3: Главный экран (Home) (4-5 дней)
**Файл:** `04_HOME.md`
- Header с уровнем, монетами, счетчиком
- Персонажи (маскот и злодей)
- Основные кнопки действий
- Модал незаконченного задания
- Последние попытки
- Нижняя навигация

### Phase 4: Поток "Помоги разобраться" (5-7 дней)
**Файл:** `05_HELP.md`
- Выбор источника изображения
- Загрузка и crop изображения
- Проверка качества
- Экран обработки (loading, long-wait)
- Экран результата с подсказками
- Навигация по подсказкам
- Отправка ответа

### Phase 5: Поток "Проверка ДЗ" (5-7 дней)
**Файл:** `06_CHECK.md`
- Выбор сценария (1 фото / 2 фото)
- Загрузка фото задания и ответа
- Проверка качества
- Экран обработки
- Экран результата
- Отображение ошибок
- Исправление и повторная отправка

### Phase 6: Достижения (2-3 дня)
**Файл:** `07_ACHIEVEMENTS.md`
- Список достижений (полки)
- Карточка достижения
- Разблокированные/заблокированные состояния
- Анимация получения достижения

### Phase 7: Друзья и Реферальная система (2-3 дня)
**Файл:** `08_FRIENDS.md`
- Прогресс приглашений
- Генерация реферальной ссылки
- Копирование и отправка ссылки
- Список приглашенных друзей

### Phase 8: Профиль и Настройки (3-4 дня)
**Файл:** `09_PROFILE.md`
- Информация о профиле
- История попыток
- Фильтрация истории
- Детальная карточка попытки
- Отчет родителю
- Поддержка
- Настройки подписки
- Paywall

### Phase 9: Злодей и Игровая механика (3-4 дня)
**Файл:** `10_VILLAIN.md`
- Отображение злодея и здоровья
- Экран злодея с репликами
- Битва и урон
- Экран победы
- Награды за победу

### Phase 10: Аналитика (2-3 дня)
**Файл:** `11_ANALYTICS.md`
- Интеграция событий из реестра
- Отправка событий на каждом экране
- Отправка user properties
- Отладка и валидация событий

### Phase 11: Тестирование и Оптимизация (3-5 дней)
**Файл:** `12_TESTING.md`
- E2E тесты критических потоков
- Unit тесты компонентов
- Тестирование на реальных устройствах VK
- Оптимизация производительности
- Проверка доступности

### Phase 12: Адаптация под Max и Telegram (5-7 дней)
**Файл:** `13_PLATFORMS.md`
- Абстракция платформенных SDK
- Адаптация UI под Max
- Адаптация UI под Telegram
- Тестирование на всех платформах

### Phase 13: Монетизация (4-5 дней)
**Файл:** `14_MONETIZATION.md`
- Интеграция VK Pay
- VK Ads SDK
- Подписки и управление доступом

### Phase 14: Модерация VK (2-3 дня)
**Файл:** `15_VK_MODERATION.md`
- Подготовка к модерации
- Чеклист готовности
- Типичные ошибки

### Phase 15: Backend REST API (11-12 дней)
**Файл:** `16_BACKEND_API.md`
- Миграция с Telegram Bot на REST API
- Все endpoints для миниаппа
- WebSocket для real-time уведомлений
- Адаптация существующей бизнес-логики

---

## Технологический стек

### Frontend
- **Framework:** React 18+ с TypeScript
- **Build Tool:** Vite
- **State Management:** Zustand / Jotai
- **Routing:** React Router v6
- **UI Components:** Custom UI Kit + VK UI (опционально)
- **Styling:** CSS Modules / Styled Components / Tailwind
- **Forms:** React Hook Form + Zod
- **HTTP Client:** Axios / Fetch
- **Image Processing:** Browser-image-compression
- **Analytics:** Custom analytics service

### Platform SDKs
- **VK:** @vkontakte/vk-bridge, @vkontakte/vkui (optional)
- **Max:** Max SDK (TBD)
- **Telegram:** @twa-dev/sdk

### Dev Tools
- **Linting:** ESLint + Prettier
- **Type Checking:** TypeScript strict mode
- **Testing:** Vitest + React Testing Library + Playwright
- **Deployment:** Static hosting (Vercel / Cloudflare Pages)

---

## Ключевые компоненты архитектуры

### 1. Роутинг
```
/onboarding
/home
/help
  /upload
  /crop
  /processing
  /result
/check
  /scenario
  /upload
  /processing
  /result
/achievements
/friends
/profile
  /history
  /settings
  /support
/villain
/paywall
```

### 2. State Management
- **Global State:** user profile, auth, platform info
- **Screen State:** local state для каждого экрана
- **Server State:** React Query для кеширования API данных

### 3. API Integration
```typescript
interface APIClient {
  // Profile
  getProfile(childProfileId: string): Promise<ChildProfile>
  updateProfile(data: Partial<ChildProfile>): Promise<void>

  // Attempts
  createAttempt(mode: 'help' | 'check'): Promise<Attempt>
  uploadImage(attemptId: string, image: Blob): Promise<AttemptImage>
  processAttempt(attemptId: string): Promise<void>
  getAttemptResult(attemptId: string): Promise<AttemptResult>

  // Achievements
  getAchievements(): Promise<Achievement[]>

  // Villain
  getActiveVillain(): Promise<Villain>

  // Referrals
  getReferralInfo(): Promise<ReferralInfo>

  // Analytics
  sendEvent(event: AnalyticsEvent): void
}
```

### 4. Платформенная абстракция
```typescript
interface PlatformBridge {
  init(): Promise<void>
  getUser(): Promise<PlatformUser>
  shareLink(url: string): Promise<void>
  copyToClipboard(text: string): Promise<void>
  openURL(url: string): Promise<void>
  hapticFeedback(type: 'light' | 'medium' | 'heavy'): void
  requestPhotoAccess(): Promise<boolean>
  showPopup(params: PopupParams): Promise<void>
}
```

---

## Критерии готовности (Definition of Done)

Каждая фаза считается завершенной, когда:
- ✅ Все экраны реализованы согласно дизайну
- ✅ Интеграция с Backend API завершена
- ✅ Аналитические события отправляются
- ✅ Обработка ошибок реализована
- ✅ Адаптивность под мобильные экраны
- ✅ Код отревьюен и протестирован
- ✅ Документация обновлена

---

## Риски и зависимости

### Риски
1. **Backend API готовность** - фронт зависит от готовых эндпоинтов
2. **Платформенные ограничения VK/Max** - могут быть лимиты на функциональность
3. **Производительность обработки изображений** - большие файлы
4. **Согласование UX между платформами** - разные HIG

### Митигация
- Mock API для независимой разработки
- Feature flags для платформенных отличий
- Оптимизация изображений на клиенте
- Платформенные UI Kit'ы

---

---

## Оценка сроков (обновлено)

### Минимальный MVP (Phase 0-5)
**Срок:** 18-25 дней
- Онбординг, главный экран
- Поток помощи и проверки
- Базовая аналитика

### Полная версия (Phase 0-11)
**Срок:** 30-42 дня
- Все MVP функции
- Достижения, друзья, профиль
- Злодей и игровая механика
- Полная аналитика, тестирование

### VK Production Ready (Phase 0-14)
**Срок:** 36-50 дней
- Все функции полной версии
- VK Pay, Ads SDK, подписки
- Готовность к модерации VK

### Кроссплатформенная версия с Backend (Phase 0-15)
**Срок:** 47-62 дня
- Все функции VK версии
- Адаптация под Max и Telegram
- **Backend REST API миграция (11-12 дней)**

---

## Следующие шаги

1. Читай `01_SETUP.md` для старта проекта
2. Согласуй API контракты с backend командой
3. Подготовь дизайн-токены и UI Kit
4. Настрой CI/CD pipeline
