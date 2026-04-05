// src/pages/Check/CheckResultPage.tsx
import { useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { Panel, PanelHeader, Group, Div, Title, Text } from '@vkontakte/vkui';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';
import type { CheckResult, CheckError } from '@/types/check';
import styles from './CheckResultPage.module.css';

export default function CheckResultPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const attemptId = location.state?.attemptId as string;
  const result = location.state?.result as CheckResult;

  useEffect(() => {
    if (!attemptId || !result) {
      navigate(ROUTES.CHECK);
      return;
    }

    analytics.trackEvent('check_result_opened', {
      child_profile_id: profile?.child_profile_id,
      attempt_id: attemptId,
      result_status: result.status,
      errors_count: result.errors?.length || 0,
    });
  }, [attemptId, result]);

  const handleReturnHome = () => {
    navigate(ROUTES.HOME);
  };

  const handleResubmit = () => {
    navigate('/check/resubmit', {
      state: { attemptId, result },
    });
  };

  if (!result) return null;

  return (
    <Panel id="check-result">
      <PanelHeader>Результат</PanelHeader>

      <Group>
        <Div>
          {result.status === 'success' && (
            <div className={styles.successSection}>
              <div className={styles.checkmark}>✓</div>
              <Title level="1" weight="1" style={{ marginBottom: '8px' }}>
                Отлично! Всё правильно
              </Title>
              <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
                Задание выполнено без ошибок
              </Text>
            </div>
          )}

          {result.status === 'error' && result.errors && result.errors.length > 0 && (
            <div className={styles.errorSection}>
              <div className={styles.errorIcon}>✕</div>
              <Title level="1" weight="1" style={{ marginBottom: '8px' }}>
                Найдены ошибки
              </Title>
              <Text style={{ color: 'var(--vkui--color_text_secondary)', marginBottom: '16px' }}>
                Исправь их и попробуй снова
              </Text>

              <div className={styles.errors}>
                {result.errors.map((error: CheckError, index: number) => (
                  <Card key={error.id} className={styles.errorCard}>
                    <div className={styles.errorHeader}>
                      <span className={styles.errorNumber}>Ошибка {index + 1}</span>
                      <span className={styles.errorSeverity}>
                        {error.severity === 'error' ? '🔴' : '🟡'}
                      </span>
                    </div>
                    {error.stepNumber && (
                      <Text style={{ fontSize: '14px', marginBottom: '4px' }}>
                        Шаг {error.stepNumber}
                      </Text>
                    )}
                    {error.lineReference && (
                      <Text style={{ fontSize: '14px', marginBottom: '4px' }}>
                        Строка: {error.lineReference}
                      </Text>
                    )}
                    <Text>{error.description}</Text>
                  </Card>
                ))}
              </div>

              <Button
                mode="primary"
                size="l"
                stretched
                onClick={handleResubmit}
                style={{ marginTop: '16px' }}
              >
                Исправить и отправить снова
              </Button>
            </div>
          )}

          {/* Rewards */}
          <Card className={styles.rewardsCard}>
            <Text style={{ fontSize: '14px', color: 'var(--vkui--color_text_secondary)' }}>
              За это задание ты получишь:
            </Text>
            <div className={styles.rewards}>
              <div className={styles.reward}>
                🪙 {result.coinsEarned} монет
              </div>
              <div className={styles.reward}>
                ⚔️ {result.damageDealt} урона злодею
              </div>
            </div>
          </Card>

          <Button
            mode="outline"
            size="l"
            stretched
            onClick={handleReturnHome}
            style={{ marginTop: '16px' }}
          >
            Вернуться на главную
          </Button>
        </Div>
      </Group>
    </Panel>
  );
}
