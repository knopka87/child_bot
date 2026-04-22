// src/pages/Home/components/UnfinishedAttemptModal.tsx
import { motion, AnimatePresence } from 'framer-motion';

interface UnfinishedAttemptModalProps {
  isOpen: boolean;
  onClose: () => void;
  onContinue: () => void;
  onNewTask: () => void;
}

export function UnfinishedAttemptModal({
  isOpen,
  onClose,
  onContinue,
  onNewTask,
}: UnfinishedAttemptModalProps) {
  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-6"
          onClick={onClose}
        >
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.9, opacity: 0 }}
            className="bg-white rounded-3xl p-6 w-full max-w-[340px] text-center"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="text-[48px] mb-2">📝</div>
            <h2 className="text-[#6C5CE7] text-[20px] font-semibold mb-2">
              Незаконченное задание
            </h2>
            <p className="text-[#636e72] text-[14px] mb-5">
              У тебя есть незаконченное задание. Хочешь продолжить?
            </p>
            <div className="flex flex-col gap-3">
              <button
                onClick={onContinue}
                className="w-full py-3 bg-[#6C5CE7] text-white rounded-2xl font-medium active:scale-[0.98] transition-transform"
              >
                Продолжить
              </button>
              <button
                onClick={onNewTask}
                className="w-full py-3 border-2 border-[#6C5CE7] text-[#6C5CE7] rounded-2xl font-medium active:scale-[0.98] transition-transform"
              >
                Новое задание
              </button>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
