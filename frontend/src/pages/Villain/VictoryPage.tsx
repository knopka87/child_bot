// src/pages/Villain/VictoryPage.tsx
import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Panel, PanelHeader, Group, Div } from '@vkontakte/vkui';
import { Button } from '@/components/ui/Button';
import { Spinner } from '@/components/ui/Spinner';
import { VictoryAnimation } from './components/VictoryAnimation';
import { RewardCard } from './components/RewardCard';
import { useVictoryData } from './hooks/useVictoryData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { ROUTES } from '@/config/routes';
import { ASSETS } from '@/config/assets';
import styles from './VictoryPage.module.css';

export function VictoryPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);
  const { villainId } = (location.state as { villainId?: string }) || {};

  // Validate villainId
  if (!villainId) {
    return (
      <Panel id="victory">
        <PanelHeader>Ошибка</PanelHeader>
        <Group>
          <Div className={styles.errorContainer}>
            <p>Неверный ID злодея</p>
            <button onClick={() => navigate(ROUTES.HOME)}>На главную</button>
          </Div>
        </Group>
      </Panel>
    );
  }

  const { victory, isLoading, error } = useVictoryData(villainId);

  useEffect(() => {
    if (victory && villainId) {
      analytics.trackEvent('victory_screen_opened', {
        child_profile_id: profile?.child_profile_id,
        villain_id: villainId,
        attempt_id: '',
      });

      victory.rewards.forEach((reward) => {
        analytics.trackEvent('victory_reward_viewed', {
          child_profile_id: profile?.child_profile_id,
          villain_id: villainId,
          reward_type: reward.type,
          reward_id: reward.id,
        });
      });
    }
  }, [victory, villainId, analytics, profile]);

  const handleContinue = () => {
    if (villainId) {
      analytics.trackEvent('victory_continue_clicked', {
        child_profile_id: profile?.child_profile_id,
        villain_id: villainId,
      });
    }

    navigate(ROUTES.HOME);
  };

  if (isLoading) {
    return (
      <Panel id="victory">
        <PanelHeader>Победа!</PanelHeader>
        <Group>
          <Div className={styles.loadingContainer}>
            <Spinner size="lg" />
          </Div>
        </Group>
      </Panel>
    );
  }

  if (error || !victory) {
    return (
      <Panel id="victory">
        <PanelHeader>Победа!</PanelHeader>
        <Group>
          <Div className={styles.errorContainer}>
            <p>Не удалось загрузить данные победы</p>
            <button onClick={() => navigate(ROUTES.HOME)}>На главную</button>
          </Div>
        </Group>
      </Panel>
    );
  }

  return (
    <Panel id="victory">
      <VictoryAnimation />

      <PanelHeader>Победа!</PanelHeader>

      <Group>
        <Div>
          <div className={styles.header}>
            <h1 className={styles.title}>🎉 Победа!</h1>
            <p className={styles.subtitle}>Ты победил злодея!</p>
          </div>

          <div className={styles.villainDefeated}>
            <img
              src={ASSETS.images.villainDefeated}
              alt="Defeated villain"
              className={styles.villainImage}
            />
          </div>

          <div className={styles.stats}>
            <div className={styles.statItem}>
              <span className={styles.statLabel}>Урон нанесён</span>
              <span className={styles.statValue}>{victory.total_damage}</span>
            </div>
            <div className={styles.statItem}>
              <span className={styles.statLabel}>Задач выполнено</span>
              <span className={styles.statValue}>{victory.tasks_completed}</span>
            </div>
          </div>

          <div className={styles.rewards}>
            <h3 className={styles.rewardsTitle}>Твои награды</h3>
            {victory.rewards.map((reward) => (
              <RewardCard key={reward.id} reward={reward} />
            ))}
          </div>

          <Button mode="primary" size="l" stretched onClick={handleContinue}>
            Продолжить учиться
          </Button>
        </Div>
      </Group>
    </Panel>
  );
}
