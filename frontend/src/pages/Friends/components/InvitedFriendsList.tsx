// src/pages/Friends/components/InvitedFriendsList.tsx
import { Card } from '@/components/ui/Card';
import type { InvitedFriend } from '@/types/referral';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
import styles from './InvitedFriendsList.module.css';

interface InvitedFriendsListProps {
  friends: InvitedFriend[];
}

export function InvitedFriendsList({ friends }: InvitedFriendsListProps) {
  if (friends.length === 0) {
    return (
      <Card className={styles.emptyCard}>
        <div className={styles.emptyIcon}>👥</div>
        <p className={styles.emptyText}>Пока нет приглашённых друзей</p>
        <p className={styles.emptySubtext}>
          Поделись своей ссылкой, чтобы начать
        </p>
      </Card>
    );
  }

  return (
    <div className={styles.container}>
      <h3 className={styles.title}>Приглашённые друзья ({friends.length})</h3>

      <div className={styles.list}>
        {friends.map((friend) => (
          <Card key={friend.id} className={styles.friendCard}>
            <div className={styles.avatar}>
              {friend.avatarUrl ? (
                <img src={friend.avatarUrl} alt={friend.displayName} />
              ) : (
                <div className={styles.avatarPlaceholder}>
                  {friend.displayName.charAt(0).toUpperCase()}
                </div>
              )}
            </div>

            <div className={styles.info}>
              <div className={styles.name}>{friend.displayName}</div>
              <div className={styles.meta}>
                <span className={styles.status}>
                  {getStatusText(friend.status)}
                </span>
                <span className={styles.time}>
                  {formatDistanceToNow(new Date(friend.invitedAt), {
                    addSuffix: true,
                    locale: ru,
                  })}
                </span>
              </div>
            </div>

            <div className={styles.statusIcon}>
              {friend.status === 'completed_first_task' && '✓'}
              {friend.status === 'active' && '👤'}
              {friend.status === 'pending' && '⏳'}
            </div>
          </Card>
        ))}
      </div>
    </div>
  );
}

function getStatusText(status: InvitedFriend['status']): string {
  switch (status) {
    case 'completed_first_task':
      return 'Выполнил задание';
    case 'active':
      return 'Активен';
    case 'pending':
      return 'Ожидание';
  }
}
