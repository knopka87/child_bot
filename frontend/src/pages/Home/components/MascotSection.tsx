// src/pages/Home/components/MascotSection.tsx
import { Card } from '@/components/ui/Card';
import { ProgressBar } from '@/components/ui/ProgressBar';
import type { Villain } from '@/types/domain';
import styles from './MascotSection.module.css';

interface MascotData {
  id: string;
  state: 'idle' | 'happy' | 'thinking' | 'celebrating' | 'encouraging';
  imageUrl: string;
  message: string;
}

interface MascotSectionProps {
  mascot: MascotData;
  villain: Villain | null;
  onMascotClick: () => void;
  onVillainClick: () => void;
}

export function MascotSection({
  mascot,
  villain,
  onMascotClick,
  onVillainClick,
}: MascotSectionProps) {
  return (
    <div className={styles.section}>
      {/* Mascot */}
      <div className={styles.mascotWrapper} onClick={onMascotClick}>
        <img
          src={mascot.imageUrl}
          alt="Mascot"
          className={styles.mascotImage}
        />
        <div className={styles.speechBubble}>
          <p className={styles.message}>{mascot.message}</p>
        </div>
      </div>

      {/* Battle indicator */}
      {villain && (
        <div className={styles.battleIndicator}>
          <span className={styles.swordIcon}>⚔️</span>
        </div>
      )}

      {/* Villain */}
      {villain && (
        <Card
          className={styles.villainCard}
          variant="bordered"
          onClick={onVillainClick}
        >
          <img
            src={villain.imageUrl}
            alt={villain.name}
            className={styles.villainImage}
          />
          <div className={styles.villainInfo}>
            <h3 className={styles.villainName}>{villain.name}</h3>
            <ProgressBar
              value={villain.healthPercent}
              variant="error"
              size="sm"
              showLabel
              label={`${villain.healthPercent}%`}
            />
          </div>
        </Card>
      )}
    </div>
  );
}
