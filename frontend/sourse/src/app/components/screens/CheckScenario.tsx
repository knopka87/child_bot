import { useNavigate } from "react-router";
import { ArrowLeft, Image, Images, BookOpen } from "lucide-react";
import { motion } from "motion/react";

const scenarios = [
  {
    id: "single",
    icon: <Image size={28} />,
    title: "1 фото",
    desc: "Задание и ответ на одном фото",
    color: "from-primary to-secondary",
  },
  {
    id: "two-separate",
    icon: <Images size={28} />,
    title: "2 фото",
    desc: "Задание отдельно, ответ отдельно",
    color: "from-[#00B894] to-[#55EFC4]",
  },
  {
    id: "two-pages",
    icon: <BookOpen size={28} />,
    title: "2 фото задания",
    desc: "Задание на нескольких страницах",
    color: "from-[#FDCB6E] to-[#F9CA24]",
  },
];

export function CheckScenario() {
  const navigate = useNavigate();

  const handleSelect = (id: string) => {
    if (id === "single") {
      navigate("/check/upload?scenario=single");
    } else if (id === "two-separate") {
      navigate("/check/upload?scenario=two-separate");
    } else {
      navigate("/check/upload?scenario=two-pages");
    }
  };

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate("/")}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Назад</span>
      </button>

      <div className="text-center mb-8">
        <h2 className="text-primary">Проверка ДЗ</h2>
        <p className="text-muted-foreground text-[14px] mt-1">
          Выбери, как выглядит твоё задание
        </p>
      </div>

      <div className="flex flex-col gap-4">
        {scenarios.map((s, i) => (
          <motion.button
            key={s.id}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.1 }}
            onClick={() => handleSelect(s.id)}
            className={`bg-gradient-to-r ${s.color} text-white rounded-3xl p-5 flex items-center gap-4 shadow-lg active:scale-[0.98] transition-transform`}
          >
            <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center flex-shrink-0">
              {s.icon}
            </div>
            <div className="text-left">
              <h3 className="text-white">{s.title}</h3>
              <p className="text-white/80 text-[13px] mt-0.5">{s.desc}</p>
            </div>
          </motion.button>
        ))}
      </div>
    </div>
  );
}
