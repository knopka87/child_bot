// src/pages/Check/CheckResultPage.tsx
import { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { ArrowLeft, CheckCircle, AlertTriangle } from 'lucide-react';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { ROUTES } from '@/config/routes';
import type { CheckResult, CheckError } from '@/types/check';

export default function CheckResultPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();

  const attemptId = location.state?.attemptId as string;
  const result = location.state?.result as CheckResult;

  useEffect(() => {
    const trackOpen = async () => {
      if (!attemptId || !result) {
        navigate(ROUTES.CHECK);
        return;
      }

      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      analytics.trackEvent('check_result_opened', {
        child_profile_id: childProfileId,
        attempt_id: attemptId,
        result_status: result.status,
        errors_count: result.errors?.length || 0,
      });
    };

    trackOpen();
  }, [attemptId, result, analytics, navigate]);

  if (!result) return null;

  // Определяем вердикт
  const verdict: 'correct' | 'errors' | 'review' =
    result.status === 'success'
      ? 'correct'
      : result.errors && result.errors.length > 0
      ? 'errors'
      : 'review';

  // Вариант 1: Всё правильно
  if (verdict === 'correct') {
    return (
      <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#E8FFF8] to-white">
        <button
          onClick={() => navigate(ROUTES.HOME)}
          className="flex items-center gap-2 text-[#6C5CE7] mb-4 active:opacity-70 transition-opacity"
        >
          <ArrowLeft size={20} />
          <span className="text-[14px] font-medium">На главную</span>
        </button>

        <div className="flex-1 flex flex-col items-center justify-center text-center">
          <div
            className="mb-4"
            style={{ animation: 'scaleIn 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275)' }}
          >
            <CheckCircle size={64} className="text-[#00B894] mx-auto" />
          </div>

          <h1 className="text-[32px] font-bold text-[#00B894] mb-2">Верно!</h1>
          <p className="text-[16px] text-[#636e72] mb-6">Молодец! Задание выполнено без ошибок</p>

          {/* Урон злодею */}
          {result.damageDealt > 0 && (
            <div className="bg-white rounded-2xl p-4 shadow-sm w-full mb-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 bg-[#FFE8F0] rounded-full flex items-center justify-center flex-shrink-0">
                  <span className="text-2xl">👾</span>
                </div>
                <div className="flex-1 text-left">
                  <p className="text-[12px] text-[#636e72] mb-1">Урон злодею!</p>
                  <div className="flex gap-1.5">
                    {[1, 2, 3].map((i) => (
                      <div
                        key={i}
                        className={`flex-1 h-3 rounded-full ${
                          i <= 3 - result.damageDealt ? 'bg-[#FF6B6B]' : 'bg-[#E0E0E0]'
                        }`}
                      />
                    ))}
                  </div>
                  <p className="text-[14px] font-bold text-[#2D3436] mt-1">
                    -{result.damageDealt} HP
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Награда */}
          {result.coinsEarned > 0 && (
            <div className="bg-[#FFF9E8] rounded-2xl p-4 w-full mb-4">
              <div className="flex items-center justify-center gap-2">
                <span className="text-3xl">🪙</span>
                <span className="text-[18px] font-bold text-[#2D3436]">
                  +{result.coinsEarned} монет
                </span>
              </div>
            </div>
          )}
        </div>

        <button
          onClick={() => navigate('/check/scenario')}
          className="w-full py-4 bg-[#6C5CE7] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#6C5CE7]/20 active:scale-[0.98] transition-transform"
        >
          Новое задание
        </button>

        <style>{`
          @keyframes scaleIn {
            from {
              transform: scale(0);
            }
            to {
              transform: scale(1);
            }
          }
        `}</style>
      </div>
    );
  }

  // Вариант 2: Есть ошибки
  if (verdict === 'errors') {
    return (
      <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#FFF9E8] to-white">
        <button
          onClick={() => navigate(ROUTES.HOME)}
          className="flex items-center gap-2 text-[#6C5CE7] mb-4 active:opacity-70 transition-opacity"
        >
          <ArrowLeft size={20} />
          <span className="text-[14px] font-medium">На главную</span>
        </button>

        <div className="flex-1 flex flex-col">
          <div
            className="bg-[#FFF9E8] rounded-3xl p-6 text-center mb-4"
            style={{ animation: 'fadeIn 0.3s ease-out' }}
          >
            <AlertTriangle size={48} className="text-[#FDCB6E] mx-auto mb-3" />
            <h2 className="text-[24px] font-bold text-[#E17055]">Есть ошибки</h2>
            <p className="text-[14px] text-[#636e72] mt-1">Не переживай, попробуй ещё раз!</p>
          </div>

          {/* Список ошибок */}
          <div className="space-y-3 mb-4">
            {result.errors?.map((error: CheckError, index: number) => (
              <div
                key={error.id}
                className="bg-[#FFF0F0] rounded-2xl p-4 border border-[#FFD0D0]"
              >
                <p className="text-[12px] text-[#FF6B6B] mb-1 font-semibold">
                  Ошибка {index + 1}
                  {error.stepNumber && ` • Шаг ${error.stepNumber}`}
                </p>
                <p className="text-[14px] text-[#2D3436]">{error.description}</p>
                {error.lineReference && (
                  <p className="text-[12px] text-[#636e72] mt-1">Строка: {error.lineReference}</p>
                )}
              </div>
            ))}
          </div>
        </div>

        <div className="flex flex-col gap-3">
          <button
            onClick={() => navigate('/check/upload', { state: { attemptId, result } })}
            className="w-full py-4 bg-[#6C5CE7] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#6C5CE7]/20 active:scale-[0.98] transition-transform"
          >
            Исправил(а) — проверить снова
          </button>
          <button
            onClick={() => navigate('/check/scenario')}
            className="w-full py-3 border-2 border-[#6C5CE7] text-[#6C5CE7] text-[16px] font-semibold rounded-2xl active:scale-[0.98] transition-transform"
          >
            Новое задание
          </button>
        </div>

        <style>{`
          @keyframes fadeIn {
            from {
              opacity: 0;
              transform: scale(0.9);
            }
            to {
              opacity: 1;
              transform: scale(1);
            }
          }
        `}</style>
      </div>
    );
  }

  // Вариант 3: Посмотри ещё раз (нет конкретных ошибок)
  return (
    <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#FFF0F0] to-white">
      <button
        onClick={() => navigate(ROUTES.HOME)}
        className="flex items-center gap-2 text-[#6C5CE7] mb-4 active:opacity-70 transition-opacity"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px] font-medium">На главную</span>
      </button>

      <div className="flex-1 flex flex-col items-center justify-center text-center">
        <div style={{ animation: 'fadeIn 0.3s ease-out' }}>
          <AlertTriangle size={48} className="text-[#FDCB6E] mx-auto mb-3" />
          <h2 className="text-[24px] font-bold text-[#E17055] mb-2">Посмотри ещё раз</h2>
        </div>

        <div className="bg-white rounded-2xl p-4 shadow-sm w-full mb-4">
          <p className="text-[14px] text-[#2D3436]">
            Ничего страшного! Ошибки — это нормально. Попробуй ещё раз, ты справишься!
          </p>
        </div>

        <div className="text-6xl mb-4">🦉</div>
        <p className="text-[14px] text-[#636e72]">Не сдавайся!</p>
      </div>

      <div className="flex flex-col gap-3">
        <button
          onClick={() => navigate('/check/upload', { state: { attemptId } })}
          className="w-full py-4 bg-[#6C5CE7] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#6C5CE7]/20 active:scale-[0.98] transition-transform"
        >
          Попробовать ещё раз
        </button>
        <button
          onClick={() => navigate('/check/scenario')}
          className="w-full py-3 border-2 border-[#6C5CE7] text-[#6C5CE7] text-[16px] font-semibold rounded-2xl active:scale-[0.98] transition-transform"
        >
          Новое задание
        </button>
      </div>

      <style>{`
        @keyframes fadeIn {
          from {
            opacity: 0;
            transform: scale(0.9);
          }
          to {
            opacity: 1;
            transform: scale(1);
          }
        }
      `}</style>
    </div>
  );
}
