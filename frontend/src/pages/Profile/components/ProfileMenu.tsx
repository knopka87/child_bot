// src/pages/Profile/components/ProfileMenu.tsx
import { SimpleCell } from '@vkontakte/vkui';
import {
  Icon24ChevronRight,
  Icon24HistoryBackwardOutline,
  Icon24ReportOutline,
  Icon24HelpOutline,
  Icon24StarsOutline,
} from '@vkontakte/icons';
import styles from './ProfileMenu.module.css';

interface ProfileMenuProps {
  onHistoryClick: () => void;
  onReportClick: () => void;
  onSupportClick: () => void;
  onSubscriptionClick: () => void;
}

export function ProfileMenu({
  onHistoryClick,
  onReportClick,
  onSupportClick,
  onSubscriptionClick,
}: ProfileMenuProps) {
  return (
    <div className={styles.menu}>
      <SimpleCell
        before={<Icon24HistoryBackwardOutline />}
        after={<Icon24ChevronRight />}
        onClick={onHistoryClick}
      >
        История попыток
      </SimpleCell>

      <SimpleCell
        before={<Icon24ReportOutline />}
        after={<Icon24ChevronRight />}
        onClick={onReportClick}
      >
        Отчёт родителю
      </SimpleCell>

      <SimpleCell
        before={<Icon24HelpOutline />}
        after={<Icon24ChevronRight />}
        onClick={onSupportClick}
      >
        Поддержка
      </SimpleCell>

      <SimpleCell
        before={<Icon24StarsOutline />}
        after={<Icon24ChevronRight />}
        onClick={onSubscriptionClick}
      >
        Подписка
      </SimpleCell>
    </div>
  );
}
