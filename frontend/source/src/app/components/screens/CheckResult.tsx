import { useState } from "react";
import { useNavigate } from "react-router";
import { ArrowLeft, CheckCircle, AlertTriangle } from "lucide-react";
import { motion } from "motion/react";
import { Mascot } from "../Mascot";
import { Villain } from "../Villain";

type Verdict = "correct" | "errors" | "review";

export function CheckResult() {
  const navigate = useNavigate();
  const [verdict] = useState<Verdict>(() => {
    const v: Verdict[] = ["correct", "errors", "review"];
    return v[Math.floor(Math.random() * 3)];
  });
  const [villainHealth] = useState(2);

  if (verdict === "correct") {
    return (
      <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#E8FFF8] to-background">
        <button onClick={() => navigate("/")} className="flex items-center gap-2 text-primary mb-4">
          <ArrowLeft size={20} /><span className="text-[14px]">На главную</span>
        </button>

        <div className="flex-1 flex flex-col items-center justify-center text-center">
          <motion.div initial={{ scale: 0 }} animate={{ scale: 1 }} transition={{ type: "spring" }}>
            <CheckCircle size={64} className="text-[#00B894] mx-auto mb-4" />
          </motion.div>
          <h1 className="text-[#00B894] mb-2">Верно!</h1>
          <Mascot size="sm" message="Молодец!" className="mb-4" />

          {/* Villain damage */}
          <div className="bg-white rounded-2xl p-4 shadow-sm w-full mb-4">
            <div className="flex items-center gap-3">
              <Villain size="sm" />
              <div className="flex-1">
                <p className="text-[12px] text-muted-foreground mb-1">Урон нанесён!</p>
                <div className="flex gap-1.5">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className={`flex-1 h-3 rounded-full ${i <= villainHealth - 1 ? "bg-destructive" : "bg-muted"}`} />
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>

        <button
          onClick={() => navigate("/check/scenario")}
          className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform"
        >
          Новое задание
        </button>
      </div>
    );
  }

  if (verdict === "errors") {
    return (
      <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#FFF9E8] to-background">
        <button onClick={() => navigate("/")} className="flex items-center gap-2 text-primary mb-4">
          <ArrowLeft size={20} /><span className="text-[14px]">На главную</span>
        </button>

        <div className="flex-1 flex flex-col">
          <motion.div initial={{ scale: 0.9, opacity: 0 }} animate={{ scale: 1, opacity: 1 }} className="bg-[#FFF9E8] rounded-3xl p-6 text-center mb-4">
            <AlertTriangle size={48} className="text-[#FDCB6E] mx-auto mb-3" />
            <h2 className="text-[#E17055]">Есть ошибки</h2>
          </motion.div>

          <div className="bg-white rounded-2xl p-4 shadow-sm mb-4">
            <Mascot size="sm" message="Не переживай!" className="mb-3" />
            <p className="text-[14px] text-foreground">
              Проверь этот шаг — кажется, ошибка в последнем действии. Попробуй пересчитать!
            </p>
          </div>

          {/* Error highlight */}
          <div className="bg-[#FFF0F0] rounded-2xl p-4 border border-[#FFD0D0] mb-4">
            <p className="text-[12px] text-destructive mb-1">Обрати внимание:</p>
            <p className="text-[14px] text-foreground">
              Шаг 3: проверь вычитание — возможно, ты перепутал числа.
            </p>
          </div>
        </div>

        <div className="flex flex-col gap-3">
          <button
            onClick={() => navigate("/check/upload")}
            className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform"
          >
            Исправил(а) — проверить снова
          </button>
          <button
            onClick={() => navigate("/check/scenario")}
            className="w-full py-3 border border-primary text-primary rounded-2xl active:scale-[0.98] transition-transform"
          >
            Новое задание
          </button>
        </div>
      </div>
    );
  }

  // "review" - посмотри ещё раз
  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#FFF0F0] to-background">
      <button onClick={() => navigate("/")} className="flex items-center gap-2 text-primary mb-4">
        <ArrowLeft size={20} /><span className="text-[14px]">На главную</span>
      </button>

      <div className="flex-1 flex flex-col items-center justify-center text-center">
        <motion.div initial={{ scale: 0.9, opacity: 0 }} animate={{ scale: 1, opacity: 1 }}>
          <AlertTriangle size={48} className="text-[#FDCB6E] mx-auto mb-3" />
          <h2 className="text-[#E17055] mb-2">Посмотри ещё раз</h2>
        </motion.div>

        <div className="bg-white rounded-2xl p-4 shadow-sm w-full mb-4">
          <p className="text-[14px] text-foreground">
            Ничего страшного! Ошибки — это нормально. Попробуй ещё раз, ты справишься!
          </p>
        </div>

        <Mascot size="sm" message="Не сдавайся!" />
      </div>

      <div className="flex flex-col gap-3">
        <button
          onClick={() => navigate("/check/upload")}
          className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform"
        >
          Попробовать ещё раз
        </button>
        <button
          onClick={() => navigate("/check/scenario")}
          className="w-full py-3 border border-primary text-primary rounded-2xl active:scale-[0.98] transition-transform"
        >
          Новое задание
        </button>
      </div>
    </div>
  );
}