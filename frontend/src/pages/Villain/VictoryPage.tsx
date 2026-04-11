// src/pages/Villain/VictoryPage.tsx
import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Spinner } from '@/components/ui/Spinner';
import { useVictoryData } from './hooks/useVictoryData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { ROUTES } from '@/config/routes';

export function VictoryPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const analytics = useAnalytics();
  const { villainId } = (location.state as { villainId?: string }) || {};

  // Validate villainId
  if (!villainId) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen px-5 bg-gradient-to-b from-[#FFF9E8] to-[#E8FFF8]">
        <p className="text-[#636e72] mb-4">Неверный ID злодея</p>
        <button
          onClick={() => navigate(ROUTES.HOME)}
          className="px-6 py-3 bg-[#6C5CE7] text-white rounded-2xl active:scale-[0.98] transition-transform"
        >
          На главную
        </button>
      </div>
    );
  }

  const { victory, isLoading, error } = useVictoryData(villainId);

  useEffect(() => {
    const trackVictory = async () => {
      if (victory && villainId) {
        const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);

        analytics.trackEvent('victory_screen_opened', {
          child_profile_id: childProfileId,
          villain_id: villainId,
          attempt_id: '',
        });

        victory.rewards.forEach((reward) => {
          analytics.trackEvent('victory_reward_viewed', {
            child_profile_id: childProfileId,
            villain_id: villainId,
            reward_type: reward.type,
            reward_id: reward.id,
          });
        });
      }
    };

    trackVictory();
  }, [victory, villainId, analytics]);

  const handleContinue = async () => {
    if (villainId) {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      analytics.trackEvent('victory_continue_clicked', {
        child_profile_id: childProfileId,
        villain_id: villainId,
      });
    }

    navigate(ROUTES.HOME);
  };

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-b from-[#FFF9E8] to-[#E8FFF8]">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !victory) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen px-5 bg-gradient-to-b from-[#FFF9E8] to-[#E8FFF8]">
        <p className="text-[#636e72] mb-4">Не удалось загрузить данные победы</p>
        <button
          onClick={() => navigate(ROUTES.HOME)}
          className="px-6 py-3 bg-[#6C5CE7] text-white rounded-2xl active:scale-[0.98] transition-transform"
        >
          На главную
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#FFF9E8] to-[#E8FFF8]">
      <div className="flex-1 flex flex-col items-center justify-center text-center">
        {/* Анимация победы */}
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: 'spring', duration: 0.6 }}
          className="text-[80px] mb-2"
        >
          🎉
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <h1 className="text-[32px] font-bold text-[#6C5CE7] mb-2">Победа!</h1>
          <p className="text-[#636e72] text-[14px] mb-6">
            Ты победил злодея!
          </p>
        </motion.div>

        {/* Персонажи */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="flex items-end gap-6 mb-6"
        >
          <div className="text-center">
            <div className="text-6xl mb-2">🦉</div>
            <p className="text-[12px] text-[#636e72]">Мы победили!</p>
          </div>
          <div className="text-center">
            <div className="text-6xl mb-2 opacity-40">👾</div>
            <p className="text-[12px] text-[#636e72]">Побеждён</p>
          </div>
        </motion.div>

        {/* Статистика */}
        {(victory.total_damage > 0 || victory.tasks_completed > 0) && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.6 }}
            className="bg-white rounded-2xl p-4 w-full mb-6 shadow-sm"
          >
            <div className="grid grid-cols-2 gap-4">
              {victory.total_damage > 0 && (
                <div className="text-center">
                  <p className="text-[24px] font-bold text-[#2D3436]">
                    {victory.total_damage}
                  </p>
                  <p className="text-[12px] text-[#636e72]">Урон нанесён</p>
                </div>
              )}
              {victory.tasks_completed > 0 && (
                <div className="text-center">
                  <p className="text-[24px] font-bold text-[#2D3436]">
                    {victory.tasks_completed}
                  </p>
                  <p className="text-[12px] text-[#636e72]">Задач выполнено</p>
                </div>
              )}
            </div>
          </motion.div>
        )}

        {/* Награды */}
        {victory.rewards && victory.rewards.length > 0 && (
          <motion.div
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.7 }}
            className="bg-white rounded-3xl p-6 shadow-sm w-full max-w-[280px]"
          >
            <p className="text-[13px] text-[#636e72] mb-3">Твоя награда</p>
            <div className="flex items-center justify-center gap-4">
              {victory.rewards.map((reward, index) => (
                <div key={reward.id || index} className="text-center">
                  <div className="text-[40px]">
                    {reward.type === 'sticker' && '⭐'}
                    {reward.type === 'achievement' && '🏆'}
                    {reward.type === 'coins' && '🪙'}
                    {reward.type === 'avatar' && '👤'}
                  </div>
                  <p className="text-[11px] text-[#2D3436] mt-1">
                    {reward.type === 'sticker' && 'Редкий стикер'}
                    {reward.type === 'achievement' && 'Достижение'}
                    {reward.type === 'coins' && `${reward.amount} монет`}
                    {reward.type === 'avatar' && 'Аватар'}
                  </p>
                </div>
              ))}
            </div>
          </motion.div>
        )}
      </div>

      <button
        onClick={handleContinue}
        className="w-full py-4 bg-[#6C5CE7] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#6C5CE7]/20 active:scale-[0.98] transition-transform"
      >
        Продолжить учиться
      </button>
    </div>
  );
}
