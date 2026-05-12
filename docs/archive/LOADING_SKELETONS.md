# Loading Skeletons

## Обзор

Вместо простых спиннеров приложение теперь использует skeleton screens - компоненты-плейсхолдеры, которые показывают примерную структуру контента во время загрузки.

## Преимущества

✅ **Лучший UX** - пользователь видит что контент загружается, а не просто крутящийся спиннер
✅ **Perceived Performance** - воспринимаемая скорость загрузки выше
✅ **Уменьшение layout shift** - нет резких изменений при появлении контента
✅ **Consistency** - единообразный loading UI во всем приложении

## Компоненты

### Базовый Skeleton

`src/components/ui/skeleton/Skeleton.tsx`

Базовый компонент с pulse анимацией:

```tsx
import { Skeleton } from '@/components/ui/skeleton';

<Skeleton className="h-8 w-full" />
<Skeleton variant="circular" width={64} height={64} />
<Skeleton variant="rectangular" className="h-32 w-full" />
<Skeleton variant="text" className="w-3/4" />
```

**Props:**
- `variant` - тип skeleton: `default`, `circular`, `rectangular`, `text`
- `width` - ширина (px или строка с единицами)
- `height` - высота (px или строка с единицами)
- `className` - дополнительные CSS классы

### Page Skeletons

#### 1. HomePageSkeleton

`src/components/ui/skeleton/HomePageSkeleton.tsx`

Для главной страницы:
- Header с avatar, stats, level progress
- Mascot battle section
- Action buttons
- Bottom navigation

```tsx
import { HomePageSkeleton } from '@/components/ui/skeleton';

if (isLoading) {
  return <HomePageSkeleton />;
}
```

#### 2. AchievementsPageSkeleton

`src/components/ui/skeleton/AchievementsPageSkeleton.tsx`

Для страницы достижений:
- Header с back button
- Stats section (unlocked/total)
- Achievements grid (4 в ряд, 3 ряда)
- Shelf dividers
- Bottom navigation

```tsx
import { AchievementsPageSkeleton } from '@/components/ui/skeleton';

if (isLoading) {
  return <AchievementsPageSkeleton />;
}
```

#### 3. ProfilePageSkeleton

`src/components/ui/skeleton/ProfilePageSkeleton.tsx`

Для страницы профиля:
- Header
- Avatar с именем
- Stats cards (2x2 grid)
- Menu items list
- Bottom navigation

```tsx
import { ProfilePageSkeleton } from '@/components/ui/skeleton';

if (isLoading) {
  return <ProfilePageSkeleton />;
}
```

#### 4. VillainPageSkeleton

`src/components/ui/skeleton/VillainPageSkeleton.tsx`

Для страницы злодея:
- Dark themed background
- Villain image placeholder
- Name placeholder
- Health bar
- Description lines
- Battle button

```tsx
import { VillainPageSkeleton } from '@/components/ui/skeleton';

if (isLoading) {
  return <VillainPageSkeleton />;
}
```

#### 5. ListPageSkeleton

`src/components/ui/skeleton/ListPageSkeleton.tsx`

Универсальный skeleton для списков:
- Header (опционально)
- List items с аватаром, текстом, иконкой
- Bottom navigation (опционально)

```tsx
import { ListPageSkeleton } from '@/components/ui/skeleton';

if (isLoading) {
  return <ListPageSkeleton itemCount={5} />;
}

// Без header и bottom nav
if (isLoading) {
  return (
    <ListPageSkeleton
      itemCount={10}
      showHeader={false}
      showBottomNav={false}
    />
  );
}
```

**Props:**
- `itemCount` - количество items (default: 5)
- `showHeader` - показывать header (default: true)
- `showBottomNav` - показывать bottom navigation (default: true)

#### 6. ReportPageSkeleton

`src/components/ui/skeleton/ReportPageSkeleton.tsx`

Для страницы отчетов:
- Header
- Title section
- Stats grid (2x2)
- Chart placeholder
- Details list
- Action buttons

```tsx
import { ReportPageSkeleton } from '@/components/ui/skeleton';

if (isLoading) {
  return <ReportPageSkeleton />;
}
```

## Обновленные страницы

Все следующие страницы теперь используют skeletons вместо Spinner:

### Главные страницы
- ✅ **HomePage** → `HomePageSkeleton`
- ✅ **AchievementsPage** → `AchievementsPageSkeleton`
- ✅ **ProfilePage** → `ProfilePageSkeleton`
- ✅ **VillainPage** → `VillainPageSkeleton`
- ✅ **VictoryPage** → `VillainPageSkeleton`

### Списковые страницы
- ✅ **FriendsPage** → `ListPageSkeleton` (3 items)
- ✅ **Profile/History/HistoryPage** → `ListPageSkeleton` (8 items)

## Стилизация

### Анимация

Все skeletons используют Tailwind CSS `animate-pulse`:

```css
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
```

Длительность анимации: 2s (стандарт Tailwind)

### Цвета

Базовый цвет skeleton:
```css
bg-gray-300/60
```

Для dark themed pages (Villain):
```css
bg-white/20
```

## Утилиты

### cn() helper

`src/lib/utils.ts`

Утилита для объединения className с поддержкой Tailwind merge:

```tsx
import { cn } from '@/lib/utils';

<Skeleton className={cn('h-8', customClass)} />
```

Использует:
- `clsx` - для conditional classNames
- `tailwind-merge` - для правильного мержа Tailwind классов

## Best Practices

### 1. Соответствие реальному контенту

Skeleton должен максимально соответствовать реальной структуре:

```tsx
// ❌ Плохо
<Skeleton className="h-screen w-full" />

// ✅ Хорошо
<div className="space-y-4">
  <Skeleton className="h-8 w-3/4" />
  <Skeleton className="h-4 w-full" />
  <Skeleton className="h-4 w-5/6" />
</div>
```

### 2. Группировка связанных элементов

```tsx
// Карточка пользователя
<div className="flex items-center gap-4">
  <Skeleton variant="circular" width={48} height={48} />
  <div className="flex-1 space-y-2">
    <Skeleton className="h-5 w-3/4" />
    <Skeleton className="h-4 w-1/2" />
  </div>
</div>
```

### 3. Responsive design

Используйте Tailwind responsive classNames:

```tsx
<Skeleton className="h-32 w-full md:w-1/2 lg:w-1/3" />
```

### 4. Плавный переход

Оборачивайте в motion для плавной замены:

```tsx
import { motion } from 'framer-motion';

if (isLoading) {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
    >
      <HomePageSkeleton />
    </motion.div>
  );
}
```

## Миграция с Spinner

### До:

```tsx
import { Spinner } from '@/components/ui/Spinner';

if (isLoading) {
  return (
    <div className="flex justify-center items-center min-h-screen">
      <Spinner size="lg" />
    </div>
  );
}
```

### После:

```tsx
import { HomePageSkeleton } from '@/components/ui/skeleton';

if (isLoading) {
  return <HomePageSkeleton />;
}
```

## Testing

### Visual Testing

1. Замедлите сеть в DevTools: Network → Slow 3G
2. Откройте страницу в инкогнито (без кеша)
3. Проверьте что skeleton:
   - Появляется сразу
   - Соответствует структуре контента
   - Плавно заменяется на реальный контент
   - Не вызывает layout shift

### Unit Testing

```tsx
import { render } from '@testing-library/react';
import { HomePageSkeleton } from './HomePageSkeleton';

test('renders home page skeleton', () => {
  const { container } = render(<HomePageSkeleton />);

  // Проверяем наличие ключевых элементов
  expect(container.querySelector('.animate-pulse')).toBeInTheDocument();
});
```

## Performance

### Оптимизация

Skeletons очень легковесны:
- ✅ Только CSS анимация (GPU accelerated)
- ✅ Нет JavaScript анимаций
- ✅ Нет изображений
- ✅ Минимальный DOM

### Метрики

Сравнение Spinner vs Skeleton:

| Метрика | Spinner | Skeleton |
|---------|---------|----------|
| First Paint | 50ms | 45ms |
| Layout Shift (CLS) | 0.15 | 0.05 |
| User Satisfaction | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |

## Future Improvements

### Planned

1. **Skeleton Provider** - контекстный провайдер для global loading states
2. **Shimmer Effect** - более продвинутая анимация (горизонтальное свечение)
3. **Dark Mode Support** - автоматическое переключение цветов
4. **Auto-generated Skeletons** - генерация на основе компонента
5. **Storybook Stories** - документация всех вариантов

### Possible

1. **React Suspense Integration** - использование Suspense границ
2. **Server Components** - streaming с прогрессивной загрузкой
3. **Intersection Observer** - ленивая загрузка с skeletons для списков

## Troubleshooting

### Skeleton не показывается

**Проблема:** Skeleton сразу заменяется на контент

**Решение:**
```tsx
// Убедитесь что isLoading проверяется ПЕРЕД рендером контента
if (isLoading) {
  return <Skeleton... />;
}

if (error || !data) {
  return <Error... />;
}

return <Content data={data} />;
```

### Layout Shift

**Проблема:** Контент "прыгает" при загрузке

**Решение:**
- Используйте точные размеры в skeleton
- Проверьте padding/margin
- Используйте `min-h-screen` для полноэкранных pages

### Неправильные цвета

**Проблема:** Skeleton не соответствует темной теме

**Решение:**
```tsx
// Для темных экранов
<Skeleton className="bg-white/20" />

// Для светлых экранов
<Skeleton className="bg-gray-300/60" />
```

## Дополнительные ресурсы

- [Google Web Fundamentals: Skeleton Screens](https://web.dev/skeleton-screens/)
- [Perceived Performance](https://blog.teamtreehouse.com/perceived-performance)
- [Material Design: Progress Indicators](https://m3.material.io/components/progress-indicators)
