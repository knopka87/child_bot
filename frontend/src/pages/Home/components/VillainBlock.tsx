// src/pages/Home/components/VillainBlock.tsx
import { ChevronRight } from 'lucide-react';
import { motion } from 'framer-motion';
import type { Villain } from '@/types/domain';

interface VillainBlockProps {
  villain: Villain;
  onVillainClick: () => void;
}

export function VillainBlock({ villain, onVillainClick }: VillainBlockProps) {
  const healthPercent = Math.max(0, villain.healthPercent);

  return (
    <motion.button
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.25 }}
      onClick={onVillainClick}
      className="bg-white rounded-2xl p-4 flex items-center gap-3 shadow-sm mb-4 active:scale-[0.98] transition-transform text-left w-full"
    >
      <span className="text-[36px]">👾</span>
      <div className="flex-1">
        <p className="text-[13px] text-[#2D3436]">{villain.name}</p>
        <p className="text-[11px] text-[#636e72] mb-1.5">«Ты меня не победишь!»</p>
        <div className="flex items-center gap-2">
          <div className="flex-1">
            <div className="w-full h-2.5 bg-[#DFE6E9] rounded-full overflow-hidden">
              <div 
                className="h-full bg-[#FF6B6B] transition-all duration-500 ease-out rounded-full"
                style={{ width: `${healthPercent}%` }}
              />
            </div>
          </div>
          <p className="text-[10px] text-[#636e72] font-medium whitespace-nowrap">
            {Math.round(healthPercent)}%
          </p>
        </div>
      </div>
      <ChevronRight size={18} className="text-[#636e72]" />
    </motion.button>
  );
}
