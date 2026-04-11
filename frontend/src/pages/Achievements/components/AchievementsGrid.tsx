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
  // Сортируем по приоритету и разбиваем на группы по 4
  const sortedAchievements = [...achievements].sort((a, b) => a.priority - b.priority);

  const achievementsPerShelf = 4;
  const shelves: Achievement[][] = [];
  for (let i = 0; i < sortedAchievements.length; i += achievementsPerShelf) {
    shelves.push(sortedAchievements.slice(i, i + achievementsPerShelf));
  }

  return (
    <div className={styles.container}>
      {shelves.map((shelfItems, shelfIndex) => (
        <div key={shelfIndex} className={styles.shelf}>
          <div className={styles.grid}>
            {shelfItems.map((achievement) => (
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
