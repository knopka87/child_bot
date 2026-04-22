import { useState } from "react";
import { motion, AnimatePresence } from "motion/react";
import { ArrowLeft } from "lucide-react";
import { useNavigate } from "react-router";

interface Achievement {
  id: number;
  sticker: string;
  title: string;
  desc: string;
  unlocked: boolean;
  howToGet?: string;
  color: string;
}

const achievements: Achievement[] = [
  { id: 1, sticker: "🔥", title: "5 дней подряд", desc: "Занимался 5 дней подряд", unlocked: true, color: "#FF6B6B" },
  { id: 2, sticker: "✅", title: "10 проверок ДЗ", desc: "Проверил 10 домашних заданий", unlocked: true, color: "#00B894" },
  { id: 3, sticker: "⭐", title: "5 ошибок исправлено", desc: "Исправил 5 ошибок после проверки", unlocked: true, color: "#FDCB6E" },
  { id: 4, sticker: "🎯", title: "Первое задание", desc: "За первое решённое задание", unlocked: false, howToGet: "Загрузи и реши своё первое задание", color: "#E17055" },
  { id: 5, sticker: "⚡", title: "Скоростной решатель", desc: "За быструю работу", unlocked: false, howToGet: "Реши задание менее чем за 60 секунд", color: "#A29BFE" },
  { id: 6, sticker: "🏆", title: "Победитель злодеев", desc: "За победу над злодеями", unlocked: false, howToGet: "Победи 3 злодеев подряд, решая задания верно", color: "#FDCB6E" },
  { id: 7, sticker: "🦉", title: "Мудрая сова", desc: "За использование подсказок", unlocked: false, howToGet: "Открой все подсказки в 5 заданиях", color: "#6C5CE7" },
  { id: 8, sticker: "💎", title: "Коллекционер", desc: "За сбор стикеров", unlocked: false, howToGet: "Получи 10 любых стикеров", color: "#00CEC9" },
  { id: 9, sticker: "🚀", title: "Ракета знаний", desc: "За 20 решённых заданий", unlocked: false, howToGet: "Реши 20 заданий в любом режиме", color: "#E84393" },
  { id: 10, sticker: "🌟", title: "Суперзвезда", desc: "За 10 безошибочных проверок", unlocked: false, howToGet: "Получи 10 проверок без единой ошибки", color: "#F9CA24" },
  { id: 11, sticker: "🎪", title: "Марафонец", desc: "За 7 дней подряд", unlocked: false, howToGet: "Занимайся 7 дней подряд без пропусков", color: "#FF7675" },
  { id: 12, sticker: "🧠", title: "Гений", desc: "За 50 решённых заданий", unlocked: false, howToGet: "Реши 50 заданий — ты настоящий гений!", color: "#74B9FF" },
];

const shelfRows = [
  achievements.slice(0, 4),
  achievements.slice(4, 8),
  achievements.slice(8, 12),
];

export function AchievementsScreen() {
  const [selectedAch, setSelectedAch] = useState<Achievement | null>(null);
  const navigate = useNavigate();

  const unlockedCount = achievements.filter((a) => a.unlocked).length;

  return (
    <div className="flex flex-col min-h-full bg-gradient-to-b from-[#F5F0FF] via-[#FFF5F9] to-[#F0FFF4] relative">
      {/* Header */}
      <div className="flex items-center gap-3 px-5 pt-6 pb-4">
        <button
          onClick={() => navigate(-1)}
          className="w-10 h-10 bg-white rounded-2xl flex items-center justify-center shadow-sm active:scale-95 transition-transform"
        >
          <ArrowLeft size={20} className="text-foreground" />
        </button>
        <div>
          <h1 className="text-foreground">Мои награды</h1>
          <p className="text-muted-foreground text-[13px]">
            Собрано {unlockedCount} из {achievements.length}
          </p>
        </div>
      </div>

      {/* Room background feel */}
      <div className="flex-1 px-4 pb-6 flex flex-col gap-2">
        {shelfRows.map((row, rowIdx) => (
          <div key={rowIdx} className="relative">
            {/* Sticker items on shelf */}
            <div className="grid grid-cols-4 gap-2 px-2 pb-2 pt-3">
              {row.map((ach, i) => (
                <motion.button
                  key={ach.id}
                  initial={{ opacity: 0, scale: 0.8 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ delay: (rowIdx * 4 + i) * 0.05, type: "spring", stiffness: 300 }}
                  onClick={() => setSelectedAch(ach)}
                  className="flex flex-col items-center gap-1.5 active:scale-[0.9] transition-transform"
                >
                  <div
                    className={`w-[68px] h-[68px] rounded-[20px] flex items-center justify-center text-[32px] transition-all ${
                      ach.unlocked
                        ? "shadow-lg"
                        : "grayscale opacity-[0.35]"
                    }`}
                    style={{
                      background: ach.unlocked
                        ? `linear-gradient(135deg, ${ach.color}22, ${ach.color}44)`
                        : "rgba(0,0,0,0.04)",
                      boxShadow: ach.unlocked
                        ? `0 4px 12px ${ach.color}30`
                        : "none",
                    }}
                  >
                    {ach.sticker}
                  </div>
                  <p className={`text-[10px] text-center leading-tight line-clamp-2 w-full ${
                    ach.unlocked ? "text-foreground" : "text-muted-foreground/60"
                  }`}>
                    {ach.title}
                  </p>
                </motion.button>
              ))}
            </div>
            {/* Wooden shelf */}
            <div className="mx-1 h-[10px] rounded-b-lg bg-gradient-to-b from-[#C9B896] to-[#B8A580] shadow-[0_3px_6px_rgba(0,0,0,0.1)]">
              <div className="h-[3px] bg-gradient-to-r from-transparent via-white/30 to-transparent rounded-full mx-2 mt-[2px]" />
            </div>
          </div>
        ))}
      </div>

      {/* Achievement detail modal */}
      <AnimatePresence>
        {selectedAch && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/40 backdrop-blur-sm flex items-center justify-center z-50 px-6"
            onClick={() => setSelectedAch(null)}
          >
            <motion.div
              initial={{ scale: 0.85, opacity: 0, y: 20 }}
              animate={{ scale: 1, opacity: 1, y: 0 }}
              exit={{ scale: 0.85, opacity: 0, y: 20 }}
              transition={{ type: "spring", stiffness: 350, damping: 25 }}
              className="bg-white rounded-[28px] p-7 w-full max-w-[320px] text-center shadow-2xl"
              onClick={(e) => e.stopPropagation()}
            >
              <div
                className={`w-24 h-24 rounded-[24px] flex items-center justify-center text-[48px] mx-auto mb-4 ${
                  selectedAch.unlocked ? "" : "grayscale opacity-50"
                }`}
                style={{
                  background: `linear-gradient(135deg, ${selectedAch.color}22, ${selectedAch.color}44)`,
                  boxShadow: selectedAch.unlocked ? `0 6px 20px ${selectedAch.color}30` : "none",
                }}
              >
                {selectedAch.sticker}
              </div>
              <h2 className="text-foreground mb-1.5">{selectedAch.title}</h2>
              <p className="text-muted-foreground text-[14px] mb-4">
                {selectedAch.unlocked
                  ? selectedAch.desc
                  : selectedAch.howToGet || "Продолжай учиться, чтобы получить!"}
              </p>
              {selectedAch.unlocked ? (
                <span className="inline-block text-[#00B894] bg-[#E8FFF8] px-4 py-1.5 rounded-full text-[13px] mb-5">
                  ✓ Получено
                </span>
              ) : (
                <span className="inline-block text-muted-foreground bg-muted px-4 py-1.5 rounded-full text-[13px] mb-5">
                  Ещё не получено
                </span>
              )}
              <button
                onClick={() => setSelectedAch(null)}
                className="w-full py-3.5 bg-primary text-white rounded-2xl active:scale-[0.97] transition-transform"
              >
                Понятно!
              </button>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
