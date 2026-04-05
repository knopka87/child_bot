// src/pages/Achievements/AchievementsPage.tsx
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import { BottomNav } from '@/components/layout/BottomNav';
import { Spinner } from '@/components/ui/Spinner';
import { AchievementDetailModal } from './components/AchievementDetailModal';
import { useAchievements } from './hooks/useAchievements';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import type { Achievement } from '@/types/achievements';

export function AchievementsPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const { achievements, stats, isLoading, error } = useAchievements();
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
    }
  }, [analytics, childProfileId]);

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

  // Группируем достижения по полкам
  const achievementsByShelf = achievements.reduce(
    (acc, achievement) => {
      const shelf = achievement.shelf_order;
      if (!acc[shelf]) {
        acc[shelf] = [];
      }
      acc[shelf].push(achievement);
      return acc;
    },
    {} as Record<number, Achievement[]>
  );

  // Сортируем полки
  const shelves = Object.keys(achievementsByShelf)
    .map(Number)
    .sort((a, b) => a - b);

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
    return (
      <div className="flex flex-col min-h-screen bg-[#F5F3FF]">
        <div className="flex items-center justify-center flex-1">
          <Spinner size="lg" />
        </div>
        <BottomNav />
      </div>
    );
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
        {shelves.map((shelfNum, shelfIndex) => {
          const shelfAchievements = achievementsByShelf[shelfNum].sort(
            (a, b) => a.position_in_shelf - b.position_in_shelf
          );

          return (
            <div key={shelfNum} className="mb-8">
              {/* Achievement Cards */}
              <div className="grid grid-cols-4 gap-3 mb-4">
                {shelfAchievements.map((achievement, index) => (
                  <button
                    key={achievement.id}
                    onClick={() => handleAchievementClick(achievement)}
                    className="flex flex-col items-center active:scale-95 transition-transform"
                  >
                    {/* Icon Circle */}
                    <div
                      className={`w-[70px] h-[70px] rounded-full flex items-center justify-center text-3xl mb-2 shadow-sm ${getIconBackground(
                        achievement,
                        index
                      )} ${!achievement.is_unlocked ? 'opacity-40' : ''}`}
                    >
                      {achievement.icon}
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

              {/* Shelf Separator (except last) */}
              {shelfIndex < shelves.length - 1 && (
                <div className="h-[2px] bg-gradient-to-r from-[#C9A969] via-[#E8D5A3] to-[#C9A969] rounded-full shadow-md" />
              )}
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

      <BottomNav />
    </div>
  );
}

export default AchievementsPage;
