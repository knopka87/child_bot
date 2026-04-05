// src/components/ui/Modal/Modal.tsx
import { ModalRoot, ModalPage, ModalPageHeader } from '@vkontakte/vkui';
import { ReactNode } from 'react';

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
}

/**
 * Wrapper над VKUI Modal
 */
export function Modal({ isOpen, onClose, title, children }: ModalProps) {
  if (!isOpen) return null;

  return (
    <ModalRoot activeModal={isOpen ? 'modal' : null} onClose={onClose}>
      <ModalPage
        id="modal"
        onClose={onClose}
        header={title && <ModalPageHeader>{title}</ModalPageHeader>}
      >
        {children}
      </ModalPage>
    </ModalRoot>
  );
}
