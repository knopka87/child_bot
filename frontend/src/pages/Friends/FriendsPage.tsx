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
      console.log('[FriendsPage] ===== INVITE DEBUG =====');
      console.log('[FriendsPage] Referral code:', data.referralCode);
      console.log('[FriendsPage] Platform:', navigator.platform);
      console.log('[FriendsPage] User Agent:', navigator.userAgent);
      console.log('[FriendsPage] VK Bridge ready:', typeof bridge !== 'undefined');
      console.log('[FriendsPage] ========================');

      // Пробуем VKWebAppShowInviteBox - работает на мобильных приложениях VK
      // requestKey передаётся как vk_request_key в Launch Params
      const result = await bridge.send('VKWebAppShowInviteBox', {
        requestKey: data.referralCode,
      } as any);

      console.log('[FriendsPage] ✅ Invite success:', result);

      if (childProfileId) {
        analytics.trackEvent('referral_invite_sent', {
          child_profile_id: childProfileId,
          referral_code: data.referralCode,
          share_channel: 'vk_invite_box',
        });

        await referralAPI.trackInviteSent(childProfileId, 'vk');
      }
    } catch (error: any) {
      console.error('[FriendsPage] ❌ VKWebAppShowInviteBox failed:', {
        error,
        errorType: error?.error_type,
        errorCode: error?.error_data?.error_code,
        errorReason: error?.error_data?.error_reason,
        errorMessage: error?.message,
        fullError: JSON.stringify(error, null, 2),
      });

      // Если пользователь отменил диалог
      if (error?.error_data?.error_code === 4 && error?.error_data?.error_reason === 'User denied') {
        console.log('[FriendsPage] User cancelled invite dialog');
        return;
      }

      // Fallback: Показываем UI с реферальным кодом для ручной отправки
      console.log('[FriendsPage] Showing manual code sharing UI');
      setShowCopyFallback(true);
    }
  };

  const handleCopyCode = async () => {
    if (!data) return;

    const code = data.referralCode;

    // Пробуем современный Clipboard API
    if (navigator.clipboard && window.isSecureContext) {
      try {
        await navigator.clipboard.writeText(code);
        alert('Код скопирован! Отправь его другу в ВК.\n\nДруг должен будет ввести код при регистрации.');

        if (childProfileId) {
          analytics.trackEvent('referral_code_copied', {
            child_profile_id: childProfileId,
            referral_code: code,
          });
        }
        return;
      } catch (err) {
        console.warn('[FriendsPage] Clipboard API failed, trying fallback:', err);
      }
    }

    // Fallback: Старый метод через document.execCommand
    try {
      const textArea = document.createElement('textarea');
      textArea.value = code;
      textArea.style.position = 'fixed';
      textArea.style.left = '-999999px';
      textArea.style.top = '-999999px';
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();

      const successful = document.execCommand('copy');
      document.body.removeChild(textArea);

      if (successful) {
        alert('Код скопирован! Отправь его другу в ВК.\n\nДруг должен будет ввести код при регистрации.');

        if (childProfileId) {
          analytics.trackEvent('referral_code_copied', {
            child_profile_id: childProfileId,
            referral_code: code,
          });
        }
      } else {
        alert('Не удалось скопировать автоматически.\n\nТвой реферальный код: ' + code + '\n\nСкопируй его вручную.');
      }
    } catch (err) {
      console.error('[FriendsPage] All copy methods failed:', err);
      alert('Не удалось скопировать автоматически.\n\nТвой реферальный код: ' + code + '\n\nСкопируй его вручную.');
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
          {showCopyFallback
            ? 'Скопируй код и отправь другу в ВК. Друг должен ввести его при регистрации.'
            : 'Нажми кнопку ниже, выбери друзей из списка и получите оба бонусные стикеры!'}
        </p>

        {!showCopyFallback ? (
          <button onClick={handleInvite} className={styles.inviteButton}>
            <Send size={20} />
            <span>Пригласить друга</span>
          </button>
        ) : (
          <>
            <div className={styles.codeDisplay}>
              <span className={styles.codeLabel}>Твой реферальный код:</span>
              <div className={styles.codeValue}>{data.referralCode}</div>
            </div>
            <button onClick={handleCopyCode} className={styles.inviteButton}>
              <Send size={20} />
              <span>Скопировать код</span>
            </button>
          </>
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
