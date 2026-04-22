// src/pages/Villain/components/VictoryAnimation.tsx
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
