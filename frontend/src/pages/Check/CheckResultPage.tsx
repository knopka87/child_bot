// src/pages/Check/CheckResultPage.tsx
import { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { ArrowLeft, CheckCircle, AlertTriangle } from 'lucide-react';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { ROUTES } from '@/config/routes';
import type { CheckError } from '@/types/check';

export default function CheckResultPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();

  const attemptId = location.state?.attemptId as string;
  const result = location.state?.result as any;
  const scenario = location.state?.scenario as string;

  useEffect(() => {
    const resultAny = result as any;
    console.log('[CheckResultPage] Received result data:', JSON.stringify(result, null, 2));
    console.log('[CheckResultPage] attemptId:', attemptId);
    console.log('[CheckResultPage] result.status:', result?.status);
    console.log('[CheckResultPage] result.is_correct:', resultAny?.is_correct);
    console.log('[CheckResultPage] result.errors:', result?.errors);
    console.log('[CheckResultPage] hasErrors:', hasErrors);
    
    const trackOpen = async () => {
      if (!attemptId || !result) {
        navigate(ROUTES.CHECK_SCENARIO);
        return;
      }

      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      analytics.trackEvent('check_result_opened', {
        child_profile_id: childProfileId,
        attempt_id: attemptId,
        result_status: result.status,
        errors_count: result.errors?.length || 0,
      });

      // Если есть ошибки, тречим просмотр фидбека
      if (result.errors && result.errors.length > 0) {
        analytics.trackEvent('check_error_feedback_viewed', {
          child_profile_id: childProfileId,
          attempt_id: attemptId,
          error_count: result.errors.length,
        });
      }
    };

    trackOpen();
  }, [attemptId, result, analytics, navigate]);

  if (!result) return null;

  const hasErrors = result.errors && result.errors.length > 0;
  const isFailed = result.status === 'failed';
  const isInternalError = result.status === 'error' && !hasErrors && result.is_correct === null;

  // Вариант 0: Обработка завершилась ошибкой (LLM недоступен и т.д.)
  if (isFailed) {
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
            <div className="text-6xl mb-4">⚠️</div>
            <h2 className="text-[24px] font-bold text-[#E17055] mb-2">Ошибка обработки</h2>
          </div>

          <div className="bg-white rounded-2xl p-4 shadow-sm w-full mb-4">
            <p className="text-[14px] text-[#2D3436]">
              Не удалось проверить задание. Попробуй ещё раз или выбери другое задание.
            </p>
          </div>

          <div className="text-6xl mb-4">🦉</div>
          <p className="text-[14px] text-[#636e72]">Не переживай, попробуем снова!</p>
        </div>

        <div className="flex flex-col gap-3">
          <button
            onClick={() => navigate(`/check/upload-images?scenario=${scenario}`)}
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

  // Вариант 0.5: Внутренняя ошибка LLM (не удалось проверить)
  if (isInternalError) {
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
            <div className="text-6xl mb-4">⚠️</div>
            <h2 className="text-[24px] font-bold text-[#E17055] mb-2">Не удалось проверить задание</h2>
          </div>

          <div className="bg-white rounded-2xl p-4 shadow-sm w-full mb-4">
            <p className="text-[14px] text-[#2D3436]">
              Произошла ошибка при проверке. Попробуй загрузить фото ещё раз или выбери другое задание.
            </p>
          </div>

          <div className="text-6xl mb-4">🦉</div>
          <p className="text-[14px] text-[#636e72]">Не переживай, попробуем снова!</p>
        </div>

        <div className="flex flex-col gap-3">
          <button
            onClick={() => navigate(`/check/upload-images?scenario=${scenario}`)}
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

  // Вариант 1: Всё правильно
  if (result.status === 'success' && !hasErrors) {
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
                  {/* Единая прогресс-бар полоска здоровья */}
                  <div className="w-full h-3 bg-[#E0E0E0] rounded-full overflow-hidden">
                    <div 
                      className="h-full bg-[#FF6B6B] transition-all duration-500 ease-out rounded-full"
                      style={{ width: `${Math.max(0, 100 - result.damageDealt * 20)}%` }}
                    />
                  </div>
                  <p className="text-[14px] font-bold text-[#2D3436] mt-1">
                    -20 HP
                  </p>
                </div>
              </div>
              <p className="text-[11px] text-[#636e72] mt-2 text-center">
                💡 Новый злодей появляется каждый день в полночь!
              </p>
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
  if (hasErrors) {
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

          {/* Сообщение от LLM - основной фидбек */}
          {(result.feedback || (result.errors && result.errors.length > 0)) && (
            <div className="bg-white rounded-2xl p-4 mb-4 border-l-4 border-[#6C5CE7]">
              {result.feedback ? (
                <p className="text-[14px] text-[#2D3436] leading-relaxed">
                  {result.feedback}
                </p>
              ) : (
                <div className="space-y-2">
                  {result.errors?.map((error: CheckError, index: number) => (
                    <p key={error.id} className="text-[14px] text-[#2D3436] leading-relaxed">
                      <span className="font-semibold text-[#FF6B6B]">Ошибка {index + 1}:</span>{' '}
                      {error.description}
                      {error.lineReference && ` (строка ${error.lineReference})`}
                    </p>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>

        <div className="flex flex-col gap-3">
          <button
            onClick={() => {
              // Пытаемся получить изображение задания из sessionStorage
              let taskImage: string | null = null;
              const lastTaskImageStr = sessionStorage.getItem('check_last_task_image');
              if (lastTaskImageStr) {
                try {
                  const lastTaskImage = JSON.parse(lastTaskImageStr);
                  taskImage = lastTaskImage.base64;
                } catch (e) {
                  console.error('Failed to parse last task image:', e);
                }
              }
              
              // Если нет в sessionStorage, пробуем из результата
              if (!taskImage) {
                taskImage = result.taskImage || result.result?.task_image;
              }
              
              navigate(`/check/upload-images?scenario=${scenario}`, { 
                state: { 
                  attemptId, 
                  existingTaskImage: taskImage,
                  mode: 'fix_errors'
                } 
              });
            }}
            className="w-full py-4 bg-[#6C5CE7] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#6C5CE7]/20 active:scale-[0.98] transition-transform"
          >
            Исправил(а) — проверить снова
          </button>
          <button
            onClick={() => {
              // Очищаем данные предыдущей попытки при выборе нового задания
              sessionStorage.removeItem('check_last_task_image');
              sessionStorage.removeItem('check_task_photo');
              sessionStorage.removeItem('check_answer_photo');
              sessionStorage.removeItem('check_single_photo_data');
              
              navigate('/check/scenario');
            }}
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
          onClick={() => {
            // Очищаем данные предыдущей попытки при выборе новой попытки
            sessionStorage.removeItem('check_last_task_image');
            sessionStorage.removeItem('check_task_photo');
            sessionStorage.removeItem('check_answer_photo');
            sessionStorage.removeItem('check_single_photo_data');
            
            navigate(`/check/upload-images?scenario=${scenario}`);
          }}
          className="w-full py-4 bg-[#6C5CE7] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#6C5CE7]/20 active:scale-[0.98] transition-transform"
        >
          Попробовать ещё раз
        </button>
        <button
          onClick={() => {
            // Очищаем данные предыдущей попытки при выборе нового задания
            sessionStorage.removeItem('check_last_task_image');
            sessionStorage.removeItem('check_task_photo');
            sessionStorage.removeItem('check_answer_photo');
            sessionStorage.removeItem('check_single_photo_data');
            
            navigate('/check/scenario');
          }}
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
