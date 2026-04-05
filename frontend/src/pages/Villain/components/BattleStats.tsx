// src/pages/Villain/components/BattleStats.tsx
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
            <span className={styles.statValue}>{battle.battle_stats.total_damage_dealt}</span>
          </div>
        </div>

        <div className={styles.stat}>
          <span className={styles.statIcon}>🎯</span>
          <div className={styles.statInfo}>
            <span className={styles.statLabel}>Правильных ответов</span>
            <span className={styles.statValue}>{battle.battle_stats.correct_tasks_count}</span>
          </div>
        </div>
      </div>
    </Card>
  );
}
