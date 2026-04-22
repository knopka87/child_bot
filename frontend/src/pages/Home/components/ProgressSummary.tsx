// src/pages/Home/components/ProgressSummary.tsx
import { motion } from 'framer-motion';

interface ProgressSummaryProps {
  tasksTotal: number;
  tasksCorrect: number;
  villainsDefeated: number;
}

export function ProgressSummary({
  tasksTotal,
  tasksCorrect,
  villainsDefeated,
}: ProgressSummaryProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.35 }}
      className="bg-white rounded-2xl p-4 shadow-sm mb-4"
    >
      <p className="text-[13px] text-[#636e72] mb-2">Прогресс</p>
      <div className="flex gap-4">
        <div className="flex-1 text-center">
          <p className="text-[#6C5CE7] text-[20px] font-semibold">{tasksTotal}</p>
          <p className="text-[11px] text-[#636e72]">Заданий</p>
        </div>
        <div className="flex-1 text-center">
          <p className="text-[#00B894] text-[20px] font-semibold">{tasksCorrect}</p>
          <p className="text-[11px] text-[#636e72]">Верно</p>
        </div>
        <div className="flex-1 text-center">
          <p className="text-[#FDCB6E] text-[20px] font-semibold">{villainsDefeated}</p>
          <p className="text-[11px] text-[#636e72]">Злодея побеждено</p>
        </div>
      </div>
    </motion.div>
  );
}
