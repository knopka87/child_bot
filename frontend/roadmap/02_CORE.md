# Phase 1: Core Infrastructure

**Длительность:** 3-5 дней
**Приоритет:** Критический
**Зависимости:** 01_SETUP.md

---

## ⚠️ КРИТИЧНО: Следуй многослойной архитектуре!

**Перед началом обязательно прочитай:** **[COMPONENT_ARCHITECTURE.md](./COMPONENT_ARCHITECTURE.md)**

Все компоненты должны следовать **6-слойной структуре**:
1. **Design Tokens** → цвета, отступы, шрифты (src/styles/tokens.css)
2. **UI Kit** → Button, Input, Card (src/components/ui/)
3. **Composite** → Header, BottomNav (src/components/composite/)
4. **Sections** → MascotSection (src/components/sections/)
5. **Templates** → MainLayout (src/components/templates/)
6. **Pages** → HomePage (src/pages/)

**Ключевой принцип:** Изменение в одном месте → применяется везде!

---

## Цель

Создать базовую инфраструктуру приложения по многослойной архитектуре: Design Tokens (Layer 1), UI Kit компоненты (Layer 2), API client, Error Boundary, Platform Bridge абстракцию.

---

## Часть 1: VKUI Integration

**ВАЖНО:** Вместо создания Custom UI Kit с нуля, используем **VKUI** - официальную библиотеку компонентов VK.

### Преимущества VKUI:
- ✅ Соответствие VK design guidelines
- ✅ Автоматическая адаптивность под iOS/Android/Desktop
- ✅ Темная/светлая тема из коробки
- ✅ 2500+ готовых иконок
- ✅ Accessibility из коробки
- ✅ Экономия 2-3 дней разработки

---

### 1.1. Настройка VKUI

**Файл:** `src/App.tsx`

```typescript
import { useEffect, useState } from 'react';
import { BrowserRouter } from 'react-router-dom';
import {
  ConfigProvider,
  AdaptivityProvider,
  AppRoot,
  Platform,
  Appearance,
} from '@vkontakte/vkui';
import '@vkontakte/vkui/dist/vkui.css';

import { detectPlatform } from '@/lib/platform/platform-detection';
import { AppRoutes } from '@/routes';

export function App() {
  const [platform, setPlatform] = useState<Platform>('android');
  const [appearance, setAppearance] = useState<Appearance>('light');

  useEffect(() => {
    // Определяем платформу
    detectPlatform().then((info) => {
      setPlatform(info.platform);
    });

    // Определяем тему (можно через VK Bridge)
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    setAppearance(isDark ? 'dark' : 'light');
  }, []);

  return (
    <ConfigProvider appearance={appearance} platform={platform}>
      <AdaptivityProvider>
        <AppRoot>
          <BrowserRouter>
            <AppRoutes />
          </BrowserRouter>
        </AppRoot>
      </AdaptivityProvider>
    </ConfigProvider>
  );
}
```

---

### 1.2. Thin Wrappers над VKUI (опционально)

Создаем тонкие обертки для кастомизации при необходимости:

**Файл:** `src/components/ui/Button/Button.tsx`

```typescript
import { Button as VKUIButton, ButtonProps as VKUIButtonProps } from '@vkontakte/vkui';
import { useHaptics } from '@/lib/platform/haptics';
import { ReactNode } from 'react';

export interface ButtonProps extends VKUIButtonProps {
  children: ReactNode;
  enableHaptics?: boolean;
}

/**
 * Wrapper над VKUI Button с haptic feedback
 */
export function Button({
  children,
  onClick,
  enableHaptics = true,
  ...props
}: ButtonProps) {
  const { onButtonClick } = useHaptics();

  const handleClick = (e: React.MouseEvent<HTMLElement>) => {
    if (enableHaptics) {
      onButtonClick();
    }
    onClick?.(e);
  };

  return (
    <VKUIButton onClick={handleClick} {...props}>
      {children}
    </VKUIButton>
  );
}

// Convenience wrappers
export function PrimaryButton(props: ButtonProps) {
  return <Button mode="primary" size="l" stretched {...props} />;
}

export function SecondaryButton(props: ButtonProps) {
  return <Button mode="secondary" size="l" stretched {...props} />;
}

export function OutlineButton(props: ButtonProps) {
  return <Button mode="outline" size="l" stretched {...props} />;
}
```

**Использование:**

```typescript
import { Button, PrimaryButton } from '@/components/ui/Button';
import { Icon24Add } from '@vkontakte/icons';

// Обычная кнопка
<Button mode="primary" size="l" before={<Icon24Add />}>
  Помоги разобраться
</Button>

// С wrapper
<PrimaryButton before={<Icon24Add />}>
  Помоги разобраться
</PrimaryButton>
```

---

### 1.3. VKUI Card Component

Используем готовый компонент Card из VKUI:

```typescript
import { Card, CardGrid, Div, Title, Text } from '@vkontakte/vkui';

function MascotCard() {
  return (
    <CardGrid size="l">
      <Card mode="shadow">
        <Div>
          <Title level="2" weight="2">
            Маскот
          </Title>
          <Text>Ваш верный помощник!</Text>
        </Div>
      </Card>
    </CardGrid>
  );
}
```

---

### 1.4. VKUI Input Component

```typescript
import { FormItem, Input } from '@vkontakte/vkui';
import { Icon24Cancel } from '@vkontakte/icons';

function EmailInput() {
  const [value, setValue] = useState('');

  return (
    <FormItem top="Email">
      <Input
        type="email"
        placeholder="example@mail.com"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        after={
          value && (
            <Icon24Cancel
              onClick={() => setValue('')}
              style={{ cursor: 'pointer' }}
            />
          )
        }
      />
    </FormItem>
  );
}
```

---

### 1.5. VKUI Modal Component

```typescript
import { ModalRoot, ModalPage, ModalPageHeader } from '@vkontakte/vkui';
import { useRouteNavigator } from '@vkontakte/vk-mini-apps-router';

function Modals({ activeModal, onClose }) {
  return (
    <ModalRoot activeModal={activeModal} onClose={onClose}>
      <ModalPage
        id="hint-modal"
        header={<ModalPageHeader>Подсказка</ModalPageHeader>}
      >
        <Div>
          <Text>Это подсказка уровня 1...</Text>
        </Div>
      </ModalPage>
    </ModalRoot>
  );
}
```

---

### 1.6. VKUI Progress Component

```typescript
import { Progress } from '@vkontakte/vkui';

function LevelProgress({ value }: { value: number }) {
  return (
    <Progress
      value={value}
      appearance="accent"
      aria-label={`Прогресс: ${value}%`}
    />
  );
}
```

---

### 1.7. VKUI Spinner Component

```typescript
import { Spinner } from '@vkontakte/vkui';

function LoadingScreen() {
  return (
    <Div style={{ display: 'flex', justifyContent: 'center', padding: 40 }}>
      <Spinner size="large" />
    </Div>
  );
}
```

---

### 1.8. VKUI Skeleton Loader (вместо Spinner)

**Лучше использовать Skeleton для улучшения UX:**

```typescript
import { Skeleton, Cell, Avatar } from '@vkontakte/vkui';

function AchievementsSkeleton() {
  return (
    <>
      {[1, 2, 3].map((i) => (
        <Cell
          key={i}
          before={<Skeleton><Avatar size={48} /></Skeleton>}
          subtitle={<Skeleton width="60%" height={12} />}
        >
          <Skeleton width="80%" height={16} />
        </Cell>
      ))}
    </>
  );
}

// Использование
function AchievementsList() {
  const { data, isLoading } = useAchievements();

  if (isLoading) {
    return <AchievementsSkeleton />;
  }

  return (
    <>
      {data.achievements.map((achievement) => (
        <AchievementCard key={achievement.id} achievement={achievement} />
      ))}
    </>
  );
}
```

---

### 1.9. VKUI Pull-to-Refresh

```typescript
import { PullToRefresh } from '@vkontakte/vkui';
import { useState } from 'react';

function HomePage() {
  const [isFetching, setIsFetching] = useState(false);
  const { refetch } = useProfile();

  const handleRefresh = async () => {
    setIsFetching(true);
    await refetch();
    setIsFetching(false);
  };

  return (
    <PullToRefresh onRefresh={handleRefresh} isFetching={isFetching}>
      {/* Контент страницы */}
    </PullToRefresh>
  );
}
```

---

### 1.10. VKUI Snackbar (Toast уведомления)

```typescript
import { Snackbar, Avatar } from '@vkontakte/vkui';
import { Icon24Done, Icon24Error } from '@vkontakte/icons';

function useToast() {
  const [snackbar, setSnackbar] = useState<React.ReactNode>(null);

  const showSuccess = (message: string) => {
    setSnackbar(
      <Snackbar
        onClose={() => setSnackbar(null)}
        before={<Avatar size={24}><Icon24Done fill="var(--vkui--color_icon_positive)" /></Avatar>}
      >
        {message}
      </Snackbar>
    );
  };

  const showError = (message: string) => {
    setSnackbar(
      <Snackbar
        onClose={() => setSnackbar(null)}
        before={<Avatar size={24}><Icon24Error fill="var(--vkui--color_icon_negative)" /></Avatar>}
      >
        {message}
      </Snackbar>
    );
  };

  return { snackbar, showSuccess, showError };
}

// Использование
function SomeComponent() {
  const { snackbar, showSuccess } = useToast();

  const handleSubmit = async () => {
    await api.submit();
    showSuccess('Ответ отправлен!');
  };

  return (
    <Group>
      <Button onClick={handleSubmit}>Отправить</Button>
      {snackbar}
    </Group>
  );
}
```

---

### 1.11. VKUI Placeholder (Empty States)

```typescript
import { Placeholder } from '@vkontakte/vkui';
import { Icon56ErrorOutline } from '@vkontakte/icons';

function NoHistoryPlaceholder() {
  return (
    <Placeholder
      icon={<Icon56ErrorOutline />}
      header="История пуста"
      action={<Button size="m">Начать решать</Button>}
    >
      Здесь будут отображаться ваши попытки
    </Placeholder>
  );
}
```

---

### 1.2. Удалены старые Input/Card/Modal компоненты

**Используйте VKUI компоненты напрямую!**

**Все стандартные UI компоненты (Input, Card, Modal, Progress, Spinner) заменены на VKUI аналоги!**

---

## Часть 2: Layout Компоненты

### 2.1. Header Component

**Файл:** `src/components/layout/Header/Header.tsx`

```typescript
import { ProgressBar } from '@/components/ui/ProgressBar';
import styles from './Header.module.css';

export interface HeaderProps {
  level: number;
  levelProgress: number; // 0-100
  coins: number;
  tasksCount: number;
}

export function Header({ level, levelProgress, coins, tasksCount }: HeaderProps) {
  return (
    <header className={styles.header}>
      <div className={styles.levelSection}>
        <div className={styles.levelBadge}>
          <span className={styles.levelNumber}>{level}</span>
          <span className={styles.levelLabel}>Уровень</span>
        </div>
        <ProgressBar
          value={levelProgress}
          variant="success"
          size="sm"
          className={styles.progressBar}
        />
      </div>

      <div className={styles.stats}>
        <div className={styles.stat}>
          <span className={styles.statIcon}>💰</span>
          <span className={styles.statValue}>{coins}</span>
        </div>
        <div className={styles.stat}>
          <span className={styles.statIcon}>✅</span>
          <span className={styles.statValue}>{tasksCount}</span>
        </div>
      </div>
    </header>
  );
}
```

---

### 2.2. BottomNav Component

**Файл:** `src/components/layout/BottomNav/BottomNav.tsx`

```typescript
import { NavLink } from 'react-router-dom';
import { ROUTES } from '@/config/routes';
import clsx from 'clsx';
import styles from './BottomNav.module.css';

export interface NavItem {
  path: string;
  label: string;
  icon: string; // Emoji or icon component
}

const NAV_ITEMS: NavItem[] = [
  { path: ROUTES.HOME, label: 'Главная', icon: '🏠' },
  { path: ROUTES.ACHIEVEMENTS, label: 'Достижения', icon: '🏆' },
  { path: ROUTES.FRIENDS, label: 'Друзья', icon: '👥' },
  { path: ROUTES.PROFILE, label: 'Профиль', icon: '👤' },
];

export function BottomNav() {
  return (
    <nav className={styles.nav}>
      {NAV_ITEMS.map((item) => (
        <NavLink
          key={item.path}
          to={item.path}
          className={({ isActive }) =>
            clsx(styles.navItem, {
              [styles.active]: isActive,
            })
          }
        >
          <span className={styles.icon}>{item.icon}</span>
          <span className={styles.label}>{item.label}</span>
        </NavLink>
      ))}
    </nav>
  );
}
```

---

### 2.3. Container Component

**Файл:** `src/components/layout/Container/Container.tsx`

```typescript
import { ReactNode } from 'react';
import clsx from 'clsx';
import styles from './Container.module.css';

export interface ContainerProps {
  children: ReactNode;
  maxWidth?: 'sm' | 'md' | 'lg' | 'full';
  padding?: boolean;
  className?: string;
}

export function Container({
  children,
  maxWidth = 'md',
  padding = true,
  className,
}: ContainerProps) {
  return (
    <div
      className={clsx(
        styles.container,
        styles[maxWidth],
        {
          [styles.padding]: padding,
        },
        className
      )}
    >
      {children}
    </div>
  );
}
```

---

## Часть 3: API Client

### 3.1. Базовый HTTP клиент

**Файл:** `src/api/client.ts`

```typescript
import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

export class APIClient {
  private client: AxiosInstance;

  constructor(baseURL: string, timeout: number = 30000) {
    this.client = axios.create({
      baseURL,
      timeout,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        // Add auth token if available
        const token = localStorage.getItem('auth_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }

        console.log('[API] Request:', config.method?.toUpperCase(), config.url);
        return config;
      },
      (error) => {
        console.error('[API] Request error:', error);
        return Promise.reject(error);
      }
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => {
        console.log('[API] Response:', response.status, response.config.url);
        return response;
      },
      (error) => {
        console.error('[API] Response error:', error.response?.status, error.message);

        // Handle specific error codes
        if (error.response?.status === 401) {
          // Handle unauthorized
          localStorage.removeItem('auth_token');
          window.location.href = '/onboarding';
        }

        return Promise.reject(error);
      }
    );
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.get<T>(url, config);
    return response.data;
  }

  async post<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.post<T>(url, data, config);
    return response.data;
  }

  async put<T>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.put<T>(url, data, config);
    return response.data;
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.client.delete<T>(url, config);
    return response.data;
  }

  async upload<T>(
    url: string,
    file: Blob,
    onProgress?: (progress: number) => void
  ): Promise<T> {
    const formData = new FormData();
    formData.append('file', file);

    const response = await this.client.post<T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = (progressEvent.loaded / progressEvent.total) * 100;
          onProgress(Math.round(progress));
        }
      },
    });

    return response.data;
  }
}

// Singleton instance
export const apiClient = new APIClient(
  import.meta.env.VITE_API_BASE_URL,
  Number(import.meta.env.VITE_API_TIMEOUT)
);
```

---

### 3.2. API Endpoints

**Файл:** `src/api/endpoints.ts`

```typescript
export const API_ENDPOINTS = {
  // Profile
  PROFILE: '/child-profile/:id',
  UPDATE_PROFILE: '/child-profile/:id',

  // Attempts
  CREATE_ATTEMPT: '/attempts',
  GET_ATTEMPT: '/attempts/:id',
  UPLOAD_IMAGE: '/attempts/:attemptId/images',
  PROCESS_ATTEMPT: '/attempts/:attemptId/process',
  GET_RESULT: '/attempts/:attemptId/result',
  UNFINISHED_ATTEMPT: '/attempts/unfinished',
  RECENT_ATTEMPTS: '/attempts/recent',

  // Achievements
  ACHIEVEMENTS: '/achievements',
  ACHIEVEMENT_DETAIL: '/achievements/:id',

  // Villain
  ACTIVE_VILLAIN: '/villain/active',
  VILLAIN_DETAIL: '/villain/:id',

  // Referrals
  REFERRAL_INFO: '/referrals',
  GENERATE_REFERRAL_CODE: '/referrals/generate',

  // Analytics
  ANALYTICS_EVENT: '/analytics/event',
} as const;

// Helper to replace path params
export function buildEndpoint(
  endpoint: string,
  params: Record<string, string | number>
): string {
  let url = endpoint;
  for (const [key, value] of Object.entries(params)) {
    url = url.replace(`:${key}`, String(value));
  }
  return url;
}
```

---

### 3.3. Profile API

**Файл:** `src/api/profile.ts`

```typescript
import { apiClient } from './client';
import { API_ENDPOINTS, buildEndpoint } from './endpoints';
import type { ChildProfile, UpdateProfileRequest } from '@/types/api';

export const profileAPI = {
  async getProfile(childProfileId: string): Promise<ChildProfile> {
    const url = buildEndpoint(API_ENDPOINTS.PROFILE, { id: childProfileId });
    return apiClient.get<ChildProfile>(url);
  },

  async updateProfile(
    childProfileId: string,
    data: UpdateProfileRequest
  ): Promise<void> {
    const url = buildEndpoint(API_ENDPOINTS.UPDATE_PROFILE, {
      id: childProfileId,
    });
    return apiClient.put<void>(url, data);
  },
};
```

---

## Часть 4: TypeScript Типы

### 4.1. Domain Types

**Файл:** `src/types/domain.ts`

```typescript
// Child Profile
export interface ChildProfile {
  id: string;
  parentUserId: string;
  displayName: string;
  avatarId: string;
  grade: number; // 1-11
  level: number;
  levelProgress: number; // 0-100
  coinsBalance: number;
  tasksSolvedCorrectCount: number;
  winsCount: number;
  checksCorrectCount: number;
  currentStreakDays: number;
  hasUnfinishedAttempt: boolean;
  activeVillainId: string | null;
  activeVillainHealthPercent: number;
  invitedCountTotal: number;
  achievementsUnlockedCount: number;
  mascotId: string;
  mascotState: MascotState;
  createdAt: string;
  updatedAt: string;
}

export type MascotState = 'idle' | 'happy' | 'thinking' | 'celebrating';

// Attempt
export interface Attempt {
  id: string;
  childProfileId: string;
  mode: 'help' | 'check';
  status: AttemptStatus;
  scenarioType: ScenarioType | null;
  createdAt: string;
  updatedAt: string;
}

export type AttemptStatus =
  | 'created'
  | 'uploading'
  | 'uploaded'
  | 'processing'
  | 'completed'
  | 'failed';

export type ScenarioType = 'single_photo' | 'two_photo';

// Attempt Image
export interface AttemptImage {
  id: string;
  attemptId: string;
  imageRole: 'task' | 'answer';
  url: string;
  thumbnailUrl: string;
  uploadedAt: string;
}

// Hint
export interface Hint {
  id: string;
  attemptId: string;
  level: number; // 1, 2, 3
  content: string;
  order: number;
}

// Achievement
export interface Achievement {
  id: string;
  title: string;
  description: string;
  icon: string;
  isUnlocked: boolean;
  unlockedAt: string | null;
  requirement: string;
  rewardType: 'coins' | 'sticker' | 'badge';
  rewardId: string;
}

// Villain
export interface Villain {
  id: string;
  name: string;
  imageUrl: string;
  healthPercent: number; // 0-100
  maxHealth: number;
  currentHealth: number;
  taunt: string;
  defeatedAt: string | null;
}

// Referral
export interface ReferralInfo {
  referralCode: string;
  invitedCount: number;
  targetCount: number;
  rewardType: 'sticker' | 'coins';
  rewardId: string;
  isRewardUnlocked: boolean;
}
```

---

### 4.2. API Types

**Файл:** `src/types/api.ts`

```typescript
import type {
  ChildProfile,
  Attempt,
  AttemptImage,
  Hint,
  Achievement,
  Villain,
  ReferralInfo,
} from './domain';

// Request types
export interface CreateAttemptRequest {
  childProfileId: string;
  mode: 'help' | 'check';
  scenarioType?: 'single_photo' | 'two_photo';
}

export interface UpdateProfileRequest {
  displayName?: string;
  avatarId?: string;
  grade?: number;
}

export interface UploadImageRequest {
  imageRole: 'task' | 'answer';
}

// Response types
export interface AttemptResultResponse {
  attemptId: string;
  status: 'success' | 'error';
  hints?: Hint[];
  errors?: AttemptError[];
  coinsEarned?: number;
  damageDealt?: number;
}

export interface AttemptError {
  id: string;
  stepNumber: number | null;
  lineReference: string | null;
  description: string;
  locationType: 'step' | 'line' | 'general';
}

// API Response wrapper
export interface APIResponse<T> {
  data: T;
  success: boolean;
  error?: {
    code: string;
    message: string;
  };
}
```

---

## Часть 5: Error Boundary

**Файл:** `src/components/ErrorBoundary.tsx`

```typescript
import { Component, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('[ErrorBoundary] Caught error:', error, errorInfo);
    // TODO: Send to error tracking service (Sentry, etc.)
  }

  render() {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return (
        <div style={{ padding: '20px', textAlign: 'center' }}>
          <h2>Что-то пошло не так</h2>
          <p>Попробуйте обновить страницу</p>
          <button
            onClick={() => window.location.reload()}
            style={{
              marginTop: '16px',
              padding: '12px 24px',
              background: '#5181B8',
              color: 'white',
              border: 'none',
              borderRadius: '8px',
              cursor: 'pointer',
            }}
          >
            Обновить
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}
```

---

## Часть 6: Zustand Stores

### 6.1. Profile Store

**Файл:** `src/stores/profileStore.ts`

```typescript
import { create } from 'zustand';
import type { ChildProfile } from '@/types/domain';
import { profileAPI } from '@/api/profile';

interface ProfileState {
  profile: ChildProfile | null;
  isLoading: boolean;
  error: Error | null;
  fetchProfile: (childProfileId: string) => Promise<void>;
  updateProfile: (data: Partial<ChildProfile>) => void;
  reset: () => void;
}

export const useProfileStore = create<ProfileState>((set, get) => ({
  profile: null,
  isLoading: false,
  error: null,

  fetchProfile: async (childProfileId: string) => {
    set({ isLoading: true, error: null });
    try {
      const profile = await profileAPI.getProfile(childProfileId);
      set({ profile, isLoading: false });
    } catch (error) {
      set({ error: error as Error, isLoading: false });
    }
  },

  updateProfile: (data: Partial<ChildProfile>) => {
    const currentProfile = get().profile;
    if (currentProfile) {
      set({ profile: { ...currentProfile, ...data } });
    }
  },

  reset: () => {
    set({ profile: null, isLoading: false, error: null });
  },
}));
```

---

## Чеклист задач

### VKUI Integration
- [ ] Настроить ConfigProvider с platform detection
- [ ] Импортировать VKUI styles в App.tsx
- [ ] Создать thin wrappers для Button (с haptics)
- [ ] Настроить AdaptivityProvider
- [ ] Тестировать VKUI компоненты на разных платформах
- [ ] Создать Snackbar hook для toast уведомлений
- [ ] Создать Skeleton loaders для загрузки

### Haptics & Platform Services
- [ ] Интегрировать haptics в Button wrapper
- [ ] Добавить haptics для важных действий
- [ ] Протестировать haptics на реальном устройстве
- [ ] Настроить platform detection hook

### Layout
- [ ] Создать Header с уровнем, монетами, счетчиком (VKUI)
- [ ] Создать BottomNav с навигацией (Tabbar)
- [ ] Настроить Panel и View структуру VKUI
- [ ] Добавить PullToRefresh на главную

### API
- [ ] Создать базовый API client с interceptors
- [ ] Интегрировать VK Storage для токенов
- [ ] Определить API endpoints
- [ ] Создать Profile API методы
- [ ] Создать типы для API запросов/ответов
- [ ] Добавить retry logic для failed requests

### Infrastructure
- [ ] Создать ErrorBoundary компонент
- [ ] Создать Zustand stores (profile, attempt, platform)
- [ ] Создать доменные TypeScript типы
- [ ] Добавить логирование и error tracking
- [ ] Настроить React Query для server state

### Testing
- [ ] Протестировать VKUI wrappers
- [ ] Протестировать API client
- [ ] Протестировать stores
- [ ] Протестировать haptics
- [ ] Проверить адаптивность на iOS/Android/Desktop

---

## Следующий этап

После завершения Core Infrastructure переходи к **04_HOME.md** для создания главного экрана.
