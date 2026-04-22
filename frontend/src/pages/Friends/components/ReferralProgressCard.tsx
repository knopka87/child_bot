// src/pages/Friends/components/ReferralProgressCard.tsx
import { Card } from '@/components/ui/Card';
import { ProgressBar } from '@/components/ui/ProgressBar';
import type { ReferralGoal } from '@/types/referral';
import styles from './ReferralProgressCard.module.css';

interface ReferralProgressCardProps {
  currentGoal: ReferralGoal;
  invitedCount: number;
  progressPercent: number;
}

export function ReferralProgressCard({
  currentGoal,
  invitedCount,
  progressPercent,
}: ReferralProgressCardProps) {
  return (
    <Card className={styles.card}>
      <div className={styles.header}>
        <h3 className={styles.title}>Текущая цель</h3>
        <div className={styles.goalBadge}>
          {currentGoal.isCompleted ? '✓ Выполнено' : 'В процессе'}
        </div>
      </div>

      <p className={styles.description}>
        Пригласи {currentGoal.targetCount} друзей
      </p>

      <div className={styles.progressIndicator}>
        {Array.from({ length: currentGoal.targetCount }).map((_, index) => (
          <span
            key={index}
            className={`${styles.dot} ${index < invitedCount ? styles.filled : ''}`}
          >
            {index < invitedCount ? '✓' : index + 1}
          </span>
        ))}
      </div>

      <div className={styles.progressText}>
        {invitedCount} из {currentGoal.targetCount}
      </div>

      <ProgressBar
        value={progressPercent}
        size="md"
        variant={currentGoal.isCompleted ? 'success' : 'default'}
      />

      <div className={styles.reward}>
        <div className={styles.rewardLabel}>Награда:</div>
        <div className={styles.rewardContent}>
          {currentGoal.reward.type === 'sticker' && currentGoal.reward.imageUrl && (
            <img
              src={currentGoal.reward.imageUrl}
              alt={currentGoal.reward.name}
              className={styles.rewardImage}
            />
          )}
          {currentGoal.reward.type === 'coins' && (
            <span className={styles.rewardEmoji}>🪙</span>
          )}
          <div className={styles.rewardInfo}>
            <div className={styles.rewardName}>{currentGoal.reward.name}</div>
            <div className={styles.rewardDescription}>
              {currentGoal.reward.description}
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
}
