// src/components/ui/Button/Button.tsx
import { Button as VKUIButton, ButtonProps as VKUIButtonProps } from '@vkontakte/vkui';
import { useHaptics } from '@/lib/platform/haptics';
import { ReactNode } from 'react';

export interface ButtonProps extends Omit<VKUIButtonProps, 'getRootRef'> {
  children: ReactNode;
  enableHaptics?: boolean;
}

/**
 * Wrapper над VKUI Button с haptic feedback
 */
export function Button({
  children,
  onClick,
  enableHaptics = true,
  ...props
}: ButtonProps) {
  const { onButtonClick } = useHaptics();

  const handleClick = (e: React.MouseEvent<HTMLElement>) => {
    if (enableHaptics) {
      onButtonClick();
    }
    onClick?.(e);
  };

  return (
    <VKUIButton onClick={handleClick} {...props}>
      {children}
    </VKUIButton>
  );
}

// Convenience wrappers
export function PrimaryButton(props: ButtonProps) {
  return <Button mode="primary" size="l" stretched {...props} />;
}

export function SecondaryButton(props: ButtonProps) {
  return <Button mode="secondary" size="l" stretched {...props} />;
}

export function OutlineButton(props: ButtonProps) {
  return <Button mode="outline" size="l" stretched {...props} />;
}
