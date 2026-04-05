import { useState } from "react";
import { useNavigate } from "react-router";
import { ArrowLeft, ChevronDown, ChevronUp, Lock } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";

const hints = [
  {
    id: "Л1",
    title: "Подсказка 1 — Направление",
    content: "Прочитай задание ещё раз внимательно. Какое действие нужно выполнить? Это сложение или вычитание?",
  },
  {
    id: "Л2",
    title: "Подсказка 2 — Разбор",
    content: "Попробуй разделить задачу на шаги: сначала найди, что дано, затем — что нужно найти.",
  },
  {
    id: "Л3",
    title: "Подсказка 3 — Решение",
    content: "Вычисли: 15 + 23 = 38. Потом вычти: 38 - 10 = 28. Ответ: 28 яблок.",
  },
];

export function HelpResult() {
  const navigate = useNavigate();
  const [unlockedLevel, setUnlockedLevel] = useState(1);
  const [openHints, setOpenHints] = useState<string[]>(["Л1"]);
  const [villainHealth, setVillainHealth] = useState(3);

  const toggleHint = (id: string) => {
    const level = parseInt(id.replace("Л", ""));
    if (level > unlockedLevel) return;
    setOpenHints((prev) =>
      prev.includes(id) ? prev.filter((h) => h !== id) : [...prev, id]
    );
  };

  const unlockNext = () => {
    if (unlockedLevel < 3) {
      const next = unlockedLevel + 1;
      setUnlockedLevel(next);
      setOpenHints((prev) => [...prev, `Л${next}`]);
    }
  };

  const handleSubmitAnswer = () => {
    const newHealth = villainHealth - 1;
    setVillainHealth(newHealth);
    if (newHealth <= 0) {
      navigate("/victory");
    }
  };

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate("/")}
        className="flex items-center gap-2 text-primary mb-4"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">На главную</span>
      </button>

      <h2 className="text-primary mb-2">Подсказки</h2>

      {/* Villain health bar */}
      <div className="bg-white rounded-2xl p-3 mb-4 flex items-center gap-3">
        <span className="text-[20px]">👾</span>
        <div className="flex-1">
          <p className="text-[12px] text-muted-foreground mb-1">Здоровье злодея</p>
          <div className="flex gap-1.5">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className={`flex-1 h-3 rounded-full transition-all ${
                  i <= villainHealth ? "bg-destructive" : "bg-muted"
                }`}
              />
            ))}
          </div>
        </div>
      </div>

      {/* Hint accordion - sequential unlock */}
      <div className="flex flex-col gap-3 mb-6">
        {hints.map((hint) => {
          const level = parseInt(hint.id.replace("Л", ""));
          const isLocked = level > unlockedLevel;
          const isOpen = openHints.includes(hint.id);
          return (
            <div key={hint.id} className={`bg-white rounded-2xl overflow-hidden shadow-sm ${isLocked ? "opacity-60" : ""}`}>
              <button
                onClick={() => toggleHint(hint.id)}
                className="w-full flex items-center justify-between px-4 py-3"
                disabled={isLocked}
              >
                <div className="flex items-center gap-2">
                  <span className="w-7 h-7 bg-primary/10 rounded-lg flex items-center justify-center text-primary text-[12px]">
                    {hint.id}
                  </span>
                  <span className="text-[14px] text-foreground">{hint.title}</span>
                </div>
                {isLocked ? (
                  <Lock size={16} className="text-muted-foreground" />
                ) : isOpen ? (
                  <ChevronUp size={18} className="text-muted-foreground" />
                ) : (
                  <ChevronDown size={18} className="text-muted-foreground" />
                )}
              </button>
              <AnimatePresence>
                {isOpen && !isLocked && (
                  <motion.div
                    initial={{ height: 0, opacity: 0 }}
                    animate={{ height: "auto", opacity: 1 }}
                    exit={{ height: 0, opacity: 0 }}
                    className="overflow-hidden"
                  >
                    <div className="px-4 pb-4 text-[14px] text-muted-foreground border-t border-border/50 pt-3">
                      {hint.content}
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          );
        })}
      </div>

      {unlockedLevel < 3 && (
        <button
          onClick={unlockNext}
          className="w-full py-3 bg-white border border-border text-foreground rounded-2xl mb-3 active:scale-[0.98] transition-transform"
        >
          Открыть следующую подсказку
        </button>
      )}

      {/* Buttons */}
      <div className="flex flex-col gap-3 mt-auto">
        <button
          onClick={handleSubmitAnswer}
          className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform"
        >
          Отправить ответ
        </button>
        <button
          onClick={() => navigate("/help/upload")}
          className="w-full py-3 border border-primary text-primary rounded-2xl active:scale-[0.98] transition-transform"
        >
          Новое задание
        </button>
      </div>
    </div>
  );
}
