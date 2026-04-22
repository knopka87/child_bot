// src/pages/Profile/ProfilePage.tsx
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { History, FileText, HelpCircle, ChevronRight, CreditCard } from 'lucide-react';
import { motion } from 'framer-motion';
import { BottomNav } from '@/components/layout/BottomNav';
import { ProfilePageSkeleton } from '@/components/ui/skeleton';
import { useProfileData } from './hooks/useProfileData';
import { useAnalytics } from '@/hooks/useAnalytics';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { ROUTES } from '@/config/routes';
import styles from './ProfilePage.module.css';

const menuItems = [
  {
    id: 'history',
    path: ROUTES.PROFILE_HISTORY,
    icon: History,
    label: 'История',
    color: styles.colorPrimary,
  },
  {
    id: 'report',
    path: '/profile/report',
    icon: FileText,
    label: 'Отчёт родителю',
    color: styles.colorGreen,
  },
  {
    id: 'help',
    path: '/profile/help',
    icon: HelpCircle,
    label: 'Помощь',
    color: styles.colorBlue,
  },
  {
    id: 'subscription',
    path: '/payment',
    icon: CreditCard,
    label: 'Подписка',
    color: styles.colorOrange,
  },
];

export function ProfilePage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const { data, isLoading, error } = useProfileData();
  const [childProfileId, setChildProfileId] = useState<string | null>(null);

  // Загружаем child_profile_id для аналитики
  useEffect(() => {
    const loadProfileId = async () => {
      const profileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
      setChildProfileId(profileId);

      if (profileId) {
        analytics.trackEvent('profile_opened', {
          child_profile_id: profileId,
        });
      }
    };

    loadProfileId();
  }, [analytics]);

  const handleMenuClick = (id: string, path: string) => {
    if (!childProfileId) return;

    // Track analytics based on menu item ID
    switch (id) {
      case 'history':
        analytics.trackEvent('profile_history_clicked', {
          child_profile_id: childProfileId,
        });
        break;
      case 'report':
        analytics.trackEvent('profile_report_clicked', {
          child_profile_id: childProfileId,
        });
        break;
      case 'help':
        analytics.trackEvent('profile_support_clicked', {
          child_profile_id: childProfileId,
        });
        break;
      case 'subscription':
        analytics.trackEvent('profile_subscription_clicked', {
          child_profile_id: childProfileId,
        });
        break;
    }
    navigate(path);
  };

  const getSubscriptionText = () => {
    if (!data) return '';

    switch (data.subscription.status) {
      case 'trial':
        return `Пробный период — ${data.subscription.trialDaysRemaining} дней`;
      case 'active':
        return data.subscription.planName || 'Активная подписка';
      case 'expired':
        return 'Подписка истекла';
      case 'cancelled':
        return 'Подписка отменена';
      default:
        return '';
    }
  };

  if (isLoading) {
    return <ProfilePageSkeleton />;
  }

  if (error || !data) {
    return (
      <div className={styles.container}>
        <div className={styles.errorContainer}>
          <p>Не удалось загрузить профиль</p>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Profile card */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className={styles.profileCard}
      >
        <div className={styles.avatar}>
          <span className={styles.avatarEmoji}>
            {data.avatarUrl || '🦊'}
          </span>
        </div>
        <div className={styles.profileInfo}>
          <h2 className={styles.profileName}>{data.displayName}</h2>
          <p className={styles.profileGrade}>{data.grade} класс</p>
          <p className={styles.profileSubscription}>
            {getSubscriptionText()}
          </p>
        </div>
      </motion.div>

      {/* Menu items */}
      <div className={styles.menuList}>
        {menuItems.map((item, i) => {
          const Icon = item.icon;
          return (
            <motion.button
              key={item.path}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.05 }}
              onClick={() => handleMenuClick(item.id, item.path)}
              className={styles.menuItem}
            >
              <div className={`${styles.menuIcon} ${item.color}`}>
                <Icon size={20} />
              </div>
              <span className={styles.menuLabel}>{item.label}</span>
              <ChevronRight size={18} className={styles.menuChevron} />
            </motion.button>
          );
        })}
      </div>

      <BottomNav />
    </div>
  );
}
