// src/pages/Home/components/ActionButtons.tsx
import { BookOpen, ClipboardCheck } from 'lucide-react';
import { motion } from 'framer-motion';

interface ActionButtonsProps {
  onHelpClick: () => void;
  onCheckClick: () => void;
}

export function ActionButtons({ onHelpClick, onCheckClick }: ActionButtonsProps) {
  return (
    <div className="flex flex-col gap-4">
      <motion.button
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.5 }}
        onClick={onHelpClick}
        className="bg-gradient-to-r from-[#6C5CE7] to-[#A29BFE] text-white rounded-[28px] p-6 flex items-center gap-4 shadow-lg active:scale-[0.98] transition-transform"
      >
        <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center flex-shrink-0">
          <BookOpen size={28} strokeWidth={2.5} />
        </div>
        <div className="text-left flex-1">
          <h3 className="text-white text-[18px] font-bold leading-tight">
            Помоги разобраться
          </h3>
          <p className="text-white/90 text-[14px] mt-1 leading-tight">
            Загрузи фото задания
          </p>
        </div>
      </motion.button>

      <motion.button
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.6 }}
        onClick={onCheckClick}
        className="bg-gradient-to-r from-[#00D9A5] to-[#5DEFC4] text-white rounded-[28px] p-6 flex items-center gap-4 shadow-lg active:scale-[0.98] transition-transform"
      >
        <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center flex-shrink-0">
          <ClipboardCheck size={28} strokeWidth={2.5} />
        </div>
        <div className="text-left flex-1">
          <h3 className="text-white text-[18px] font-bold leading-tight">Проверка ДЗ</h3>
          <p className="text-white/90 text-[14px] mt-1 leading-tight">
            Проверю твою работу
          </p>
        </div>
      </motion.button>
    </div>
  );
}
