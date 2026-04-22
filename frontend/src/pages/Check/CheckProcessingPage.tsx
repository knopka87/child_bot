// src/pages/Check/CheckProcessingPage.tsx
import { useEffect, useState, useRef } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { ArrowLeft, Save, RefreshCw, X } from 'lucide-react';
import { motion } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { checkAPI } from '@/api/check';
import { ROUTES } from '@/config/routes';

const LONG_WAIT_THRESHOLD = 30000; // 30 seconds
const POLLING_INTERVAL = 2000; // 2 seconds
const MAX_POLL_ATTEMPTS = 150; // 5 minutes max

export default function CheckProcessingPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();

  const [status, setStatus] = useState<'processing' | 'long_wait' | 'completed' | 'failed'>('processing');
  const [seconds, setSeconds] = useState(0);
  const [dots, setDots] = useState('');
  const attemptId = location.state?.attemptId as string;
  const scenario = location.state?.scenario as string;
  const pollingRef = useRef<NodeJS.Timeout | null>(null);
  const hasNavigatedRef = useRef(false); // Prevent multiple navigations
  const startTimeRef = useRef(Date.now());

  // Счётчик секунд
  useEffect(() => {
    const timer = setInterval(() => setSeconds((s) => s + 1), 1000);
    return () => clearInterval(timer);
  }, []);

  // Анимация точек
  useEffect(() => {
    const dotTimer = setInterval(() => {
      setDots((d) => (d.length >= 3 ? '' : d + '.'));
    }, 500);
    return () => clearInterval(dotTimer);
  }, []);

  // Показать опции долгого ожидания после порога
  useEffect(() => {
    if (seconds * 1000 >= LONG_WAIT_THRESHOLD && status === 'processing') {
      setStatus('long_wait');

      const trackLongWait = async () => {
        const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
        analytics.trackEvent('check_long_wait_shown', {
          child_profile_id: childProfileId,
          attempt_id: attemptId,
          duration_seconds: seconds,
        });
      };
      trackLongWait();
    }
  }, [seconds, status, attemptId, analytics]);

  useEffect(() => {
    if (!attemptId) {
      navigate(ROUTES.CHECK_SCENARIO);
      return;
    }

    processAttempt();

    return () => {
      if (pollingRef.current) {
        clearInterval(pollingRef.current);
      }
    };
  }, [attemptId]);

  const processAttempt = async () => {
    try {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);

      analytics.trackEvent('check_processing_started', {
        child_profile_id: childProfileId,
        attempt_id: attemptId,
        scenario,
      });

      console.log('[CheckProcessingPage] Starting process for attempt:', attemptId);

      // Start processing
      const processResponse = await checkAPI.processAttempt(attemptId);
      console.log('[CheckProcessingPage] Process response:', processResponse);

      // Poll for result
      startPolling(attemptId);
    } catch (error: any) {
      console.error('[CheckProcessingPage] Processing failed:', error);
      console.error('[CheckProcessingPage] Error response:', error.response?.data);
      console.error('[CheckProcessingPage] Error status:', error.response?.status);
      setStatus('failed');
    }
  };

  const startPolling = async (attemptId: string) => {
    let pollCount = 0;

    pollingRef.current = setInterval(async () => {
      // Check if we already navigated (prevent race conditions)
      if (hasNavigatedRef.current) {
        if (pollingRef.current) {
          clearInterval(pollingRef.current);
          pollingRef.current = null;
        }
        return;
      }

      try {
        pollCount++;
        console.log(`[CheckProcessingPage] Poll attempt ${pollCount}/${MAX_POLL_ATTEMPTS} for attempt:`, attemptId);

        if (pollCount > MAX_POLL_ATTEMPTS) {
          console.error('[CheckProcessingPage] Polling timeout exceeded');
          hasNavigatedRef.current = true; // Prevent further polling
          if (pollingRef.current) {
            clearInterval(pollingRef.current);
            pollingRef.current = null;
          }
          setStatus('failed');
          return;
        }

        const result: any = await checkAPI.getResult(attemptId);

        // Double-check navigation flag after async call
        if (hasNavigatedRef.current) {
          if (pollingRef.current) {
            clearInterval(pollingRef.current);
            pollingRef.current = null;
          }
          return;
        }

        console.log('[CheckProcessingPage] Poll result:', result);

        // Backend returns: { status: 'completed', result: { status: 'success', is_correct: true, ... } }
        // We need to check the nested result.status ('success' or 'error')
        const checkResultStatus = result.result?.status || result.status;
        const hasResult = checkResultStatus && 
          (checkResultStatus === 'success' || checkResultStatus === 'error' ||
           checkResultStatus === 'completed' || checkResultStatus === 'failed');

        if (result && hasResult) {
          // Set flag IMMEDIATELY before any async operations
          hasNavigatedRef.current = true;
          
          // Clear polling IMMEDIATELY
          if (pollingRef.current) {
            clearInterval(pollingRef.current);
            pollingRef.current = null;
          }
          
          console.log('[CheckProcessingPage] Processing finished with check status:', checkResultStatus);

          setStatus('completed');

          const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);

          // Extract the actual check result from the nested structure
          const checkResult = result.result || result;

          analytics.trackEvent('check_processing_completed', {
            child_profile_id: childProfileId,
            attempt_id: attemptId,
            result_status: checkResult.status,
            errors_count: checkResult.errors?.length || 0,
          });

          // Navigate to result page
          navigate(ROUTES.CHECK_PROCESSING.replace('/check/processing', '/check/result'), {
            state: {
              attemptId,
              result: checkResult,
              scenario,
            },
          });
        }
      } catch (error: any) {
        // Still processing - continue polling
        if (pollCount % 10 === 0) {
          console.log(`[CheckProcessingPage] Still polling... (${pollCount}/${MAX_POLL_ATTEMPTS})`);
        }
        if (error.response?.status && error.response?.status !== 404 && error.response?.status !== 200) {
          console.error('[CheckProcessingPage] Polling error:', error.response?.data);
        }
      }
    }, POLLING_INTERVAL);
  };

  const handleSaveAndWait = async () => {
    try {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);

      analytics.trackEvent('check_save_and_wait_clicked', {
        child_profile_id: childProfileId,
        attempt_id: attemptId,
      });

      await checkAPI.saveAndWait(attemptId);

      analytics.trackEvent('check_saved_for_later', {
        child_profile_id: childProfileId,
        attempt_id: attemptId,
      });

      navigate(ROUTES.HOME);
    } catch (error) {
      console.error('[CheckProcessingPage] Save and wait failed:', error);
    }
  };

  const handleRetry = () => {
    const childProfileId = vkStorage.getItem(storageKeys.PROFILE_ID);
    analytics.trackEvent('check_retry_clicked', {
      child_profile_id: childProfileId,
      attempt_id: attemptId,
    });

    setSeconds(0);
    setStatus('processing');
    startTimeRef.current = Date.now();

    if (pollingRef.current) {
      clearInterval(pollingRef.current);
    }

    processAttempt();
  };

  const handleCancel = async () => {
    try {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);

      analytics.trackEvent('check_cancel_clicked', {
        child_profile_id: childProfileId,
        attempt_id: attemptId,
      });

      if (pollingRef.current) {
        clearInterval(pollingRef.current);
      }

      await checkAPI.cancelAttempt(attemptId);
      navigate(ROUTES.HOME);
    } catch (error) {
      console.error('[CheckProcessingPage] Cancel failed:', error);
    }
  };

  const messages = [
    { time: 0, text: 'Проверяю задание' },
    { time: 8, text: 'Нужно чуть больше времени' },
    { time: 15, text: 'Почти готово, подожди ещё немного' },
    { time: 25, text: 'Скоро закончу' },
  ];

  const currentMessage =
    [...messages].reverse().find((m) => seconds >= m.time)?.text || messages[0].text;

  return (
    <div className="flex flex-col min-h-screen px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF]">
      <button
        onClick={() => navigate(ROUTES.HOME)}
        className="flex items-center gap-2 text-[#6C5CE7] mb-6 active:opacity-70 transition-opacity"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px] font-medium">Назад</span>
      </button>

      {status === 'failed' ? (
        <div className="flex-1 flex flex-col items-center justify-center text-center">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="text-6xl mb-4"
          >
            ✕
          </motion.div>
          <h2 className="text-[24px] font-bold text-[#2D3436] mb-2">Ошибка обработки</h2>
          <p className="text-[14px] text-[#636e72] mb-6">
            Не удалось проверить изображение
          </p>
          <button
            onClick={handleRetry}
            className="w-full py-4 bg-[#6C5CE7] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#6C5CE7]/20 active:scale-[0.98] transition-transform mb-3"
          >
            Попробовать снова
          </button>
          <button
            onClick={() => navigate(ROUTES.HOME)}
            className="w-full py-3 text-[#636e72] text-[14px] active:opacity-70 transition-opacity"
          >
            На главную
          </button>
        </div>
      ) : (
        <div className="flex-1 flex flex-col items-center justify-center">
          {/* Mascot с сообщением */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="text-center mb-6"
          >
            <div className="text-[120px] leading-none mb-4">🦉</div>
            <div className="bg-white rounded-2xl py-3 px-5 shadow-sm inline-block">
              <p className="text-[#6C5CE7] text-[16px] font-medium">
                {currentMessage}{dots}
              </p>
            </div>
          </motion.div>

          {/* Progress bar */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.3 }}
            className="w-full max-w-[240px] mb-6"
          >
            <div className="h-2 bg-white rounded-full overflow-hidden">
              <motion.div
                animate={{ x: ['-100%', '100%'] }}
                transition={{ repeat: Infinity, duration: 1.5, ease: 'easeInOut' }}
                className="w-full h-full bg-gradient-to-r from-[#6C5CE7] to-[#A29BFE] rounded-full"
              />
            </div>
            <p className="text-[#636e72] text-[12px] text-center mt-3">
              Обработка продолжится, даже если ты выйдешь
            </p>
          </motion.div>

          {/* Long wait actions */}
          {status === 'long_wait' && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              className="w-full flex flex-col gap-3"
            >
              <button
                onClick={handleSaveAndWait}
                className="w-full py-3 bg-white border border-[#DFE6E9] text-[#2D3436] rounded-2xl flex items-center justify-center gap-2 active:scale-[0.98] transition-transform"
              >
                <Save size={16} className="text-[#6C5CE7]" />
                <span className="text-[14px] font-medium">Сохранить и подождать</span>
              </button>
              <button
                onClick={handleRetry}
                className="w-full py-3 bg-white border border-[#DFE6E9] text-[#2D3436] rounded-2xl flex items-center justify-center gap-2 active:scale-[0.98] transition-transform"
              >
                <RefreshCw size={16} className="text-[#6C5CE7]" />
                <span className="text-[14px] font-medium">Повторить</span>
              </button>
              <button
                onClick={handleCancel}
                className="w-full py-3 text-[#636e72] rounded-2xl flex items-center justify-center gap-2 active:opacity-70 transition-opacity"
              >
                <X size={16} />
                <span className="text-[14px]">Отменить</span>
              </button>
            </motion.div>
          )}
        </div>
      )}
    </div>
  );
}
