// src/pages/Check/ScenarioSelection.tsx
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Image, Images } from 'lucide-react';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import type { CheckScenario } from '@/types/check';

const scenarios = [
  {
    id: 'single_photo' as CheckScenario,
    icon: Image,
    title: '1 фото',
    desc: 'Задание и ответ на одном фото',
    color: 'from-[#6C5CE7] to-[#A29BFE]', // фиолетовый градиент
  },
  {
    id: 'two_photo' as CheckScenario,
    icon: Images,
    title: '2 фото',
    desc: 'Задание отдельно, ответ отдельно',
    color: 'from-[#00B894] to-[#55EFC4]', // зелёный градиент
  },
];

export default function ScenarioSelection() {
  const navigate = useNavigate();
  const analytics = useAnalytics();

  useEffect(() => {
    const trackOpen = async () => {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      analytics.trackEvent('check_scenario_selection_opened', {
        child_profile_id: childProfileId,
      });
    };
    trackOpen();
  }, [analytics]);

  const handleScenarioSelect = async (scenario: CheckScenario) => {
    const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
    analytics.trackEvent('check_scenario_selected', {
      child_profile_id: childProfileId,
      scenario_type: scenario,
    });

    // Маппинг внутренних сценариев на URL-параметры
    const scenarioParam = scenario === 'single_photo' ? 'single_photo' : 'two_photo';
    navigate(`/check/upload-images?scenario=${scenarioParam}`);
  };

  return (
    <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-white">
      {/* Кнопка назад */}
      <button
        onClick={() => navigate('/')}
        className="flex items-center gap-2 text-[#6C5CE7] mb-6 active:opacity-70 transition-opacity"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px] font-medium">Назад</span>
      </button>

      {/* Заголовок */}
      <div className="text-center mb-8">
        <h2 className="text-[28px] font-bold text-[#2D3436] mb-1">Проверка ДЗ</h2>
        <p className="text-[#636e72] text-[14px]">
          Выбери, как выглядит твоё задание
        </p>
      </div>

      {/* Карточки сценариев */}
      <div className="flex flex-col gap-4 flex-1">
        {scenarios.map((scenario, index) => {
          const IconComponent = scenario.icon;
          return (
            <button
              key={scenario.id}
              onClick={() => handleScenarioSelect(scenario.id)}
              className={`bg-gradient-to-r ${scenario.color} text-white rounded-3xl p-6 flex items-center gap-4 shadow-lg active:scale-[0.98] transition-transform flex-1`}
              style={{
                animation: `fadeInUp 0.3s ease-out ${index * 0.1}s both`,
              }}
            >
              {/* Иконка */}
              <div className="w-16 h-16 bg-white/20 rounded-2xl flex items-center justify-center flex-shrink-0">
                <IconComponent size={32} />
              </div>

              {/* Текст */}
              <div className="text-left">
                <h3 className="text-white text-[20px] font-bold">{scenario.title}</h3>
                <p className="text-white/80 text-[13px] mt-0.5">{scenario.desc}</p>
              </div>
            </button>
          );
        })}
      </div>

      <style>{`
        @keyframes fadeInUp {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
      `}</style>
    </div>
  );
}
