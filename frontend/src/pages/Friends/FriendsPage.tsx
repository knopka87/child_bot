// src/pages/Friends/FriendsPage.tsx
import { useEffect, useState } from 'react';
import { Copy, Send, Users, CheckCircle, Gift } from 'lucide-react';
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
  const [copied, setCopied] = useState(false);
  const [childProfileId, setChildProfileId] = useState<string | null>(null);

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

  const handleCopy = () => {
    if (!data) return;

    navigator.clipboard?.writeText(data.referralLink).catch(() => {});
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);

    if (childProfileId) {
      analytics.trackEvent('referral_link_copied', {
        child_profile_id: childProfileId,
        referral_code: data.referralCode,
      });
    }
  };

  const handleShare = async () => {
    if (!data) return;

    try {
      // Используем официальный VK механизм приглашений
      // Документация: https://dev.vk.com/ru/bridge/VKWebAppShowInviteBox
      // requestKey передаётся приглашённому как vk_request_key в Launch Params
      const result = await bridge.send('VKWebAppShowInviteBox', {
        requestKey: data.referralCode,
      } as any);

      console.log('[FriendsPage] Invite result:', result);

      if (childProfileId) {
        analytics.trackEvent('referral_invite_sent', {
          child_profile_id: childProfileId,
          referral_code: data.referralCode,
          share_channel: 'vk_invite_box',
        });

        await referralAPI.trackInviteSent(childProfileId, 'vk');
      }
    } catch (vkError: any) {
      console.log('[FriendsPage] VK invite failed:', vkError);

      // Если пользователь отменил приглашение
      if (vkError?.error_data?.error_code === 4 && vkError?.error_data?.error_reason === 'User denied') {
        console.log('[FriendsPage] User cancelled invite dialog');
        return;
      }

      // Fallback: если VKWebAppShowInviteBox не поддерживается, используем VKWebAppShare
      try {
        await bridge.send('VKWebAppShare', {
          link: data.referralLink,
        });

        if (childProfileId) {
          analytics.trackEvent('referral_link_shared', {
            child_profile_id: childProfileId,
            referral_code: data.referralCode,
            share_channel: 'vk_share_fallback',
          });

          await referralAPI.trackInviteSent(childProfileId, 'vk');
        }
      } catch (shareError) {
        console.log('[FriendsPage] VKWebAppShare also failed:', shareError);

        // Последний fallback - копируем ссылку
        handleCopy();
      }
    }
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
          Отправь ссылку другу и получите оба бонусные стикеры!
        </p>

        <div className={styles.linkContainer}>{data.referralLink}</div>

        <div className={styles.buttonsRow}>
          <button onClick={handleCopy} className={`${styles.button} ${styles.buttonPrimary} ${copied ? styles.buttonSuccess : ''}`}>
            {copied ? <CheckCircle size={18} /> : <Copy size={18} />}
            {copied ? 'Скопировано' : 'Скопировать'}
          </button>
          <button onClick={handleShare} className={`${styles.button} ${styles.buttonOutline}`}>
            <Send size={18} />
            Отправить
          </button>
        </div>
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
