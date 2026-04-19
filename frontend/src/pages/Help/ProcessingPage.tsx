// src/pages/Help/ProcessingPage.tsx
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { ArrowLeft, Save, RefreshCw, X } from 'lucide-react';
import { motion } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';

export default function ProcessingPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();

  const [status, setStatus] = useState<'processing' | 'long_wait' | 'completed' | 'failed'>('processing');
  const [seconds, setSeconds] = useState(0);
  const [dots, setDots] = useState('');
  const attemptId = location.state?.attemptId as string;
  const [imageUrl, setImageUrl] = useState<string | null>(location.state?.imageUrl as string);

  // Загружаем изображение из sessionStorage если нет в state
  useEffect(() => {
    if (!imageUrl) {
      const helpPhotoStr = sessionStorage.getItem('help_photo_data');
      if (helpPhotoStr) {
        try {
          const helpPhoto = JSON.parse(helpPhotoStr);
          setImageUrl(helpPhoto.base64);
        } catch (e) {
          console.error('[ProcessingPage] Failed to parse help photo data:', e);
        }
      }
    }
  }, [imageUrl]);

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

  // Показать опции долгого ожидания после 10 секунд
  useEffect(() => {
    if (seconds >= 10 && status === 'processing') {
      setStatus('long_wait');
    }
  }, [seconds, status]);

  useEffect(() => {
    if (!attemptId) {
      navigate(ROUTES.HELP);
      return;
    }

    processAttempt();
  }, [attemptId]);

  const processAttempt = async () => {
    try {
      const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);

      analytics.trackEvent('help_processing_started', {
        child_profile_id: childProfileId,
        attempt_id: attemptId,
      });

      // Start processing
      await helpAPI.processAttempt(attemptId);

      // Poll for result
      const result = await pollForResult(attemptId);

      if (result) {
        navigate(`/help/result/${attemptId}`, {
          state: { result },
        });
      }
    } catch (error) {
      console.error('[ProcessingPage] Processing failed:', error);
      setStatus('failed');
    }
  };

  const pollForResult = async (attemptId: string, maxAttempts = 60): Promise<any> => {
    for (let i = 0; i < maxAttempts; i++) {
      try {
        const result = await helpAPI.getResult(attemptId);

        if (result.hints && result.hints.length > 0) {
          return result;
        }

        await new Promise((resolve) => setTimeout(resolve, 2000));
      } catch (error) {
        await new Promise((resolve) => setTimeout(resolve, 2000));
      }
    }
    throw new Error('Polling timeout');
  };

  const handleSaveAndWait = async () => {
    try {
      await helpAPI.saveAndWait(attemptId);
      navigate(ROUTES.HOME);
    } catch (error) {
      console.error('[ProcessingPage] Save and wait failed:', error);
    }
  };

  const handleRetry = () => {
    setSeconds(0);
    setStatus('processing');
    processAttempt();
  };

  const messages = [
    { time: 0, text: 'Думаю' },
    { time: 8, text: 'Нужно чуть больше времени' },
    { time: 15, text: 'Почти готово, подожди ещё немного' },
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
          <div className="text-6xl mb-4">✕</div>
          <h2 className="text-[24px] font-bold text-[#2D3436] mb-2">Ошибка обработки</h2>
          <p className="text-[14px] text-[#636e72] mb-6">
            Не удалось обработать изображение
          </p>
          <button
            onClick={() => navigate(ROUTES.HELP)}
            className="w-full py-4 bg-[#6C5CE7] text-white text-[16px] font-semibold rounded-2xl shadow-lg shadow-[#6C5CE7]/20 active:scale-[0.98] transition-transform"
          >
            Попробовать снова
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

          {/* Превью изображения */}
          {imageUrl && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.5 }}
              className="w-full max-w-[280px] mb-6"
            >
              <img
                src={imageUrl}
                alt="Загруженное изображение"
                className="w-full h-auto rounded-2xl shadow-lg"
              />
            </motion.div>
          )}

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
                onClick={() => navigate(ROUTES.HOME)}
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
