// src/components/ui/Button/PurpleButton.tsx
import { useHaptics } from '@/lib/platform/haptics';
import styles from './PurpleButton.module.css';

interface PurpleButtonProps {
  onClick: () => void;
  children: React.ReactNode;
  disabled?: boolean;
  style?: React.CSSProperties;
}

/**
 * Фиолетовая кнопка как в дизайне
 */
export function PurpleButton({ onClick, children, disabled, style }: PurpleButtonProps) {
  const { onButtonClick } = useHaptics();

  const handleClick = () => {
    onButtonClick();
    onClick();
  };

  return (
    <button
      className={styles.purpleButton}
      onClick={handleClick}
      disabled={disabled}
      style={style}
    >
      {children}
    </button>
  );
}
