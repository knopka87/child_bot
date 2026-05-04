// src/pages/Friends/FriendsPage.tsx
import { useEffect, useState } from 'react';
import { Send, Users, Gift } from 'lucide-react';
import { motion } from 'framer-motion';
import { BottomNav } from '@/components/layout/BottomNav';
import { ListPageSkeleton } from '@/components/ui/skeleton';
import { useReferralData } from './hooks/useReferralData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { referralAPI } from '@/api/referral';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import bridge from '@/lib/platform/bridge';
import styles from './FriendsPage.module.css';

export function FriendsPage() {
  const analytics = useAnalytics();
  const { data, isLoading, error } = useReferralData();
  const [childProfileId, setChildProfileId] = useState<string | null>(null);
  const [showCopyFallback, setShowCopyFallback] = useState(false);

  // Загружаем child_profile_id для аналитики
  useEffect(() => {
    const loadProfileId = async () => {
      const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      setChildProfileId(profileId);

      if (profileId) {
        analytics.trackEvent('friends_opened', {
          child_profile_id: profileId,
        });
      }
    };

    loadProfileId();
  }, [analytics]);

  const handleInvite = async () => {
    console.log('[FriendsPage] handleInvite called');

    if (!data) {
      console.warn('[FriendsPage] No data available');
      return;
    }

    try {
      console.log('[FriendsPage] Starting invite with code:', data.referralCode);
      console.log('[FriendsPage] Using link:', `https://vk.com/app54517931#start=${data.referralCode}`);

      // Используем VKWebAppShare для отправки приглашения
      // Для Mini Apps это основной способ пригласить друзей
      const result = await bridge.send('VKWebAppShare', {
        link: `https://vk.com/app54517931#start=${data.referralCode}`,
      });

      console.log('[FriendsPage] Share result:', result);

      if (childProfileId) {
        analytics.trackEvent('referral_invite_sent', {
          child_profile_id: childProfileId,
          referral_code: data.referralCode,
          share_channel: 'vk_share',
        });

        await referralAPI.trackInviteSent(childProfileId, 'vk');
      }
    } catch (error: any) {
      console.error('[FriendsPage] Share failed:', {
        error,
        errorType: error?.error_type,
        errorCode: error?.error_data?.error_code,
        errorReason: error?.error_data?.error_reason,
        errorMessage: error?.message,
      });

      // Если пользователь отменил диалог
      if (error?.error_data?.error_code === 4) {
        console.log('[FriendsPage] User cancelled share dialog');
        return;
      }

      // Показываем fallback с возможностью скопировать ссылку
      console.log('[FriendsPage] Showing copy fallback');
      setShowCopyFallback(true);
    }
  };

  const handleCopyLink = () => {
    if (!data) return;

    const link = `https://vk.com/app54517931#start=${data.referralCode}`;
    navigator.clipboard?.writeText(link)
      .then(() => {
        alert('Ссылка скопирована! Отправь её другу в ВК');

        if (childProfileId) {
          analytics.trackEvent('referral_link_copied', {
            child_profile_id: childProfileId,
            referral_code: data.referralCode,
          });
        }
      })
      .catch(() => {
        alert('Не удалось скопировать. Попробуй ещё раз.');
      });
  };

  if (isLoading) {
    return <ListPageSkeleton itemCount={3} />;
  }

  if (error || !data) {
    return (
      <div className={styles.container}>
        <div className={styles.errorContainer}>
          <p>Не удалось загрузить данные</p>
        </div>
      </div>
    );
  }

  const invitedCount = data.invitedCount; // Относительный прогресс (от предыдущего уровня)
  const targetCount = data.targetCount; // Относительная цель (от предыдущего уровня)
  const totalInvited = data.totalInvited; // Абсолютное количество всех приглашённых

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Друзья</h1>
      <p className={styles.subtitle}>Пригласи друзей и учитесь вместе!</p>

      {/* Карточка прогресса с наградой */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className={styles.rewardCard}
      >
        <div className={styles.rewardHeader}>
          <Gift size={20} className={styles.giftIcon} />
          <p className={styles.rewardText}>
            Пригласи ещё {targetCount} {targetCount === 1 ? 'друга' : 'друзей'} — получи редкий стикер!
          </p>
        </div>

        <div className={styles.progressDots}>
          {Array.from({ length: targetCount }).map((_, i) => (
            <div
              key={i}
              className={`${styles.progressDot} ${
                i < invitedCount ? styles.progressDotActive : ''
              }`}
            >
              {i < invitedCount ? '✓' : i + 1}
            </div>
          ))}
        </div>

        <p className={styles.progressCount}>
          {invitedCount} из {targetCount}
        </p>

        {/* Превью награды */}
        <div className={styles.rewardPreview}>
          <span className={styles.rewardEmoji}>⭐</span>
          <span className={styles.rewardName}>
            {data.currentGoal.reward.name}
          </span>
        </div>
      </motion.div>

      {/* Карточка приглашения */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className={styles.inviteCard}
      >
        <div className={styles.inviteIcon}>
          <Users size={30} className={styles.usersIcon} />
        </div>

        <h3 className={styles.inviteTitle}>Пригласи друга</h3>
        <p className={styles.inviteDescription}>
          {showCopyFallback
            ? 'Скопируй ссылку и отправь её другу в ВК!'
            : 'Нажми кнопку ниже, выбери друзей из списка и получите оба бонусные стикеры!'}
        </p>

        {!showCopyFallback ? (
          <button onClick={handleInvite} className={styles.inviteButton}>
            <Send size={20} />
            <span>Пригласить друга</span>
          </button>
        ) : (
          <button onClick={handleCopyLink} className={styles.inviteButton}>
            <Send size={20} />
            <span>Скопировать ссылку</span>
          </button>
        )}
      </motion.div>

      {/* Статистика */}
      <div className={styles.statsCard}>
        <div className={styles.statsRow}>
          <span className={styles.statsLabel}>Приглашено друзей</span>
          <span className={styles.statsValue}>{totalInvited}</span>
        </div>
      </div>

      <BottomNav />
    </div>
  );
}
