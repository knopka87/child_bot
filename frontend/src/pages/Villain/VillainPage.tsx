// src/pages/Villain/VillainPage.tsx
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import { motion } from 'framer-motion';
import { useVillainData } from './hooks/useVillainData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { ROUTES } from '@/config/routes';
import { VillainPageSkeleton } from '@/components/ui/skeleton';

export function VillainPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const { villain, battle, isLoading, error } = useVillainData();
  const [childProfileId, setChildProfileId] = useState<string | null>(null);
  const [imageFailed, setImageFailed] = useState(false);

  // Вычисляем количество задач для победы (только если villain загружен)
  const tasksNeeded = villain 
    ? (battle?.battle_stats?.damage_per_task 
      ? Math.ceil(villain.max_hp / battle.battle_stats.damage_per_task)
      : Math.ceil(villain.max_hp / 20)) // Fallback на 20 если данных нет
    : 5; // Default пока загружается

  useEffect(() => {
    vkStorage.getItem(storageKeys.PROFILE_ID).then(setChildProfileId);
  }, []);

  useEffect(() => {
    if (villain && childProfileId) {
      analytics.trackEvent('villain_screen_opened', {
        child_profile_id: childProfileId,
        villain_id: villain.id,
      });

      analytics.trackEvent('villain_taunt_viewed', {
        child_profile_id: childProfileId,
        villain_id: villain.id,
      });
    }
  }, [villain?.id, childProfileId]);

  useEffect(() => {
    if (villain?.is_defeated) {
      navigate(ROUTES.VILLAIN_VICTORY, {
        state: { villainId: villain.id },
      });
    }
  }, [villain?.is_defeated, navigate]);

  if (isLoading) {
    return <VillainPageSkeleton />;
  }

  if (error || !villain) {
    // Если злодея нет - значит он уже побеждён сегодня
    return (
      <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#E8FFF8] to-[#F0F4FF]">
        <button
          onClick={() => navigate(ROUTES.HOME)}
          className="flex items-center gap-2 text-primary mb-6"
        >
          <ArrowLeft size={20} />
          <span className="text-[14px]">На главную</span>
        </button>

        <div className="flex-1 flex flex-col items-center justify-center text-center">
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: 'spring', duration: 0.6 }}
            className="text-[80px] mb-4"
          >
            🎉
          </motion.div>

          <h1 className="text-[#00B894] mb-2 text-[32px] font-bold">Злодей побеждён!</h1>
          <p className="text-[#636e72] text-[16px] mb-6">
            Ты молодец! Следующий злодей появится завтра.
          </p>

          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3 }}
            className="bg-white rounded-2xl p-6 shadow-sm w-full max-w-[320px] mb-6"
          >
            <div className="text-6xl mb-4">🦉</div>
            <p className="text-[#2D3436] text-[16px] font-medium mb-2">Отличная работа!</p>
            <p className="text-[#636e72] text-[14px]">
              Возвращайся завтра за новым злодеем и новыми наградами!
            </p>
          </motion.div>
        </div>

        <button
          onClick={() => navigate(ROUTES.HOME)}
          className="w-full py-4 bg-[#00B894] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#00B894]/20 active:scale-[0.98] transition-transform"
        >
          На главную
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-[#FFE8E8]">
      <button
        onClick={() => navigate(ROUTES.HOME)}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Назад</span>
      </button>

      <div className="flex-1 flex flex-col items-center justify-center text-center">
        <motion.div
          animate={{ y: [0, -8, 0] }}
          transition={{ repeat: Infinity, duration: 2, ease: 'easeInOut' }}
          className="mb-4 w-[150px] h-[150px] flex items-center justify-center"
        >
          {!imageFailed && villain.image_url ? (
            <img 
              src={villain.image_url} 
              alt={villain.name} 
              className="w-full h-full object-contain"
              onError={() => setImageFailed(true)}
            />
          ) : (
            <span className="text-[100px]">👾</span>
          )}
        </motion.div>

        <h2 className="text-foreground mb-1 text-[24px] font-semibold">
          {villain.name || 'Злодей Кракозябра'}
        </h2>
        <p className="text-muted-foreground text-[14px] mb-6">
          {villain.taunt || '«Ха-ха! Попробуй-ка реши задачки!»'}
        </p>

        {/* Health bar - единая прогресс-бар полоска */}
        <div className="w-full max-w-[260px] mb-6">
          <p className="text-[12px] text-muted-foreground mb-2">
            Здоровье злодея: {villain.hp} / {villain.max_hp} HP
          </p>
          <div className="w-full h-4 bg-gray-200 rounded-full overflow-hidden">
            <div 
              className="h-full bg-red-500 transition-all duration-500 ease-out rounded-full"
              style={{ width: `${Math.max(0, (villain.hp / villain.max_hp) * 100)}%` }}
            />
          </div>
          <p className="text-[11px] text-muted-foreground mt-1">
            {Math.round((villain.hp / villain.max_hp) * 100)}% здоровья осталось
          </p>
        </div>

        {/* Info about daily villain */}
        <div className="bg-[#FFF9E8] rounded-2xl p-4 shadow-sm w-full mb-6">
          <p className="text-[13px] text-[#2D3436]">
            📅 Новый злодей появляется каждый день!
          </p>
          <p className="text-[12px] text-[#636e72] mt-1">
            Здоровье восстанавливается в полночь. Успей победить!
          </p>
        </div>

        {/* Explanation */}
        <div className="bg-white rounded-2xl p-5 shadow-sm w-full mb-6">
          <h3 className="text-foreground mb-2 text-[18px] font-semibold">Как победить?</h3>
          <div className="flex flex-col gap-2 text-left text-[13px] text-muted-foreground">
            <p>• Решай задания правильно — каждый верный ответ снимает {battle?.battle_stats?.damage_per_task || 20} HP</p>
            <p>• Реши {tasksNeeded} {tasksNeeded === 1 ? 'задачу' : tasksNeeded < 5 ? 'задачи' : 'заданий'} правильно, чтобы победить!</p>
            <p>• За победу ты можешь получишь редкий стикер и достижение</p>
          </div>
        </div>
      </div>

      <button
        onClick={() => navigate(ROUTES.HELP)}
        className="w-full py-4 bg-blue-600 text-white rounded-2xl shadow-lg shadow-blue-600/20 active:scale-[0.98] transition-transform font-semibold"
      >
        Продолжить учиться
      </button>
    </div>
  );
}
