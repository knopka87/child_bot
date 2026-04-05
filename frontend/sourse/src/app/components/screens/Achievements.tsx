import { useState } from "react";
import { Trophy, Star, Flame, CheckCircle, Target, Zap } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";

const achievements = [
  { id: 1, icon: <Flame size={28} />, title: "5 дней подряд", desc: "Занимался 5 дней подряд", unlocked: true, color: "from-[#FF6B6B] to-[#FD79A8]" },
  { id: 2, icon: <CheckCircle size={28} />, title: "10 проверок ДЗ", desc: "Проверил 10 домашних заданий", unlocked: true, color: "from-[#00B894] to-[#55EFC4]" },
  { id: 3, icon: <Star size={28} />, title: "5 ошибок исправлено", desc: "Исправил 5 ошибок после проверки", unlocked: true, color: "from-[#FDCB6E] to-[#F9CA24]" },
  { id: 4, icon: <Target size={28} />, title: "Первое задание", desc: "Загрузи первое задание", unlocked: false, color: "from-gray-300 to-gray-400" },
  { id: 5, icon: <Zap size={28} />, title: "Скоростной решатель", desc: "Реши задание за 1 минуту", unlocked: false, color: "from-gray-300 to-gray-400" },
  { id: 6, icon: <Trophy size={28} />, title: "Победитель злодеев", desc: "Победи 3 злодеев", unlocked: false, color: "from-gray-300 to-gray-400" },
];

export function AchievementsScreen() {
  const [showAll, setShowAll] = useState(false);
  const [showPopup, setShowPopup] = useState(false);

  const displayedAchievements = showAll ? achievements : achievements.slice(0, 4);

  return (
    <div className="flex flex-col min-h-full px-5 pt-8 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background relative">
      <h1 className="text-primary mb-1">Достижения</h1>
      <p className="text-muted-foreground text-[14px] mb-6">
        {achievements.filter((a) => a.unlocked).length} из {achievements.length} получено
      </p>

      {/* Achievement cards */}
      <div className="grid grid-cols-2 gap-3 mb-6">
        {displayedAchievements.map((ach, i) => (
          <motion.button
            key={ach.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.05 }}
            onClick={() => ach.unlocked && setShowPopup(true)}
            className={`rounded-2xl p-4 flex flex-col items-center gap-2 text-center shadow-sm transition-transform active:scale-[0.97] ${
              ach.unlocked ? "bg-white" : "bg-muted/50 opacity-60"
            }`}
          >
            <div
              className={`w-14 h-14 rounded-2xl flex items-center justify-center text-white bg-gradient-to-br ${ach.color}`}
            >
              {ach.icon}
            </div>
            <p className="text-[13px] text-foreground">{ach.title}</p>
            <p className="text-[11px] text-muted-foreground">{ach.desc}</p>
            {ach.unlocked && (
              <span className="text-[11px] text-[#00B894] bg-[#E8FFF8] px-2 py-0.5 rounded-full">Получено</span>
            )}
          </motion.button>
        ))}
      </div>

      {!showAll && (
        <button
          onClick={() => setShowAll(true)}
          className="w-full py-3 border border-primary text-primary rounded-2xl mb-4"
        >
          Посмотреть всё
        </button>
      )}

      {/* Stats */}
      <div className="bg-white rounded-2xl p-4 shadow-sm mb-4">
        <div className="flex gap-4">
          <div className="flex-1 text-center">
            <p className="text-primary text-[20px]">12</p>
            <p className="text-[11px] text-muted-foreground">Заданий</p>
          </div>
          <div className="flex-1 text-center">
            <p className="text-[#00B894] text-[20px]">3</p>
            <p className="text-[11px] text-muted-foreground">Злодеев</p>
          </div>
          <div className="flex-1 text-center">
            <p className="text-[#FDCB6E] text-[20px]">5</p>
            <p className="text-[11px] text-muted-foreground">Дней подряд</p>
          </div>
        </div>
      </div>

      {/* Achievement popup */}
      <AnimatePresence>
        {showPopup && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 px-6"
            onClick={() => setShowPopup(false)}
          >
            <motion.div
              initial={{ scale: 0.8, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.8, opacity: 0 }}
              className="bg-white rounded-3xl p-8 w-full max-w-[340px] text-center"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="text-[56px] mb-3">🏆</div>
              <h2 className="text-primary mb-2">Достижение получено!</h2>
              <p className="text-muted-foreground text-[14px] mb-4">
                Ты молодец! Продолжай в том же духе!
              </p>
              <button
                onClick={() => setShowPopup(false)}
                className="w-full py-3 bg-primary text-white rounded-2xl"
              >
                Отлично!
              </button>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}