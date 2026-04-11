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

        {/* Информация о следующем уровне для серийных наград */}
        {achievement.next_level && (
          <div className="mt-4 p-3 bg-[#F0F4FF] rounded-xl">
            <p className="text-[12px] text-[#636e72] mb-1">Следующая награда:</p>
            <p className="text-[14px] text-[#2D3436] font-medium">
              {achievement.next_level.description}
            </p>
            <div className="mt-2 flex items-center gap-2">
              <div className="flex-1 h-2 bg-white rounded-full overflow-hidden">
                <div
                  className="h-full bg-[#6C5CE7] rounded-full transition-all"
                  style={{
                    width: `${Math.min(100, (achievement.progress.current / achievement.next_level.requirement_value) * 100)}%`,
                  }}
                />
              </div>
              <span className="text-[12px] text-[#636e72] font-medium min-w-[60px] text-right">
                {achievement.progress.current} / {achievement.next_level.requirement_value}
              </span>
            </div>
          </div>
        )}

        {/* Кнопка */}
        <PurpleButton onClick={onClose} style={{ marginTop: '24px' }}>
          Понятно!
        </PurpleButton>
      </div>
    </CustomModal>
  );
}
