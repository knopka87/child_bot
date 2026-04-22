// src/components/layout/Container/Container.tsx
import { ReactNode } from 'react';
import styles from './Container.module.css';

export interface ContainerProps {
  children: ReactNode;
  className?: string;
}

export function Container({ children, className }: ContainerProps) {
  const classes = [styles.container, className].filter(Boolean).join(' ');
  return <div className={classes}>{children}</div>;
}
