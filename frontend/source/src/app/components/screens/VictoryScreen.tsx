import { useNavigate } from "react-router";
import { motion } from "motion/react";
import { Mascot } from "../Mascot";
import { Villain } from "../Villain";

export function VictoryScreen() {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#FFF9E8] to-[#E8FFF8]">
      <div className="flex-1 flex flex-col items-center justify-center text-center">
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: "spring", duration: 0.6 }}
          className="text-[80px] mb-2"
        >
          🎉
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <h1 className="text-primary mb-2">Победа!</h1>
          <p className="text-muted-foreground text-[14px] mb-6">
            Ты победил злодея Кракозябру!
          </p>
        </motion.div>

        {/* Characters */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
          className="flex items-end gap-6 mb-6"
        >
          <Mascot size="sm" message="Мы победили!" />
          <Villain size="sm" defeated />
        </motion.div>

        {/* Reward */}
        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: 0.7 }}
          className="bg-white rounded-3xl p-6 shadow-sm w-full max-w-[280px]"
        >
          <p className="text-[13px] text-muted-foreground mb-3">Твоя награда</p>
          <div className="flex items-center justify-center gap-4">
            <div className="text-center">
              <div className="text-[40px]">⭐</div>
              <p className="text-[11px] text-foreground mt-1">Редкий стикер</p>
            </div>
            <div className="text-center">
              <div className="text-[40px]">🏆</div>
              <p className="text-[11px] text-foreground mt-1">Достижение</p>
            </div>
          </div>
        </motion.div>
      </div>

      <button
        onClick={() => navigate("/")}
        className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform"
      >
        Продолжить
      </button>
    </div>
  );
}
