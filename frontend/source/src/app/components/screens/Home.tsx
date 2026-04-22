import { useState } from "react";
import { useNavigate } from "react-router";
import { BookOpen, ClipboardCheck, Trophy, Coins } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";
import { Mascot } from "../Mascot";
import { Villain } from "../Villain";

export function HomeScreen() {
  const navigate = useNavigate();
  const [showUnfinishedModal, setShowUnfinishedModal] = useState(false);
  const [pendingMode, setPendingMode] = useState<"help" | "check" | null>(null);

  const hasUnfinished = true;
  const villainHealth = 2;
  const level = 5;
  const levelProgress = 65;
  const tasksSolved = 12;
  const coins = 340;

  const handleModeClick = (mode: "help" | "check") => {
    if (hasUnfinished) {
      setPendingMode(mode);
      setShowUnfinishedModal(true);
    } else {
      navigate(mode === "help" ? "/help/upload" : "/check/scenario");
    }
  };

  const handleContinue = () => {
    setShowUnfinishedModal(false);
    navigate("/help/result");
  };

  const handleNewTask = () => {
    setShowUnfinishedModal(false);
    navigate(pendingMode === "help" ? "/help/upload" : "/check/scenario");
  };

  return (
    <div className="flex flex-col min-h-full px-4 pt-5 pb-4 bg-gradient-to-b from-[#F0F4FF] to-background">
      {/* Game UI top bar */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="flex items-center gap-2 mb-3"
      >
        {/* Level badge */}
        <div className="bg-gradient-to-br from-primary to-secondary rounded-2xl px-3 py-2 flex items-center gap-2 shadow-md shadow-primary/20">
          <div className="w-8 h-8 bg-white/25 rounded-xl flex items-center justify-center">
            <span className="text-white text-[15px]">{level}</span>
          </div>
          <div className="flex flex-col">
            <span className="text-white/70 text-[9px] leading-none">Уровень</span>
            <div className="w-14 h-2 bg-white/25 rounded-full mt-1 overflow-hidden">
              <div className="h-full bg-white rounded-full transition-all" style={{ width: `${levelProgress}%` }} />
            </div>
          </div>
        </div>

        <div className="flex-1" />

        {/* Cups counter */}
        <div className="bg-white rounded-2xl py-2.5 flex items-center gap-1.5 shadow-sm p-[12px]">
          <Trophy size={16} className="text-[#FDCB6E]" />
          <span className="text-[14px] text-foreground">{tasksSolved}</span>
        </div>

        {/* Coins counter */}
        <div className="bg-white rounded-2xl py-2.5 flex items-center gap-1.5 shadow-sm p-[12px]">
          <Coins size={16} className="text-[#FDCB6E]" />
          <span className="text-[14px] text-foreground">{coins}</span>
        </div>
      </motion.div>

      {/* Characters block — 3x enlarged */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="flex items-end justify-between px-2"
      >
        {/* Mascot - left with speech bubble ABOVE */}
        <button
          onClick={() => navigate("/achievements")}
          className="flex flex-col items-center active:scale-[0.95] transition-transform"
        >
          <Mascot size="xl" message="Мы справимся! 💪" />
        </button>

        {/* VS */}
        <div className="flex items-center pb-24">
          <span className="text-[24px]">⚔️</span>
        </div>

        {/* Villain - right */}
        <button
          onClick={() => navigate("/villain")}
          className="flex flex-col items-center active:scale-[0.95] transition-transform"
        >
          <Villain size="2xl" />
          <div className="flex gap-1.5 mt-0.5">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className={`w-7 h-2.5 rounded-full transition-all ${
                  i <= villainHealth ? "bg-destructive" : "bg-muted"
                }`}
              />
            ))}
          </div>
        </button>
      </motion.div>

      {/* Spacer */}
      <div className="flex-1" />

      {/* Main CTA buttons — thumb zone */}
      <div className="flex flex-col gap-3">
        <motion.button
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          onClick={() => handleModeClick("help")}
          className="bg-gradient-to-r from-primary to-secondary text-white rounded-3xl p-5 flex items-center gap-4 shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform"
        >
          <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center flex-shrink-0">
            <BookOpen size={28} />
          </div>
          <div className="text-left">
            <h3 className="text-white">Помоги разобраться</h3>
            <p className="text-white/80 text-[13px] mt-0.5">Загрузи фото задания</p>
          </div>
        </motion.button>

        <motion.button
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.25 }}
          onClick={() => handleModeClick("check")}
          className="bg-gradient-to-r from-[#00B894] to-[#55EFC4] text-white rounded-3xl p-5 flex items-center gap-4 shadow-lg shadow-[#00B894]/20 active:scale-[0.98] transition-transform"
        >
          <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center flex-shrink-0">
            <ClipboardCheck size={28} />
          </div>
          <div className="text-left">
            <h3 className="text-white">Проверка ДЗ</h3>
            <p className="text-white/80 text-[13px] mt-0.5">Проверю твою работу</p>
          </div>
        </motion.button>
      </div>

      {/* Unfinished task modal */}
      <AnimatePresence>
        {showUnfinishedModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-6"
            onClick={() => setShowUnfinishedModal(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              className="bg-white rounded-3xl p-6 w-full max-w-[340px] text-center"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="text-[48px] mb-2">📝</div>
              <h2 className="text-primary mb-2">Незаконченное задание</h2>
              <p className="text-muted-foreground text-[14px] mb-5">
                У тебя есть незаконченное задание. Хочешь продолжить?
              </p>
              <div className="flex flex-col gap-3">
                <button onClick={handleContinue} className="w-full py-3 bg-primary text-white rounded-2xl">
                  Продолжить
                </button>
                <button onClick={handleNewTask} className="w-full py-3 border border-primary text-primary rounded-2xl">
                  Новое задание
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}