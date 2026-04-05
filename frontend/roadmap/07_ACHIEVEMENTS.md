# Phase 6: Достижения (Achievements)

**Длительность:** 3-4 дня
**Приоритет:** Средний
**Зависимости:** 02_CORE.md, 04_HOME.md

---

## ⚠️ ВАЖНО: Динамические данные с Backend

**ВСЕ достижения приходят с бекенда!**

Список достижений на дизайне (🔥 5 дней подряд, ✅ 10 проверок ДЗ и т.д.) - это **ПРИМЕРЫ**.

**Реальные данные полностью определяет бекенд:**
- ✅ Список всех достижений
- ✅ Условия разблокировки (requirements)
- ✅ Иконки/эмодзи
- ✅ Названия и описания
- ✅ Текущий прогресс пользователя
- ✅ Размер наград (монеты, стикеры)
- ✅ Порядок на полках (shelf_order)

**Frontend НЕ хардкодит ничего!**

Вместо этого frontend:
- Получает весь список через `GET /api/v1/achievements`
- Рендерит любое количество достижений (не только 12)
- Показывает любые иконки (не только захардкоженные)
- Адаптируется под новые типы достижений без изменения кода

## Цель

Создать универсальную систему отображения достижений, которая рендерит любые достижения, приходящие с бекенда, с прогрессом, детальными модалами и анимацией разблокировки.

---

## Архитектура достижений

### Структура компонентов

```
AchievementsPage
├── Header ("Мои награды", "Собрано 3 из 12")
├── AchievementsGrid
│   ├── AchievementShelf[] (3 полки, по 4 элемента)
│   │   └── AchievementCard[]
│   │       ├── Icon/Image
│   │       ├── Title
│   │       ├── LockedOverlay (если закрыто)
│   │       └── ProgressIndicator (если частично)
├── AchievementDetailModal
│   ├── LargeIcon
│   ├── Title + Description
│   ├── Reward (стикер, монеты, аватар)
│   └── ProgressBar (для незавершённых)
└── UnlockAnimation
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/achievements.ts`

```typescript
// ⚠️ ВАЖНО: Не хардкодим типы! Backend может добавить новые.
// Achievement ID - это string, не enum
export type AchievementID = string // uuid или slug

// Категории тоже приходят с бекенда
export type AchievementCategory =
  | 'streak'
  | 'tasks'
  | 'fixes'
  | 'milestones'
  | 'mastery'
  | string // Backend может добавить новые категории
  | 'speed_solver' // ⚡ Скоростной решатель
  | 'villain_defeater' // 🏆 Победитель злодеев
  | 'wise_owl' // 🦉 Мудрая сова
  | 'collector' // 💎 Коллекционер
  | 'knowledge_rocket' // 🚀 Ракета знаний
  | 'superstar' // 🌟 Суперзвезда
  | 'marathoner' // 🎪 Марафонец
  | 'genius'; // 🧠 Гений

export interface Achievement {
  id: string;
  type: AchievementType;
  title: string;
  description: string;
  icon: string; // emoji or image URL
  isUnlocked: boolean;
  unlockedAt?: string;
  progress: {
    current: number;
    total: number;
    percent: number;
  };
  reward: AchievementReward;
  shelfOrder: number; // 1, 2, 3
  positionInShelf: number; // 0-3
}

export interface AchievementReward {
  type: 'sticker' | 'coins' | 'avatar' | 'badge';
  id: string;
  name: string;
  imageUrl?: string;
  amount?: number; // для coins
}

export interface AchievementsStats {
  unlockedCount: number;
  totalCount: number;
  progressPercent: number;
}
```

---

## Часть 2: API Integration

### 2.1. Achievements API

**Файл:** `src/api/achievements.ts`

```typescript
import { apiClient } from './client';
import type { Achievement, AchievementsStats } from '@/types/achievements';

export const achievementsAPI = {
  /**
   * Получить все достижения
   */
  async getAchievements(childProfileId: string): Promise<Achievement[]> {
    return apiClient.get<Achievement[]>(`/achievements/${childProfileId}`);
  },

  /**
   * Получить статистику достижений
   */
  async getAchievementsStats(childProfileId: string): Promise<AchievementsStats> {
    return apiClient.get<AchievementsStats>(`/achievements/${childProfileId}/stats`);
  },

  /**
   * Получить детали достижения
   */
  async getAchievementDetail(
    childProfileId: string,
    achievementId: string
  ): Promise<Achievement> {
    return apiClient.get<Achievement>(
      `/achievements/${childProfileId}/${achievementId}`
    );
  },

  /**
   * Получить награду за достижение
   */
  async claimAchievementReward(
    childProfileId: string,
    achievementId: string
  ): Promise<void> {
    return apiClient.post<void>(
      `/achievements/${childProfileId}/${achievementId}/claim`
    );
  },
};
```

---

## Часть 3: Компоненты

### 3.1. AchievementsPage Container

**Файл:** `src/pages/Achievements/AchievementsPage.tsx`

```typescript
import { useEffect, useState } from 'react';
import { Container } from '@/components/layout/Container';
import { BottomNav } from '@/components/layout/BottomNav';
import { AchievementsGrid } from './components/AchievementsGrid';
import { AchievementDetailModal } from './components/AchievementDetailModal';
import { UnlockAnimation } from './components/UnlockAnimation';
import { useAchievements } from './hooks/useAchievements';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { Spinner } from '@/components/ui/Spinner';
import type { Achievement } from '@/types/achievements';
import styles from './AchievementsPage.module.css';

export default function AchievementsPage() {
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const { achievements, stats, isLoading, error, refetch } = useAchievements();
  const [selectedAchievement, setSelectedAchievement] = useState<Achievement | null>(
    null
  );
  const [showUnlockAnimation, setShowUnlockAnimation] = useState(false);
  const [newlyUnlocked, setNewlyUnlocked] = useState<Achievement | null>(null);

  useEffect(() => {
    // Analytics: achievements_opened
    analytics.trackEvent('achievements_opened', {
      child_profile_id: profile?.id,
    });
  }, [analytics, profile]);

  const handleAchievementClick = (achievement: Achievement) => {
    // Analytics: achievement_clicked
    analytics.trackEvent('achievement_clicked', {
      child_profile_id: profile?.id,
      achievement_id: achievement.id,
      is_unlocked: achievement.isUnlocked,
    });

    setSelectedAchievement(achievement);
  };

  const handleCloseModal = () => {
    setSelectedAchievement(null);
  };

  const handleShelfView = (shelfOrder: number) => {
    // Analytics: achievement_shelf_viewed
    analytics.trackEvent('achievement_shelf_viewed', {
      child_profile_id: profile?.id,
      shelf_order: shelfOrder,
    });
  };

  useEffect(() => {
    // Check for newly unlocked achievements
    // This would come from WebSocket or polling
    // For now, we'll just demo the animation
  }, []);

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !achievements) {
    return (
      <div className={styles.errorContainer}>
        <p>Не удалось загрузить достижения</p>
        <button onClick={() => refetch()}>Попробовать снова</button>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <Container className={styles.content}>
        <div className={styles.header}>
          <h1 className={styles.title}>Мои награды</h1>
          {stats && (
            <p className={styles.stats}>
              Собрано {stats.unlockedCount} из {stats.totalCount}
            </p>
          )}
        </div>

        <AchievementsGrid
          achievements={achievements}
          onAchievementClick={handleAchievementClick}
          onShelfView={handleShelfView}
        />
      </Container>

      <BottomNav />

      {selectedAchievement && (
        <AchievementDetailModal
          achievement={selectedAchievement}
          isOpen={!!selectedAchievement}
          onClose={handleCloseModal}
        />
      )}

      {showUnlockAnimation && newlyUnlocked && (
        <UnlockAnimation
          achievement={newlyUnlocked}
          onComplete={() => {
            setShowUnlockAnimation(false);
            setNewlyUnlocked(null);
          }}
        />
      )}
    </div>
  );
}
```

**Файл:** `src/pages/Achievements/AchievementsPage.module.css`

```css
.page {
  min-height: 100vh;
  background: #f5f5f5;
}

.content {
  padding: 20px 16px 100px 16px;
}

.header {
  margin-bottom: 24px;
  text-align: center;
}

.title {
  font-size: 28px;
  font-weight: 700;
  color: #000;
  margin: 0 0 8px 0;
}

.stats {
  font-size: 16px;
  color: #666;
  margin: 0;
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

### 3.2. AchievementsGrid Component

**Файл:** `src/pages/Achievements/components/AchievementsGrid.tsx`

```typescript
import { useEffect, useRef } from 'react';
import { AchievementCard } from './AchievementCard';
import type { Achievement } from '@/types/achievements';
import styles from './AchievementsGrid.module.css';

interface AchievementsGridProps {
  achievements: Achievement[];
  onAchievementClick: (achievement: Achievement) => void;
  onShelfView: (shelfOrder: number) => void;
}

export function AchievementsGrid({
  achievements,
  onAchievementClick,
  onShelfView,
}: AchievementsGridProps) {
  const shelfRefs = useRef<(HTMLDivElement | null)[]>([]);

  useEffect(() => {
    // Setup intersection observer for shelf views
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const shelfIndex = shelfRefs.current.indexOf(
              entry.target as HTMLDivElement
            );
            if (shelfIndex !== -1) {
              onShelfView(shelfIndex + 1);
            }
          }
        });
      },
      { threshold: 0.5 }
    );

    shelfRefs.current.forEach((ref) => {
      if (ref) observer.observe(ref);
    });

    return () => observer.disconnect();
  }, [onShelfView]);

  // Group achievements by shelf
  const shelves = [1, 2, 3].map((shelfOrder) =>
    achievements
      .filter((a) => a.shelfOrder === shelfOrder)
      .sort((a, b) => a.positionInShelf - b.positionInShelf)
  );

  return (
    <div className={styles.container}>
      {shelves.map((shelf, index) => (
        <div
          key={index}
          ref={(el) => (shelfRefs.current[index] = el)}
          className={styles.shelf}
        >
          <div className={styles.shelfGrid}>
            {shelf.map((achievement) => (
              <AchievementCard
                key={achievement.id}
                achievement={achievement}
                onClick={() => onAchievementClick(achievement)}
              />
            ))}
          </div>
          <div className={styles.shelfDivider} />
        </div>
      ))}
    </div>
  );
}
```

**Файл:** `src/pages/Achievements/components/AchievementsGrid.module.css`

```css
.container {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.shelf {
  position: relative;
}

.shelfGrid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  padding: 0 8px;
  margin-bottom: 16px;
}

.shelfDivider {
  height: 4px;
  background: linear-gradient(
    to bottom,
    rgba(0, 0, 0, 0.1),
    rgba(0, 0, 0, 0.05)
  );
  border-radius: 2px;
}
```

---

### 3.3. AchievementCard Component

**Файл:** `src/pages/Achievements/components/AchievementCard.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import { ProgressBar } from '@/components/ui/ProgressBar';
import type { Achievement } from '@/types/achievements';
import styles from './AchievementCard.module.css';

interface AchievementCardProps {
  achievement: Achievement;
  onClick: () => void;
}

export function AchievementCard({ achievement, onClick }: AchievementCardProps) {
  const { icon, title, isUnlocked, progress } = achievement;

  return (
    <Card
      className={styles.card}
      variant={isUnlocked ? 'default' : 'bordered'}
      onClick={onClick}
    >
      <div className={styles.iconWrapper}>
        <span className={styles.icon}>{icon}</span>
        {!isUnlocked && <div className={styles.lockedOverlay}>🔒</div>}
      </div>

      <p className={styles.title}>{title}</p>

      {!isUnlocked && progress.total > 0 && (
        <ProgressBar
          value={progress.percent}
          size="xs"
          className={styles.progress}
          showLabel={false}
        />
      )}
    </Card>
  );
}
```

**Файл:** `src/pages/Achievements/components/AchievementCard.module.css`

```css
.card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 12px 8px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  min-height: 99px;
}

.card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.card:active {
  transform: scale(0.95);
}

.iconWrapper {
  position: relative;
  margin-bottom: 8px;
}

.icon {
  font-size: 32px;
  line-height: 1;
}

.lockedOverlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.8);
  font-size: 20px;
  border-radius: 50%;
}

.title {
  font-size: 12px;
  font-weight: 500;
  text-align: center;
  color: #000;
  margin: 0;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.progress {
  margin-top: 8px;
  width: 100%;
}
```

---

### 3.4. AchievementDetailModal Component

**Файл:** `src/pages/Achievements/components/AchievementDetailModal.tsx`

```typescript
import { Modal } from '@/components/ui/Modal';
import { Button } from '@/components/ui/Button';
import { ProgressBar } from '@/components/ui/ProgressBar';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import type { Achievement } from '@/types/achievements';
import styles from './AchievementDetailModal.module.css';

interface AchievementDetailModalProps {
  achievement: Achievement;
  isOpen: boolean;
  onClose: () => void;
}

export function AchievementDetailModal({
  achievement,
  isOpen,
  onClose,
}: AchievementDetailModalProps) {
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  useEffect(() => {
    if (isOpen) {
      // Analytics: achievement_detail_opened
      analytics.trackEvent('achievement_detail_opened', {
        child_profile_id: profile?.id,
        achievement_id: achievement.id,
        is_unlocked: achievement.isUnlocked,
      });

      if (achievement.isUnlocked) {
        // Analytics: achievement_reward_viewed
        analytics.trackEvent('achievement_reward_viewed', {
          child_profile_id: profile?.id,
          achievement_id: achievement.id,
        });
      } else {
        // Analytics: locked_achievement_requirement_viewed
        analytics.trackEvent('locked_achievement_requirement_viewed', {
          child_profile_id: profile?.id,
          achievement_id: achievement.id,
        });
      }
    }
  }, [isOpen, achievement, analytics, profile]);

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className={styles.content}>
        <div className={styles.iconLarge}>
          {achievement.icon}
          {!achievement.isUnlocked && (
            <div className={styles.lockedBadge}>🔒</div>
          )}
        </div>

        <h2 className={styles.title}>{achievement.title}</h2>
        <p className={styles.description}>{achievement.description}</p>

        {achievement.isUnlocked ? (
          <div className={styles.rewardSection}>
            <h3 className={styles.rewardTitle}>Награда</h3>
            <div className={styles.reward}>
              {achievement.reward.imageUrl && (
                <img
                  src={achievement.reward.imageUrl}
                  alt={achievement.reward.name}
                  className={styles.rewardImage}
                />
              )}
              <p className={styles.rewardName}>{achievement.reward.name}</p>
              {achievement.reward.amount && (
                <p className={styles.rewardAmount}>
                  +{achievement.reward.amount}{' '}
                  {achievement.reward.type === 'coins' ? 'монет' : ''}
                </p>
              )}
            </div>
            {achievement.unlockedAt && (
              <p className={styles.unlockedDate}>
                Получено {new Date(achievement.unlockedAt).toLocaleDateString('ru-RU')}
              </p>
            )}
          </div>
        ) : (
          <div className={styles.progressSection}>
            <h3 className={styles.progressTitle}>Прогресс</h3>
            <ProgressBar
              value={achievement.progress.percent}
              size="md"
              showLabel
              label={`${achievement.progress.current} / ${achievement.progress.total}`}
            />
            <p className={styles.progressHint}>
              Ещё {achievement.progress.total - achievement.progress.current}!
            </p>
          </div>
        )}

        <Button variant="primary" isFullWidth onClick={onClose}>
          Закрыть
        </Button>
      </div>
    </Modal>
  );
}
```

**Файл:** `src/pages/Achievements/components/AchievementDetailModal.module.css`

```css
.content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 16px;
}

.iconLarge {
  position: relative;
  font-size: 80px;
  margin-bottom: 16px;
}

.lockedBadge {
  position: absolute;
  top: 0;
  right: -8px;
  font-size: 24px;
}

.title {
  font-size: 24px;
  font-weight: 700;
  color: #000;
  margin: 0 0 8px 0;
  text-align: center;
}

.description {
  font-size: 16px;
  color: #666;
  margin: 0 0 24px 0;
  text-align: center;
  line-height: 1.5;
}

.rewardSection,
.progressSection {
  width: 100%;
  margin-bottom: 24px;
}

.rewardTitle,
.progressTitle {
  font-size: 18px;
  font-weight: 600;
  color: #000;
  margin: 0 0 12px 0;
  text-align: center;
}

.reward {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background: #f5f5f5;
  border-radius: 12px;
}

.rewardImage {
  width: 64px;
  height: 64px;
  object-fit: contain;
  margin-bottom: 8px;
}

.rewardName {
  font-size: 16px;
  font-weight: 600;
  color: #000;
  margin: 0;
}

.rewardAmount {
  font-size: 14px;
  color: #007aff;
  margin: 4px 0 0 0;
}

.unlockedDate {
  font-size: 14px;
  color: #999;
  margin: 12px 0 0 0;
  text-align: center;
}

.progressHint {
  font-size: 14px;
  color: #666;
  margin: 8px 0 0 0;
  text-align: center;
}
```

---

### 3.5. UnlockAnimation Component

**Файл:** `src/pages/Achievements/components/UnlockAnimation.tsx`

```typescript
import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import type { Achievement } from '@/types/achievements';
import styles from './UnlockAnimation.module.css';

interface UnlockAnimationProps {
  achievement: Achievement;
  onComplete: () => void;
}

export function UnlockAnimation({
  achievement,
  onComplete,
}: UnlockAnimationProps) {
  const [phase, setPhase] = useState<'enter' | 'show' | 'exit'>('enter');

  useEffect(() => {
    // Animation timeline:
    // 0-500ms: enter (fade in + scale up)
    // 500-3000ms: show (display achievement)
    // 3000-3500ms: exit (fade out)

    const enterTimer = setTimeout(() => setPhase('show'), 500);
    const showTimer = setTimeout(() => setPhase('exit'), 3000);
    const exitTimer = setTimeout(onComplete, 3500);

    return () => {
      clearTimeout(enterTimer);
      clearTimeout(showTimer);
      clearTimeout(exitTimer);
    };
  }, [onComplete]);

  return createPortal(
    <div className={`${styles.overlay} ${styles[phase]}`}>
      <div className={styles.content}>
        <div className={styles.confetti}>
          <span>🎉</span>
          <span>✨</span>
          <span>🎊</span>
          <span>⭐</span>
        </div>

        <div className={styles.achievement}>
          <div className={styles.icon}>{achievement.icon}</div>
          <h2 className={styles.title}>Достижение разблокировано!</h2>
          <p className={styles.name}>{achievement.title}</p>

          <div className={styles.reward}>
            <p className={styles.rewardLabel}>Награда:</p>
            <p className={styles.rewardName}>{achievement.reward.name}</p>
          </div>
        </div>
      </div>
    </div>,
    document.body
  );
}
```

**Файл:** `src/pages/Achievements/components/UnlockAnimation.module.css`

```css
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
  opacity: 0;
  transition: opacity 0.5s ease;
}

.overlay.show {
  opacity: 1;
}

.overlay.exit {
  opacity: 0;
}

.content {
  position: relative;
  transform: scale(0.5);
  animation: scaleUp 0.5s ease forwards;
}

@keyframes scaleUp {
  to {
    transform: scale(1);
  }
}

.confetti {
  position: absolute;
  top: -50px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 20px;
  font-size: 32px;
  animation: confettiFall 2s ease-out;
}

@keyframes confettiFall {
  0% {
    transform: translateX(-50%) translateY(-50px) rotate(0deg);
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: translateX(-50%) translateY(50px) rotate(360deg);
    opacity: 0;
  }
}

.achievement {
  background: white;
  border-radius: 20px;
  padding: 40px 32px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  min-width: 320px;
}

.icon {
  font-size: 80px;
  margin-bottom: 16px;
  animation: bounce 0.6s ease infinite alternate;
}

@keyframes bounce {
  from {
    transform: translateY(0);
  }
  to {
    transform: translateY(-10px);
  }
}

.title {
  font-size: 20px;
  font-weight: 700;
  color: #000;
  margin: 0 0 8px 0;
}

.name {
  font-size: 24px;
  font-weight: 700;
  color: #007aff;
  margin: 0 0 24px 0;
}

.reward {
  padding: 16px;
  background: #f5f5f5;
  border-radius: 12px;
}

.rewardLabel {
  font-size: 14px;
  color: #666;
  margin: 0 0 4px 0;
}

.rewardName {
  font-size: 16px;
  font-weight: 600;
  color: #000;
  margin: 0;
}
```

---

## Часть 4: Custom Hooks

### 4.1. useAchievements Hook

**Файл:** `src/pages/Achievements/hooks/useAchievements.ts`

```typescript
import { useEffect, useState } from 'react';
import { achievementsAPI } from '@/api/achievements';
import type { Achievement, AchievementsStats } from '@/types/achievements';
import { useProfileStore } from '@/stores/profileStore';

export function useAchievements() {
  const [achievements, setAchievements] = useState<Achievement[] | null>(null);
  const [stats, setStats] = useState<AchievementsStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  const profile = useProfileStore((state) => state.profile);

  const fetchAchievements = async () => {
    if (!profile?.id) {
      setIsLoading(false);
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const [achievementsData, statsData] = await Promise.all([
        achievementsAPI.getAchievements(profile.id),
        achievementsAPI.getAchievementsStats(profile.id),
      ]);

      setAchievements(achievementsData);
      setStats(statsData);
    } catch (err) {
      setError(err as Error);
      console.error('[useAchievements] Failed to fetch:', err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchAchievements();
  }, [profile?.id]);

  return {
    achievements,
    stats,
    isLoading,
    error,
    refetch: fetchAchievements,
  };
}
```

---

## Часть 5: Аналитические события

### События достижений

| Event Name | Когда срабатывает | Параметры |
|------------|-------------------|-----------|
| `achievements_opened` | Открыт экран достижений | `child_profile_id` |
| `achievement_shelf_viewed` | Просмотрена полка достижений | `child_profile_id`, `shelf_order` |
| `achievement_clicked` | Клик по достижению | `child_profile_id`, `achievement_id`, `is_unlocked` |
| `achievement_detail_opened` | Открыта детальная карточка | `child_profile_id`, `achievement_id`, `is_unlocked` |
| `achievement_reward_viewed` | Просмотрена награда открытого достижения | `child_profile_id`, `achievement_id` |
| `locked_achievement_requirement_viewed` | Просмотрено условие закрытого достижения | `child_profile_id`, `achievement_id` |
| `achievement_unlocked` | Достижение открыто (backend) | `child_profile_id`, `achievement_id`, `unlock_reason` |

---

## Часть 6: API Endpoints

### Backend эндпоинты

```typescript
// GET /api/v1/achievements/:childProfileId
// Response:
interface AchievementsResponse {
  achievements: Array<{
    id: string;
    type: string;
    title: string;
    description: string;
    icon: string;
    isUnlocked: boolean;
    unlockedAt?: string;
    progress: {
      current: number;
      total: number;
      percent: number;
    };
    reward: {
      type: 'sticker' | 'coins' | 'avatar' | 'badge';
      id: string;
      name: string;
      imageUrl?: string;
      amount?: number;
    };
    shelfOrder: number;
    positionInShelf: number;
  }>;
}

// GET /api/v1/achievements/:childProfileId/stats
// Response:
interface AchievementsStatsResponse {
  unlockedCount: number;
  totalCount: number;
  progressPercent: number;
}

// POST /api/v1/achievements/:childProfileId/:achievementId/claim
// Response: 204 No Content
```

---

## Чеклист задач

### Типы и API
- [ ] Создать типы achievements.ts
- [ ] Реализовать achievementsAPI.getAchievements()
- [ ] Реализовать achievementsAPI.getAchievementsStats()
- [ ] Реализовать achievementsAPI.getAchievementDetail()

### Компоненты
- [ ] Создать AchievementsPage container
- [ ] Создать AchievementsGrid (3 полки по 4 элемента)
- [ ] Создать AchievementCard с locked overlay
- [ ] Создать AchievementDetailModal
- [ ] Создать UnlockAnimation с confetti

### Аналитика
- [ ] Добавить achievements_opened
- [ ] Добавить achievement_shelf_viewed (intersection observer)
- [ ] Добавить achievement_clicked
- [ ] Добавить achievement_detail_opened
- [ ] Добавить achievement_reward_viewed
- [ ] Добавить locked_achievement_requirement_viewed

### UI/UX
- [ ] Добавить анимацию разблокировки
- [ ] Добавить прогресс-бары для частично открытых
- [ ] Добавить эффект полки (shadow/divider)
- [ ] Адаптировать под разные размеры экранов
- [ ] Добавить haptic feedback при клике

### Тестирование
- [ ] Протестировать отображение всех 12 достижений
- [ ] Протестировать locked/unlocked состояния
- [ ] Протестировать анимацию разблокировки
- [ ] Протестировать детальный модал

---

## Следующий этап

После завершения Achievements переходи к **08_FRIENDS.md** для создания реферальной системы.
