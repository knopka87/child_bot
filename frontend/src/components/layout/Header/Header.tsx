// src/components/layout/Header/Header.tsx
import { Trophy, Coins } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { ROUTES } from '@/config/routes';
import { useHaptics } from '@/lib/platform/haptics';

export interface HeaderProps {
  level: number;
  levelProgress: number;
  coins: number;
  tasksCount: number;
  showCoins?: boolean;
  showTasks?: boolean;
  hasNewAchievements?: boolean;
}

export function Header({
  level,
  levelProgress,
  coins,
  tasksCount,
  showCoins = true,
  showTasks = true,
  hasNewAchievements = false,
}: HeaderProps) {
  const navigate = useNavigate();
  const { onButtonClick } = useHaptics();
  return (
    <div className="flex items-center justify-between gap-3 px-4 pt-4 pb-2">
      {/* Level Card */}
      <div className="bg-gradient-to-r from-[#6C5CE7] to-[#A29BFE] rounded-3xl px-5 py-3 flex items-center gap-3 shadow-lg">
        <div className="text-white">
          <div className="text-2xl font-bold leading-none">{level}</div>
          <div className="text-xs opacity-90 mt-0.5">Уровень</div>
        </div>
        <div className="flex-1 min-w-[40px]">
          <div className="h-1.5 bg-white/30 rounded-full overflow-hidden">
            <div
              className="h-full bg-white rounded-full transition-all"
              style={{ width: `${levelProgress}%` }}
            />
          </div>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="flex gap-3">
        {showTasks && (
          <button
            onClick={() => {
              onButtonClick();
              navigate(ROUTES.ACHIEVEMENTS);
            }}
            className="relative bg-white rounded-2xl px-4 py-2.5 flex items-center gap-2 shadow-sm active:scale-[0.98] transition-transform"
          >
            <Trophy size={20} className="text-[#FDCB6E]" />
            <span className="text-[#2D3436] font-semibold">{tasksCount}</span>
            {hasNewAchievements && (
              <div className="absolute -top-0.5 -right-0.5 w-2 h-2 bg-red-500 rounded-full" />
            )}
          </button>
        )}
        {showCoins && (
          <div className="bg-white rounded-2xl px-4 py-2.5 flex items-center gap-2 shadow-sm">
            <Coins size={20} className="text-[#FDCB6E]" />
            <span className="text-[#2D3436] font-semibold">{coins}</span>
          </div>
        )}
      </div>
    </div>
  );
}
