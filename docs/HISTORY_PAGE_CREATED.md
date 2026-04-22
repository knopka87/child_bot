# История попыток — Страница создана

**Дата:** 2026-04-04

## Что сделано

Создана страница истории попыток пользователя с полной интеграцией с бэкенд API.

---

## Созданные файлы

### 1. Компоненты

#### HistoryPage.tsx
**Путь:** `frontend/src/pages/Profile/History/HistoryPage.tsx`

**Функциональность:**
- Отображение списка всех попыток пользователя (help и check)
- Фильтрация по режиму (все/помощь/проверка)
- Фильтрация по статусу (все/верно/ошибки)
- Кнопка "Назад" в профиль
- Трекинг аналитики для всех действий
- Анимация появления карточек попыток
- Модальное окно с деталями попытки

**Основные функции:**
- `handleBack()` - возврат в профиль с трекингом
- `handleCardClick()` - открытие деталей попытки
- `handleFilterModeChange()` - смена фильтра режима
- `handleFilterStatusChange()` - смена фильтра статуса
- `formatDate()` - форматирование даты в читаемый вид
- `getStatusConfig()` - получение конфигурации статуса (иконка, цвет, текст)

#### HistoryDetailModal.tsx
**Путь:** `frontend/src/pages/Profile/History/components/HistoryDetailModal.tsx`

**Функциональность:**
- Модальное окно с полными деталями попытки
- Отображение изображений (условие/решение)
- Отображение использованных подсказок
- Список найденных ошибок с описанием
- Кнопки действий:
  - "Исправить и проверить" (для ошибочных попыток)
  - "Повторить" (новая попытка того же типа)
  - "Закрыть"
- Анимация появления снизу вверх
- Drag handle для визуального понимания

#### HistoryPage.module.css
**Путь:** `frontend/src/pages/Profile/History/HistoryPage.module.css`

**Стили для:**
- Градиентный фон страницы
- Фильтры с горизонтальной прокруткой
- Карточки попыток с тенями и анимацией
- Бейджи статусов (успех/ошибка/в процессе)
- Состояния загрузки/ошибки/пустой истории

#### HistoryDetailModal.module.css
**Путь:** `frontend/src/pages/Profile/History/components/HistoryDetailModal.module.css`

**Стили для:**
- Overlay затемнения
- Модальное окно с rounded corners
- Drag handle
- Сетка изображений
- Карточки ошибок
- Кнопки действий (primary/secondary/text)

---

### 2. Хуки

#### useHistory.ts
**Путь:** `frontend/src/pages/Profile/History/hooks/useHistory.ts`

**Функциональность:**
- Загрузка истории попыток с бэкенда
- Автоматический refetch при изменении фильтров
- Обработка состояний загрузки и ошибок
- Возможность ручного обновления (refetch)

**Параметры:**
- `childProfileId: string | null` - ID профиля ребёнка
- `filters?: HistoryFilters` - фильтры (mode, status, даты)

**Возвращает:**
```typescript
{
  data: HistoryAttempt[];      // список попыток
  isLoading: boolean;          // флаг загрузки
  error: Error | null;         // ошибка загрузки
  refetch: () => Promise<void>; // функция перезагрузки
}
```

---

### 3. API интеграция

#### Обновлён profileAPI в profile.ts

**Метод `getHistory()`**

Эндпоинт: `GET /api/v1/attempts/history`

**Query параметры:**
- `child_profile_id` - ID профиля ребёнка
- `limit` - количество записей (50)
- `offset` - смещение (0)
- `filter` - фильтр ('all', 'help', 'check')

**Конвертация данных:**
- Backend формат `attempt_id` → Frontend формат `id`
- Backend формат `result_status` → Frontend формат `status`
- Создание минимального объекта `images` с thumbnail
- Создание объекта `result` с ошибками

**Метод `getHistoryDetail()`**

Эндпоинт: `GET /api/v1/attempts/:attemptId/result`

**Конвертация:**
- Преобразование ошибок в формат `ErrorFeedback`
- Подсчёт использованных подсказок
- Создание полного объекта `HistoryAttempt`

---

### 4. Роутинг

#### Обновлён routes.tsx

**Изменения:**
```typescript
// Добавлен импорт
import { HistoryPage } from '@/pages/Profile/History/HistoryPage';

// Изменён роут
<Route path={ROUTES.PROFILE_HISTORY} element={<HistoryPage />} />
```

**Роут:** `/profile/history`

---

## Аналитика

### События

Страница отправляет следующие события в analytics:

1. **history_opened**
   - Когда: при открытии страницы
   - Параметры: `child_profile_id`

2. **history_back_clicked**
   - Когда: клик на кнопку "Назад"
   - Параметры: `child_profile_id`

3. **history_item_clicked**
   - Когда: клик на карточку попытки
   - Параметры: `child_profile_id`, `attempt_id`, `mode`, `status`

4. **history_filter_changed**
   - Когда: изменение фильтра
   - Параметры: `child_profile_id`, `filter_type` ('mode'/'status'), `filter_value`

5. **history_retry_clicked**
   - Когда: клик на "Повторить" в модалке
   - Параметры: `child_profile_id`, `attempt_id`, `mode`

6. **history_fix_errors_clicked**
   - Когда: клик на "Исправить и проверить"
   - Параметры: `child_profile_id`, `attempt_id`

---

## UX детали

### Фильтры

**Режим:**
- Все
- Помощь
- Проверка

**Статус:**
- Все статусы
- Верно
- Ошибки

Фильтры имеют горизонтальную прокрутку на маленьких экранах.

### Статусы попыток

| Статус | Иконка | Цвет | Текст |
|--------|--------|------|-------|
| success | CheckCircle | Зелёный (#00B894) | "Решено верно" |
| error | XCircle | Красный (#DC3545) | "Есть ошибки" |
| in_progress | Clock | Жёлтый (#FDCB6E) | "В процессе" |

### Модальное окно

**Секции:**
- Заголовок с названием режима и датой
- Бейдж статуса
- Информация о подсказках (если использовались)
- Сетка изображений (условие, решение)
- Результат (summary)
- Список ошибок с номерами шагов
- Кнопки действий

**Анимация:**
- Появление снизу вверх (slide up)
- Overlay fade in/out
- Spring transition для плавности

---

## Состояния

### Загрузка
- Центрированный спиннер
- Текст отсутствует

### Ошибка
- Текст: "Не удалось загрузить историю"
- Кнопка "Повторить" для refetch

### Пустая история
- Текст: "История пуста"
- Подсказка: "Попробуйте решить первую задачу!"

### Данные загружены
- Список карточек с анимацией
- Каждая карточка появляется с задержкой 0.05s * index

---

## Backend требования

### Эндпоинт истории

**URL:** `GET /api/v1/attempts/history`

**Query параметры:**
- `child_profile_id` (required) - ID профиля ребёнка
- `limit` (optional, default: 20) - количество записей
- `offset` (optional, default: 0) - смещение для пагинации
- `filter` (optional) - фильтр ('all', 'help', 'check', 'correct', 'errors')

**Response:**
```typescript
{
  attempts: [
    {
      attempt_id: string;
      mode: 'help' | 'check';
      attempt_status: 'completed' | 'failed' | 'cancelled';
      result_status?: 'correct' | 'has_errors' | 'wrong';
      error_count?: number;
      thumbnail_url: string;
      scenario_type?: 'single_photo' | 'two_photo';
      created_at: string; // ISO 8601
      completed_at: string; // ISO 8601
    }
  ];
  total: number;
  has_more: boolean;
}
```

### Эндпоинт деталей попытки

**URL:** `GET /api/v1/attempts/:attemptId/result`

**Response:**
```typescript
{
  attempt_id: string;
  attempt_status: 'completed' | 'processing' | 'failed';
  result_status?: 'correct' | 'has_errors' | 'wrong';
  error_count?: number;
  errors?: Array<{
    error_block_id: string;
    step_number?: number;
    line_reference?: string;
    error_type: string;
    error_message: string;
    error_hint: string;
    location_type: 'step' | 'line' | 'general';
  }>;
  hints?: Array<{
    hint_level: 1 | 2 | 3;
    hint_text: string;
    hint_images?: string[];
    unlocked: boolean;
  }>;
  used_hints_count?: number;
}
```

---

## Навигация

### Из профиля в историю
```typescript
// ProfilePage.tsx
navigate(ROUTES.PROFILE_HISTORY); // /profile/history
```

### Из истории в профиль
```typescript
// HistoryPage.tsx - кнопка "Назад"
navigate(ROUTES.PROFILE); // /profile
```

### Из модалки деталей
```typescript
// "Повторить" для help
navigate(ROUTES.HELP_UPLOAD);

// "Повторить" для check
navigate(ROUTES.CHECK_SCENARIO);

// "Исправить и проверить"
navigate(ROUTES.CHECK);
```

---

## Тестирование

### Ручное тестирование

1. **Открытие страницы:**
   ```bash
   # Перейти в профиль
   open http://localhost:5173/profile

   # Кликнуть на "История"
   # Должна открыться страница с историей
   ```

2. **Фильтрация:**
   ```bash
   # Кликнуть "Помощь"
   # Должны остаться только help попытки

   # Кликнуть "Верно"
   # Должны остаться только успешные попытки
   ```

3. **Детали попытки:**
   ```bash
   # Кликнуть на карточку попытки
   # Должно открыться модальное окно
   # Проверить отображение всех данных
   ```

4. **Действия:**
   ```bash
   # В модалке для ошибочной попытки
   # Кликнуть "Исправить и проверить"
   # Должна открыться страница проверки

   # Кликнуть "Повторить"
   # Должна открыться страница upload/scenario
   ```

### Состояния для проверки

- ✅ Пустая история (новый пользователь)
- ✅ Несколько попыток с разными статусами
- ✅ Фильтрация работает
- ✅ Модалка открывается/закрывается
- ✅ Кнопка "Назад" ведёт в профиль
- ✅ Analytics события отправляются

---

## Зависимости

### Пакеты (уже установлены)
- `framer-motion` - анимации
- `lucide-react` - иконки
- `react-router-dom` - роутинг

### Внутренние зависимости
- `@/components/layout/BottomNav` - нижняя навигация
- `@/components/ui/Spinner` - спиннер загрузки
- `@/hooks/useAnalytics` - хук аналитики
- `@/lib/platform/vk-storage` - хранилище VK
- `@/api/profile` - API методы
- `@/types/profile` - TypeScript типы

---

## Следующие шаги (опционально)

### Пагинация
- Добавить бесконечную прокрутку
- Загружать следующие 20 попыток при скролле

### Дополнительные фильтры
- Фильтр по дате (сегодня, неделя, месяц)
- Комбинированные фильтры

### Детали попытки
- Полная информация об изображениях
- Увеличение изображений по клику
- Шаринг результата

### Производительность
- Мемоизация списка попыток
- Виртуализация длинного списка
- Ленивая загрузка изображений

---

## Результат

✅ Страница истории создана
✅ Интеграция с бэкенд API
✅ Фильтрация по режиму и статусу
✅ Модальное окно с деталями
✅ Полная аналитика
✅ Адаптивный дизайн
✅ Анимации и переходы

**URL:** http://localhost:5173/profile/history

---

**Статус:** ✅ Готово к тестированию
**Требует:** Backend API endpoints `/api/v1/attempts/history` и `/api/v1/attempts/:id/result`
