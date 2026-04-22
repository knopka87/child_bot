# Детальный анализ дизайна (Figma AI экспорт)

## Типичные проблемы Figma AI генерации

### 1. Magic Numbers везде

Figma AI генерирует "пиксель-перфект" значения вместо системных:

```tsx
// ❌ Проблема: hardcoded значения
className="text-[14px]"           // вместо text-sm
className="text-[13px]"           // нестандартный размер
className="text-[11px]"           // нестандартный размер
className="text-[12px]"           // нестандартный размер
className="text-[28px]"           // нестандартный размер
className="text-[20px]"           // нестандартный размер
className="text-[36px]"           // вместо text-4xl
className="text-[48px]"           // вместо text-5xl
className="rounded-3xl"           // OK
className="rounded-xl"            // OK
className="gap-3.5"               // нестандартный gap
className="mb-1.5"                // нестандартный margin
className="py-0.5"                // нестандартный padding
className="mt-0.5"                // нестандартный margin
```

**Рекомендация**: Создать typography scale:
```css
/* theme.css */
--text-xs: 11px;
--text-sm: 13px;
--text-base: 14px;
--text-lg: 16px;
--text-xl: 20px;
--text-2xl: 24px;
--text-3xl: 28px;
```

### 2. Hardcoded цвета вместо токенов

```tsx
// ❌ Проблема: прямые HEX вместо CSS переменных
className="bg-[#F0F4FF]"          // есть var(--background)
className="bg-[#E8E4FF]"          // есть var(--muted)
className="text-[#2D3436]"        // есть var(--foreground)
className="text-[#00B894]"        // есть var(--success)
className="text-[#FDCB6E]"        // есть var(--warning)
className="from-[#E8FFF8]"        // нет токена
className="from-[#FFF0F0]"        // нет токена

// ✅ Правильно:
className="bg-background"
className="text-foreground"
className="text-success"
```

**Рекомендация**: Добавить недостающие токены:
```css
--success-background: #E8FFF8;
--error-background: #FFF0F0;
--warning-background: #FFF9E8;
```

### 3. Дублирование UI библиотек

```json
// package.json — избыточность
"@mui/material": "7.3.5",         // 500KB+ — НЕ НУЖЕН
"@mui/icons-material": "7.3.5",   // НЕ НУЖЕН (есть Lucide)
"@emotion/react": "11.14.0",      // НЕ НУЖЕН без MUI
"@emotion/styled": "11.14.1",     // НЕ НУЖЕН без MUI
"react-slick": "0.31.0",          // ДУБЛЬ (есть Embla)
"react-popper": "2.3.0",          // НЕ НУЖЕН (есть Radix)
"@popperjs/core": "2.11.8",       // НЕ НУЖЕН
```

**Экономия**: ~600KB после tree-shaking, ~2MB до.

### 4. Компоненты UI не используются

50+ Radix компонентов, но в screens используются inline стили:

```tsx
// ❌ Проблема: кастомные кнопки вместо готовых
<button className="w-full py-4 bg-primary text-white rounded-2xl...">

// ✅ Правильно: использовать Button компонент
import { Button } from "../ui/button";
<Button size="lg" className="w-full">Продолжить</Button>
```

### 5. Отсутствует система spacing

```tsx
// ❌ Проблема: разные отступы везде
className="px-5"    // 20px
className="px-6"    // 24px
className="pt-8"    // 32px
className="pt-4"    // 16px
className="pb-4"    // 16px
className="pb-6"    // 24px
className="mb-5"    // 20px
className="mb-6"    // 24px
className="mb-4"    // 16px
className="gap-3"   // 12px
className="gap-4"   // 16px
```

**Рекомендация**: Единая система:
- `px-4` (16px) — стандартный горизонтальный padding
- `py-4` (16px) — стандартный вертикальный padding
- `gap-3` (12px) — между элементами
- `gap-4` (16px) — между секциями

---

## Критические UX проблемы

### 1. Нет loading states

```tsx
// Payment.tsx:117 — просто меняет state
onClick={() => {
  setState(Math.random() > 0.2 ? "success" : "error");
}}
```

Нужно:
- Spinner во время загрузки
- Disable кнопки
- Skeleton для контента

### 2. Нет error boundaries

Если компонент упадёт — белый экран. Нужно:
```tsx
<ErrorBoundary fallback={<ErrorScreen />}>
  <App />
</ErrorBoundary>
```

### 3. Accessibility проблемы

```tsx
// ❌ Кнопки без aria-label
<button onClick={() => navigate("/villain")}>
  <span className="text-[36px]">👾</span>
  ...
</button>

// ❌ Checkbox без связи с label
<button onClick={() => setAdultConsent(!adultConsent)}>
  <div className={...}>{adultConsent && <Check />}</div>
  <span>текст</span>
</button>
```

### 4. Touch targets слишком маленькие

```tsx
// Layout.tsx — табы 44px минимум для touch
className="py-2 pt-3"  // ~40px — мало для пальца

// Рекомендация: минимум 48px
className="py-3"
```

---

## Проблемы для VK/MAX Mini Apps

### 1. Нет интеграции с платформой

Для VK Mini Apps нужно:
```bash
npm install @vkontakte/vk-bridge @vkontakte/vkui
```

Для MAX Mini Apps — их SDK (когда выйдет публичный).

### 2. Фиксированная ширина 390px

```tsx
// Layout.tsx:25
<div className="w-full max-w-[390px]...">
```

VK Mini Apps работают на разных устройствах. Нужно:
```tsx
<div className="w-full max-w-md"> // 448px — более гибко
```

### 3. Safe area не учтена

```tsx
className="safe-area-bottom"  // только снизу
```

Для VK/MAX нужно также `safe-area-top` и учёт нативного header'а.

### 4. Навигация не учитывает платформу

VK Mini Apps имеют свой back button в header'е. Текущие кнопки "Назад" дублируют функционал.

---

## Что хорошо сделано

| Аспект | Оценка |
|--------|--------|
| Цветовая палитра | Продуманная, детская, дружелюбная |
| Анимации | Плавные, не перегружены |
| Иерархия информации | Чёткая на главном экране |
| CTA кнопки | Контрастные, заметные |
| Маскот персонаж | Добавляет характер |
| Rounded corners | Консистентные (2xl, 3xl) |

---

## Рекомендации по улучшению

### Приоритет 1 — Обязательно

| Задача | Причина |
|--------|---------|
| Убрать MUI, emotion, react-slick | -600KB bundle size |
| Заменить hardcoded цвета на токены | Поддерживаемость |
| Использовать готовые UI компоненты | Консистентность |
| Добавить loading/error states | UX |
| Увеличить touch targets до 48px | Accessibility |

### Приоритет 2 — Желательно

| Задача | Причина |
|--------|---------|
| Создать typography scale | Консистентность |
| Унифицировать spacing систему | Поддерживаемость |
| Добавить skeleton loaders | Perceived performance |
| Вынести константы в отдельный файл | Чистота кода |

### Приоритет 3 — Для VK/MAX

| Задача | Причина |
|--------|---------|
| Интеграция VK Bridge | Авторизация, платежи |
| Адаптация под VKUI guidelines | Нативный вид |
| Изучить MAX Mini Apps SDK | Когда станет доступен |
| Убрать фиксированную ширину | Адаптивность |

---

## Вывод

**Figma AI** генерирует визуально корректный код, но с проблемами:
- Не использует дизайн-систему
- Дублирует библиотеки
- Hardcoded значения вместо токенов
- Игнорирует готовые компоненты

**Для production** требуется рефакторинг: замена magic numbers, очистка зависимостей, интеграция с платформой (VK/MAX).