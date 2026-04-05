// src/components/ui/Card/Card.tsx
import { Card as VKUICard, CardProps as VKUICardProps } from '@vkontakte/vkui';
import { ReactNode } from 'react';
import styles from './Card.module.css';

export interface CardProps extends Omit<VKUICardProps, 'getRootRef'> {
  children: ReactNode;
  variant?: 'default' | 'bordered';
  className?: string;
  onClick?: () => void;
}

/**
 * Wrapper над VKUI Card
 */
export function Card({
  children,
  variant = 'default',
  className,
  onClick,
  ...props
}: CardProps) {
  const classes = [
    styles.card,
    variant === 'bordered' && styles.bordered,
    onClick && styles.clickable,
    className,
  ]
    .filter(Boolean)
    .join(' ');

  return (
    <VKUICard mode="shadow" className={classes} onClick={onClick} {...props}>
      {children}
    </VKUICard>
  );
}
