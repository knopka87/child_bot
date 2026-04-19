// src/pages/Home/HomePage.tsx
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Header } from '@/components/layout/Header';
import { BottomNav } from '@/components/layout/BottomNav';
import { ActionButtons } from './components/ActionButtons';
import { MascotBattle } from './components/MascotBattle';
import { UnfinishedAttemptModal } from './components/UnfinishedAttemptModal';
import { LevelUpAnimation } from '@/components/ui/LevelUpAnimation';
import { useHomeData } from './hooks/useHomeData';
import { useNewAchievements } from '@/hooks/useNewAchievements';
import { useAnalytics } from '@/hooks/useAnalytics';
import { homeAPI } from '@/api/home';
import { ROUTES } from '@/config/routes';
import { Spinner } from '@/components/ui/Spinner';

export default function HomePage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const { data, isLoading, error, refetch } = useHomeData();
  const { hasNew: hasNewAchievements } = useNewAchievements();
  const [showUnfinishedModal, setShowUnfinishedModal] = useState(false);
  const [showLevelUp, setShowLevelUp] = useState(false);
  const [levelUpData, setLevelUpData] = useState({ level: 0, coins: 0 });
  const [previousLevel, setPreviousLevel] = useState<number | null>(null);

  // Проверяем повышение уровня
  useEffect(() => {
    if (!data) return;

    const currentLevel = data.profile.level;
    
    // Если это первая загрузка, запоминаем уровень
    if (previousLevel === null) {
      setPreviousLevel(currentLevel);
      return;
    }

    // Если уровень повысился, показываем анимацию
    if (currentLevel > previousLevel) {
      setLevelUpData({
        level: currentLevel,
        coins: 100, // Константа из бэкенда
      });
      setShowLevelUp(true);
    }

    setPreviousLevel(currentLevel);
  }, [data, previousLevel]);

  useEffect(() => {
    // Отправляем события только когда данные загружены
    if (!data) return;

    // Analytics: home_opened
    analytics.trackEvent('home_opened', {
      child_profile_id: data.profile.id,
      entry_point: 'direct',
    });

    // Analytics: level_bar_viewed
    analytics.trackEvent('level_bar_viewed', {
      child_profile_id: data.profile.id,
      level: data.profile.level,
      level_progress_percent: data.profile.levelProgress,
    });

    analytics.trackEvent('coins_balance_viewed', {
      child_profile_id: data.profile.id,
      coins_balance: data.profile.coinsBalance,
    });

    analytics.trackEvent('tasks_correct_count_viewed', {
      child_profile_id: data.profile.id,
      tasks_solved_correct_count: data.profile.tasksSolvedCorrectCount,
    });
  }, [data, analytics]);

  const handleHelpClick = () => {
    // Analytics: home_help_clicked
    analytics.trackEvent('home_help_clicked', {
      child_profile_id: data?.profile.id,
    });

    // Проверяем наличие незавершённой попытки ТОЛЬКО типа help
    if (data?.unfinishedAttempt && data.unfinishedAttempt.mode === 'help') {
      setShowUnfinishedModal(true);

      // Analytics: unfinished_attempt_modal_shown
      analytics.trackEvent('unfinished_attempt_modal_shown', {
        child_profile_id: data.profile.id,
        attempt_id: data.unfinishedAttempt.id,
        mode: data.unfinishedAttempt.mode,
      });
      return;
    }

    navigate(ROUTES.HELP_UPLOAD);
  };

  const handleCheckClick = () => {
    // Analytics: home_check_clicked
    analytics.trackEvent('home_check_clicked', {
      child_profile_id: data?.profile.id,
    });

    // Проверяем наличие незавершённой попытки ТОЛЬКО типа check
    if (data?.unfinishedAttempt && data.unfinishedAttempt.mode === 'check') {
      setShowUnfinishedModal(true);

      // Analytics: unfinished_attempt_modal_shown
      analytics.trackEvent('unfinished_attempt_modal_shown', {
        child_profile_id: data.profile.id,
        attempt_id: data.unfinishedAttempt.id,
        mode: data.unfinishedAttempt.mode,
      });
      return;
    }

    navigate(ROUTES.CHECK_SCENARIO);
  };

  const handleVillainClick = () => {
    if (!data?.villain) return;

    // Analytics: villain_clicked
    analytics.trackEvent('villain_clicked', {
      child_profile_id: data.profile.id,
      villain_id: data.villain.id,
      villain_state: 'active',
    });

    navigate(ROUTES.VILLAIN);
  };

  const handleContinueAttempt = () => {
    if (!data?.unfinishedAttempt) return;

    // Analytics: unfinished_attempt_continue_clicked
    analytics.trackEvent('unfinished_attempt_continue_clicked', {
      child_profile_id: data.profile.id,
      attempt_id: data.unfinishedAttempt.id,
      mode: data.unfinishedAttempt.mode,
    });

    setShowUnfinishedModal(false);

    // Navigate to appropriate flow based on attempt mode and status
    const attempt = data.unfinishedAttempt;
    
    if (attempt.mode === 'help') {
      // For help attempts, check if processing is complete
      if (attempt.status === 'completed') {
        // Help processing is done, go to result page
        navigate(`/help/result/${attempt.id}`, {
          state: { attemptId: attempt.id },
        });
      } else {
        // Still processing, go to processing page
        navigate(ROUTES.HELP_PROCESSING, {
          state: { attemptId: attempt.id },
        });
      }
    } else {
      // Check mode
      if (attempt.status === 'completed') {
        // Check processing is done, go to result page
        navigate('/check/result', {
          state: { attemptId: attempt.id },
        });
      } else {
        // Still processing, go to processing page
        navigate(ROUTES.CHECK_PROCESSING, {
          state: { attemptId: attempt.id },
        });
      }
    }
  };

  const handleNewTask = async () => {
    if (!data?.unfinishedAttempt) return;

    // Analytics: unfinished_attempt_new_task_clicked
    analytics.trackEvent('unfinished_attempt_new_task_clicked', {
      child_profile_id: data.profile.id,
      attempt_id: data.unfinishedAttempt.id,
      mode: data.unfinishedAttempt.mode,
    });

    // Delete unfinished attempt
    try {
      await homeAPI.deleteAttempt(data.unfinishedAttempt.id);
      setShowUnfinishedModal(false);
      refetch(); // Refresh data
    } catch (error) {
      console.error('[HomePage] Failed to delete attempt:', error);
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-[#E8E4FF]">
        <Spinner size="lg" />
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="flex flex-col justify-center items-center min-h-screen px-6 text-center bg-[#E8E4FF]">
        <p className="text-[#2D3436] text-[16px] mb-4">Не удалось загрузить данные</p>
        <button
          onClick={() => refetch()}
          className="py-3 px-6 bg-[#6C5CE7] text-white rounded-2xl font-medium active:scale-[0.98] transition-transform"
        >
          Попробовать снова
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen bg-[#E8E4FF]">
      <Header
        level={data.profile.level}
        xpTotal={data.profile.xpTotal}
        xpForNextLevel={data.profile.xpForNextLevel}
        levelProgress={data.profile.levelProgress}
        coins={data.profile.coinsBalance}
        tasksCount={data.profile.tasksSolvedCorrectCount}
        hasNewAchievements={hasNewAchievements}
      />

      <div className="flex-1 flex flex-col justify-between pb-20">
        {/* Mascot Battle Section */}
        <MascotBattle villain={data.villain} onVillainClick={handleVillainClick} />

        {/* Action Buttons */}
        <div className="px-4">
          <ActionButtons onHelpClick={handleHelpClick} onCheckClick={handleCheckClick} />
        </div>
      </div>

      <BottomNav hasNewAchievements={hasNewAchievements} />

      <UnfinishedAttemptModal
        isOpen={showUnfinishedModal}
        onClose={() => setShowUnfinishedModal(false)}
        onContinue={handleContinueAttempt}
        onNewTask={handleNewTask}
      />

      <LevelUpAnimation
        show={showLevelUp}
        newLevel={levelUpData.level}
        coinsReward={levelUpData.coins}
        onComplete={() => setShowLevelUp(false)}
      />
    </div>
  );
}
