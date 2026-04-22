// src/pages/Villain/components/RewardCard.tsx
import { Card } from '@/components/ui/Card';
import type { VillainReward } from '@/types/villain';
import styles from './RewardCard.module.css';

interface RewardCardProps {
  reward: VillainReward;
}

export function RewardCard({ reward }: RewardCardProps) {
  const getRewardTypeLabel = () => {
    switch (reward.type) {
      case 'sticker':
        return '🎨 Стикер';
      case 'achievement':
        return '🏆 Достижение';
      case 'coins':
        return '💰 Монеты';
      case 'avatar':
        return '👤 Аватар';
      default:
        return reward.type;
    }
  };

  return (
    <Card className={styles.card} variant="bordered">
      <div className={styles.content}>
        {reward.image_url && (
          <div className={styles.imageWrapper}>
            <img
              src={reward.image_url}
              alt={reward.name}
              className={styles.rewardImage}
            />
          </div>
        )}

        <div className={styles.info}>
          <div className={styles.type}>{getRewardTypeLabel()}</div>
          <h4 className={styles.name}>{reward.name}</h4>

          {reward.amount !== undefined && (
            <div className={styles.amount}>+{reward.amount}</div>
          )}
        </div>
      </div>
    </Card>
  );
}
