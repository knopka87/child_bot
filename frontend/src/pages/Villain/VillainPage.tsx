// src/pages/Villain/VillainPage.tsx
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft } from 'lucide-react';
import { motion } from 'framer-motion';
import { useVillainData } from './hooks/useVillainData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';

export function VillainPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const { villain, isLoading, error } = useVillainData();

  useEffect(() => {
    if (villain) {
      analytics.trackEvent('villain_screen_opened', {
        child_profile_id: profile?.child_profile_id,
        villain_id: villain.id,
      });

      analytics.trackEvent('villain_taunt_viewed', {
        child_profile_id: profile?.child_profile_id,
        villain_id: villain.id,
      });
    }
  }, [villain?.id]);

  useEffect(() => {
    if (villain?.is_defeated) {
      navigate(ROUTES.VILLAIN_VICTORY, {
        state: { villainId: villain.id },
      });
    }
  }, [villain?.is_defeated, navigate]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Загрузка...</p>
        </div>
      </div>
    );
  }

  if (error || !villain) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen px-5">
        <p className="text-gray-600 mb-4">Нет активного злодея</p>
        <button
          onClick={() => navigate(ROUTES.HOME)}
          className="px-6 py-3 bg-blue-600 text-white rounded-2xl"
        >
          На главную
        </button>
      </div>
    );
  }

  // Вычисляем количество заполненных сегментов (1-3)
  const healthSegments = Math.ceil((villain.hp / villain.max_hp) * 3);

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
          className="text-[100px] mb-4"
        >
          👾
        </motion.div>

        <h2 className="text-foreground mb-1 text-[24px] font-semibold">
          {villain.name || 'Злодей Кракозябра'}
        </h2>
        <p className="text-muted-foreground text-[14px] mb-6">
          {villain.taunt || '«Ха-ха! Попробуй-ка реши задачки!»'}
        </p>

        {/* Health bar */}
        <div className="w-full max-w-[260px] mb-6">
          <p className="text-[12px] text-muted-foreground mb-2">Здоровье злодея</p>
          <div className="flex gap-2">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className={`flex-1 h-4 rounded-full transition-all ${
                  i <= healthSegments ? 'bg-red-500' : 'bg-gray-200'
                }`}
              />
            ))}
          </div>
        </div>

        {/* Explanation */}
        <div className="bg-white rounded-2xl p-5 shadow-sm w-full mb-6">
          <h3 className="text-foreground mb-2 text-[18px] font-semibold">Как победить?</h3>
          <div className="flex flex-col gap-2 text-left text-[13px] text-muted-foreground">
            <p>• Решай задания правильно — каждый верный ответ снимает здоровье злодея</p>
            <p>• После 3 верных ответов злодей побеждён!</p>
            <p>• За победу ты получишь редкий стикер и достижение</p>
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
