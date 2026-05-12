# Вопросы по результатам проверки Roadmap Frontend

**Дата проверки:** 2026-04-04
**Общий прогресс выполнения:** ~85%
**Готовность к релизу:** Требуется доработка критичных элементов

---

## 🔴 Критические вопросы (блокирующие релиз)

### 1. Тестирование (Phase 11)
**Статус:** ❌ Не реализовано
**Проблема:** В проекте нет ни одного теста (0 .test.tsx/.test.ts файлов)

**Вопросы:**
- Нужно ли реализовать unit тесты для критических компонентов?
- Какой минимальный coverage требуется для релиза?

**Рекомендации для минимального покрытия:**
- [ ] API client тесты
- [ ] Analytics service тесты
- [ ] Platform adapters тесты
- [ ] Key user flows (onboarding, help, check)
- [ ] Error boundary тесты

**Файлы для создания:**
```
tests/unit/api/client.test.ts
tests/unit/services/analytics.test.ts
tests/unit/platform/adapters.test.ts
tests/integration/flows/onboarding.test.tsx
tests/integration/flows/help.test.tsx
```

---

### 2. EmailVerificationWaiting компонент
**Статус:** ❌ Отсутствует
**Проблема:** Не найден экран ожидания верификации email

**Вопросы:**
- Как реализована верификация email в данный момент?
- Должен ли быть polling каждые 3 секунды как описано в roadmap?
- Какой timeout для верификации? (roadmap предлагает 2 минуты)

**Ожидаемый файл:**
```
src/pages/Onboarding/screens/EmailVerificationWaiting.tsx
```

**Требуемый функционал (из roadmap):**
- Polling `/api/v1/consent/check?childProfileId=xxx` каждые 3 секунды
- Timeout через 2 минуты
- Кнопка "Отправить письмо повторно"
- Кнопка "Пропустить" (сохранить с is_verified=false)

---

### 3. VK Pay Integration
**Статус:** ⚠️ Требует проверки
**Проблема:** API monetization.ts создан (2.4KB), но интеграция не проверена

**Вопросы:**
- Реализована ли полная интеграция с VK Pay?
- Есть ли VK_APP_ID в production environment?
- Протестирована ли покупка подписки в sandbox VK?
- Работают ли VK Bridge методы:
  - `VKWebAppOpenPayForm`
  - `VKWebAppCheckNativeAds`

**Файлы для проверки:**
```
src/api/monetization.ts
src/types/monetization.ts
.env (production)
```

---

### 4. VK Moderation Checklist
**Статус:** ⚠️ ~30% выполнено
**Проблема:** Roadmap содержит 50+ пунктов для модерации, большинство не проверено

**Вопросы:**
- Проводилась ли подготовка к модерации VK?
- Когда планируется подача заявки на модерацию?
- Готовы ли все необходимые материалы?

**Критичные пункты для проверки:**

#### Технические требования
- [ ] Sign validation работает корректно?
- [ ] Приложение запускается на iOS?
- [ ] Приложение запускается на Android?
- [ ] Приложение запускается на Desktop Web?
- [ ] Bundle size < 10MB? (текущий: 9.3MB ✅)

#### Юридические требования
- [ ] Privacy Policy добавлена и доступна?
- [ ] Terms of Service добавлены?
- [ ] Согласие на обработку персональных данных?
- [ ] Возрастные ограничения указаны?

#### Контент
- [ ] Все изображения соответствуют требованиям VK?
- [ ] Нет нарушений авторских прав?
- [ ] Описание приложения корректное?
- [ ] Скриншоты актуальные и качественные?

**Справочный файл:**
```
frontend/roadmap/15_VK_MODERATION.md
```

---

## 🟡 Средние приоритеты (улучшения)

### 5. Analytics Coverage
**Статус:** ⚠️ 42/110 событий (38%)
**Проблема:** По roadmap должно быть 110/117 событий (94%)

**Вопросы:**
- Нужно ли добавить все пропущенные analytics события?
- Какие события наиболее критичны для отслеживания?

**Отсутствующие события:**
- `mascot_stats_opened`
- `villain_stats_opened`
- `ui_error_shown` (в Error Boundary)
- `subscription_cancel_initiated`
- `subscription_cancel_confirmed`
- `subscription_resume_clicked`
- И другие (~68 событий)

**Файлы для обновления:**
```
src/components/ErrorBoundary.tsx - добавить ui_error_shown
src/pages/Home/components/MascotSection.tsx - добавить mascot_stats_opened
src/pages/Villain/*.tsx - добавить villain_stats_opened
```

---

### 6. RecentAttempts компонент
**Статус:** ⚠️ Создан, но не используется
**Проблема:** Файл существует, но не импортирован в HomePage

**Вопросы:**
- Должен ли список последних попыток отображаться на главной?
- Где именно на странице он должен располагаться?
- Сколько последних попыток показывать? (roadmap предлагает 3)

**Файл:**
```
src/pages/Home/components/RecentAttempts.tsx - создан
src/pages/Home/HomePage.tsx - нужно добавить импорт и использование
```

---

### 7. Backend REST API Migration
**Статус:** ⚠️ ~80% завершено
**Проблема:** Миграция с Telegram Bot на REST API в процессе

**Вопросы:**
- Завершена ли миграция backend полностью?
- Все ли endpoints работают стабильно?
- Есть ли production база данных?

**Критичные endpoints для проверки:**
- [ ] `GET /api/v1/home/:childProfileId`
- [ ] `POST /api/v1/attempts`
- [ ] `POST /api/v1/attempts/:id/images`
- [ ] `POST /api/v1/attempts/:id/process`
- [ ] `GET /api/v1/attempts/:id/result`
- [ ] `POST /api/v1/attempts/:id/next-hint`
- [ ] `GET /api/v1/avatars`
- [ ] `POST /api/v1/profiles/child`
- [ ] `GET /api/v1/profile`
- [ ] `GET /api/v1/achievements`
- [ ] `GET /api/v1/friends/referrals`

**Справочный файл:**
```
frontend/roadmap/16_BACKEND_API.md
api/REST_API_STATUS.md (если существует)
```

---

## 🟢 Низкие приоритеты (оптимизация)

### 8. Bundle Optimization
**Статус:** ✅ В пределах нормы (9.3MB < 10MB)
**Текущее разбиение:**
- react-vendor
- vk-vendor
- ui-vendor

**Вопросы:**
- Планируется ли дополнительная оптимизация bundle size?
- Нужно ли добавить lazy loading для страниц?

**Возможные улучшения:**
```typescript
// Route-based code splitting
const HomePage = lazy(() => import('@/pages/Home/HomePage'));
const AchievementsPage = lazy(() => import('@/pages/Achievements/AchievementsPage'));
// и т.д.
```

---

### 9. Production Environment
**Статус:** ⚠️ Требует настройки
**Проблема:** `.env.example` содержит placeholder значения

**Вопросы:**
- Настроен ли production .env с реальными credentials?
- Какие значения должны быть в production?

**Требуемые переменные:**
```bash
VITE_VK_APP_ID=your_actual_vk_app_id
VITE_API_BASE_URL=https://your-production-api.com
VITE_APP_VERSION=1.0.0
```

---

### 10. HTTPS Configuration
**Статус:** ⚠️ Отключен для dev
**Проблема:** vite.config.ts комментарий: "https отключен для локальной разработки"

**Вопросы:**
- Настроен ли HTTPS для production через nginx?
- Есть ли SSL сертификат для production домена?

**Файлы для проверки:**
```
frontend/nginx.conf
docker-compose.yml (production)
```

---

## 📊 Статистика выполнения

### По Phases

| Phase | Название | Статус | Процент |
|-------|----------|--------|---------|
| 0 | Setup | ✅ | 100% |
| 1 | Core Infrastructure | ✅ | 100% |
| 2 | Onboarding | ⚠️ | 90% |
| 3 | Home | ⚠️ | 95% |
| 4 | Help Flow | ✅ | 100% |
| 5 | Check Flow | ✅ | 100% |
| 6 | Achievements | ✅ | 100% |
| 7 | Friends/Referrals | ✅ | 100% |
| 8 | Profile | ✅ | 100% |
| 9 | Villain | ✅ | 100% |
| 10 | Analytics | ⚠️ | 38% |
| 11 | Testing | ❌ | 0% |
| 12 | Platforms | ✅ | 100% |
| 13 | Monetization | ⚠️ | 70% |
| 14 | VK Moderation | ⚠️ | 30% |
| 15 | Backend API | ⚠️ | 80% |

### Общий прогресс: ~85%

**Что готово к релизу:**
- ✅ Core функциональность (помощь с ДЗ, проверка ДЗ)
- ✅ Gamification (маскот, злодей, достижения)
- ✅ Social features (друзья, реферралы)
- ✅ Platform integration (VK Bridge, VKUI)
- ✅ Основные UI компоненты
- ✅ Routing и навигация
- ✅ API интеграция

**Что блокирует production релиз:**
- ❌ Отсутствие тестов
- ⚠️ Incomplete VK Moderation checklist
- ⚠️ Email verification flow
- ⚠️ VK Pay integration verification

---

## 🎯 Рекомендации

### Для немедленного релиза (MVP)

**Приоритет 1 (Критично):**
1. Добавить EmailVerificationWaiting экран
2. Написать минимальные unit тесты (API, Analytics)
3. Пройти VK Moderation checklist полностью
4. Проверить работу всех backend endpoints

**Приоритет 2 (Важно):**
5. Настроить production environment (.env)
6. Проверить VK Pay интеграцию
7. Настроить HTTPS для production
8. Добавить Error Boundary analytics

### Для следующей итерации

**Улучшения функциональности:**
1. Покрыть все 110 analytics событий
2. Добавить RecentAttempts на Home page
3. Реализовать полную VK Pay интеграцию
4. Добавить route-based code splitting

**Улучшения качества:**
5. Написать E2E тесты для ключевых flows
6. Увеличить test coverage до 80%+
7. Провести performance аудит
8. Оптимизировать bundle size дальше

---

## 📁 Полезные файлы для справки

**Roadmap документы:**
```
frontend/roadmap/INDEX.md - общий индекс
frontend/roadmap/FINAL_STATUS.md - текущий статус
frontend/roadmap/15_VK_MODERATION.md - чеклист модерации
frontend/roadmap/16_BACKEND_API.md - спецификация API
frontend/roadmap/ANALYTICS_COVERAGE.md - список событий
```

**Конфигурация:**
```
frontend/.env.example - шаблон переменных окружения
frontend/vite.config.ts - конфигурация сборки
frontend/nginx.conf - конфигурация web сервера
frontend/tsconfig.json - TypeScript настройки
```

**Документация:**
```
frontend/README.md - основная документация
frontend/docs/PLATFORM_INTEGRATION.md - интеграция с VK
frontend/docs/MONETIZATION.md - монетизация
frontend/docs/BACKEND_MIGRATION.md - миграция API
```

---

**Контакт для вопросов:** [Ваш email/контакт]
**Следующая проверка:** После реализации критичных пунктов
