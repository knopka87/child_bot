import { useNavigate } from "react-router";
import { ArrowLeft } from "lucide-react";
import { motion } from "motion/react";

export function VillainScreen() {
  const navigate = useNavigate();
  const villainHealth = 2;

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-[#FFE8E8]">
      <button
        onClick={() => navigate("/")}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Назад</span>
      </button>

      <div className="flex-1 flex flex-col items-center justify-center text-center">
        <motion.div
          animate={{ y: [0, -8, 0] }}
          transition={{ repeat: Infinity, duration: 2, ease: "easeInOut" }}
          className="text-[100px] mb-4"
        >
          👾
        </motion.div>

        <h2 className="text-foreground mb-1">Злодей Кракозябра</h2>
        <p className="text-muted-foreground text-[14px] mb-6">
          «Ха-ха! Попробуй-ка реши задачки!»
        </p>

        {/* Health bar */}
        <div className="w-full max-w-[260px] mb-6">
          <p className="text-[12px] text-muted-foreground mb-2">Здоровье злодея</p>
          <div className="flex gap-2">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className={`flex-1 h-4 rounded-full transition-all ${
                  i <= villainHealth ? "bg-destructive" : "bg-muted"
                }`}
              />
            ))}
          </div>
        </div>

        {/* Explanation */}
        <div className="bg-white rounded-2xl p-5 shadow-sm w-full mb-6">
          <h3 className="text-foreground mb-2">Как победить?</h3>
          <div className="flex flex-col gap-2 text-left text-[13px] text-muted-foreground">
            <p>• Решай задания правильно — каждый верный ответ снимает здоровье злодея</p>
            <p>• После 3 верных ответов злодей побеждён!</p>
            <p>• За победу ты получишь редкий стикер и достижение</p>
          </div>
        </div>
      </div>

      <button
        onClick={() => navigate("/help/upload")}
        className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform"
      >
        Продолжить учиться
      </button>
    </div>
  );
}
