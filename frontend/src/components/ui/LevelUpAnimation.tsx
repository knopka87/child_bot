// src/components/ui/LevelUpAnimation.tsx
import { motion, AnimatePresence } from 'framer-motion';
import { Star } from 'lucide-react';

interface LevelUpAnimationProps {
  show: boolean;
  newLevel: number;
  coinsReward: number;
  onComplete: () => void;
}

export function LevelUpAnimation({ show, newLevel, coinsReward, onComplete }: LevelUpAnimationProps) {
  return (
    <AnimatePresence>
      {show && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.5 }}
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/60"
          onClick={onComplete}
        >
          <motion.div
            initial={{ scale: 0, rotate: -180 }}
            animate={{ scale: 1, rotate: 0 }}
            exit={{ scale: 0, rotate: 180 }}
            transition={{ 
              type: 'spring',
              damping: 12,
              stiffness: 200,
              duration: 0.6 
            }}
            className="bg-gradient-to-br from-[#6C5CE7] to-[#A29BFE] rounded-3xl p-8 shadow-2xl max-w-[300px] w-full mx-4 text-center"
            onClick={(e) => e.stopPropagation()}
          >
            {/* Animated stars */}
            <div className="relative mb-4">
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
                className="absolute -top-4 -left-4"
              >
                <Star size={24} className="text-[#FDCB6E]" fill="#FDCB6E" />
              </motion.div>
              <motion.div
                animate={{ rotate: -360 }}
                transition={{ duration: 2, repeat: Infinity, ease: 'linear' }}
                className="absolute -top-4 -right-4"
              >
                <Star size={24} className="text-[#FDCB6E]" fill="#FDCB6E" />
              </motion.div>
              
              {/* Trophy icon */}
              <div className="text-7xl mb-2">🏆</div>
            </div>

            {/* Level up text */}
            <motion.h2
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.3 }}
              className="text-white text-3xl font-bold mb-2"
            >
              Уровень {newLevel}!
            </motion.h2>

            <motion.p
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.4 }}
              className="text-white/90 text-lg mb-4"
            >
              Поздравляем! 🎉
            </motion.p>

            {/* Coins reward */}
            <motion.div
              initial={{ y: 20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ delay: 0.5 }}
              className="bg-white/20 rounded-2xl px-6 py-3 inline-flex items-center gap-2"
            >
              <span className="text-3xl">🪙</span>
              <span className="text-white text-xl font-bold">+{coinsReward} монет</span>
            </motion.div>

            {/* Tap to continue */}
            <motion.p
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.8 }}
              className="text-white/70 text-sm mt-4"
            >
              Нажми, чтобы продолжить
            </motion.p>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
