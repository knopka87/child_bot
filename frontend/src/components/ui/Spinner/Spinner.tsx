// src/components/ui/Spinner/Spinner.tsx
import { Spinner as VKUISpinner, SpinnerProps as VKUISpinnerProps } from '@vkontakte/vkui';

export interface SpinnerProps extends Omit<VKUISpinnerProps, 'size'> {
  size?: 'sm' | 'md' | 'lg';
}

/**
 * Wrapper над VKUI Spinner
 */
export function Spinner({ size = 'md', ...props }: SpinnerProps) {
  const sizeMap = {
    sm: 'small',
    md: 'medium',
    lg: 'large',
  } as const;

  return <VKUISpinner size={sizeMap[size]} {...props} />;
}
