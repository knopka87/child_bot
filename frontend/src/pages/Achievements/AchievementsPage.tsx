// src/pages/Achievements/AchievementsPage.tsx
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import { BottomNav } from '@/components/layout/BottomNav';
import { AchievementsPageSkeleton } from '@/components/ui/skeleton';
import { AchievementDetailModal } from './components/AchievementDetailModal';
import { useAchievements } from './hooks/useAchievements';
import { useNewAchievements } from '@/hooks/useNewAchievements';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import type { Achievement } from '@/types/achievements';

export function AchievementsPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const { achievements, stats, isLoading, error } = useAchievements();
  const { hasNew, markAsViewed } = useNewAchievements();
  const [selectedAchievement, setSelectedAchievement] = useState<Achievement | null>(null);
  const [childProfileId, setChildProfileId] = useState<string | null>(null);

  // Загружаем child_profile_id из storage
  useEffect(() => {
    const loadProfileId = async () => {
      const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      setChildProfileId(profileId);
    };
    loadProfileId();
  }, []);

  useEffect(() => {
    if (childProfileId) {
      analytics.trackEvent('achievements_opened', {
        child_profile_id: childProfileId,
      });

      // Отмечаем что пользователь просмотрел страницу достижений
      if (hasNew) {
        markAsViewed();
      }
    }
  }, [analytics, childProfileId, hasNew, markAsViewed]);

  const handleAchievementClick = (achievement: Achievement) => {
    if (childProfileId) {
      analytics.trackEvent('achievement_clicked', {
        child_profile_id: childProfileId,
        achievement_id: achievement.id,
        is_unlocked: achievement.is_unlocked,
      });
    }

    setSelectedAchievement(achievement);
  };

  const handleCloseModal = () => {
    setSelectedAchievement(null);
  };

  // Бэкенд уже возвращает правильный порядок: активные первыми (новые в начале), затем неактивные по приоритету
  // Группируем достижения по "полкам" (по 4 в ряд)
  const achievementsPerShelf = 4;
  const shelves: Achievement[][] = [];
  for (let i = 0; i < achievements.length; i += achievementsPerShelf) {
    shelves.push(achievements.slice(i, i + achievementsPerShelf));
  }

  // Функция для определения цвета фона иконки
  const getIconBackground = (achievement: Achievement, index: number) => {
    if (!achievement.is_unlocked) {
      return 'bg-[#E0E0E0]'; // серый
    }

    // Цвета для разблокированных достижений
    const colors = [
      'bg-gradient-to-br from-[#FFB8B8] to-[#FF6B6B]', // розовый
      'bg-gradient-to-br from-[#A8E6CF] to-[#56C596]', // зеленый/бирюзовый
      'bg-gradient-to-br from-[#FFE17B] to-[#FFC75F]', // желтый
    ];

    return colors[index % colors.length];
  };

  if (isLoading) {
    return <AchievementsPageSkeleton />;
  }

  if (error) {
    return (
      <div className="flex flex-col min-h-screen bg-[#F5F3FF]">
        <div className="flex items-center justify-center flex-1">
          <p className="text-[#636e72]">Не удалось загрузить достижения</p>
        </div>
        <BottomNav />
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen bg-[#F5F3FF] pb-20">
      {/* Header */}
      <div className="px-4 pt-4 pb-2">
        <div className="flex items-center gap-3 mb-1">
          <button
            onClick={() => navigate(-1)}
            className="w-10 h-10 flex items-center justify-center rounded-full hover:bg-white/50 transition-colors"
          >
            <ArrowLeft size={24} className="text-[#2D3436]" />
          </button>
          <h1 className="text-[28px] font-bold text-[#2D3436]">Мои награды</h1>
        </div>
        {stats && (
          <p className="text-[15px] text-[#636e72] ml-[52px]">
            Собрано {stats.unlockedAchievements} из {stats.totalAchievements}
          </p>
        )}
      </div>

      {/* Shelves */}
      <div className="flex-1 px-4 mt-4">
        {shelves.map((shelfAchievements, shelfIndex) => {
          return (
            <div key={shelfIndex} className="mb-8">
              {/* Achievement Cards */}
              <div className="grid grid-cols-4 gap-3 mb-4">
                {shelfAchievements.map((achievement, index) => (
                  <button
                    key={achievement.id}
                    onClick={() => handleAchievementClick(achievement)}
                    className="flex flex-col items-center active:scale-95 transition-transform"
                  >
                    {/* Icon Circle with Badge */}
                    <div className="relative">
                      <div
                        className={`w-[70px] h-[70px] rounded-full flex items-center justify-center text-3xl mb-2 shadow-sm ${getIconBackground(
                          achievement,
                          index
                        )} ${!achievement.is_unlocked ? 'opacity-40' : ''}`}
                      >
                        {achievement.icon}
                      </div>
                      {/* Badge for Serial Stickers (Дружба, Стрик, Проверки ДЗ, Злодеи, Исправленные ошибки) */}
                      {achievement.reward.type === 'sticker' &&
                        ['Дружба', 'Стрик', 'Проверки ДЗ', 'Победитель злодеев', 'Исправленные ошибки'].includes(achievement.reward.name) &&
                        achievement.progress.total && (
                          <div className="absolute -top-1 -right-1 w-6 h-6 bg-[#FF6B6B] text-white text-xs font-bold rounded-full flex items-center justify-center shadow-md">
                            {achievement.progress.total}
                          </div>
                        )}
                    </div>
                    {/* Title */}
                    <p
                      className={`text-[11px] text-center leading-tight ${
                        achievement.is_unlocked ? 'text-[#2D3436]' : 'text-[#B2BEC3]'
                      }`}
                    >
                      {achievement.title}
                    </p>
                  </button>
                ))}
              </div>
            </div>
          );
        })}
      </div>

      {/* Detail Modal */}
      {selectedAchievement && (
        <AchievementDetailModal
          achievement={selectedAchievement}
          isOpen={!!selectedAchievement}
          onClose={handleCloseModal}
        />
      )}

      <BottomNav hasNewAchievements={hasNew} />
    </div>
  );
}

export default AchievementsPage;
