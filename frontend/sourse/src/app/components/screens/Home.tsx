import { useState } from "react";
import { useNavigate } from "react-router";
import { BookOpen, ClipboardCheck, ChevronRight, Clock, CheckCircle, XCircle } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";
import { Mascot } from "../Mascot";

const jokes = [
  "Почему учебник грустит? Потому что у него слишком много проблем!",
  "Что сказал ноль восьмёрке? Красивый пояс!",
  "Какой предмет самый вкусный? ИЗО-бражение торта!",
  "Почему карандаш пошёл в школу? Чтобы стать острее!",
];

const recentAttempts = [
  { id: 1, title: "Математика — задача про яблоки", status: "correct" as const, date: "Сегодня" },
  { id: 2, title: "Русский язык — упражнение 45", status: "almost" as const, date: "Вчера" },
  { id: 3, title: "Математика — примеры на сложение", status: "incorrect" as const, date: "Вчера" },
];

const statusConfig = {
  correct: { label: "Решено верно", icon: <CheckCircle size={14} />, color: "text-[#00B894]", bg: "bg-[#E8FFF8]" },
  almost: { label: "Почти верно", icon: <Clock size={14} />, color: "text-[#FDCB6E]", bg: "bg-[#FFF9E8]" },
  incorrect: { label: "Решено неверно", icon: <XCircle size={14} />, color: "text-destructive", bg: "bg-red-50" },
};

export function HomeScreen() {
  const navigate = useNavigate();
  const [showUnfinishedModal, setShowUnfinishedModal] = useState(false);
  const [pendingMode, setPendingMode] = useState<"help" | "check" | null>(null);
  const joke = jokes[Math.floor(Math.random() * jokes.length)];

  // Simulated unfinished session
  const hasUnfinished = true;
  const villainHealth = 2; // out of 3

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
    <div className="flex flex-col min-h-full px-5 pt-8 pb-4 bg-gradient-to-b from-[#F0F4FF] to-background">
      {/* Header */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="mb-5"
      >
        <h1 className="text-primary">Привет! 👋</h1>
        <p className="text-muted-foreground text-[14px]">Чем займёмся сегодня?</p>
      </motion.div>

      {/* Main CTA buttons */}
      <div className="flex flex-col gap-3 mb-5">
        <motion.button
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.1 }}
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
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.2 }}
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

      {/* Villain block */}
      <motion.button
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.25 }}
        onClick={() => navigate("/villain")}
        className="bg-white rounded-2xl p-4 flex items-center gap-3 shadow-sm mb-4 active:scale-[0.98] transition-transform text-left"
      >
        <span className="text-[36px]">👾</span>
        <div className="flex-1">
          <p className="text-[13px] text-foreground">Злодей Кракозябра</p>
          <p className="text-[11px] text-muted-foreground mb-1.5">«Ты меня не победишь!»</p>
          <div className="flex gap-1.5">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className={`flex-1 h-2.5 rounded-full transition-all ${
                  i <= villainHealth ? "bg-destructive" : "bg-muted"
                }`}
              />
            ))}
          </div>
        </div>
        <ChevronRight size={18} className="text-muted-foreground" />
      </motion.button>

      {/* Mascot joke */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.3 }}
        className="bg-[#F0F4FF] rounded-2xl p-3 flex gap-3 items-center mb-4"
      >
        <Mascot size="sm" className="flex-shrink-0" />
        <p className="text-[12px] text-[#2D3436]">{joke}</p>
      </motion.div>

      {/* Progress summary */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.35 }}
        className="bg-white rounded-2xl p-4 shadow-sm mb-4"
      >
        <p className="text-[13px] text-muted-foreground mb-2">Прогресс</p>
        <div className="flex gap-4">
          <div className="flex-1 text-center">
            <p className="text-primary text-[20px]">12</p>
            <p className="text-[11px] text-muted-foreground">Заданий</p>
          </div>
          <div className="flex-1 text-center">
            <p className="text-[#00B894] text-[20px]">8</p>
            <p className="text-[11px] text-muted-foreground">Верно</p>
          </div>
          <div className="flex-1 text-center">
            <p className="text-[#FDCB6E] text-[20px]">3</p>
            <p className="text-[11px] text-muted-foreground">Злодея побеждено</p>
          </div>
        </div>
      </motion.div>

      {/* Recent attempts */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
      >
        <p className="text-[13px] text-muted-foreground mb-2">Последние задания</p>
        <div className="flex flex-col gap-2">
          {recentAttempts.map((item) => {
            const s = statusConfig[item.status];
            return (
              <div key={item.id} className="bg-white rounded-xl p-3 shadow-sm flex items-center justify-between">
                <div className="flex-1 min-w-0">
                  <p className="text-[13px] text-foreground truncate">{item.title}</p>
                  <div className="flex items-center gap-1.5 mt-1">
                    <span className={`flex items-center gap-1 ${s.color} ${s.bg} px-2 py-0.5 rounded-full text-[10px]`}>
                      {s.icon} {s.label}
                    </span>
                    <span className="text-[10px] text-muted-foreground">{item.date}</span>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </motion.div>

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
                <button
                  onClick={handleContinue}
                  className="w-full py-3 bg-primary text-white rounded-2xl"
                >
                  Продолжить
                </button>
                <button
                  onClick={handleNewTask}
                  className="w-full py-3 border border-primary text-primary rounded-2xl"
                >
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
