# Phase 7: Друзья и реферальная система (Friends & Referral)

**Длительность:** 3-4 дня
**Приоритет:** Средний
**Зависимости:** 02_CORE.md, 04_HOME.md

---

## Цель

Создать реферальную систему с прогрессом приглашений, индикаторами прогресса, наградами за цели, генерацией реферальной ссылки и share sheet интеграцией.

---

## Архитектура Friends

### Структура компонентов

```
FriendsPage
├── ReferralProgressCard
│   ├── GoalDescription ("Пригласи 5 друзей")
│   ├── ProgressIndicator (✓✓345)
│   ├── ProgressText ("2 из 5")
│   └── RewardPreview ("Редкий стикер «Дружба»")
├── ReferralLinkSection
│   ├── ReferralCode/Link
│   ├── CopyButton
│   └── ShareButton
├── InvitedFriendsList
│   └── FriendCard[]
└── RewardUnlockModal
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/referral.ts`

```typescript
export interface ReferralData {
  referralCode: string;
  referralLink: string;
  invitedCount: number;
  targetCount: number;
  progressPercent: number;
  currentGoal: ReferralGoal;
  invitedFriends: InvitedFriend[];
}

export interface ReferralGoal {
  id: string;
  targetCount: number;
  reward: ReferralReward;
  isCompleted: boolean;
  completedAt?: string;
}

export interface ReferralReward {
  type: 'sticker' | 'coins' | 'avatar' | 'premium_feature';
  id: string;
  name: string;
  description: string;
  imageUrl?: string;
  amount?: number;
}

export interface InvitedFriend {
  id: string;
  displayName: string;
  avatarUrl?: string;
  invitedAt: string;
  status: 'pending' | 'active' | 'completed_first_task';
}

export interface ShareChannel {
  type: 'telegram' | 'vk' | 'whatsapp' | 'copy' | 'native';
  name: string;
  icon: string;
}
```

---

## Часть 2: API Integration

### 2.1. Referral API

**Файл:** `src/api/referral.ts`

```typescript
import { apiClient } from './client';
import type {
  ReferralData,
  ReferralGoal,
  InvitedFriend,
} from '@/types/referral';

export const referralAPI = {
  /**
   * Получить реферальные данные
   */
  async getReferralData(childProfileId: string): Promise<ReferralData> {
    return apiClient.get<ReferralData>(`/referrals/${childProfileId}`);
  },

  /**
   * Сгенерировать реферальный код
   */
  async generateReferralCode(
    childProfileId: string
  ): Promise<{ code: string; link: string }> {
    return apiClient.post<{ code: string; link: string }>(
      `/referrals/${childProfileId}/generate`
    );
  },

  /**
   * Получить список приглашённых друзей
   */
  async getInvitedFriends(childProfileId: string): Promise<InvitedFriend[]> {
    return apiClient.get<InvitedFriend[]>(
      `/referrals/${childProfileId}/friends`
    );
  },

  /**
   * Получить текущую цель
   */
  async getCurrentGoal(childProfileId: string): Promise<ReferralGoal> {
    return apiClient.get<ReferralGoal>(`/referrals/${childProfileId}/goal`);
  },

  /**
   * Получить награду за достижение цели
   */
  async claimReferralReward(
    childProfileId: string,
    goalId: string
  ): Promise<void> {
    return apiClient.post<void>(
      `/referrals/${childProfileId}/goal/${goalId}/claim`
    );
  },

  /**
   * Трекинг отправки приглашения
   */
  async trackShareSent(
    childProfileId: string,
    channel: string
  ): Promise<void> {
    return apiClient.post<void>(`/referrals/${childProfileId}/share`, {
      channel,
    });
  },
};
```

---

## Часть 3: Компоненты

### 3.1. FriendsPage Container

**Файл:** `src/pages/Friends/FriendsPage.tsx`

```typescript
import { useEffect, useState } from 'react';
import { Container } from '@/components/layout/Container';
import { BottomNav } from '@/components/layout/BottomNav';
import { ReferralProgressCard } from './components/ReferralProgressCard';
import { ReferralLinkSection } from './components/ReferralLinkSection';
import { InvitedFriendsList } from './components/InvitedFriendsList';
import { RewardUnlockModal } from './components/RewardUnlockModal';
import { useReferralData } from './hooks/useReferralData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { usePlatform } from '@/hooks/usePlatform';
import { useProfileStore } from '@/stores/profileStore';
import { Spinner } from '@/components/ui/Spinner';
import type { ReferralReward } from '@/types/referral';
import styles from './FriendsPage.module.css';

export default function FriendsPage() {
  const analytics = useAnalytics();
  const platform = usePlatform();
  const profile = useProfileStore((state) => state.profile);
  const { data, isLoading, error, refetch } = useReferralData();
  const [showRewardModal, setShowRewardModal] = useState(false);
  const [unlockedReward, setUnlockedReward] = useState<ReferralReward | null>(
    null
  );

  useEffect(() => {
    // Analytics: friends_opened
    analytics.trackEvent('friends_opened', {
      child_profile_id: profile?.id,
    });
  }, [analytics, profile]);

  useEffect(() => {
    if (data) {
      // Analytics: friends_reward_offer_viewed
      analytics.trackEvent('friends_reward_offer_viewed', {
        child_profile_id: profile?.id,
        target_count: data.targetCount,
        current_count: data.invitedCount,
      });

      // Analytics: referral_progress_viewed
      analytics.trackEvent('referral_progress_viewed', {
        child_profile_id: profile?.id,
        invited_count_total: data.invitedCount,
        target_count: data.targetCount,
      });
    }
  }, [data, analytics, profile]);

  const handleCopyLink = (link: string) => {
    // Analytics: referral_link_copied
    analytics.trackEvent('referral_link_copied', {
      child_profile_id: profile?.id,
      referral_code: data?.referralCode,
    });
  };

  const handleShareClick = () => {
    // Analytics: referral_share_opened
    analytics.trackEvent('referral_share_opened', {
      child_profile_id: profile?.id,
      referral_code: data?.referralCode,
    });
  };

  const handleShareSent = (channel: string) => {
    // Analytics: referral_share_sent
    analytics.trackEvent('referral_share_sent', {
      child_profile_id: profile?.id,
      referral_code: data?.referralCode,
      channel_type: channel,
    });

    // Track on backend
    if (data && profile?.id) {
      referralAPI.trackShareSent(profile.id, channel);
    }
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
      <Container className={styles.content}>
        <h1 className={styles.title}>Пригласи друзей</h1>

        <ReferralProgressCard
          invitedCount={data.invitedCount}
          targetCount={data.targetCount}
          progressPercent={data.progressPercent}
          reward={data.currentGoal.reward}
          isCompleted={data.currentGoal.isCompleted}
        />

        <ReferralLinkSection
          referralCode={data.referralCode}
          referralLink={data.referralLink}
          onCopy={handleCopyLink}
          onShare={handleShareClick}
          onShareSent={handleShareSent}
        />

        {data.invitedFriends.length > 0 && (
          <InvitedFriendsList friends={data.invitedFriends} />
        )}
      </Container>

      <BottomNav />

      {showRewardModal && unlockedReward && (
        <RewardUnlockModal
          reward={unlockedReward}
          isOpen={showRewardModal}
          onClose={() => {
            setShowRewardModal(false);
            setUnlockedReward(null);
          }}
        />
      )}
    </div>
  );
}
```

**Файл:** `src/pages/Friends/FriendsPage.module.css`

```css
.page {
  min-height: 100vh;
  background: #f5f5f5;
}

.content {
  padding: 20px 16px 100px 16px;
}

.title {
  font-size: 28px;
  font-weight: 700;
  color: #000;
  margin: 0 0 24px 0;
  text-align: center;
}

.loadingContainer,
.errorContainer {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
}

.errorContainer button {
  margin-top: 16px;
  padding: 12px 24px;
  background: #007aff;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}
```

---

### 3.2. ReferralProgressCard Component

**Файл:** `src/pages/Friends/components/ReferralProgressCard.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import { ProgressBar } from '@/components/ui/ProgressBar';
import type { ReferralReward } from '@/types/referral';
import styles from './ReferralProgressCard.module.css';

interface ReferralProgressCardProps {
  invitedCount: number;
  targetCount: number;
  progressPercent: number;
  reward: ReferralReward;
  isCompleted: boolean;
}

export function ReferralProgressCard({
  invitedCount,
  targetCount,
  progressPercent,
  reward,
  isCompleted,
}: ReferralProgressCardProps) {
  return (
    <Card className={styles.card} variant="primary">
      <h2 className={styles.goal}>
        Пригласи {targetCount} друзей — получи редкий стикер!
      </h2>

      <div className={styles.indicators}>
        {Array.from({ length: targetCount }).map((_, index) => (
          <span key={index} className={styles.indicator}>
            {index < invitedCount ? '✓' : index + 1}
          </span>
        ))}
      </div>

      <p className={styles.progress}>
        {invitedCount} из {targetCount}
      </p>

      <div className={styles.reward}>
        {reward.imageUrl && (
          <img
            src={reward.imageUrl}
            alt={reward.name}
            className={styles.rewardImage}
          />
        )}
        <div className={styles.rewardInfo}>
          <span className={styles.rewardIcon}>⭐</span>
          <span className={styles.rewardName}>{reward.name}</span>
        </div>
      </div>

      {!isCompleted && (
        <ProgressBar
          value={progressPercent}
          size="md"
          className={styles.progressBar}
          variant="success"
        />
      )}

      {isCompleted && (
        <div className={styles.completedBadge}>
          <span>🎉</span>
          <span>Цель достигнута!</span>
        </div>
      )}
    </Card>
  );
}
```

**Файл:** `src/pages/Friends/components/ReferralProgressCard.module.css`

```css
.card {
  padding: 24px;
  margin-bottom: 24px;
}

.goal {
  font-size: 18px;
  font-weight: 600;
  color: #000;
  margin: 0 0 16px 0;
  text-align: center;
  line-height: 1.4;
}

.indicators {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 12px;
}

.indicator {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border: 2px solid #007aff;
  border-radius: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #007aff;
}

.indicator:has-text('✓') {
  background: #007aff;
  color: white;
}

.progress {
  font-size: 16px;
  font-weight: 600;
  color: #000;
  margin: 0 0 16px 0;
  text-align: center;
}

.reward {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  margin-bottom: 16px;
}

.rewardImage {
  width: 64px;
  height: 64px;
  object-fit: contain;
  margin-bottom: 8px;
}

.rewardInfo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rewardIcon {
  font-size: 20px;
}

.rewardName {
  font-size: 16px;
  font-weight: 600;
  color: #000;
}

.progressBar {
  margin-top: 16px;
}

.completedBadge {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: #34c759;
  border-radius: 8px;
  color: white;
  font-size: 16px;
  font-weight: 600;
}
```

---

### 3.3. ReferralLinkSection Component

**Файл:** `src/pages/Friends/components/ReferralLinkSection.tsx`

```typescript
import { useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { usePlatform } from '@/hooks/usePlatform';
import { useAnalytics } from '@/hooks/useAnalytics';
import { ShareSheet } from './ShareSheet';
import styles from './ReferralLinkSection.module.css';

interface ReferralLinkSectionProps {
  referralCode: string;
  referralLink: string;
  onCopy: (link: string) => void;
  onShare: () => void;
  onShareSent: (channel: string) => void;
}

export function ReferralLinkSection({
  referralCode,
  referralLink,
  onCopy,
  onShare,
  onShareSent,
}: ReferralLinkSectionProps) {
  const platform = usePlatform();
  const analytics = useAnalytics();
  const [showShareSheet, setShowShareSheet] = useState(false);
  const [copySuccess, setCopySuccess] = useState(false);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(referralLink);
      setCopySuccess(true);
      onCopy(referralLink);

      setTimeout(() => setCopySuccess(false), 2000);
    } catch (error) {
      console.error('[ReferralLink] Copy failed:', error);
    }
  };

  const handleShare = () => {
    onShare();

    // Use native share if available
    if (navigator.share) {
      navigator
        .share({
          title: 'Объяснятель ДЗ',
          text: 'Присоединяйся ко мне в приложении Объяснятель ДЗ!',
          url: referralLink,
        })
        .then(() => {
          onShareSent('native');
        })
        .catch((error) => {
          console.error('[ReferralLink] Share failed:', error);
        });
    } else {
      // Show custom share sheet
      setShowShareSheet(true);
    }
  };

  return (
    <div className={styles.section}>
      <h3 className={styles.title}>Твоя реферальная ссылка</h3>

      <Card className={styles.linkCard} variant="bordered">
        <div className={styles.linkWrapper}>
          <code className={styles.link}>{referralLink}</code>
        </div>
      </Card>

      <div className={styles.buttons}>
        <Button
          variant="outline"
          size="lg"
          onClick={handleCopy}
          className={styles.button}
        >
          {copySuccess ? '✓ Скопировано' : 'Скопировать'}
        </Button>

        <Button
          variant="primary"
          size="lg"
          onClick={handleShare}
          className={styles.button}
        >
          Отправить
        </Button>
      </div>

      {showShareSheet && (
        <ShareSheet
          link={referralLink}
          isOpen={showShareSheet}
          onClose={() => setShowShareSheet(false)}
          onShare={onShareSent}
        />
      )}
    </div>
  );
}
```

**Файл:** `src/pages/Friends/components/ReferralLinkSection.module.css`

```css
.section {
  margin-bottom: 24px;
}

.title {
  font-size: 16px;
  font-weight: 600;
  color: #000;
  margin: 0 0 12px 0;
}

.linkCard {
  padding: 16px;
  margin-bottom: 12px;
}

.linkWrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.link {
  display: block;
  font-family: 'SF Mono', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  color: #007aff;
  word-break: break-all;
  white-space: nowrap;
}

.buttons {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.button {
  height: 47px;
}
```

---

### 3.4. ShareSheet Component

**Файл:** `src/pages/Friends/components/ShareSheet.tsx`

```typescript
import { Modal } from '@/components/ui/Modal';
import type { ShareChannel } from '@/types/referral';
import styles from './ShareSheet.module.css';

interface ShareSheetProps {
  link: string;
  isOpen: boolean;
  onClose: () => void;
  onShare: (channel: string) => void;
}

export function ShareSheet({ link, isOpen, onClose, onShare }: ShareSheetProps) {
  const channels: ShareChannel[] = [
    { type: 'telegram', name: 'Telegram', icon: '✈️' },
    { type: 'vk', name: 'VKontakte', icon: '🔵' },
    { type: 'whatsapp', name: 'WhatsApp', icon: '📱' },
    { type: 'copy', name: 'Копировать', icon: '📋' },
  ];

  const handleShare = (channel: ShareChannel) => {
    switch (channel.type) {
      case 'telegram':
        window.open(
          `https://t.me/share/url?url=${encodeURIComponent(link)}&text=${encodeURIComponent('Присоединяйся!')}`,
          '_blank'
        );
        break;
      case 'vk':
        window.open(
          `https://vk.com/share.php?url=${encodeURIComponent(link)}`,
          '_blank'
        );
        break;
      case 'whatsapp':
        window.open(
          `https://wa.me/?text=${encodeURIComponent(link)}`,
          '_blank'
        );
        break;
      case 'copy':
        navigator.clipboard.writeText(link);
        break;
    }

    onShare(channel.type);
    onClose();
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Поделиться">
      <div className={styles.channels}>
        {channels.map((channel) => (
          <button
            key={channel.type}
            className={styles.channelButton}
            onClick={() => handleShare(channel)}
          >
            <span className={styles.channelIcon}>{channel.icon}</span>
            <span className={styles.channelName}>{channel.name}</span>
          </button>
        ))}
      </div>
    </Modal>
  );
}
```

---

### 3.5. InvitedFriendsList Component

**Файл:** `src/pages/Friends/components/InvitedFriendsList.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import type { InvitedFriend } from '@/types/referral';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
import styles from './InvitedFriendsList.module.css';

interface InvitedFriendsListProps {
  friends: InvitedFriend[];
}

export function InvitedFriendsList({ friends }: InvitedFriendsListProps) {
  return (
    <div className={styles.section}>
      <div className={styles.header}>
        <h3 className={styles.title}>Приглашённые друзья</h3>
        <span className={styles.count}>{friends.length}</span>
      </div>

      <div className={styles.list}>
        {friends.map((friend) => (
          <Card key={friend.id} className={styles.friendCard} variant="bordered">
            {friend.avatarUrl && (
              <img
                src={friend.avatarUrl}
                alt={friend.displayName}
                className={styles.avatar}
              />
            )}
            <div className={styles.friendInfo}>
              <p className={styles.friendName}>{friend.displayName}</p>
              <p className={styles.friendStatus}>
                {getStatusText(friend.status)} •{' '}
                {formatDistanceToNow(new Date(friend.invitedAt), {
                  addSuffix: true,
                  locale: ru,
                })}
              </p>
            </div>
            <span className={styles.statusIcon}>
              {getStatusIcon(friend.status)}
            </span>
          </Card>
        ))}
      </div>
    </div>
  );
}

function getStatusText(status: InvitedFriend['status']): string {
  switch (status) {
    case 'pending':
      return 'Ожидает регистрации';
    case 'active':
      return 'Зарегистрирован';
    case 'completed_first_task':
      return 'Выполнил первое задание';
  }
}

function getStatusIcon(status: InvitedFriend['status']): string {
  switch (status) {
    case 'pending':
      return '⏳';
    case 'active':
      return '✅';
    case 'completed_first_task':
      return '🎉';
  }
}
```

---

## Часть 4: Аналитические события

### События реферальной системы

| Event Name | Когда срабатывает | Параметры |
|------------|-------------------|-----------|
| `friends_opened` | Открыт экран друзей | `child_profile_id` |
| `friends_reward_offer_viewed` | Просмотрен оффер награды | `child_profile_id`, `target_count`, `current_count` |
| `invite_friend_clicked` | Нажата кнопка приглашения | `child_profile_id`, `referral_code` |
| `referral_link_copied` | Скопирована ссылка | `child_profile_id`, `referral_code` |
| `referral_share_opened` | Открыт share sheet | `child_profile_id`, `referral_code` |
| `referral_share_sent` | Отправлено приглашение | `child_profile_id`, `referral_code`, `channel_type` |
| `referral_progress_viewed` | Просмотрен прогресс | `child_profile_id`, `invited_count_total`, `target_count` |
| `referral_reward_unlocked` | Разблокирована награда (backend) | `child_profile_id`, `reward_type`, `reward_id` |
| `referral_reward_claimed` | Получена награда | `child_profile_id`, `reward_type`, `reward_id` |

---

## Чеклист задач

### Типы и API
- [ ] Создать типы referral.ts
- [ ] Реализовать referralAPI.getReferralData()
- [ ] Реализовать referralAPI.generateReferralCode()
- [ ] Реализовать referralAPI.getInvitedFriends()
- [ ] Реализовать referralAPI.trackShareSent()

### Компоненты
- [ ] Создать FriendsPage container
- [ ] Создать ReferralProgressCard с индикаторами (✓✓345)
- [ ] Создать ReferralLinkSection с Copy/Share
- [ ] Создать ShareSheet модал
- [ ] Создать InvitedFriendsList

### Интеграции
- [ ] Добавить Clipboard API для копирования
- [ ] Добавить Web Share API для native share
- [ ] Добавить ссылки для Telegram/VK/WhatsApp

### Аналитика
- [ ] Добавить friends_opened
- [ ] Добавить referral_link_copied
- [ ] Добавить referral_share_opened/sent
- [ ] Добавить referral_progress_viewed

### UI/UX
- [ ] Добавить индикаторы прогресса
- [ ] Добавить анимацию копирования
- [ ] Добавить список приглашённых друзей
- [ ] Адаптировать под платформы

### Тестирование
- [ ] Протестировать копирование ссылки
- [ ] Протестировать share на разных платформах
- [ ] Протестировать отображение прогресса
- [ ] Протестировать список друзей

---

## Следующий этап

После завершения Friends переходи к **09_PROFILE.md** для создания профиля и настроек.
