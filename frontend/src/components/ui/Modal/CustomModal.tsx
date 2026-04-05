// src/components/ui/Modal/CustomModal.tsx
import { ReactNode, useEffect } from 'react';
import styles from './CustomModal.module.css';

export interface CustomModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: ReactNode;
}

/**
 * Кастомное модальное окно с блюром и закругленными углами
 */
export function CustomModal({ isOpen, onClose, children }: CustomModalProps) {
  // Блокируем скролл body когда модалка открыта
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        {children}
      </div>
    </div>
  );
}
