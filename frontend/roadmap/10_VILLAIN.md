# Phase 9: Злодей и игровая механика (Villain Game Mechanics)

**Длительность:** 3-4 дня
**Приоритет:** Средний
**Зависимости:** 02_CORE.md, 04_HOME.md

---

## Цель

Создать игровую механику битвы со злодеем Кракозябра: экран злодея с репликами, индикатор здоровья, механику урона от правильных ответов, экран победы и награды.

---

## Архитектура Villain

### Структура компонентов

```
VillainPage
├── VillainCard
│   ├── VillainImage (белый персонаж с короной)
│   ├── VillainName ("Кракозябра")
│   ├── HealthBar
│   └── TauntBubble ("Ха-ха! Попробуй-ка реши задачки!")
├── BattleStats
│   ├── DamageDealt
│   ├── AttacksCount
│   └── ProgressToVictory
└── ActionButton ("К заданиям")

VictoryPage
├── VictoryAnimation (confetti, celebration)
├── VillainDefeatedImage
├── RewardCard
│   ├── RewardImage (редкий стикер)
│   ├── RewardName
│   └── AchievementBadge
└── ContinueButton ("Продолжить учиться")
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/villain.ts`

```typescript
export interface Villain {
  id: string;
  name: string;
  description: string;
  imageUrl: string;
  healthPercent: number;
  currentHealth: number;
  maxHealth: number;
  taunt: string;
  isActive: boolean;
  isDefeated: boolean;
}

export interface VillainBattle {
  villainId: string;
  childProfileId: string;
  damageDealt: number;
  attacksCount: number;
  startedAt: string;
  defeatedAt?: string;
}

export interface VillainAttack {
  attemptId: string;
  damage: number;
  reason: 'correct_answer' | 'error_fixed' | 'hint_completed';
  villainHealthBefore: number;
  villainHealthAfter: number;
}

export interface VillainVictory {
  villainId: string;
  defeatedAt: string;
  totalDamage: number;
  totalAttempts: number;
  rewards: VillainReward[];
}

export interface VillainReward {
  type: 'sticker' | 'achievement' | 'coins' | 'avatar';
  id: string;
  name: string;
  description: string;
  imageUrl?: string;
  amount?: number;
  rarity?: 'common' | 'rare' | 'epic' | 'legendary';
}

export type VillainTaunt =
  | 'Ха-ха! Попробуй-ка реши задачки!'
  | 'Думаешь, справишься?'
  | 'Я непобедим!'
  | 'Ещё немного, и ты сдашься!'
  | 'Ну давай, удиви меня!';
```

---

## Часть 2: API Integration

### 2.1. Villain API

**Файл:** `src/api/villain.ts`

```typescript
import { apiClient } from './client';
import type {
  Villain,
  VillainBattle,
  VillainAttack,
  VillainVictory,
} from '@/types/villain';

export const villainAPI = {
  /**
   * Получить активного злодея
   */
  async getActiveVillain(childProfileId: string): Promise<Villain | null> {
    return apiClient.get<Villain | null>(`/villains/${childProfileId}/active`);
  },

  /**
   * Получить детали битвы
   */
  async getBattleDetails(
    childProfileId: string,
    villainId: string
  ): Promise<VillainBattle> {
    return apiClient.get<VillainBattle>(
      `/villains/${childProfileId}/battle/${villainId}`
    );
  },

  /**
   * Получить детали победы
   */
  async getVictoryDetails(
    childProfileId: string,
    villainId: string
  ): Promise<VillainVictory> {
    return apiClient.get<VillainVictory>(
      `/villains/${childProfileId}/victory/${villainId}`
    );
  },

  /**
   * Получить новую реплику злодея
   */
  async getRandomTaunt(villainId: string): Promise<{ taunt: string }> {
    return apiClient.get<{ taunt: string }>(`/villains/${villainId}/taunt`);
  },
};
```

---

## Часть 3: Компоненты

### 3.1. VillainPage Container

**Файл:** `src/pages/Villain/VillainPage.tsx`

```typescript
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { VillainCard } from './components/VillainCard';
import { BattleStats } from './components/BattleStats';
import { Button } from '@/components/ui/Button';
import { useVillainData } from './hooks/useVillainData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';
import { Spinner } from '@/components/ui/Spinner';
import styles from './VillainPage.module.css';

export default function VillainPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const { villain, battle, isLoading, error, refetch } = useVillainData();

  useEffect(() => {
    if (villain) {
      // Analytics: villain_screen_opened
      analytics.trackEvent('villain_screen_opened', {
        child_profile_id: profile?.id,
        villain_id: villain.id,
      });

      // Analytics: villain_taunt_viewed
      analytics.trackEvent('villain_taunt_viewed', {
        child_profile_id: profile?.id,
        villain_id: villain.id,
      });
    }
  }, [villain, analytics, profile]);

  const handleGoToTasks = () => {
    navigate(ROUTES.HOME);
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !villain) {
    return (
      <div className={styles.errorContainer}>
        <p>Нет активного злодея</p>
        <button onClick={() => navigate(ROUTES.HOME)}>На главную</button>
      </div>
    );
  }

  // If villain is defeated, redirect to victory page
  if (villain.isDefeated) {
    navigate(ROUTES.VILLAIN_VICTORY, {
      state: { villainId: villain.id },
    });
    return null;
  }

  return (
    <div className={styles.page}>
      <Container className={styles.content}>
        <button className={styles.backButton} onClick={() => navigate(-1)}>
          ← Назад
        </button>

        <VillainCard villain={villain} />

        {battle && <BattleStats battle={battle} />}

        <div className={styles.actions}>
          <Button variant="primary" size="lg" isFullWidth onClick={handleGoToTasks}>
            К заданиям
          </Button>
          <p className={styles.hint}>
            Решай задания правильно, чтобы нанести урон злодею!
          </p>
        </div>
      </Container>
    </div>
  );
}
```

**Файл:** `src/pages/Villain/VillainPage.module.css`

```css
.page {
  min-height: 100vh;
  background: linear-gradient(180deg, #fff 0%, #f5f5f5 100%);
}

.content {
  padding: 20px 16px;
}

.backButton {
  background: none;
  border: none;
  color: #666;
  font-size: 16px;
  cursor: pointer;
  padding: 8px;
  margin-bottom: 16px;
}

.actions {
  margin-top: 32px;
}

.hint {
  font-size: 14px;
  color: #666;
  text-align: center;
  margin: 12px 0 0 0;
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

### 3.2. VillainCard Component

**Файл:** `src/pages/Villain/components/VillainCard.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import { ProgressBar } from '@/components/ui/ProgressBar';
import type { Villain } from '@/types/villain';
import styles from './VillainCard.module.css';

interface VillainCardProps {
  villain: Villain;
}

export function VillainCard({ villain }: VillainCardProps) {
  return (
    <Card className={styles.card} variant="default">
      <div className={styles.imageWrapper}>
        <img
          src={villain.imageUrl}
          alt={villain.name}
          className={styles.villainImage}
        />
      </div>

      <div className={styles.info}>
        <h2 className={styles.name}>{villain.name}</h2>
        <p className={styles.description}>{villain.description}</p>

        <div className={styles.health}>
          <div className={styles.healthLabel}>
            <span>Здоровье</span>
            <span>
              {villain.currentHealth} / {villain.maxHealth}
            </span>
          </div>
          <ProgressBar
            value={villain.healthPercent}
            variant="error"
            size="lg"
            showLabel={false}
          />
        </div>
      </div>

      <div className={styles.tauntBubble}>
        <p className={styles.taunt}>{villain.taunt}</p>
        <div className={styles.tauntArrow} />
      </div>
    </Card>
  );
}
```

**Файл:** `src/pages/Villain/components/VillainCard.module.css`

```css
.card {
  padding: 24px;
  position: relative;
  overflow: visible;
  margin-bottom: 24px;
}

.imageWrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.villainImage {
  width: 200px;
  height: 200px;
  object-fit: contain;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.info {
  text-align: center;
}

.name {
  font-size: 28px;
  font-weight: 700;
  color: #000;
  margin: 0 0 8px 0;
}

.description {
  font-size: 16px;
  color: #666;
  margin: 0 0 20px 0;
}

.health {
  margin-top: 16px;
}

.healthLabel {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
  font-weight: 600;
}

.tauntBubble {
  position: absolute;
  top: -20px;
  right: 20px;
  background: white;
  border: 2px solid #000;
  border-radius: 16px;
  padding: 12px 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  max-width: 200px;
  animation: bounce 2s ease-in-out infinite;
}

@keyframes bounce {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.taunt {
  font-size: 14px;
  font-weight: 600;
  color: #000;
  margin: 0;
  line-height: 1.4;
}

.tauntArrow {
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 0;
  border-left: 8px solid transparent;
  border-right: 8px solid transparent;
  border-top: 8px solid #000;
}

.tauntArrow::before {
  content: '';
  position: absolute;
  top: -10px;
  left: -7px;
  width: 0;
  height: 0;
  border-left: 7px solid transparent;
  border-right: 7px solid transparent;
  border-top: 7px solid white;
}
```

---

### 3.3. BattleStats Component

**Файл:** `src/pages/Villain/components/BattleStats.tsx`

```typescript
import { Card } from '@/components/ui/Card';
import type { VillainBattle } from '@/types/villain';
import styles from './BattleStats.module.css';

interface BattleStatsProps {
  battle: VillainBattle;
}

export function BattleStats({ battle }: BattleStatsProps) {
  return (
    <Card className={styles.card} variant="bordered">
      <h3 className={styles.title}>Статистика битвы</h3>

      <div className={styles.stats}>
        <div className={styles.stat}>
          <span className={styles.statIcon}>⚔️</span>
          <div className={styles.statInfo}>
            <span className={styles.statLabel}>Урон нанесён</span>
            <span className={styles.statValue}>{battle.damageDealt}</span>
          </div>
        </div>

        <div className={styles.stat}>
          <span className={styles.statIcon}>🎯</span>
          <div className={styles.statInfo}>
            <span className={styles.statLabel}>Атак совершено</span>
            <span className={styles.statValue}>{battle.attacksCount}</span>
          </div>
        </div>
      </div>
    </Card>
  );
}
```

**Файл:** `src/pages/Villain/components/BattleStats.module.css`

```css
.card {
  padding: 20px;
}

.title {
  font-size: 18px;
  font-weight: 600;
  color: #000;
  margin: 0 0 16px 0;
  text-align: center;
}

.stats {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f5f5f5;
  border-radius: 8px;
}

.statIcon {
  font-size: 24px;
}

.statInfo {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.statLabel {
  font-size: 14px;
  color: #666;
}

.statValue {
  font-size: 20px;
  font-weight: 700;
  color: #000;
}
```

---

### 3.4. VictoryPage Container

**Файл:** `src/pages/Villain/VictoryPage.tsx`

```typescript
import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Container } from '@/components/layout/Container';
import { Button } from '@/components/ui/Button';
import { VictoryAnimation } from './components/VictoryAnimation';
import { RewardCard } from './components/RewardCard';
import { useVictoryData } from './hooks/useVictoryData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';
import { Spinner } from '@/components/ui/Spinner';
import styles from './VictoryPage.module.css';

export default function VictoryPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const { villainId } = (location.state as any) || {};
  const { victory, isLoading, error } = useVictoryData(villainId);

  useEffect(() => {
    if (victory) {
      // Analytics: victory_screen_opened
      analytics.trackEvent('victory_screen_opened', {
        child_profile_id: profile?.id,
        villain_id: villainId,
        attempt_id: '', // Would come from context
      });

      // Analytics: victory_reward_viewed for each reward
      victory.rewards.forEach((reward) => {
        analytics.trackEvent('victory_reward_viewed', {
          child_profile_id: profile?.id,
          villain_id: villainId,
          reward_type: reward.type,
          reward_id: reward.id,
        });
      });
    }
  }, [victory, villainId, analytics, profile]);

  const handleContinue = () => {
    // Analytics: victory_continue_clicked
    analytics.trackEvent('victory_continue_clicked', {
      child_profile_id: profile?.id,
      villain_id: villainId,
    });

    navigate(ROUTES.HOME);
  };

  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !victory) {
    return (
      <div className={styles.errorContainer}>
        <p>Не удалось загрузить данные победы</p>
        <button onClick={() => navigate(ROUTES.HOME)}>На главную</button>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <VictoryAnimation />

      <Container className={styles.content}>
        <div className={styles.header}>
          <h1 className={styles.title}>🎉 Победа!</h1>
          <p className={styles.subtitle}>Ты победил злодея!</p>
        </div>

        <div className={styles.villainDefeated}>
          <img
            src="/assets/villain-defeated.png"
            alt="Defeated villain"
            className={styles.villainImage}
          />
        </div>

        <div className={styles.stats}>
          <div className={styles.statItem}>
            <span className={styles.statLabel}>Урон нанесён</span>
            <span className={styles.statValue}>{victory.totalDamage}</span>
          </div>
          <div className={styles.statItem}>
            <span className={styles.statLabel}>Попыток</span>
            <span className={styles.statValue}>{victory.totalAttempts}</span>
          </div>
        </div>

        <div className={styles.rewards}>
          <h3 className={styles.rewardsTitle}>Твои награды</h3>
          {victory.rewards.map((reward) => (
            <RewardCard key={reward.id} reward={reward} />
          ))}
        </div>

        <Button variant="primary" size="lg" isFullWidth onClick={handleContinue}>
          Продолжить учиться
        </Button>
      </Container>
    </div>
  );
}
```

**Файл:** `src/pages/Villain/VictoryPage.module.css`

```css
.page {
  min-height: 100vh;
  background: linear-gradient(180deg, #ffd700 0%, #fff 50%);
  position: relative;
}

.content {
  padding: 40px 16px;
}

.header {
  text-align: center;
  margin-bottom: 32px;
}

.title {
  font-size: 36px;
  font-weight: 700;
  color: #000;
  margin: 0 0 8px 0;
}

.subtitle {
  font-size: 20px;
  color: #666;
  margin: 0;
}

.villainDefeated {
  display: flex;
  justify-content: center;
  margin-bottom: 32px;
}

.villainImage {
  width: 200px;
  height: 200px;
  object-fit: contain;
  filter: grayscale(100%);
  opacity: 0.6;
}

.stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 32px;
}

.statItem {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.statLabel {
  font-size: 14px;
  color: #666;
  margin-bottom: 4px;
}

.statValue {
  font-size: 28px;
  font-weight: 700;
  color: #000;
}

.rewards {
  margin-bottom: 32px;
}

.rewardsTitle {
  font-size: 20px;
  font-weight: 600;
  color: #000;
  margin: 0 0 16px 0;
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
```

---

### 3.5. VictoryAnimation Component

**Файл:** `src/pages/Villain/components/VictoryAnimation.tsx`

```typescript
import { createPortal } from 'react-dom';
import styles from './VictoryAnimation.module.css';

export function VictoryAnimation() {
  return createPortal(
    <div className={styles.container}>
      <div className={styles.confetti}>
        {Array.from({ length: 50 }).map((_, i) => (
          <span
            key={i}
            className={styles.confettiPiece}
            style={{
              left: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 3}s`,
              animationDuration: `${2 + Math.random() * 2}s`,
            }}
          >
            {['🎉', '✨', '⭐', '🎊', '🏆'][Math.floor(Math.random() * 5)]}
          </span>
        ))}
      </div>
    </div>,
    document.body
  );
}
```

**Файл:** `src/pages/Villain/components/VictoryAnimation.module.css`

```css
.container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 9999;
}

.confetti {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.confettiPiece {
  position: absolute;
  top: -50px;
  font-size: 24px;
  animation: fall linear infinite;
}

@keyframes fall {
  to {
    transform: translateY(100vh) rotate(360deg);
  }
}
```

---

## Часть 4: Аналитические события

### События Villain

| Event Name | Когда срабатывает | Параметры |
|------------|-------------------|-----------|
| `villain_screen_opened` | Открыт экран злодея | `child_profile_id`, `villain_id` |
| `villain_taunt_viewed` | Показана реплика | `child_profile_id`, `villain_id` |
| `villain_health_changed` | Изменилось здоровье (backend) | `child_profile_id`, `villain_id`, `health_before`, `health_after`, `damage_amount`, `reason` |
| `villain_victory_triggered` | Злодей побеждён (backend) | `child_profile_id`, `villain_id`, `attempt_id` |
| `victory_screen_opened` | Открыт экран победы | `child_profile_id`, `villain_id`, `attempt_id` |
| `victory_reward_viewed` | Просмотрена награда | `child_profile_id`, `villain_id`, `reward_type`, `reward_id` |
| `victory_continue_clicked` | Нажато продолжение | `child_profile_id`, `villain_id` |

---

## Чеклист задач

### Типы и API
- [ ] Создать типы villain.ts
- [ ] Реализовать villainAPI.getActiveVillain()
- [ ] Реализовать villainAPI.getBattleDetails()
- [ ] Реализовать villainAPI.getVictoryDetails()

### Компоненты Villain
- [ ] Создать VillainPage container
- [ ] Создать VillainCard с изображением и health bar
- [ ] Создать TauntBubble с репликами
- [ ] Создать BattleStats

### Компоненты Victory
- [ ] Создать VictoryPage container
- [ ] Создать VictoryAnimation (confetti)
- [ ] Создать RewardCard
- [ ] Добавить defeated villain image

### Игровая механика
- [ ] Реализовать урон от правильных ответов
- [ ] Реализовать изменение здоровья в реальном времени
- [ ] Добавить анимации атак
- [ ] Добавить звуковые эффекты (опционально)

### Аналитика
- [ ] Добавить villain_screen_opened
- [ ] Добавить villain_taunt_viewed
- [ ] Добавить victory_screen_opened
- [ ] Добавить victory_reward_viewed
- [ ] Добавить victory_continue_clicked

### Тестирование
- [ ] Протестировать экран злодея
- [ ] Протестировать изменение здоровья
- [ ] Протестировать экран победы
- [ ] Протестировать анимации

---

## Следующий этап

После завершения Villain переходи к **11_ANALYTICS.md** для интеграции аналитики.
