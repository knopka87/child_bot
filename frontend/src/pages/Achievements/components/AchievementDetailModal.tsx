// src/pages/Achievements/components/AchievementDetailModal.tsx
import { CustomModal } from '@/components/ui/Modal/CustomModal';
import { PurpleButton } from '@/components/ui/Button';
import type { Achievement } from '@/types/achievements';
import styles from './AchievementDetailModal.module.css';

interface AchievementDetailModalProps {
  achievement: Achievement | null;
  isOpen: boolean;
  onClose: () => void;
}

export function AchievementDetailModal({
  achievement,
  isOpen,
  onClose,
}: AchievementDetailModalProps) {
  if (!achievement) return null;

  return (
    <CustomModal isOpen={isOpen} onClose={onClose}>
      <div className={styles.content}>
        {/* Иконка на розовом фоне */}
        <div className={styles.iconContainer}>
          <div className={styles.iconBackground}>
            <span className={styles.iconEmoji}>{achievement.icon}</span>
          </div>
        </div>

        {/* Заголовок */}
        <h2 className={styles.title}>{achievement.title}</h2>

        {/* Описание */}
        <p className={styles.description}>{achievement.description}</p>

        {/* Статус "Получено" если разблокировано */}
        {achievement.is_unlocked && (
          <div className={styles.statusBadge}>✓ Получено</div>
        )}

        {/* Кнопка */}
        <PurpleButton onClick={onClose} style={{ marginTop: '24px' }}>
          Понятно!
        </PurpleButton>
      </div>
    </CustomModal>
  );
}
