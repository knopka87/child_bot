// src/pages/Villain/components/VillainCard.tsx
import { Card } from '@/components/ui/Card';
import { ProgressBar } from '@/components/ui/ProgressBar';
import type { Villain } from '@/types/villain';
import styles from './VillainCard.module.css';

interface VillainCardProps {
  villain: Villain;
}

export function VillainCard({ villain }: VillainCardProps) {
  const healthPercent = (villain.hp / villain.max_hp) * 100;

  return (
    <Card className={styles.card} variant="default">
      <div className={styles.imageWrapper}>
        <img
          src={villain.image_url}
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
              {villain.hp} / {villain.max_hp}
            </span>
          </div>
          <ProgressBar
            value={healthPercent}
            variant="error"
            size="lg"
            showLabel={false}
          />
        </div>
      </div>

      {villain.taunt && (
        <div className={styles.tauntBubble}>
          <p className={styles.taunt}>{villain.taunt}</p>
          <div className={styles.tauntArrow} />
        </div>
      )}
    </Card>
  );
}
