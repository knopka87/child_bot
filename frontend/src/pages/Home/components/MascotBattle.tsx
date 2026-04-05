// src/pages/Home/components/MascotBattle.tsx
import { motion } from 'framer-motion';
import type { Villain } from '@/types/domain';

interface MascotBattleProps {
  villain: Villain | null;
  onVillainClick: () => void;
}

export function MascotBattle({ villain, onVillainClick }: MascotBattleProps) {
  const healthBars = 3;
  // Рассчитываем сколько полосок заполнено на основе healthPercent
  // healthPercent = 66% → 2 полоски из 3
  const filledBars = villain ? Math.round((villain.healthPercent / 100) * healthBars) : 0;

  return (
    <div className="relative w-full h-[320px] flex items-center justify-center">
      {/* Battle Arena Container */}
      <div className="flex items-center justify-center gap-4">
        {/* Left side - Mascot */}
        <div className="relative">
          {/* Speech Bubble above mascot */}
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="absolute -top-16 left-2 bg-white rounded-2xl px-4 py-2 shadow-lg z-10"
            style={{ width: 'fit-content' }}
          >
            <p className="text-[13px] text-[#2D3436] font-medium whitespace-nowrap">
              Мы справимся! 💪
            </p>
            <div className="absolute -bottom-2 left-4 w-3 h-3 bg-white rotate-45 shadow-lg" />
          </motion.div>

          {/* Mascot */}
          <motion.img
            initial={{ x: -50, opacity: 0 }}
            animate={{ x: 0, opacity: 1 }}
            transition={{ delay: 0.2 }}
            src="/images/mascot.png"
            alt="Mascot"
            className="w-28 h-auto object-contain"
          />
        </div>

        {/* Center - Sword */}
        {villain && (
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1, rotate: [0, -10, 10, -10, 0] }}
            transition={{
              delay: 0.4,
              rotate: { repeat: Infinity, duration: 2, ease: 'easeInOut' },
            }}
            className="text-3xl self-center -mb-6"
            style={{ filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.15))' }}
          >
            ⚔️
          </motion.div>
        )}

        {/* Right side - Villain */}
        {villain && (
          <div className="relative flex flex-col">
            <motion.button
              initial={{ x: 50, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ delay: 0.3 }}
              onClick={onVillainClick}
              className="relative flex flex-col items-center active:scale-95 transition-transform"
            >
              {/* Villain image */}
              <motion.img
                src="/images/villain.png"
                alt={villain.name}
                className="w-44 h-auto object-contain"
                animate={{ y: [0, -6, 0] }}
                transition={{ repeat: Infinity, duration: 2.5, ease: 'easeInOut' }}
              />
            </motion.button>

            {/* Health Bars below villain - 3 divisions: 2 dark red + 1 light red */}
            <div className="flex gap-1.5 justify-center mt-2">
              {Array.from({ length: healthBars }).map((_, i) => {
                let bgColor = 'bg-[#FFB8B8]'; // Светло-красная (пустая)
                if (i < filledBars) {
                  bgColor = 'bg-[#FF6B6B]'; // Ярко-красная (заполненная)
                }
                return (
                  <div
                    key={i}
                    className={`w-10 h-2.5 rounded-full transition-all ${bgColor}`}
                  />
                );
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
