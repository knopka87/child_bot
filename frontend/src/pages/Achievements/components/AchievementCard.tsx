// src/pages/Achievements/components/AchievementCard.tsx
import { Card } from '@/components/ui/Card';
import { ProgressBar } from '@/components/ui/ProgressBar';
import type { Achievement } from '@/types/achievements';
import styles from './AchievementCard.module.css';

interface AchievementCardProps {
  achievement: Achievement;
  onClick: (achievement: Achievement) => void;
}

export function AchievementCard({ achievement, onClick }: AchievementCardProps) {
  const isLocked = !achievement.is_unlocked;
  const isPartial = !isLocked && achievement.progress.percent < 100;

  return (
    <Card
      className={`${styles.card} ${isLocked ? styles.locked : ''}`}
      onClick={() => onClick(achievement)}
    >
      <div className={styles.iconWrapper}>
        {typeof achievement.icon === 'string' && achievement.icon.startsWith('http') ? (
          <img src={achievement.icon} alt={achievement.title} className={styles.icon} />
        ) : (
          <span className={styles.emoji}>{achievement.icon}</span>
        )}
        {isLocked && <div className={styles.lockedOverlay}>🔒</div>}
      </div>

      <div className={styles.title}>{achievement.title}</div>

      {isPartial && (
        <ProgressBar
          value={achievement.progress.percent}
          size="sm"
          variant="default"
        />
      )}
    </Card>
  );
}
