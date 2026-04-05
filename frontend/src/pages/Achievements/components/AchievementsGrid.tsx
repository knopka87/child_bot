// src/pages/Achievements/components/AchievementsGrid.tsx
import { AchievementCard } from './AchievementCard';
import type { Achievement } from '@/types/achievements';
import styles from './AchievementsGrid.module.css';

interface AchievementsGridProps {
  achievements: Achievement[];
  onAchievementClick: (achievement: Achievement) => void;
  // onShelfView?: (shelfOrder: number) => void; // TODO: track shelf views
}

export function AchievementsGrid({
  achievements,
  onAchievementClick,
}: AchievementsGridProps) {
  // Group by shelf
  const shelves = [1, 2, 3];
  const achievementsByShelf = shelves.map((shelfOrder) => ({
    shelfOrder,
    items: achievements.filter((a) => a.shelf_order === shelfOrder),
  }));

  return (
    <div className={styles.container}>
      {achievementsByShelf.map(({ shelfOrder, items }) => (
        <div key={shelfOrder} className={styles.shelf}>
          <div className={styles.shelfLabel}>Полка {shelfOrder}</div>
          <div className={styles.grid}>
            {items.map((achievement) => (
              <AchievementCard
                key={achievement.id}
                achievement={achievement}
                onClick={onAchievementClick}
              />
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
