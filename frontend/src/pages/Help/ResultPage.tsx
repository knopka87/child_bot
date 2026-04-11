// src/pages/Help/ResultPage.tsx
import { useEffect, useState } from 'react';
import { useLocation, useNavigate, useParams } from 'react-router-dom';
import { ArrowLeft, ChevronDown, ChevronUp, Lock } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';
import type { HelpResult } from '@/types/help';

export default function ResultPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const params = useParams();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  // Получаем attemptId из URL параметров или из state (для обратной совместимости)
  const attemptId = params.attemptId || (location.state?.attemptId as string);

  const [resultData, setResultData] = useState<HelpResult | null>(location.state?.result || null);
  const [isLoading, setIsLoading] = useState(
    !location.state?.result || (location.state?.result && location.state.result.hints.length === 0)
  );

  const [unlockedLevel, setUnlockedLevel] = useState(1);
  const [openHints, setOpenHints] = useState<string[]>([]);
  const [villainHealth] = useState(3);
  const [hasInitialized, setHasInitialized] = useState(false);

  useEffect(() => {
    if (!attemptId) {
      navigate(ROUTES.HELP);
      return;
    }

    // Загружаем данные если:
    // 1. Их вообще нет в state
    // 2. Или hints пустой (ProcessingPage получил данные до завершения обработки)
    if (!resultData || (resultData && resultData.hints.length === 0)) {
      console.log('[ResultPage] Need to load data, resultData:', resultData);
      loadResultData();
    }
  }, [attemptId]);

  useEffect(() => {
    if (resultData && !hasInitialized) {
      console.log('[ResultPage] Initializing hints...');
      setHasInitialized(true);

      analytics.trackEvent('help_result_opened', {
        child_profile_id: profile?.child_profile_id,
        attempt_id: attemptId,
        hints_count: resultData.hints.length,
      });

      // Автоматически открываем первую подсказку и вызываем API для инкремента
      if (resultData.hints.length > 0) {
        setOpenHints([resultData.hints[0].id]);

        // Вызываем API для первой подсказки (level 0 -> level 1)
        helpAPI
          .getNextHint(attemptId, 0)
          .then(() => {
            console.log('[ResultPage] First hint API call successful');
          })
          .catch((error) => {
            console.error('[ResultPage] Failed to call API for first hint:', error);
          });
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [resultData, hasInitialized]);

  const loadResultData = async () => {
    try {
      console.log('[ResultPage] Loading result for attempt:', attemptId);
      const result = await helpAPI.getResult(attemptId);
      console.log('[ResultPage] Result loaded:', result);
      setResultData(result);
      setIsLoading(false);
    } catch (error) {
      console.error('[ResultPage] Failed to load result:', error);
      navigate(ROUTES.HELP);
    }
  };

  // Если данных нет и идёт загрузка, показываем спиннер
  if (!resultData && isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Загрузка результатов...</p>
        </div>
      </div>
    );
  }

  // Если данных нет после загрузки, редиректим
  if (!attemptId || !resultData) {
    return null;
  }

  const toggleHint = (hintId: string, level: number) => {
    console.log('[ResultPage] toggleHint called:', { hintId, level, unlockedLevel });

    if (level > unlockedLevel) {
      console.log('[ResultPage] Hint is locked, cannot toggle');
      return;
    }

    setOpenHints((prev) =>
      prev.includes(hintId) ? prev.filter((h) => h !== hintId) : [...prev, hintId]
    );
  };

  const unlockNext = async () => {
    console.log('[ResultPage] unlockNext called');
    console.log('[ResultPage] Current unlockedLevel:', unlockedLevel);
    console.log('[ResultPage] Total hints:', resultData?.hints.length);

    if (unlockedLevel < 3 && resultData && resultData.hints.length > unlockedLevel) {
      const next = unlockedLevel + 1;
      console.log('[ResultPage] Unlocking level:', next);

      try {
        // Вызываем API для инкремента счётчика подсказок
        await helpAPI.getNextHint(attemptId, unlockedLevel);
        console.log('[ResultPage] API getNextHint called successfully');

        setUnlockedLevel(next);

        analytics.trackEvent('help_hint_requested', {
          child_profile_id: profile?.child_profile_id,
          attempt_id: attemptId,
          hint_level: next,
        });

        const nextHint = resultData.hints.find(h => h.level === next);
        console.log('[ResultPage] Next hint found:', nextHint);

        if (nextHint) {
          setOpenHints((prev) => {
            const newOpenHints = [...prev, nextHint.id];
            console.log('[ResultPage] New openHints:', newOpenHints);
            return newOpenHints;
          });
        }
      } catch (error) {
        console.error('[ResultPage] Failed to get next hint:', error);
        // Не блокируем UI если API упал
        setUnlockedLevel(next);
        const nextHint = resultData.hints.find(h => h.level === next);
        if (nextHint) {
          setOpenHints((prev) => [...prev, nextHint.id]);
        }
      }
    } else {
      console.log('[ResultPage] Cannot unlock - condition failed');
    }
  };

  const handleSubmitAnswer = () => {
    analytics.trackEvent('help_answer_submitted', {
      child_profile_id: profile?.child_profile_id,
      attempt_id: attemptId,
      hints_used: unlockedLevel,
    });

    // Переходим на главную (задание выполнено)
    navigate(ROUTES.HOME);
  };

  const handleNewTask = () => {
    navigate(ROUTES.HELP);
  };

  return (
    <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate(ROUTES.HOME)}
        className="flex items-center gap-2 text-primary mb-4"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">На главную</span>
      </button>

      <h2 className="text-primary mb-2 text-[24px] font-semibold">Подсказки</h2>

      {/* Villain health bar */}
      <div className="bg-white rounded-2xl p-3 mb-4 flex items-center gap-3 shadow-sm">
        <span className="text-[20px]">👾</span>
        <div className="flex-1">
          <p className="text-[12px] text-gray-500 mb-1">Здоровье злодея</p>
          <div className="flex gap-1.5">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className={`flex-1 h-3 rounded-full transition-all ${
                  i <= villainHealth ? 'bg-red-500' : 'bg-gray-200'
                }`}
              />
            ))}
          </div>
        </div>
      </div>

      {/* Hint accordion - sequential unlock */}
      <div className="flex flex-col gap-3 mb-6">
        {resultData.hints.map((hint) => {
          const isLocked = hint.level > unlockedLevel;
          const isOpen = openHints.includes(hint.id);

          return (
            <div
              key={hint.id}
              className={`bg-white rounded-2xl overflow-hidden shadow-sm ${
                isLocked ? 'opacity-60' : ''
              }`}
            >
              <button
                onClick={() => toggleHint(hint.id, hint.level)}
                className="w-full flex items-center justify-between px-4 py-3"
                disabled={isLocked}
              >
                <div className="flex items-center gap-2">
                  <span className="w-7 h-7 bg-blue-100 rounded-lg flex items-center justify-center text-blue-600 text-[12px] font-medium">
                    Л{hint.level}
                  </span>
                  <span className="text-[14px] text-gray-900 font-medium">{hint.title}</span>
                </div>
                {isLocked ? (
                  <Lock size={16} className="text-gray-400" />
                ) : isOpen ? (
                  <ChevronUp size={18} className="text-gray-400" />
                ) : (
                  <ChevronDown size={18} className="text-gray-400" />
                )}
              </button>
              <AnimatePresence>
                {isOpen && !isLocked && (
                  <motion.div
                    initial={{ height: 0, opacity: 0 }}
                    animate={{ height: 'auto', opacity: 1 }}
                    exit={{ height: 0, opacity: 0 }}
                    className="overflow-hidden"
                  >
                    <div className="px-4 pb-4 text-[14px] text-gray-600 border-t border-gray-100 pt-3">
                      {hint.content}
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          );
        })}
      </div>

      {/* Unlock next hint button */}
      {unlockedLevel < 3 && resultData.hints.length > unlockedLevel && (
        <button
          onClick={unlockNext}
          className="w-full py-3 bg-white border-2 border-blue-600 text-blue-600 rounded-2xl mb-3 active:scale-[0.98] transition-transform shadow-sm font-medium hover:bg-blue-50"
        >
          🔓 Открыть следующую подсказку (Л{unlockedLevel + 1})
        </button>
      )}

      {/* Action buttons */}
      <div className="flex flex-col gap-3 mt-auto">
        <button
          onClick={handleSubmitAnswer}
          className="w-full py-4 bg-blue-600 text-white rounded-2xl shadow-lg shadow-blue-600/20 active:scale-[0.98] transition-transform font-semibold"
        >
          Завершить
        </button>
        <button
          onClick={handleNewTask}
          className="w-full py-3 border border-blue-600 text-blue-600 rounded-2xl active:scale-[0.98] transition-transform font-medium"
        >
          Новое задание
        </button>
      </div>
    </div>
  );
}
