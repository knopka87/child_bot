# Phase 3: Главный экран (Home)

**Длительность:** 4-5 дней
**Приоритет:** Критический
**Зависимости:** 02_CORE.md

---

## Цель

Создать главный экран приложения со всеми интерактивными элементами, динамическими данными, аналитикой и API интеграцией.

---

## Архитектура экрана

### Структура компонентов

```
HomePage
├── Header (level, coins, tasks count)
├── MascotSection
│   ├── MascotCard (mascot image, state, speech)
│   └── VillainCard (villain image, health bar)
├── ActionButtons
│   ├── HelpButton ("Помоги разобраться")
│   └── CheckButton ("Проверка ДЗ")
├── RecentAttempts (last 3 attempts)
│   └── AttemptCard[]
├── UnfinishedAttemptModal
└── BottomNav
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/home.ts`

```typescript
import type { Attempt, Villain } from './domain';

export interface HomeData {
  profile: {
    id: string;
    displayName: string;
    level: number;
    levelProgress: number; // 0-100
    coinsBalance: number;
    tasksSolvedCorrectCount: number;
  };
  mascot: {
    id: string;
    state: MascotState;
    imageUrl: string;
    message: string;
  };
  villain: Villain | null;
  unfinishedAttempt: Attempt | null;
  recentAttempts: RecentAttempt[];
}

export type MascotState = 'idle' | 'happy' | 'thinking' | 'celebrating' | 'encouraging';

export interface RecentAttempt {
  id: string;
  mode: 'help' | 'check';
  status: 'success' | 'error' | 'in_progress';
  createdAt: string;
  thumbnail?: string;
  resultSummary?: string;
}
```

---

## Часть 2: API Integration

### 2.1. Home API

**Файл:** `src/api/home.ts`

```typescript
import { apiClient } from './client';
import type { HomeData, RecentAttempt } from '@/types/home';
import type { Attempt } from '@/types/domain';

export const homeAPI = {
  /**
   * Получить все данные для главного экрана
   */
  async getHomeData(childProfileId: string): Promise<HomeData> {
    return apiClient.get<HomeData>(`/home/${childProfileId}`);
  },

  /**
   * Получить незавершенную попытку
   */
  async getUnfinishedAttempt(childProfileId: string): Promise<Attempt | null> {
    return apiClient.get<Attempt | null>(`/attempts/unfinished`, {
      params: { childProfileId },
    });
  },

  /**
   * Получить последние попытки
   */
  async getRecentAttempts(
    childProfileId: string,
    limit: number = 3
  ): Promise<RecentAttempt[]> {
    return apiClient.get<RecentAttempt[]>(`/attempts/recent`, {
      params: { childProfileId, limit },
    });
  },

  /**
   * Удалить незавершенную попытку
   */
  async deleteAttempt(attemptId: string): Promise<void> {
    return apiClient.delete<void>(`/attempts/${attemptId}`);
  },
};
```

---

## Часть 3: Компоненты

### 3.1. HomePage Container

**Файл:** `src/pages/Home/HomePage.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Header } from '@/components/layout/Header';
import { BottomNav } from '@/components/layout/BottomNav';
import { Container } from '@/components/layout/Container';
import { MascotSection } from './components/MascotSection';
import { ActionButtons } from './components/ActionButtons';
import { RecentAttempts } from './components/RecentAttempts';
import { UnfinishedAttemptModal } from './components/UnfinishedAttemptModal';
import { useHomeData } from './hooks/useHomeData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { ROUTES } from '@/config/routes';
import { Spinner } from '@/components/ui/Spinner';
import styles from './HomePage.module.css';

export default function HomePage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const { data, isLoading, error, refetch } = useHomeData();
  const [showUnfinishedModal, setShowUnfinishedModal] = useState(false);

  useEffect(() => {
    // Analytics: home_opened
    analytics.trackEvent('home_opened', {
      child_profile_id: data?.profile.id,
      entry_point: 'direct', // TODO: Get from navigation state
    });

    // Analytics: level_bar_viewed
    if (data) {
      analytics.trackEvent('level_bar_viewed', {
        child_profile_id: data.profile.id,
        level: data.profile.level,
        level_progress_percent: data.profile.levelProgress,
      });

      analytics.trackEvent('coins_balance_viewed', {
        child_profile_id: data.profile.id,
        coins_balance: data.profile.coinsBalance,
      });

      analytics.trackEvent('tasks_correct_count_viewed', {
        child_profile_id: data.profile.id,
        tasks_solved_correct_count: data.profile.tasksSolvedCorrectCount,
      });
    }
  }, [data, analytics]);

  useEffect(() => {
    // Show unfinished attempt modal if exists
    if (data?.unfinishedAttempt) {
      setShowUnfinishedModal(true);

      // Analytics: unfinished_attempt_modal_shown
      analytics.trackEvent('unfinished_attempt_modal_shown', {
        child_profile_id: data.profile.id,
        attempt_id: data.unfinishedAttempt.id,
        mode: data.unfinishedAttempt.mode,
      });
    }
  }, [data?.unfinishedAttempt, analytics]);

  const handleHelpClick = () => {
    // Analytics: home_help_clicked
    analytics.trackEvent('home_help_clicked', {
      child_profile_id: data?.profile.id,
    });

    navigate(ROUTES.HELP_UPLOAD);
  };

  const handleCheckClick = () => {
    // Analytics: home_check_clicked
    analytics.trackEvent('home_check_clicked', {
      child_profile_id: data?.profile.id,
    });

    navigate(ROUTES.CHECK_SCENARIO);
  };

  const handleMascotClick = () => {
    if (!data) return;

    // Analytics: mascot_clicked
    analytics.trackEvent('mascot_clicked', {
      child_profile_id: data.profile.id,
      mascot_id: data.mascot.id,
      mascot_state: data.mascot.state,
    });

    // TODO: Show mascot stats or joke
  };

  const handleVillainClick = () => {
    if (!data?.villain) return;

    // Analytics: villain_clicked
    analytics.trackEvent('villain_clicked', {
      child_profile_id: data.profile.id,
      villain_id: data.villain.id,
      villain_state: 'active',
    });

    navigate(ROUTES.VILLAIN);
  };

  const handleContinueAttempt = () => {
    if (!data?.unfinishedAttempt) return;

    // Analytics: unfinished_attempt_continue_clicked
    analytics.trackEvent('unfinished_attempt_continue_clicked', {
      child_profile_id: data.profile.id,
      attempt_id: data.unfinishedAttempt.id,
      mode: data.unfinishedAttempt.mode,
    });

    setShowUnfinishedModal(false);

    // Navigate to appropriate flow
    const route =
      data.unfinishedAttempt.mode === 'help'
        ? ROUTES.HELP_PROCESSING
        : ROUTES.CHECK_PROCESSING;

    navigate(route, {
      state: { attemptId: data.unfinishedAttempt.id },
    });
  };

  const handleNewTask = async () => {
    if (!data?.unfinishedAttempt) return;

    // Analytics: unfinished_attempt_new_task_clicked
    analytics.trackEvent('unfinished_attempt_new_task_clicked', {
      child_profile_id: data.profile.id,
      attempt_id: data.unfinishedAttempt.id,
      mode: data.unfinishedAttempt.mode,
    });

    // Delete unfinished attempt
    try {
      await homeAPI.deleteAttempt(data.unfinishedAttempt.id);
      setShowUnfinishedModal(false);
      refetch(); // Refresh data
    } catch (error) {
      console.error('[HomePage] Failed to delete attempt:', error);
    }
  };

  const handleRecentAttemptClick = (attempt: RecentAttempt) => {
    // Analytics: recent_attempt_clicked
    analytics.trackEvent('recent_attempt_clicked', {
      child_profile_id: data?.profile.id,
      attempt_id: attempt.id,
      history_status: attempt.status,
    });

    navigate(ROUTES.PROFILE_HISTORY, {
      state: { attemptId: attempt.id },
    });
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className={styles.errorContainer}>
        <p>Не удалось загрузить данные</p>
        <button onClick={() => refetch()}>Попробовать снова</button>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <Header
        level={data.profile.level}
        levelProgress={data.profile.levelProgress}
        coins={data.profile.coinsBalance}
        tasksCount={data.profile.tasksSolvedCorrectCount}
      />

      <Container className={styles.content}>
        <MascotSection
          mascot={data.mascot}
          villain={data.villain}
          onMascotClick={handleMascotClick}
          onVillainClick={handleVillainClick}
        />

        <ActionButtons
          onHelpClick={handleHelpClick}
          onCheckClick={handleCheckClick}
        />

        {data.recentAttempts.length > 0 && (
          <RecentAttempts
            attempts={data.recentAttempts}
            onAttemptClick={handleRecentAttemptClick}
            onViewAllClick={() => {
              analytics.trackEvent('recent_attempts_view_all_clicked', {
                child_profile_id: data.profile.id,
              });
              navigate(ROUTES.PROFILE_HISTORY);
            }}
          />
        )}
      </Container>

      <BottomNav />

      <UnfinishedAttemptModal
        isOpen={showUnfinishedModal}
        onClose={() => setShowUnfinishedModal(false)}
        onContinue={handleContinueAttempt}
        onNewTask={handleNewTask}
      />
    </div>
  );
}
```

---

### 3.2. MascotSection Component

**Файл:** `src/pages/Home/components/MascotSection.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import { ProgressBar } from '@/components/ui/ProgressBar';
import type { Villain } from '@/types/domain';
import styles from './MascotSection.module.css';

interface MascotData {
  id: string;
  state: 'idle' | 'happy' | 'thinking' | 'celebrating' | 'encouraging';
  imageUrl: string;
  message: string;
}

interface MascotSectionProps {
  mascot: MascotData;
  villain: Villain | null;
  onMascotClick: () => void;
  onVillainClick: () => void;
}

export function MascotSection({
  mascot,
  villain,
  onMascotClick,
  onVillainClick,
}: MascotSectionProps) {
  return (
    <div className={styles.section}>
      {/* Mascot */}
      <div className={styles.mascotWrapper} onClick={onMascotClick}>
        <img
          src={mascot.imageUrl}
          alt="Mascot"
          className={styles.mascotImage}
        />
        <div className={styles.speechBubble}>
          <p className={styles.message}>{mascot.message}</p>
        </div>
      </div>

      {/* Battle indicator */}
      {villain && (
        <div className={styles.battleIndicator}>
          <span className={styles.swordIcon}>⚔️</span>
        </div>
      )}

      {/* Villain */}
      {villain && (
        <Card
          className={styles.villainCard}
          variant="bordered"
          onClick={onVillainClick}
        >
          <img
            src={villain.imageUrl}
            alt={villain.name}
            className={styles.villainImage}
          />
          <div className={styles.villainInfo}>
            <h3 className={styles.villainName}>{villain.name}</h3>
            <ProgressBar
              value={villain.healthPercent}
              variant="error"
              size="sm"
              showLabel
              label={`${villain.currentHealth} / ${villain.maxHealth}`}
            />
          </div>
        </Card>
      )}
    </div>
  );
}
```

**Файл:** `src/pages/Home/components/MascotSection.module.css`

```css
.section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin: 24px 0;
  position: relative;
}

.mascotWrapper {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  cursor: pointer;
  transition: transform 0.2s;
}

.mascotWrapper:active {
  transform: scale(0.98);
}

.mascotImage {
  width: 133px;
  height: 229px;
  object-fit: contain;
}

.speechBubble {
  position: relative;
  background: white;
  border-radius: 12px;
  padding: 12px 16px;
  margin-top: -20px;
  margin-left: 20px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  max-width: 250px;
}

.speechBubble::before {
  content: '';
  position: absolute;
  top: -8px;
  left: 20px;
  width: 0;
  height: 0;
  border-left: 8px solid transparent;
  border-right: 8px solid transparent;
  border-bottom: 8px solid white;
}

.message {
  font-size: 14px;
  line-height: 1.4;
  color: #000;
  margin: 0;
}

.battleIndicator {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 32px;
  animation: pulse 1.5s ease-in-out infinite;
  z-index: 10;
}

@keyframes pulse {
  0%,
  100% {
    transform: translate(-50%, -50%) scale(1);
    opacity: 1;
  }
  50% {
    transform: translate(-50%, -50%) scale(1.2);
    opacity: 0.7;
  }
}

.swordIcon {
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}

.villainCard {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.villainCard:hover {
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.villainCard:active {
  transform: scale(0.98);
}

.villainImage {
  width: 80px;
  height: 80px;
  object-fit: contain;
  flex-shrink: 0;
}

.villainInfo {
  flex: 1;
  min-width: 0;
}

.villainName {
  font-size: 16px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: #000;
}
```

---

### 3.3. ActionButtons Component

**Файл:** `src/pages/Home/components/ActionButtons.tsx`

```typescript
import { Button } from '@/components/ui/Button';
import styles from './ActionButtons.module.css';

interface ActionButtonsProps {
  onHelpClick: () => void;
  onCheckClick: () => void;
}

export function ActionButtons({ onHelpClick, onCheckClick }: ActionButtonsProps) {
  return (
    <div className={styles.buttons}>
      <Button
        variant="primary"
        size="lg"
        isFullWidth
        onClick={onHelpClick}
        className={styles.helpButton}
      >
        <div className={styles.buttonContent}>
          <span className={styles.buttonTitle}>Помоги разобраться</span>
          <span className={styles.buttonSubtitle}>Загрузи фото задания</span>
        </div>
      </Button>

      <Button
        variant="secondary"
        size="lg"
        isFullWidth
        onClick={onCheckClick}
        className={styles.checkButton}
      >
        <div className={styles.buttonContent}>
          <span className={styles.buttonTitle}>Проверка ДЗ</span>
          <span className={styles.buttonSubtitle}>Проверю твою работу</span>
        </div>
      </Button>
    </div>
  );
}
```

**Файл:** `src/pages/Home/components/ActionButtons.module.css`

```css
.buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin: 24px 0;
}

.helpButton,
.checkButton {
  height: 96px;
  padding: 16px;
}

.buttonContent {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.buttonTitle {
  font-size: 18px;
  font-weight: 600;
}

.buttonSubtitle {
  font-size: 14px;
  opacity: 0.8;
}
```

---

### 3.4. RecentAttempts Component

**Файл:** `src/pages/Home/components/RecentAttempts.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import type { RecentAttempt } from '@/types/home';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
import styles from './RecentAttempts.module.css';

interface RecentAttemptsProps {
  attempts: RecentAttempt[];
  onAttemptClick: (attempt: RecentAttempt) => void;
  onViewAllClick: () => void;
}

export function RecentAttempts({
  attempts,
  onAttemptClick,
  onViewAllClick,
}: RecentAttemptsProps) {
  return (
    <div className={styles.section}>
      <div className={styles.header}>
        <h2 className={styles.title}>Последние попытки</h2>
        <button className={styles.viewAllButton} onClick={onViewAllClick}>
          Смотреть все
        </button>
      </div>

      <div className={styles.list}>
        {attempts.map((attempt) => (
          <Card
            key={attempt.id}
            className={styles.attemptCard}
            variant="bordered"
            onClick={() => onAttemptClick(attempt)}
          >
            {attempt.thumbnail && (
              <img
                src={attempt.thumbnail}
                alt="Task"
                className={styles.thumbnail}
              />
            )}
            <div className={styles.attemptInfo}>
              <div className={styles.attemptMeta}>
                <span className={styles.attemptMode}>
                  {attempt.mode === 'help' ? '💡 Помощь' : '✅ Проверка'}
                </span>
                <span className={styles.attemptStatus}>
                  {getStatusIcon(attempt.status)}
                </span>
              </div>
              <p className={styles.attemptSummary}>
                {attempt.resultSummary || 'Задание обработано'}
              </p>
              <span className={styles.attemptTime}>
                {formatDistanceToNow(new Date(attempt.createdAt), {
                  addSuffix: true,
                  locale: ru,
                })}
              </span>
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}

function getStatusIcon(status: RecentAttempt['status']): string {
  switch (status) {
    case 'success':
      return '✅';
    case 'error':
      return '❌';
    case 'in_progress':
      return '⏳';
  }
}
```

---

### 3.5. UnfinishedAttemptModal Component

**Файл:** `src/pages/Home/components/UnfinishedAttemptModal.tsx`

```typescript
import { Modal } from '@/components/ui/Modal';
import { Button } from '@/components/ui/Button';
import styles from './UnfinishedAttemptModal.module.css';

interface UnfinishedAttemptModalProps {
  isOpen: boolean;
  onClose: () => void;
  onContinue: () => void;
  onNewTask: () => void;
}

export function UnfinishedAttemptModal({
  isOpen,
  onClose,
  onContinue,
  onNewTask,
}: UnfinishedAttemptModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title="📝 Незаконченное задание">
      <div className={styles.content}>
        <p className={styles.message}>
          У тебя есть незаконченное задание. Хочешь продолжить?
        </p>

        <div className={styles.buttons}>
          <Button variant="primary" isFullWidth onClick={onContinue}>
            Продолжить
          </Button>
          <Button variant="outline" isFullWidth onClick={onNewTask}>
            Новое задание
          </Button>
        </div>
      </div>
    </Modal>
  );
}
```

**Файл:** `src/pages/Home/components/UnfinishedAttemptModal.module.css`

```css
.content {
  padding: 16px 0;
}

.message {
  font-size: 16px;
  line-height: 1.5;
  color: #000;
  margin: 0 0 24px 0;
  text-align: center;
}

.buttons {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
```

---

## Часть 4: Custom Hooks

### 4.1. useHomeData Hook

**Файл:** `src/pages/Home/hooks/useHomeData.ts`

```typescript
import { useEffect, useState } from 'react';
import { homeAPI } from '@/api/home';
import type { HomeData } from '@/types/home';
import { useProfileStore } from '@/stores/profileStore';

export function useHomeData() {
  const [data, setData] = useState<HomeData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const profile = useProfileStore((state) => state.profile);

  const fetchData = async () => {
    if (!profile?.id) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const homeData = await homeAPI.getHomeData(profile.id);
      setData(homeData);
    } catch (err) {
      setError(err as Error);
      console.error('[useHomeData] Failed to fetch home data:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [profile?.id]);

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  };
}
```

---

## Часть 5: Аналитические события

### События, отправляемые на Home экране

| Event Name | Когда срабатывает | Параметры |
|------------|-------------------|-----------|
| `home_opened` | Экран открыт | `child_profile_id`, `entry_point` |
| `level_bar_viewed` | Показан прогресс уровня | `child_profile_id`, `level`, `level_progress_percent` |
| `coins_balance_viewed` | Показан баланс монет | `child_profile_id`, `coins_balance` |
| `tasks_correct_count_viewed` | Показан счетчик заданий | `child_profile_id`, `tasks_solved_correct_count` |
| `home_help_clicked` | Нажата кнопка "Помоги разобраться" | `child_profile_id` |
| `home_check_clicked` | Нажата кнопка "Проверка ДЗ" | `child_profile_id` |
| `unfinished_attempt_modal_shown` | Показан модал незавершённой попытки | `child_profile_id`, `attempt_id`, `mode` |
| `unfinished_attempt_continue_clicked` | Нажато "Продолжить" | `child_profile_id`, `attempt_id`, `mode` |
| `unfinished_attempt_new_task_clicked` | Нажато "Новое задание" | `child_profile_id`, `attempt_id`, `mode` |
| `mascot_clicked` | Клик по маскоту | `child_profile_id`, `mascot_id`, `mascot_state` |
| `villain_clicked` | Клик по злодею | `child_profile_id`, `villain_id`, `villain_state` |
| `recent_attempt_clicked` | Клик по карточке попытки | `child_profile_id`, `attempt_id`, `history_status` |
| `recent_attempts_view_all_clicked` | Нажато "Смотреть все" | `child_profile_id` |

---

## Часть 6: API Endpoints

### Эндпоинты для Home экрана

```typescript
// GET /api/v1/home/:childProfileId
// Response:
interface HomeDataResponse {
  profile: {
    id: string;
    displayName: string;
    level: number;
    levelProgress: number;
    coinsBalance: number;
    tasksSolvedCorrectCount: number;
  };
  mascot: {
    id: string;
    state: 'idle' | 'happy' | 'thinking' | 'celebrating' | 'encouraging';
    imageUrl: string;
    message: string;
  };
  villain: {
    id: string;
    name: string;
    imageUrl: string;
    healthPercent: number;
    currentHealth: number;
    maxHealth: number;
    taunt: string;
  } | null;
  unfinishedAttempt: {
    id: string;
    mode: 'help' | 'check';
    status: 'created' | 'uploading' | 'uploaded' | 'processing';
    createdAt: string;
  } | null;
  recentAttempts: Array<{
    id: string;
    mode: 'help' | 'check';
    status: 'success' | 'error' | 'in_progress';
    createdAt: string;
    thumbnail?: string;
    resultSummary?: string;
  }>;
}

// DELETE /api/v1/attempts/:attemptId
// Response: 204 No Content
```

---

## Часть 7: Состояния и переходы

### State Machine для Home экрана

```typescript
type HomeState =
  | 'loading'
  | 'idle'
  | 'show_unfinished_modal'
  | 'error';

type HomeAction =
  | { type: 'FETCH_SUCCESS'; data: HomeData }
  | { type: 'FETCH_ERROR'; error: Error }
  | { type: 'SHOW_MODAL' }
  | { type: 'HIDE_MODAL' }
  | { type: 'RETRY' };

// Transitions:
// loading -> idle (on success)
// loading -> error (on failure)
// idle -> show_unfinished_modal (if unfinished attempt exists)
// show_unfinished_modal -> idle (on close/continue/new task)
// error -> loading (on retry)
```

---

## Часть 8: Обработка ошибок

### Типы ошибок

1. **Network Error** - нет интернета
2. **API Error 404** - профиль не найден
3. **API Error 500** - серверная ошибка
4. **Timeout Error** - превышено время ожидания

### Error Handling Strategy

```typescript
function handleHomeError(error: Error): void {
  // Log error
  console.error('[HomePage] Error:', error);

  // Send to analytics
  analytics.trackEvent('ui_error_shown', {
    screen_name: 'home',
    error_code: getErrorCode(error),
  });

  // Show user-friendly message
  if (error.message.includes('Network')) {
    showToast('Проверьте подключение к интернету');
  } else if (error.message.includes('404')) {
    showToast('Профиль не найден');
    navigate(ROUTES.ONBOARDING);
  } else {
    showToast('Что-то пошло не так. Попробуйте позже');
  }
}
```

---

## Чеклист задач

### Компоненты
- [ ] Создать HomePage container
- [ ] Создать MascotSection с маскотом и злодеем
- [ ] Создать ActionButtons (Help, Check)
- [ ] Создать RecentAttempts список
- [ ] Создать UnfinishedAttemptModal

### API
- [ ] Реализовать homeAPI.getHomeData()
- [ ] Реализовать homeAPI.getUnfinishedAttempt()
- [ ] Реализовать homeAPI.getRecentAttempts()
- [ ] Реализовать homeAPI.deleteAttempt()

### Hooks
- [ ] Создать useHomeData hook
- [ ] Добавить auto-refresh при возврате на экран

### Аналитика
- [ ] Добавить событие home_opened
- [ ] Добавить события level/coins/tasks viewed
- [ ] Добавить события для кликов по кнопкам
- [ ] Добавить события для модала незавершённой попытки
- [ ] Добавить события для кликов по маскоту/злодею

### Стили и анимации
- [ ] Добавить анимацию появления маскота
- [ ] Добавить анимацию битвы (⚔️)
- [ ] Добавить hover/active состояния для карточек
- [ ] Адаптировать под разные размеры экранов

### Тестирование
- [ ] Протестировать загрузку данных
- [ ] Протестировать модал незавершённой попытки
- [ ] Протестировать навигацию к Help/Check
- [ ] Протестировать обработку ошибок

---

## Следующий этап

После завершения Home экрана переходи к **05_HELP.md** для создания потока помощи с заданиями.
