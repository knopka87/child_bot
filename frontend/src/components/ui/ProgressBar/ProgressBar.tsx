// src/components/ui/ProgressBar/ProgressBar.tsx
import { Progress, ProgressProps as VKUIProgressProps } from '@vkontakte/vkui';
import styles from './ProgressBar.module.css';

export interface ProgressBarProps extends Omit<VKUIProgressProps, 'getRootRef'> {
  value: number;
  variant?: 'default' | 'error' | 'success';
  size?: 'sm' | 'md' | 'lg';
  showLabel?: boolean;
  label?: string;
}

/**
 * Wrapper над VKUI Progress
 */
export function ProgressBar({
  value,
  variant = 'default',
  size = 'md',
  showLabel = false,
  label,
  ...props
}: ProgressBarProps) {
  const classes = [
    styles.progress,
    styles[`variant-${variant}`],
    styles[`size-${size}`],
  ].join(' ');

  return (
    <div className={styles.container}>
      <Progress value={value} className={classes} {...props} />
      {showLabel && label && <span className={styles.label}>{label}</span>}
    </div>
  );
}
