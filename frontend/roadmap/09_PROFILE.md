# Phase 8: Профиль, история, настройки (Profile, History, Settings)

**Длительность:** 5-6 дней
**Приоритет:** Высокий
**Зависимости:** 02_CORE.md, 03_ONBOARDING.md, 04_HOME.md

---

## Цель

Создать полный профиль пользователя с карточкой профиля, историей попыток, детальными карточками, отчётом родителю, поддержкой, настройками подписки и paywall экранами.

---

## Архитектура Profile

### Структура компонентов

```
ProfilePage
├── ProfileCard (🦊 Артём, 2 класс, trial)
├── MenuSection
│   ├── HistoryMenuItem
│   ├── ReportMenuItem (parent gate)
│   ├── SupportMenuItem
│   └── SubscriptionMenuItem (paywall)
└── BottomNav

HistoryPage
├── FiltersBar (mode, status, date)
├── HistoryList
│   └── HistoryCard[]
└── HistoryDetailModal
    ├── Images
    ├── Result
    ├── Hints/Feedback
    └── Actions (retry, fix, share)

ReportSettingsPage (Parent Gate)
├── EmailSettings
├── WeeklyReportToggle
├── ArchiveToggle
└── ReportArchive

SupportPage
├── FAQ
├── ContactForm
└── ChatWidget

PaywallPage
├── FeaturesList
├── PricingPlans
├── PaymentButton
└── LaterButton
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/profile.ts`

```typescript
export interface ProfileData {
  id: string;
  displayName: string;
  avatarId: string;
  avatarUrl: string;
  grade: number;
  subscription: SubscriptionData;
}

export interface SubscriptionData {
  status: 'trial' | 'active' | 'expired' | 'cancelled';
  planId?: string;
  planName?: string;
  trialDaysRemaining?: number;
  expiresAt?: string;
}

export interface HistoryAttempt {
  id: string;
  mode: 'help' | 'check';
  status: 'success' | 'error' | 'in_progress';
  scenarioType?: 'single_photo' | 'two_photo';
  createdAt: string;
  completedAt?: string;
  images: HistoryImage[];
  result?: HistoryResult;
  hintsUsed?: number;
}

export interface HistoryImage {
  id: string;
  role: 'task' | 'answer' | 'single';
  url: string;
  thumbnailUrl: string;
}

export interface HistoryResult {
  status: 'correct' | 'has_errors' | 'processing';
  errorCount?: number;
  feedback?: ErrorFeedback[];
  summary?: string;
}

export interface ErrorFeedback {
  id: string;
  stepNumber?: number;
  lineReference?: string;
  description: string;
  locationType: 'line' | 'step' | 'general';
}

export interface HistoryFilters {
  mode?: 'help' | 'check' | 'all';
  status?: 'success' | 'error' | 'in_progress' | 'all';
  dateFrom?: string;
  dateTo?: string;
}

export interface ReportSettings {
  email: string;
  emailVerified: boolean;
  weeklyReportEnabled: boolean;
  archiveEnabled: boolean;
}

export interface WeeklyReport {
  id: string;
  periodStart: string;
  periodEnd: string;
  generatedAt: string;
  downloadUrl?: string;
  stats: {
    totalAttempts: number;
    successfulAttempts: number;
    errorsFixed: number;
    streakDays: number;
  };
}
```

---

## Часть 2: API Integration

### 2.1. Profile & History API

**Файл:** `src/api/profile.ts`

```typescript
import { apiClient } from './client';
import type {
  ProfileData,
  HistoryAttempt,
  HistoryFilters,
  ReportSettings,
  WeeklyReport,
} from '@/types/profile';

export const profileAPI = {
  /**
   * Получить профиль
   */
  async getProfile(childProfileId: string): Promise<ProfileData> {
    return apiClient.get<ProfileData>(`/profiles/child/${childProfileId}`);
  },

  /**
   * Получить историю попыток
   */
  async getHistory(
    childProfileId: string,
    filters?: HistoryFilters
  ): Promise<HistoryAttempt[]> {
    return apiClient.get<HistoryAttempt[]>(`/history/${childProfileId}`, {
      params: filters,
    });
  },

  /**
   * Получить детали попытки
   */
  async getHistoryDetail(
    childProfileId: string,
    attemptId: string
  ): Promise<HistoryAttempt> {
    return apiClient.get<HistoryAttempt>(
      `/history/${childProfileId}/${attemptId}`
    );
  },

  /**
   * Получить настройки отчётов
   */
  async getReportSettings(parentUserId: string): Promise<ReportSettings> {
    return apiClient.get<ReportSettings>(`/reports/${parentUserId}/settings`);
  },

  /**
   * Обновить настройки отчётов
   */
  async updateReportSettings(
    parentUserId: string,
    settings: Partial<ReportSettings>
  ): Promise<void> {
    return apiClient.put<void>(`/reports/${parentUserId}/settings`, settings);
  },

  /**
   * Получить архив отчётов
   */
  async getReportArchive(parentUserId: string): Promise<WeeklyReport[]> {
    return apiClient.get<WeeklyReport[]>(`/reports/${parentUserId}/archive`);
  },

  /**
   * Скачать отчёт
   */
  async downloadReport(
    parentUserId: string,
    reportId: string
  ): Promise<Blob> {
    return apiClient.get<Blob>(`/reports/${parentUserId}/${reportId}/download`, {
      responseType: 'blob',
    });
  },

  /**
   * Отправить сообщение в поддержку
   */
  async sendSupportMessage(
    parentUserId: string,
    message: string
  ): Promise<void> {
    return apiClient.post<void>(`/support/messages`, {
      parentUserId,
      message,
    });
  },
};
```

---

## Часть 3: Компоненты Profile

### 3.1. ProfilePage Container

**Файл:** `src/pages/Profile/ProfilePage.tsx`

```typescript
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { BottomNav } from '@/components/layout/BottomNav';
import { ProfileCard } from './components/ProfileCard';
import { MenuSection } from './components/MenuSection';
import { useProfileData } from './hooks/useProfileData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';
import { Spinner } from '@/components/ui/Spinner';
import styles from './ProfilePage.module.css';

export default function ProfilePage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const { data, isLoading, error, refetch } = useProfileData();

  useEffect(() => {
    // Analytics: profile_opened
    analytics.trackEvent('profile_opened', {
      child_profile_id: profile?.id,
    });
  }, [analytics, profile]);

  const handleHistoryClick = () => {
    // Analytics: profile_history_opened
    analytics.trackEvent('profile_history_opened', {
      child_profile_id: profile?.id,
    });

    navigate(ROUTES.PROFILE_HISTORY);
  };

  const handleReportClick = () => {
    // Analytics: profile_report_settings_opened
    analytics.trackEvent('profile_report_settings_opened', {
      parent_user_id: data?.parentUserId,
    });

    // Show parent gate first
    navigate(ROUTES.PARENT_GATE, {
      state: { nextRoute: ROUTES.REPORT_SETTINGS },
    });
  };

  const handleSupportClick = () => {
    // Analytics: profile_support_opened
    analytics.trackEvent('profile_support_opened', {
      parent_user_id: data?.parentUserId,
    });

    navigate(ROUTES.SUPPORT);
  };

  const handleSubscriptionClick = () => {
    navigate(ROUTES.PAYWALL, {
      state: {
        entry_point: 'profile',
        blocked_feature: 'subscription_management',
      },
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
        <p>Не удалось загрузить профиль</p>
        <button onClick={() => refetch()}>Попробовать снова</button>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <Container className={styles.content}>
        <ProfileCard profile={data} />

        <MenuSection
          onHistoryClick={handleHistoryClick}
          onReportClick={handleReportClick}
          onSupportClick={handleSupportClick}
          onSubscriptionClick={handleSubscriptionClick}
        />
      </Container>

      <BottomNav />
    </div>
  );
}
```

---

### 3.2. ProfileCard Component

**Файл:** `src/pages/Profile/components/ProfileCard.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import type { ProfileData } from '@/types/profile';
import styles from './ProfileCard.module.css';

interface ProfileCardProps {
  profile: ProfileData;
}

export function ProfileCard({ profile }: ProfileCardProps) {
  const subscriptionText = getSubscriptionText(profile.subscription);

  return (
    <Card className={styles.card} variant="default">
      <div className={styles.avatar}>
        <img
          src={profile.avatarUrl}
          alt={profile.displayName}
          className={styles.avatarImage}
        />
      </div>

      <div className={styles.info}>
        <h2 className={styles.name}>{profile.displayName}</h2>
        <p className={styles.grade}>{profile.grade} класс</p>
        <p className={styles.subscription}>{subscriptionText}</p>
      </div>
    </Card>
  );
}

function getSubscriptionText(subscription: ProfileData['subscription']): string {
  switch (subscription.status) {
    case 'trial':
      return `Пробный период — ${subscription.trialDaysRemaining} дней`;
    case 'active':
      return subscription.planName || 'Активная подписка';
    case 'expired':
      return 'Подписка истекла';
    case 'cancelled':
      return 'Подписка отменена';
  }
}
```

**Файл:** `src/pages/Profile/components/ProfileCard.module.css`

```css
.card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  margin-bottom: 24px;
}

.avatar {
  flex-shrink: 0;
}

.avatarImage {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
}

.info {
  flex: 1;
  min-width: 0;
}

.name {
  font-size: 24px;
  font-weight: 700;
  color: #000;
  margin: 0 0 4px 0;
}

.grade {
  font-size: 16px;
  color: #666;
  margin: 0 0 8px 0;
}

.subscription {
  font-size: 14px;
  color: #007aff;
  margin: 0;
  font-weight: 500;
}
```

---

### 3.3. MenuSection Component

**Файл:** `src/pages/Profile/components/MenuSection.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import styles from './MenuSection.module.css';

interface MenuSectionProps {
  onHistoryClick: () => void;
  onReportClick: () => void;
  onSupportClick: () => void;
  onSubscriptionClick: () => void;
}

export function MenuSection({
  onHistoryClick,
  onReportClick,
  onSupportClick,
  onSubscriptionClick,
}: MenuSectionProps) {
  return (
    <div className={styles.section}>
      <Card className={styles.menuItem} variant="bordered" onClick={onHistoryClick}>
        <span className={styles.menuIcon}>📋</span>
        <span className={styles.menuText}>История</span>
        <span className={styles.menuArrow}>›</span>
      </Card>

      <Card className={styles.menuItem} variant="bordered" onClick={onReportClick}>
        <span className={styles.menuIcon}>📊</span>
        <span className={styles.menuText}>Отчёт родителю</span>
        <span className={styles.menuArrow}>›</span>
      </Card>

      <Card className={styles.menuItem} variant="bordered" onClick={onSupportClick}>
        <span className={styles.menuIcon}>💬</span>
        <span className={styles.menuText}>Помощь</span>
        <span className={styles.menuArrow}>›</span>
      </Card>

      <Card
        className={styles.menuItem}
        variant="bordered"
        onClick={onSubscriptionClick}
      >
        <span className={styles.menuIcon}>⭐</span>
        <span className={styles.menuText}>Подписка</span>
        <span className={styles.menuArrow}>›</span>
      </Card>
    </div>
  );
}
```

---

## Часть 4: Компоненты History

### 4.1. HistoryPage Container

**Файл:** `src/pages/Profile/History/HistoryPage.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { FiltersBar } from './components/FiltersBar';
import { HistoryList } from './components/HistoryList';
import { HistoryDetailModal } from './components/HistoryDetailModal';
import { useHistory } from './hooks/useHistory';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import type { HistoryAttempt, HistoryFilters } from '@/types/profile';
import { Spinner } from '@/components/ui/Spinner';
import styles from './HistoryPage.module.css';

export default function HistoryPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const [filters, setFilters] = useState<HistoryFilters>({});
  const [selectedAttempt, setSelectedAttempt] = useState<HistoryAttempt | null>(
    null
  );
  const { data, isLoading, error, refetch } = useHistory(filters);

  useEffect(() => {
    // Analytics: history_opened
    analytics.trackEvent('history_opened', {
      child_profile_id: profile?.id,
    });
  }, [analytics, profile]);

  const handleAttemptClick = (attempt: HistoryAttempt) => {
    // Analytics: history_item_clicked
    analytics.trackEvent('history_item_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attempt.id,
      history_status: attempt.status,
    });

    setSelectedAttempt(attempt);
  };

  const handleFilterChange = (newFilters: HistoryFilters) => {
    // Analytics: history_filter_used
    Object.keys(newFilters).forEach((key) => {
      analytics.trackEvent('history_filter_used', {
        child_profile_id: profile?.id,
        filter_type: key,
        filter_value: newFilters[key as keyof HistoryFilters],
      });
    });

    setFilters(newFilters);
  };

  const handleRetry = (attemptId: string) => {
    // Analytics: history_retry_clicked
    analytics.trackEvent('history_retry_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    // Navigate to appropriate flow
    navigate('/help/upload');
  };

  const handleFixAndRecheck = (attemptId: string) => {
    // Analytics: history_fix_and_recheck_clicked
    analytics.trackEvent('history_fix_and_recheck_clicked', {
      child_profile_id: profile?.id,
      attempt_id: attemptId,
    });

    // Navigate to check flow
    navigate('/check/upload');
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
        <p>Не удалось загрузить историю</p>
        <button onClick={() => refetch()}>Попробовать снова</button>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <Container className={styles.content}>
        <button className={styles.backButton} onClick={() => navigate(-1)}>
          ← Назад
        </button>

        <h1 className={styles.title}>История</h1>

        <FiltersBar filters={filters} onChange={handleFilterChange} />

        <HistoryList attempts={data} onAttemptClick={handleAttemptClick} />
      </Container>

      {selectedAttempt && (
        <HistoryDetailModal
          attempt={selectedAttempt}
          isOpen={!!selectedAttempt}
          onClose={() => setSelectedAttempt(null)}
          onRetry={handleRetry}
          onFixAndRecheck={handleFixAndRecheck}
        />
      )}
    </div>
  );
}
```

---

### 4.2. HistoryDetailModal Component

**Файл:** `src/pages/Profile/History/components/HistoryDetailModal.tsx`

```typescript
import { useEffect } from 'react';
import { Modal } from '@/components/ui/Modal';
import { Button } from '@/components/ui/Button';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import type { HistoryAttempt } from '@/types/profile';
import styles from './HistoryDetailModal.module.css';

interface HistoryDetailModalProps {
  attempt: HistoryAttempt;
  isOpen: boolean;
  onClose: () => void;
  onRetry: (attemptId: string) => void;
  onFixAndRecheck: (attemptId: string) => void;
}

export function HistoryDetailModal({
  attempt,
  isOpen,
  onClose,
  onRetry,
  onFixAndRecheck,
}: HistoryDetailModalProps) {
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  useEffect(() => {
    if (isOpen) {
      // Analytics: history_detail_opened
      analytics.trackEvent('history_detail_opened', {
        child_profile_id: profile?.id,
        attempt_id: attempt.id,
      });
    }
  }, [isOpen, attempt, analytics, profile]);

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="lg">
      <div className={styles.content}>
        <h2 className={styles.title}>
          {attempt.mode === 'help' ? '💡 Помощь' : '✅ Проверка'}
        </h2>

        <div className={styles.images}>
          {attempt.images.map((image) => (
            <img
              key={image.id}
              src={image.url}
              alt={image.role}
              className={styles.image}
            />
          ))}
        </div>

        {attempt.result && (
          <div className={styles.result}>
            <h3 className={styles.resultTitle}>Результат</h3>
            {attempt.result.status === 'correct' && (
              <p className={styles.resultSuccess}>✅ Всё правильно!</p>
            )}
            {attempt.result.status === 'has_errors' && (
              <div className={styles.errors}>
                <p className={styles.errorCount}>
                  Найдено ошибок: {attempt.result.errorCount}
                </p>
                {attempt.result.feedback?.map((error) => (
                  <div key={error.id} className={styles.errorBlock}>
                    {error.stepNumber && (
                      <span className={styles.errorStep}>
                        Шаг {error.stepNumber}
                      </span>
                    )}
                    {error.lineReference && (
                      <span className={styles.errorLine}>{error.lineReference}</span>
                    )}
                    <p className={styles.errorDescription}>{error.description}</p>
                  </div>
                ))}
              </div>
            )}
            {attempt.result.summary && (
              <p className={styles.summary}>{attempt.result.summary}</p>
            )}
          </div>
        )}

        {attempt.hintsUsed !== undefined && (
          <p className={styles.hints}>Использовано подсказок: {attempt.hintsUsed}</p>
        )}

        <div className={styles.actions}>
          <Button variant="outline" onClick={() => onRetry(attempt.id)}>
            Повторить
          </Button>
          {attempt.result?.status === 'has_errors' && (
            <Button variant="primary" onClick={() => onFixAndRecheck(attempt.id)}>
              Исправить и проверить
            </Button>
          )}
        </div>
      </div>
    </Modal>
  );
}
```

---

## Часть 5: Paywall Components

### 5.1. PaywallPage Container

**Файл:** `src/pages/Paywall/PaywallPage.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Button } from '@/components/ui/Button';
import { PricingPlans } from './components/PricingPlans';
import { useAnalytics } from '@/hooks/useAnalytics';
import { usePlatform } from '@/hooks/usePlatform';
import type { BillingPlan } from '@/types/billing';
import styles from './PaywallPage.module.css';

export default function PaywallPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const analytics = useAnalytics();
  const platform = usePlatform();
  const [selectedPlan, setSelectedPlan] = useState<BillingPlan | null>(null);

  const { entry_point, blocked_feature } = (location.state as any) || {};

  useEffect(() => {
    // Analytics: paywall_opened
    analytics.trackEvent('paywall_opened', {
      parent_user_id: platform.userId,
      entry_point,
      blocked_feature,
    });
  }, [analytics, platform, entry_point, blocked_feature]);

  const handlePlanSelect = (plan: BillingPlan) => {
    setSelectedPlan(plan);

    // Analytics: plan_selected
    analytics.trackEvent('plan_selected', {
      parent_user_id: platform.userId,
      billing_plan_id: plan.id,
      billing_period: plan.period,
      price_amount: plan.price,
    });
  };

  const handlePayment = () => {
    if (!selectedPlan) return;

    // Analytics: payment_started
    analytics.trackEvent('payment_started', {
      parent_user_id: platform.userId,
      billing_plan_id: selectedPlan.id,
      amount: selectedPlan.price,
      currency: 'RUB',
    });

    // Navigate to payment flow
    // Payment success/failure tracked by backend
  };

  const handleLater = () => {
    navigate(-1);
  };

  return (
    <div className={styles.page}>
      <Container className={styles.content}>
        <button className={styles.backButton} onClick={() => navigate(-1)}>
          ← Назад
        </button>

        <h1 className={styles.title}>Доступ к сервису</h1>
        <p className={styles.subtitle}>Разблокируй все возможности!</p>

        <div className={styles.features}>
          <div className={styles.feature}>
            <span className={styles.featureIcon}>✅</span>
            <span className={styles.featureText}>Безлимитная проверка ДЗ</span>
          </div>
          <div className={styles.feature}>
            <span className={styles.featureIcon}>💡</span>
            <span className={styles.featureText}>
              Подробные объяснения заданий
            </span>
          </div>
          <div className={styles.feature}>
            <span className={styles.featureIcon}>📊</span>
            <span className={styles.featureText}>
              Еженедельные отчёты для родителей
            </span>
          </div>
          <div className={styles.feature}>
            <span className={styles.featureIcon}>🎮</span>
            <span className={styles.featureText}>Игровые элементы и достижения</span>
          </div>
        </div>

        <PricingPlans
          selectedPlan={selectedPlan}
          onSelectPlan={handlePlanSelect}
        />

        <Button
          variant="primary"
          size="lg"
          isFullWidth
          onClick={handlePayment}
          disabled={!selectedPlan}
        >
          Выбрать план
        </Button>

        <Button variant="ghost" size="lg" isFullWidth onClick={handleLater}>
          Позже
        </Button>
      </Container>
    </div>
  );
}
```

---

## Часть 6: Аналитические события

### События Profile, History, Settings

| Event Name | Когда срабатывает | Параметры |
|------------|-------------------|-----------|
| `profile_opened` | Открыт профиль | `child_profile_id` |
| `profile_history_opened` | Открыт блок истории | `child_profile_id` |
| `profile_report_settings_opened` | Открыт отчёт родителю | `parent_user_id` |
| `profile_support_opened` | Открыта поддержка | `parent_user_id` |
| `profile_parent_gate_opened` | Открыт parent gate | `parent_user_id` |
| `history_opened` | Открыта история | `child_profile_id` |
| `history_item_clicked` | Клик по карточке | `child_profile_id`, `attempt_id`, `history_status` |
| `history_filter_used` | Использован фильтр | `child_profile_id`, `filter_type`, `filter_value` |
| `history_detail_opened` | Открыта детальная карточка | `child_profile_id`, `attempt_id` |
| `history_retry_clicked` | Нажато повторить | `child_profile_id`, `attempt_id` |
| `history_fix_and_recheck_clicked` | Нажато исправить и проверить | `child_profile_id`, `attempt_id` |
| `report_settings_opened` | Открыты настройки отчёта | `parent_user_id` |
| `report_email_changed` | Изменён email | `parent_user_id`, `email_domain` |
| `weekly_report_toggled` | Включён/выключен отчёт | `parent_user_id`, `enabled` |
| `report_archive_opened` | Открыт архив | `parent_user_id` |
| `report_download_clicked` | Скачан отчёт | `parent_user_id`, `report_id` |
| `support_opened` | Открыта поддержка | `parent_user_id`, `screen_name` |
| `support_message_sent` | Отправлено сообщение | `parent_user_id`, `message_length` |
| `paywall_opened` | Открыт paywall | `parent_user_id`, `entry_point`, `blocked_feature` |
| `pricing_opened` | Открыт экран тарифов | `parent_user_id` |
| `plan_selected` | Выбран тариф | `parent_user_id`, `billing_plan_id`, `billing_period`, `price_amount` |
| `payment_started` | Старт оплаты | `parent_user_id`, `billing_plan_id`, `amount`, `currency` |

---

## Чеклист задач

### Profile
- [ ] Создать ProfilePage container
- [ ] Создать ProfileCard (avatar, name, grade, trial)
- [ ] Создать MenuSection (4 пункта меню)
- [ ] Реализовать навигацию

### History
- [ ] Создать HistoryPage container
- [ ] Создать FiltersBar (mode, status, date)
- [ ] Создать HistoryList с карточками
- [ ] Создать HistoryDetailModal
- [ ] Реализовать retry/fix actions

### Report Settings
- [ ] Создать ReportSettingsPage
- [ ] Реализовать parent gate
- [ ] Добавить email settings
- [ ] Добавить weekly report toggle
- [ ] Добавить archive

### Support
- [ ] Создать SupportPage
- [ ] Добавить FAQ
- [ ] Добавить contact form
- [ ] Интегрировать chat widget

### Paywall
- [ ] Создать PaywallPage
- [ ] Создать PricingPlans
- [ ] Интегрировать payment gateway
- [ ] Добавить success/failure экраны

### Тестирование
- [ ] Протестировать все экраны профиля
- [ ] Протестировать фильтрацию истории
- [ ] Протестировать parent gate
- [ ] Протестировать paywall flow

---

## Следующий этап

После завершения Profile переходи к **10_VILLAIN.md** для создания игровой механики.
