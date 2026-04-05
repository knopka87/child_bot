// src/pages/Friends/FriendsPage.tsx
import { useEffect, useState } from 'react';
import { Copy, Send, Users, CheckCircle, Gift } from 'lucide-react';
import { motion } from 'framer-motion';
import { BottomNav } from '@/components/layout/BottomNav';
import { Spinner } from '@/components/ui/Spinner';
import { useReferralData } from './hooks/useReferralData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { referralAPI } from '@/api/referral';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
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

    if (navigator.share) {
      try {
        await navigator.share({
          title: 'Помощник ДЗ',
          text: 'Присоединяйся!',
          url: data.referralLink,
        });

        if (childProfileId) {
          analytics.trackEvent('referral_link_shared', {
            child_profile_id: childProfileId,
            referral_code: data.referralCode,
            share_channel: 'native',
          });

          await referralAPI.trackInviteSent(childProfileId, 'native');
        }
      } catch (err) {
        // User cancelled share
      }
    }
  };

  if (isLoading) {
    return (
      <div className={styles.container}>
        <div className={styles.loadingContainer}>
          <Spinner size="lg" />
        </div>
      </div>
    );
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

  const invitedCount = data.invitedCount;
  const targetCount = data.targetCount;

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
            Пригласи {targetCount} друзей — получи редкий стикер!
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
          <span className={styles.statsValue}>{invitedCount}</span>
        </div>
      </div>

      <BottomNav />
    </div>
  );
}
