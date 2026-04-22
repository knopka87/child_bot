import { useNavigate } from "react-router";
import { ArrowLeft, HelpCircle, MessageCircle, Book } from "lucide-react";
import { motion } from "motion/react";

const helpItems = [
  { icon: <HelpCircle size={20} />, title: "Как пользоваться?", desc: "Загрузи фото задания и получи помощь" },
  { icon: <MessageCircle size={20} />, title: "Связаться с поддержкой", desc: "Напишите нам, если что-то не работает" },
  { icon: <Book size={20} />, title: "Частые вопросы", desc: "Ответы на популярные вопросы" },
];

export function HelpScreen() {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate("/profile")}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Профиль</span>
      </button>

      <h2 className="text-primary mb-4">Помощь</h2>

      <div className="flex flex-col gap-3">
        {helpItems.map((item, i) => (
          <motion.div
            key={i}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.05 }}
            className="bg-white rounded-2xl p-4 flex items-start gap-3 shadow-sm"
          >
            <div className="w-10 h-10 bg-primary/10 rounded-xl flex items-center justify-center text-primary flex-shrink-0">
              {item.icon}
            </div>
            <div>
              <p className="text-foreground text-[14px]">{item.title}</p>
              <p className="text-muted-foreground text-[12px] mt-0.5">{item.desc}</p>
            </div>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
