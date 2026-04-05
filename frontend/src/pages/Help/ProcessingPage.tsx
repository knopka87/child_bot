// src/pages/Help/ProcessingPage.tsx
import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Panel, PanelHeader, Group, Div, Title, Text } from '@vkontakte/vkui';
import { Spinner } from '@/components/ui/Spinner';
import { Button } from '@/components/ui/Button';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { helpAPI } from '@/api/help';
import { ROUTES } from '@/config/routes';
import styles from './ProcessingPage.module.css';

export default function ProcessingPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const [status, setStatus] = useState<'processing' | 'long_wait' | 'completed' | 'failed'>('processing');
  const attemptId = location.state?.attemptId as string;
  const imageUrl = location.state?.imageUrl as string;

  useEffect(() => {
    if (!attemptId) {
      navigate(ROUTES.HELP);
      return;
    }

    processAttempt();
  }, [attemptId]);

  const processAttempt = async () => {
    try {
      analytics.trackEvent('help_processing_started', {
        child_profile_id: profile?.child_profile_id,
        attempt_id: attemptId,
      });

      // Start processing
      await helpAPI.processAttempt(attemptId);

      // Poll for result
      const result = await pollForResult(attemptId);

      if (result) {
        // Передаём attemptId в URL, а result в state (для оптимизации)
        navigate(`/help/result/${attemptId}`, {
          state: {
            result,
          },
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

        console.log('[ProcessingPage] Poll attempt', i + 1);
        console.log('[ProcessingPage] Hints count:', result.hints?.length || 0);

        // Проверяем что обработка завершена и есть подсказки
        if (result.hints && result.hints.length > 0) {
          console.log('[ProcessingPage] Processing completed with hints!');
          return result;
        }

        // Если нет подсказок, продолжаем polling
        console.log('[ProcessingPage] No hints yet, continuing polling...');
        await new Promise((resolve) => setTimeout(resolve, 2000));

        if (i > 10) {
          setStatus('long_wait');
        }
      } catch (error) {
        console.log('[ProcessingPage] Poll error, retrying...', error);
        // Continue polling
        await new Promise((resolve) => setTimeout(resolve, 2000));

        if (i > 10) {
          setStatus('long_wait');
        }
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

  return (
    <Panel id="processing">
      <PanelHeader>Обработка</PanelHeader>

      <Group>
        <Div className={styles.container}>
          {status === 'processing' && (
            <div className={styles.content}>
              <Spinner size="lg" />
              <Title level="2" weight="2" style={{ marginTop: '16px' }}>
                Обрабатываем задание...
              </Title>
              <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
                Это может занять до минуты
              </Text>
              {imageUrl && (
                <img src={imageUrl} alt="Task" className={styles.preview} />
              )}
            </div>
          )}

          {status === 'long_wait' && (
            <div className={styles.content}>
              <Spinner size="lg" />
              <Title level="2" weight="2" style={{ marginTop: '16px' }}>
                Обработка занимает больше времени...
              </Title>
              <Text style={{ marginBottom: '16px', color: 'var(--vkui--color_text_secondary)' }}>
                Можешь вернуться позже, мы сохраним твой прогресс
              </Text>
              <Button mode="outline" onClick={handleSaveAndWait}>
                Вернуться позже
              </Button>
            </div>
          )}

          {status === 'failed' && (
            <div className={styles.content}>
              <div className={styles.errorIcon}>✕</div>
              <Title level="2" weight="2">
                Ошибка обработки
              </Title>
              <Text style={{ marginBottom: '16px', color: 'var(--vkui--color_text_secondary)' }}>
                Не удалось обработать изображение
              </Text>
              <Button mode="primary" onClick={() => navigate(ROUTES.HELP)}>
                Попробовать снова
              </Button>
            </div>
          )}
        </Div>
      </Group>
    </Panel>
  );
}
