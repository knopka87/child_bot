# Component Architecture — Архитектура компонентов

**Проект:** Объяснятель ДЗ MiniApp
**Дата:** 2026-03-29

---

## 🎯 Принципы архитектуры

### Многослойная блочная структура

Приложение строится по принципу **многослойной композиции**:
1. **Design Tokens** - базовые значения (цвета, размеры, отступы)
2. **UI Kit** - атомарные компоненты (Button, Input, Card)
3. **Composite Components** - составные блоки (Header, BottomNav, AchievementCard)
4. **Page Sections** - секции страниц (ProfileInfo, HistoryList)
5. **Page Templates** - шаблоны страниц (MainLayout, OnboardingLayout)
6. **Pages** - конечные страницы (HomePage, AchievementsPage)

**Ключевой принцип:** Изменение в одном слое не требует переписывания всех страниц!

---

## 📐 Слой 1: Design Tokens

**Файл:** `src/styles/tokens.css` или `src/styles/tokens.ts`

### Что это?

Design tokens - это **единственный источник истины** для всех визуальных значений.

```css
/* tokens.css */
:root {
  /* ===== COLORS ===== */

  /* Primary palette */
  --color-primary: #7B5BF2; /* Фиолетовый для основных кнопок */
  --color-primary-hover: #6A4FDB;
  --color-primary-disabled: #B8A8F5;

  /* Secondary palette */
  --color-secondary: #4ECDC4; /* Бирюзовый для проверки */
  --color-secondary-hover: #45B8B0;

  /* Neutral colors */
  --color-bg-primary: #F8F9FE; /* Основной фон */
  --color-bg-secondary: #FFFFFF; /* Карточки */
  --color-bg-tertiary: #F0F2F8; /* Альтернативный фон */

  --color-text-primary: #1F1F1F; /* Основной текст */
  --color-text-secondary: #6B6B6B; /* Вторичный текст */
  --color-text-tertiary: #A0A0A0; /* Подписи */
  --color-text-inverse: #FFFFFF; /* Текст на темном */

  /* Status colors */
  --color-success: #4ECDC4;
  --color-warning: #FFB74D;
  --color-error: #FF6B6B;
  --color-info: #7B5BF2;

  /* Coins & rewards */
  --color-coin: #FFD700;
  --color-xp: #7B5BF2;

  /* ===== SPACING ===== */

  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --space-2xl: 48px;

  /* ===== TYPOGRAPHY ===== */

  --font-family-primary: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  --font-family-display: 'SF Pro Display', -apple-system, sans-serif;

  --font-size-xs: 12px;
  --font-size-sm: 14px;
  --font-size-base: 16px;
  --font-size-lg: 18px;
  --font-size-xl: 20px;
  --font-size-2xl: 24px;
  --font-size-3xl: 32px;

  --font-weight-regular: 400;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;

  --line-height-tight: 1.2;
  --line-height-normal: 1.5;
  --line-height-relaxed: 1.75;

  /* ===== BORDER RADIUS ===== */

  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 24px;
  --radius-full: 9999px;

  /* ===== SHADOWS ===== */

  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.08);
  --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.12);
  --shadow-xl: 0 16px 48px rgba(0, 0, 0, 0.16);

  /* ===== Z-INDEX ===== */

  --z-base: 1;
  --z-dropdown: 1000;
  --z-sticky: 1100;
  --z-modal-backdrop: 1200;
  --z-modal: 1300;
  --z-toast: 1400;

  /* ===== TRANSITIONS ===== */

  --transition-fast: 150ms ease-in-out;
  --transition-base: 250ms ease-in-out;
  --transition-slow: 350ms ease-in-out;

  /* ===== BREAKPOINTS (for JS) ===== */

  --breakpoint-mobile: 390px;
  --breakpoint-tablet: 768px;
  --breakpoint-desktop: 1024px;
}
```

**TypeScript версия:**

```typescript
// src/styles/tokens.ts
export const tokens = {
  colors: {
    primary: '#7B5BF2',
    primaryHover: '#6A4FDB',
    secondary: '#4ECDC4',
    // ...
  },
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
    '2xl': 48,
  },
  // ...
} as const

export type Token = typeof tokens
```

### Почему это важно?

✅ **Изменение цвета в одном месте** → применяется везде
✅ **Консистентный дизайн** → все отступы из единого набора
✅ **Легко переключать темы** → dark mode за 5 минут
✅ **Платформенная адаптация** → разные токены для VK/Max/Telegram

---

## 🧱 Слой 2: UI Kit (Атомарные компоненты)

**Путь:** `src/components/ui/`

### Что это?

Базовые неделимые компоненты, которые используются везде.

### Список компонентов UI Kit

```
ui/
├── Button/
│   ├── Button.tsx
│   ├── Button.module.css
│   └── Button.stories.tsx (опционально)
├── Input/
├── Card/
├── Avatar/
├── Badge/
├── Spinner/
├── ProgressBar/
├── Modal/
├── Tooltip/
└── Icon/
```

### Пример: Button

```typescript
// src/components/ui/Button/Button.tsx
import styles from './Button.module.css'

export type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost'
export type ButtonSize = 'sm' | 'md' | 'lg'

interface ButtonProps {
  variant?: ButtonVariant
  size?: ButtonSize
  fullWidth?: boolean
  loading?: boolean
  disabled?: boolean
  icon?: React.ReactNode
  children: React.ReactNode
  onClick?: () => void
}

export function Button({
  variant = 'primary',
  size = 'md',
  fullWidth = false,
  loading = false,
  disabled = false,
  icon,
  children,
  onClick,
}: ButtonProps) {
  return (
    <button
      className={`
        ${styles.button}
        ${styles[variant]}
        ${styles[size]}
        ${fullWidth ? styles.fullWidth : ''}
      `}
      disabled={disabled || loading}
      onClick={onClick}
    >
      {loading && <Spinner size="sm" />}
      {!loading && icon && <span className={styles.icon}>{icon}</span>}
      <span className={styles.label}>{children}</span>
    </button>
  )
}
```

```css
/* Button.module.css */
.button {
  /* Используем токены! */
  font-family: var(--font-family-primary);
  font-weight: var(--font-weight-semibold);
  border-radius: var(--radius-md);
  transition: all var(--transition-base);

  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-sm);

  border: none;
  cursor: pointer;
}

/* Variants */
.primary {
  background: var(--color-primary);
  color: var(--color-text-inverse);
}

.primary:hover {
  background: var(--color-primary-hover);
}

.secondary {
  background: var(--color-secondary);
  color: var(--color-text-inverse);
}

/* Sizes */
.sm {
  padding: var(--space-sm) var(--space-md);
  font-size: var(--font-size-sm);
  height: 36px;
}

.md {
  padding: var(--space-md) var(--space-lg);
  font-size: var(--font-size-base);
  height: 48px;
}

.lg {
  padding: var(--space-lg) var(--space-xl);
  font-size: var(--font-size-lg);
  height: 56px;
}

.fullWidth {
  width: 100%;
}
```

### Почему это важно?

✅ **Изменение стиля Button** → меняется на всех 50+ страницах
✅ **Добавление нового варианта** → доступен везде
✅ **Консистентность** → все кнопки выглядят одинаково
✅ **Accessibility** → добавляем aria-label один раз

---

## 🧩 Слой 3: Composite Components (Составные блоки)

**Путь:** `src/components/composite/`

### Что это?

Составные компоненты из UI Kit, которые решают конкретную бизнес-задачу.

### Структура

```
composite/
├── Header/
│   ├── Header.tsx           # Общий header
│   ├── HeaderLevelBadge.tsx # Бейдж уровня
│   ├── HeaderCoins.tsx      # Счетчик монет
│   └── HeaderTasks.tsx      # Счетчик заданий
├── BottomNav/
│   ├── BottomNav.tsx
│   └── BottomNavItem.tsx
├── AchievementCard/
├── AttemptCard/
├── VillainCard/
└── MascotCard/
```

### Пример: Header (переиспользуемый!)

```typescript
// src/components/composite/Header/Header.tsx
import { Card } from '@/components/ui/Card'
import { ProgressBar } from '@/components/ui/ProgressBar'
import { HeaderLevelBadge } from './HeaderLevelBadge'
import { HeaderCoins } from './HeaderCoins'
import { HeaderTasks } from './HeaderTasks'
import styles from './Header.module.css'

interface HeaderProps {
  level: number
  levelProgress: number
  coins: number
  tasksCount: number
  showCoins?: boolean
  showTasks?: boolean
}

export function Header({
  level,
  levelProgress,
  coins,
  tasksCount,
  showCoins = true,
  showTasks = true,
}: HeaderProps) {
  return (
    <Card className={styles.header}>
      {/* Level badge - переиспользуемый компонент */}
      <HeaderLevelBadge level={level} progress={levelProgress} />

      <div className={styles.stats}>
        {/* Монеты - переиспользуемый компонент */}
        {showCoins && <HeaderCoins count={coins} />}

        {/* Задания - переиспользуемый компонент */}
        {showTasks && <HeaderTasks count={tasksCount} />}
      </div>
    </Card>
  )
}
```

```typescript
// src/components/composite/Header/HeaderLevelBadge.tsx
import { Badge } from '@/components/ui/Badge'
import { ProgressBar } from '@/components/ui/ProgressBar'
import styles from './HeaderLevelBadge.module.css'

interface HeaderLevelBadgeProps {
  level: number
  progress: number // 0-100
}

export function HeaderLevelBadge({ level, progress }: HeaderLevelBadgeProps) {
  return (
    <Badge variant="primary" className={styles.badge}>
      <span className={styles.level}>{level}</span>
      <span className={styles.label}>Уровень</span>
      <ProgressBar value={progress} size="sm" />
    </Badge>
  )
}
```

### Почему это важно?

✅ **Header используется на 10+ страницах** → меняем один раз
✅ **Можно показывать/скрывать элементы** через props
✅ **Легко добавить новый элемент** в header везде сразу
✅ **Каждый подкомпонент можно использовать отдельно**

---

## 📄 Слой 4: Page Sections (Секции страниц)

**Путь:** `src/components/sections/`

### Что это?

Крупные блоки, специфичные для определенных страниц, но переиспользуемые.

### Структура

```
sections/
├── home/
│   ├── MascotSection.tsx      # Маскот VS злодей
│   ├── ActionButtons.tsx      # Помощь + Проверка
│   └── RecentAttempts.tsx     # Последние попытки
├── profile/
│   ├── ProfileCard.tsx
│   ├── ProfileMenu.tsx
│   └── ProfileSettings.tsx
├── achievements/
│   ├── AchievementGrid.tsx
│   └── AchievementShelf.tsx
└── friends/
    ├── ReferralProgress.tsx
    └── InviteFriendCard.tsx
```

### Пример: MascotSection (используется в Home и Villain)

```typescript
// src/components/sections/home/MascotSection.tsx
import { MascotCard } from '@/components/composite/MascotCard'
import { VillainCard } from '@/components/composite/VillainCard'
import { BattleIndicator } from '@/components/composite/BattleIndicator'
import styles from './MascotSection.module.css'

interface MascotSectionProps {
  mascot: {
    id: string
    state: 'idle' | 'happy' | 'thinking' | 'celebrating'
    imageUrl: string
  }
  villain: {
    id: string
    name: string
    health: number
    maxHealth: number
    imageUrl: string
  } | null
  onMascotClick?: () => void
  onVillainClick?: () => void
}

export function MascotSection({
  mascot,
  villain,
  onMascotClick,
  onVillainClick,
}: MascotSectionProps) {
  return (
    <div className={styles.container}>
      {/* Маскот - слева */}
      <MascotCard
        mascot={mascot}
        onClick={onMascotClick}
        className={styles.mascot}
      />

      {/* Индикатор битвы - по центру */}
      {villain && (
        <BattleIndicator className={styles.battle} />
      )}

      {/* Злодей - справа */}
      {villain && (
        <VillainCard
          villain={villain}
          onClick={onVillainClick}
          className={styles.villain}
        />
      )}
    </div>
  )
}
```

### Почему это важно?

✅ **MascotSection используется в Home и Villain screen**
✅ **Изменение логики битвы** → меняется везде
✅ **Можно использовать части секции** отдельно
✅ **A/B тестирование** - легко менять layout

---

## 🗂️ Слой 5: Page Templates (Шаблоны страниц)

**Путь:** `src/components/templates/`

### Что это?

Общие layouts, которые определяют структуру страниц.

### Типы шаблонов

```
templates/
├── MainLayout.tsx           # Основной layout с header + content + bottomNav
├── OnboardingLayout.tsx     # Layout онбординга (без bottomNav)
├── ModalLayout.tsx          # Layout для модальных экранов
└── FullscreenLayout.tsx     # Fullscreen (например, для камеры)
```

### Пример: MainLayout (используется на 90% страниц!)

```typescript
// src/components/templates/MainLayout.tsx
import { Header } from '@/components/composite/Header'
import { BottomNav } from '@/components/composite/BottomNav'
import { useProfile } from '@/hooks/useProfile'
import styles from './MainLayout.module.css'

interface MainLayoutProps {
  children: React.ReactNode
  showHeader?: boolean
  showBottomNav?: boolean
  headerProps?: {
    showCoins?: boolean
    showTasks?: boolean
  }
  backgroundColor?: string
}

export function MainLayout({
  children,
  showHeader = true,
  showBottomNav = true,
  headerProps,
  backgroundColor,
}: MainLayoutProps) {
  const { profile } = useProfile()

  return (
    <div className={styles.layout} style={{ backgroundColor }}>
      {/* Header - опциональный */}
      {showHeader && profile && (
        <Header
          level={profile.level}
          levelProgress={profile.level_progress_percent}
          coins={profile.coins_balance}
          tasksCount={profile.tasks_solved_correct_count}
          {...headerProps}
        />
      )}

      {/* Content - всегда есть */}
      <main className={styles.content}>
        {children}
      </main>

      {/* Bottom Nav - опциональный */}
      {showBottomNav && (
        <BottomNav />
      )}
    </div>
  )
}
```

```css
/* MainLayout.module.css */
.layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: var(--color-bg-primary);
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-md);
  padding-bottom: calc(60px + var(--space-md)); /* Отступ для bottomNav */
}
```

### Использование в страницах

```typescript
// src/pages/Home/HomePage.tsx
import { MainLayout } from '@/components/templates/MainLayout'
import { MascotSection } from '@/components/sections/home/MascotSection'
import { ActionButtons } from '@/components/sections/home/ActionButtons'
import { RecentAttempts } from '@/components/sections/home/RecentAttempts'

export function HomePage() {
  return (
    <MainLayout>
      {/* Все секции внутри layout */}
      <MascotSection {...mascotData} />
      <ActionButtons />
      <RecentAttempts />
    </MainLayout>
  )
}

// src/pages/Achievements/AchievementsPage.tsx
export function AchievementsPage() {
  return (
    <MainLayout>
      {/* Тот же layout, другой контент! */}
      <AchievementGrid />
    </MainLayout>
  )
}

// src/pages/Onboarding/OnboardingPage.tsx
export function OnboardingPage() {
  return (
    <OnboardingLayout> {/* Другой layout для онбординга */}
      <OnboardingContent />
    </OnboardingLayout>
  )
}
```

### Почему это важно?

✅ **Изменение header/bottomNav** → меняется на всех страницах
✅ **Добавление нового элемента** в layout → везде сразу
✅ **Разные layouts** для разных flows
✅ **Легко скрыть header/nav** на конкретной странице

---

## 📱 Слой 6: Pages (Конечные страницы)

**Путь:** `src/pages/`

### Что это?

Конечные страницы, которые композируют все слои вместе.

```
pages/
├── Home/
│   └── HomePage.tsx
├── Help/
│   ├── HelpUploadPage.tsx
│   ├── HelpProcessingPage.tsx
│   └── HelpResultPage.tsx
├── Achievements/
│   └── AchievementsPage.tsx
└── Profile/
    ├── ProfilePage.tsx
    └── ProfileHistoryPage.tsx
```

### Пример: HomePage (композиция всех слоев!)

```typescript
// src/pages/Home/HomePage.tsx
import { MainLayout } from '@/components/templates/MainLayout'
import { MascotSection } from '@/components/sections/home/MascotSection'
import { ActionButtons } from '@/components/sections/home/ActionButtons'
import { RecentAttempts } from '@/components/sections/home/RecentAttempts'
import { UnfinishedAttemptModal } from '@/components/composite/UnfinishedAttemptModal'
import { useHome } from './useHome'

export function HomePage() {
  const {
    mascot,
    villain,
    unfinishedAttempt,
    recentAttempts,
    handleHelpClick,
    handleCheckClick,
    handleMascotClick,
    handleVillainClick,
  } = useHome()

  return (
    <>
      {/* Layout - слой 5 */}
      <MainLayout>
        {/* Секция маскота - слой 4 */}
        <MascotSection
          mascot={mascot}
          villain={villain}
          onMascotClick={handleMascotClick}
          onVillainClick={handleVillainClick}
        />

        {/* Кнопки действий - слой 4 */}
        <ActionButtons
          onHelpClick={handleHelpClick}
          onCheckClick={handleCheckClick}
        />

        {/* Последние попытки - слой 4 */}
        <RecentAttempts attempts={recentAttempts} />
      </MainLayout>

      {/* Модал - слой 3 */}
      {unfinishedAttempt && (
        <UnfinishedAttemptModal attempt={unfinishedAttempt} />
      )}
    </>
  )
}
```

### Почему это важно?

✅ **Страница - это композиция переиспользуемых блоков**
✅ **Минимум кода на уровне страницы** (только композиция)
✅ **Изменение блока** → меняется на всех страницах
✅ **Легко создать новую страницу** из готовых блоков

---

## 🔄 Пример изменения: Добавление темной темы

### Задача
Добавить темную тему для всего приложения.

### Решение

**Шаг 1:** Обновить токены (1 файл!)

```css
/* tokens.css */
:root {
  --color-bg-primary: #F8F9FE;
  --color-text-primary: #1F1F1F;
}

[data-theme="dark"] {
  --color-bg-primary: #1F1F1F;
  --color-text-primary: #F8F9FE;
}
```

**Шаг 2:** Всё!

Все компоненты используют токены → автоматически поддерживают темную тему!

**Результат:**
- ✅ 0 изменений в компонентах
- ✅ 0 изменений на страницах
- ✅ 1 файл изменен (tokens.css)

---

## 🔄 Пример изменения: Изменить высоту кнопок

### Задача
Все кнопки размера `md` должны быть 52px вместо 48px.

### Решение

**Изменить 1 строку в 1 файле:**

```css
/* Button.module.css */
.md {
  height: 52px; /* было 48px */
}
```

**Результат:**
- ✅ Все 150+ кнопок на всех страницах изменились
- ✅ 1 файл изменен
- ✅ 0 изменений на страницах

---

## 🔄 Пример изменения: Добавить счетчик XP в header

### Задача
Добавить счетчик XP рядом с монетами в header на всех страницах.

### Решение

**Шаг 1:** Создать компонент (слой 3)

```typescript
// src/components/composite/Header/HeaderXP.tsx
export function HeaderXP({ xp }: { xp: number }) {
  return (
    <Badge variant="info">
      <Icon name="star" />
      <span>{xp}</span>
    </Badge>
  )
}
```

**Шаг 2:** Добавить в Header (слой 3)

```typescript
// src/components/composite/Header/Header.tsx
export function Header({ level, coins, xp, tasksCount }: HeaderProps) {
  return (
    <Card>
      <HeaderLevelBadge level={level} />
      <div className={styles.stats}>
        <HeaderCoins count={coins} />
        <HeaderXP xp={xp} /> {/* Добавили */}
        <HeaderTasks count={tasksCount} />
      </div>
    </Card>
  )
}
```

**Шаг 3:** Обновить MainLayout (слой 5)

```typescript
// src/components/templates/MainLayout.tsx
<Header
  xp={profile.xp} // Добавили
  // ... остальные props
/>
```

**Результат:**
- ✅ XP появился на всех 40+ страницах с MainLayout
- ✅ 3 файла изменены
- ✅ 0 изменений на конечных страницах

---

## 📦 Структура папок проекта

```
src/
├── styles/
│   ├── tokens.css                # Слой 1: Design tokens
│   ├── global.css
│   └── reset.css
│
├── components/
│   ├── ui/                       # Слой 2: UI Kit
│   │   ├── Button/
│   │   ├── Input/
│   │   ├── Card/
│   │   ├── Modal/
│   │   ├── ProgressBar/
│   │   └── ...
│   │
│   ├── composite/                # Слой 3: Составные компоненты
│   │   ├── Header/
│   │   ├── BottomNav/
│   │   ├── AchievementCard/
│   │   ├── AttemptCard/
│   │   └── ...
│   │
│   ├── sections/                 # Слой 4: Секции страниц
│   │   ├── home/
│   │   ├── profile/
│   │   ├── achievements/
│   │   └── ...
│   │
│   └── templates/                # Слой 5: Page templates
│       ├── MainLayout.tsx
│       ├── OnboardingLayout.tsx
│       └── ModalLayout.tsx
│
├── pages/                        # Слой 6: Конечные страницы
│   ├── Home/
│   ├── Help/
│   ├── Check/
│   ├── Achievements/
│   ├── Friends/
│   └── Profile/
│
├── hooks/                        # Custom hooks
├── stores/                       # State management
├── api/                          # API client
├── types/                        # TypeScript types
└── utils/                        # Utilities
```

---

## 🎨 Платформенная адаптация

### Разные токены для разных платформ

```typescript
// src/styles/platforms.ts
export const platformTokens = {
  vk: {
    primaryColor: '#0077FF', // VK синий
    borderRadius: '8px',
  },
  max: {
    primaryColor: '#FF3347', // Max красный
    borderRadius: '12px',
  },
  telegram: {
    primaryColor: '#0088CC', // Telegram синий
    borderRadius: '10px',
  },
}

// Применение в зависимости от платформы
const platform = detectPlatform()
document.documentElement.style.setProperty(
  '--color-primary',
  platformTokens[platform].primaryColor
)
```

---

## ✅ Преимущества архитектуры

### Maintainability
✅ Изменение в одном месте → применяется везде
✅ Легко находить код (четкая структура слоев)
✅ Переиспользование → меньше дублирования

### Scalability
✅ Легко добавлять новые страницы (композиция готовых блоков)
✅ Легко добавлять новые компоненты (следуем структуре)
✅ Можно разделить разработку между командой (разные слои)

### Flexibility
✅ Можно использовать компоненты на разных уровнях
✅ Можно комбинировать по-разному
✅ Можно легко A/B тестировать layouts

### Consistency
✅ Все страницы выглядят одинаково (используют одни токены)
✅ Все кнопки одинаковые (используют один Button)
✅ Все layouts консистентные (используют templates)

---

## 🚀 Практические рекомендации

### DO ✅

1. **Используй токены ВЕЗДЕ**
   ```css
   /* ✅ ХОРОШО */
   .card {
     padding: var(--space-md);
     border-radius: var(--radius-md);
   }

   /* ❌ ПЛОХО */
   .card {
     padding: 16px;
     border-radius: 12px;
   }
   ```

2. **Композируй из готовых компонентов**
   ```typescript
   /* ✅ ХОРОШО */
   <Card>
     <Button>Click me</Button>
   </Card>

   /* ❌ ПЛОХО - создавать новый компонент с нуля */
   <div className="custom-card">
     <div className="custom-button">Click me</div>
   </div>
   ```

3. **Используй правильный слой**
   - Базовый компонент? → `ui/`
   - Составной блок? → `composite/`
   - Секция страницы? → `sections/`
   - Layout? → `templates/`
   - Конечная страница? → `pages/`

4. **Делай компоненты гибкими через props**
   ```typescript
   <Header showCoins={false} /> // Скрыть монеты
   <MainLayout showBottomNav={false} /> // Скрыть навигацию
   ```

### DON'T ❌

1. **НЕ дублируй компоненты**
   ```typescript
   /* ❌ ПЛОХО */
   function HomeButton() { /* копия Button */ }
   function ProfileButton() { /* еще одна копия */ }

   /* ✅ ХОРОШО */
   <Button variant="primary">Home</Button>
   <Button variant="secondary">Profile</Button>
   ```

2. **НЕ хардкоди значения**
   ```css
   /* ❌ ПЛОХО */
   margin: 16px;
   color: #7B5BF2;

   /* ✅ ХОРОШО */
   margin: var(--space-md);
   color: var(--color-primary);
   ```

3. **НЕ создавай страницу с нуля**
   ```typescript
   /* ❌ ПЛОХО */
   function MyPage() {
     return (
       <div>
         <header>...</header>
         <main>...</main>
         <nav>...</nav>
       </div>
     )
   }

   /* ✅ ХОРОШО */
   function MyPage() {
     return (
       <MainLayout>
         <MyContent />
       </MainLayout>
     )
   }
   ```

---

## 📋 Checklist при добавлении нового компонента

- [ ] Определил правильный слой (ui/composite/sections/templates)
- [ ] Использую токены вместо хардкода значений
- [ ] Компонент переиспользуемый (гибкие props)
- [ ] Компонент не дублирует существующий
- [ ] TypeScript типы определены
- [ ] CSS Module для стилей
- [ ] Accessibility (aria-label, role, etc.)
- [ ] Responsive (работает на мобильных)
- [ ] Документация (props, примеры)

---

## 🎓 Итоги

Многослойная архитектура позволяет:

1. **Быстро вносить изменения** → 1 файл вместо 50
2. **Переиспользовать код** → меньше багов, меньше кода
3. **Консистентный дизайн** → все используют одни токены
4. **Легко масштабировать** → добавление новых страниц за минуты
5. **Работать командой** → четкое разделение ответственности

**Следуй этой архитектуре → твой код будет поддерживаемым!**
