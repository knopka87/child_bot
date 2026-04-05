# Frontend Roadmap — Индекс

**Проект:** Объяснятель ДЗ MiniApp
**Дата создания:** 2026-03-29
**Платформы:** VK (primary), Max, Telegram (future)

---

## 📚 Содержание

Этот roadmap содержит полное техническое описание разработки frontend миниапп приложения для помощи школьникам с домашними заданиями.

### 📋 Обзорные документы

| Файл | Описание |
|------|----------|
| **[00_OVERVIEW.md](./00_OVERVIEW.md)** | Общий план проекта, этапы разработки, технологический стек |
| **[SCREEN_MAP.md](./SCREEN_MAP.md)** | Визуальная карта всех экранов и навигации |
| **[API_DATA_REQUIREMENTS.md](./API_DATA_REQUIREMENTS.md)** | Полное описание API endpoints и типов данных |
| **[SECURITY.md](./SECURITY.md)** | 🔐 Security guidelines, аутентификация, защита данных |
| **[ANALYTICS_COVERAGE.md](./ANALYTICS_COVERAGE.md)** | 📊 Покрытие аналитических событий (94%) |
| **[COMPONENT_ARCHITECTURE.md](./COMPONENT_ARCHITECTURE.md)** | 🏗️ Многослойная компонентная архитектура |
| **[REVIEW_RESULTS.md](./REVIEW_RESULTS.md)** | ✅ Результаты проверки roadmap |
| **[VK_BEST_PRACTICES.md](./VK_BEST_PRACTICES.md)** | ⚡ VK Mini Apps best practices |
| **[VK_ROADMAP_IMPROVEMENTS.md](./VK_ROADMAP_IMPROVEMENTS.md)** | 📋 Улучшения roadmap для VK |

### 🛠️ Детальные roadmap по этапам

| Этап | Файл | Длительность | Описание |
|------|------|--------------|----------|
| **Phase 0** | [01_SETUP.md](./01_SETUP.md) | 1-2 дня | Настройка проекта: Vite, TypeScript, VK Bridge, структура папок |
| **Phase 1** | [02_CORE.md](./02_CORE.md) | 3-5 дней | Базовая инфраструктура: UI Kit, API client, роутинг, state management |
| **Phase 2** | [03_ONBOARDING.md](./03_ONBOARDING.md) | 3-4 дня | Онбординг: выбор класса, аватара, имени, email verification |
| **Phase 3** | [04_HOME.md](./04_HOME.md) | 4-5 дней | Главный экран: header, персонажи, действия, незаконченные попытки |
| **Phase 4** | [05_HELP.md](./05_HELP.md) | 5-7 дней | Поток помощи: загрузка фото, обработка, подсказки, ответ |
| **Phase 5** | [06_CHECK.md](./06_CHECK.md) | 5-7 дней | Поток проверки: сценарии, загрузка, результат, ошибки |
| **Phase 6** | [07_ACHIEVEMENTS.md](./07_ACHIEVEMENTS.md) | 2-3 дня | Достижения: полки, разблокировка, награды |
| **Phase 7** | [08_FRIENDS.md](./08_FRIENDS.md) | 2-3 дня | Друзья: реферальная система, приглашения, награды |
| **Phase 8** | [09_PROFILE.md](./09_PROFILE.md) | 3-4 дня | Профиль: история, отчеты, поддержка, подписка |
| **Phase 9** | [10_VILLAIN.md](./10_VILLAIN.md) | 3-4 дня | Злодей: битва, здоровье, победа, награды |
| **Phase 10** | [11_ANALYTICS.md](./11_ANALYTICS.md) | 2-3 дня | Аналитика: интеграция событий, валидация, отправка |
| **Phase 11** | [12_TESTING.md](./12_TESTING.md) | 3-5 дней | Тестирование: unit, integration, E2E, performance |
| **Phase 12** | [13_PLATFORMS.md](./13_PLATFORMS.md) | 5-7 дней | Платформы: адаптация под Max и Telegram |
| **Phase 13** | [14_MONETIZATION.md](./14_MONETIZATION.md) | 4-5 дней | Монетизация: VK Pay, Ads SDK, подписки |
| **Phase 14** | [15_VK_MODERATION.md](./15_VK_MODERATION.md) | 2-3 дня | Подготовка к модерации VK, чеклист, типичные ошибки |
| **Phase 15** | [16_BACKEND_API.md](./16_BACKEND_API.md) | 11-12 дней | Backend REST API: миграция с Telegram Bot, все endpoints для миниаппа |

---

## 🎯 Ключевые особенности roadmap

### ✅ Что включено в каждый файл

- **Детальная архитектура** - структура компонентов, зависимости
- **TypeScript типы** - полная типизация всех данных и API
- **Примеры кода** - готовые компоненты и hooks
- **API интеграция** - endpoints, request/response types
- **Аналитика** - все события с параметрами из реестра
- **Обработка ошибок** - error boundaries, fallbacks
- **CSS стили** - примеры стилей для компонентов
- **Чеклисты** - пошаговые задачи для отслеживания прогресса

### 📊 Основано на данных

Roadmap создан на основе:
- ✅ **Анализа дизайна** (21 экран, 15+ динамических элементов)
- ✅ **Реестра аналитики** (100+ событий, user properties)
- ✅ **Backend API требований** (25+ endpoints)
- ✅ **VK miniapp best practices**

---

## 🚀 Как использовать roadmap

### Для разработчиков

1. **Начни с обзора** - прочитай `00_OVERVIEW.md` для понимания общей картины
2. **Настрой проект** - следуй `01_SETUP.md` для инициализации
3. **Работай последовательно** - выполняй phases по порядку (01 → 02 → ... → 13)
4. **Используй чеклисты** - отмечай выполненные задачи в каждом файле
5. **Изучай примеры кода** - копируй и адаптируй готовые компоненты

### Для product/project менеджеров

1. **Отслеживай прогресс** - каждый phase имеет оценку времени
2. **Проверяй чеклисты** - Definition of Done в конце каждого этапа
3. **Координируй с backend** - используй `API_DATA_REQUIREMENTS.md`
4. **Планируй тестирование** - `12_TESTING.md` содержит test plan

### Для дизайнеров

1. **Смотри карту экранов** - `SCREEN_MAP.md` показывает все flows
2. **Проверяй компоненты** - каждый roadmap содержит UI компоненты
3. **Уточняй состояния** - все состояния (loading, error, success) описаны

---

## 📈 Оценка сроков

### Минимальный MVP (Phase 0-5)
**Срок:** 18-25 дней
**Функциональность:**
- Онбординг
- Главный экран
- Поток помощи
- Поток проверки
- Базовая аналитика

### Полная версия (Phase 0-11)
**Срок:** 30-42 дня
**Функциональность:**
- Все MVP функции
- Достижения
- Друзья и рефералы
- Профиль и история
- Злодей и игровая механика
- Полная аналитика
- Тестирование

### VK Production Ready (Phase 0-14)
**Срок:** 36-50 дней
**Функциональность:**
- Все функции полной версии
- VK Pay интеграция
- VK Ads SDK
- Подписки
- Готовность к модерации VK

### Кроссплатформенная версия с Backend (Phase 0-15)
**Срок:** 47-62 дня
**Функциональность:**
- Все функции VK версии
- Адаптация под Max и Telegram
- **Backend REST API миграция (11-12 дней)**
- Полная интеграция frontend и backend

---

## 🔗 Связанные документы

### Из проекта
- **[ANALYTICS_EVENTS_REGISTRY_Obiasnyatel_DZ_MiniApp.md](../ANALYTICS_EVENTS_REGISTRY_Obiasnyatel_DZ_MiniApp.md)** - реестр аналитических событий
- **[analysis_report.json](../analysis_report.json)** - анализ дизайна (JSON)
- **[analysis_report.md](../analysis_report.md)** - анализ дизайна (Markdown)
- **[DESIGN_ANALYSIS.md](../DESIGN_ANALYSIS.md)** - анализ дизайна (если есть)

### Внешние ресурсы
- [VK Bridge Documentation](https://dev.vk.com/bridge/overview)
- [VK Mini Apps Dev Guide](https://dev.vk.com/mini-apps/development)
- [React Documentation](https://react.dev)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Vite Guide](https://vitejs.dev/guide/)

---

## 🎨 Технологический стек

### Core
- **React 18+** - UI framework
- **TypeScript 5+** - type safety
- **Vite 5+** - build tool
- **React Router v6** - routing

### State Management
- **Zustand** - global state
- **React Query** - server state
- **React Hook Form** - forms

### Platform SDKs
- **@vkontakte/vk-bridge** - VK integration
- **@twa-dev/sdk** - Telegram WebApp (future)
- **Max SDK** - Max integration (future)

### UI & Styling
- **CSS Modules** - scoped styles
- **@vkontakte/vkui** - VK UI components (optional)
- **browser-image-compression** - image optimization

### Testing
- **Vitest** - unit tests
- **React Testing Library** - component tests
- **Playwright** - E2E tests

### Dev Tools
- **ESLint** - linting
- **Prettier** - formatting
- **TypeScript ESLint** - type linting

---

## 📞 Контакты и поддержка

Если у вас есть вопросы по roadmap:
1. Проверьте соответствующий phase файл
2. Посмотрите `API_DATA_REQUIREMENTS.md` для вопросов по данным
3. Используйте `SCREEN_MAP.md` для понимания навигации

---

## 📝 История изменений

### 2026-03-29 - Initial Release
- Создан полный roadmap (13 phases)
- Добавлены API требования
- Добавлена карта экранов
- Добавлены детальные чеклисты

---

## ✅ Следующие шаги

1. **Прочитай** `00_OVERVIEW.md` - понять общую картину
2. **Настрой** проект по `01_SETUP.md`
3. **Согласуй** API с backend командой
4. **Начни** разработку с Phase 1 (`02_CORE.md`)
5. **Отслеживай** прогресс по чеклистам
