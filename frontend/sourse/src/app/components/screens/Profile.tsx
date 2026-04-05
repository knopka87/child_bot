import { useNavigate } from "react-router";
import { History, FileText, HelpCircle, ChevronRight, CreditCard } from "lucide-react";
import { motion } from "motion/react";

const menuItems = [
  { path: "/profile/history", icon: History, label: "История", color: "bg-primary/10 text-primary" },
  { path: "/profile/report", icon: FileText, label: "Отчёт родителю", color: "bg-[#00B894]/10 text-[#00B894]" },
  { path: "/profile/help", icon: HelpCircle, label: "Помощь", color: "bg-[#74B9FF]/10 text-[#0984E3]" },
  { path: "/payment", icon: CreditCard, label: "Подписка", color: "bg-[#FDCB6E]/10 text-[#E17055]" },
];

export function ProfileScreen() {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col min-h-full px-5 pt-8 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      {/* Profile card */}
      <motion.div
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-white rounded-3xl p-5 flex items-center gap-4 shadow-sm mb-6"
      >
        <div className="w-16 h-16 bg-primary rounded-2xl flex items-center justify-center text-[32px]">
          🦊
        </div>
        <div>
          <h2 className="text-foreground">Артём</h2>
          <p className="text-muted-foreground text-[13px]">2 класс</p>
          <p className="text-[#00B894] text-[12px] mt-0.5">Пробный период — 5 дней</p>
        </div>
      </motion.div>

      {/* Menu items */}
      <div className="flex flex-col gap-2">
        {menuItems.map((item, i) => {
          const Icon = item.icon;
          return (
            <motion.button
              key={item.path}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.05 }}
              onClick={() => navigate(item.path)}
              className="bg-white rounded-2xl p-4 flex items-center gap-3 shadow-sm active:scale-[0.98] transition-transform"
            >
              <div className={`w-10 h-10 rounded-xl flex items-center justify-center ${item.color}`}>
                <Icon size={20} />
              </div>
              <span className="flex-1 text-left text-foreground">{item.label}</span>
              <ChevronRight size={18} className="text-muted-foreground" />
            </motion.button>
          );
        })}
      </div>
    </div>
  );
}
