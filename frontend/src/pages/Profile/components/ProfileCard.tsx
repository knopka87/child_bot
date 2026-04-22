// src/pages/Profile/components/ProfileCard.tsx
import { Card } from '@/components/ui/Card';
import type { ProfileData } from '@/types/profile';
import styles from './ProfileCard.module.css';

interface ProfileCardProps {
  profile: ProfileData;
}

export function ProfileCard({ profile }: ProfileCardProps) {
  const getSubscriptionBadge = () => {
    switch (profile.subscription.status) {
      case 'trial':
        return {
          text: `Пробный ${profile.subscription.trialDaysRemaining} дн.`,
          color: 'var(--vkui--color_accent_blue)',
        };
      case 'active':
        return {
          text: profile.subscription.planName || 'Активна',
          color: 'var(--vkui--color_text_positive)',
        };
      case 'expired':
        return {
          text: 'Истекла',
          color: 'var(--vkui--color_text_negative)',
        };
      case 'cancelled':
        return {
          text: 'Отменена',
          color: 'var(--vkui--color_text_secondary)',
        };
    }
  };

  const badge = getSubscriptionBadge();

  return (
    <Card className={styles.card}>
      <div className={styles.content}>
        <div className={styles.avatar}>
          {profile.avatarUrl ? (
            <img src={profile.avatarUrl} alt={profile.displayName} />
          ) : (
            <div className={styles.avatarPlaceholder}>
              {profile.displayName.charAt(0).toUpperCase()}
            </div>
          )}
        </div>

        <div className={styles.info}>
          <h2 className={styles.name}>{profile.displayName}</h2>
          <div className={styles.grade}>{profile.grade} класс</div>
          <div className={styles.subscription} style={{ color: badge.color }}>
            {badge.text}
          </div>
        </div>
      </div>
    </Card>
  );
}
