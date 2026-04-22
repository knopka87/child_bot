// src/pages/Home/components/RecentAttempts.tsx
import { CheckCircle, Clock, XCircle } from 'lucide-react';
import { motion } from 'framer-motion';
import type { RecentAttempt } from '@/types/home';

interface RecentAttemptsProps {
  attempts: RecentAttempt[];
  onAttemptClick: (attempt: RecentAttempt) => void;
}

const statusConfig = {
  success: {
    label: 'Решено верно',
    icon: CheckCircle,
    color: 'text-[#00B894]',
    bg: 'bg-[#E8FFF8]',
  },
  in_progress: {
    label: 'Почти верно',
    icon: Clock,
    color: 'text-[#FDCB6E]',
    bg: 'bg-[#FFF9E8]',
  },
  error: {
    label: 'Решено неверно',
    icon: XCircle,
    color: 'text-[#FF6B6B]',
    bg: 'bg-red-50',
  },
};

export function RecentAttempts({ attempts, onAttemptClick }: RecentAttemptsProps) {
  if (!attempts || attempts.length === 0) return null;

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.4 }}
    >
      <p className="text-[13px] text-[#636e72] mb-2">Последние задания</p>
      <div className="flex flex-col gap-2">
        {attempts.map((item) => {
          const status = statusConfig[item.status] || statusConfig.error;
          const Icon = status.icon;

          return (
            <button
              key={item.id}
              onClick={() => onAttemptClick(item)}
              className="bg-white rounded-xl p-3 shadow-sm flex items-center justify-between text-left active:scale-[0.98] transition-transform"
            >
              <div className="flex-1 min-w-0">
                <p className="text-[13px] text-[#2D3436] truncate">
                  {item.mode === 'help' ? '💡 Помощь' : '✅ Проверка'} —{' '}
                  {item.resultSummary || 'Задание обработано'}
                </p>
                <div className="flex items-center gap-1.5 mt-1">
                  <span
                    className={`flex items-center gap-1 ${status.color} ${status.bg} px-2 py-0.5 rounded-full text-[10px]`}
                  >
                    <Icon size={12} /> {status.label}
                  </span>
                  <span className="text-[10px] text-[#636e72]">
                    {formatDate(item.createdAt)}
                  </span>
                </div>
              </div>
            </button>
          );
        })}
      </div>
    </motion.div>
  );
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays === 0) return 'Сегодня';
  if (diffDays === 1) return 'Вчера';
  return `${diffDays} дн. назад`;
}
