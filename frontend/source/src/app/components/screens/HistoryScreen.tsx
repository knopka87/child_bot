import { useState } from "react";
import { useNavigate } from "react-router";
import { ArrowLeft, BookOpen, ClipboardCheck, CheckCircle, XCircle, Clock, Lightbulb, Loader, RefreshCw, Pencil } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";

type Status = "correct" | "errors" | "unfinished" | "hints" | "processing";

const historyItems = [
  { id: 1, date: "10 марта 2026", mode: "help" as const, title: "Математика — задача про яблоки", status: "correct" as Status },
  { id: 2, date: "9 марта 2026", mode: "check" as const, title: "Русский — упражнение 45", status: "errors" as Status },
  { id: 3, date: "9 марта 2026", mode: "help" as const, title: "Математика — примеры", status: "unfinished" as Status },
  { id: 4, date: "8 марта 2026", mode: "help" as const, title: "Окружающий мир — вопросы", status: "hints" as Status },
  { id: 5, date: "7 марта 2026", mode: "check" as const, title: "Математика — задание 12", status: "processing" as Status },
];

const statusConfig: Record<Status, { label: string; icon: React.ReactNode; color: string; bg: string }> = {
  correct: { label: "Решено верно", icon: <CheckCircle size={14} />, color: "text-[#00B894]", bg: "bg-[#E8FFF8]" },
  errors: { label: "Есть ошибки", icon: <XCircle size={14} />, color: "text-[#E17055]", bg: "bg-[#FFF0F0]" },
  unfinished: { label: "Незакончено", icon: <Clock size={14} />, color: "text-[#FDCB6E]", bg: "bg-[#FFF9E8]" },
  hints: { label: "Использованы подсказки", icon: <Lightbulb size={14} />, color: "text-primary", bg: "bg-primary/10" },
  processing: { label: "В обработке", icon: <Loader size={14} />, color: "text-muted-foreground", bg: "bg-muted" },
};

export function HistoryScreen() {
  const navigate = useNavigate();
  const [selectedId, setSelectedId] = useState<number | null>(null);

  const selectedItem = historyItems.find((h) => h.id === selectedId);

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate("/profile")}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Профиль</span>
      </button>

      <h2 className="text-primary mb-4">История</h2>

      <div className="flex flex-col gap-3">
        {historyItems.map((item, i) => {
          const s = statusConfig[item.status];
          return (
            <motion.button
              key={item.id}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.05 }}
              onClick={() => setSelectedId(item.id)}
              className="bg-white rounded-2xl p-4 shadow-sm text-left active:scale-[0.98] transition-transform"
            >
              <div className="flex items-center justify-between mb-1">
                <div className="flex items-center gap-2">
                  {item.mode === "help" ? (
                    <BookOpen size={16} className="text-primary" />
                  ) : (
                    <ClipboardCheck size={16} className="text-[#00B894]" />
                  )}
                  <span className="text-[13px] text-foreground truncate max-w-[200px]">{item.title}</span>
                </div>
                <span className="text-[11px] text-muted-foreground flex-shrink-0">{item.date}</span>
              </div>
              <div className={`inline-flex items-center gap-1.5 ${s.color} ${s.bg} px-2.5 py-1 rounded-full`}>
                {s.icon}
                <span className="text-[11px]">{s.label}</span>
              </div>
            </motion.button>
          );
        })}
      </div>

      {/* Detail modal */}
      <AnimatePresence>
        {selectedItem && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-end justify-center z-50"
            onClick={() => setSelectedId(null)}
          >
            <motion.div
              initial={{ y: "100%" }}
              animate={{ y: 0 }}
              exit={{ y: "100%" }}
              className="bg-white rounded-t-3xl p-6 w-full max-w-[390px]"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="w-10 h-1 bg-muted rounded-full mx-auto mb-4" />
              <h3 className="text-foreground mb-1">{selectedItem.title}</h3>
              <p className="text-[12px] text-muted-foreground mb-3">{selectedItem.date}</p>
              <div className={`inline-flex items-center gap-1.5 ${statusConfig[selectedItem.status].color} ${statusConfig[selectedItem.status].bg} px-2.5 py-1 rounded-full mb-5`}>
                {statusConfig[selectedItem.status].icon}
                <span className="text-[11px]">{statusConfig[selectedItem.status].label}</span>
              </div>

              <div className="flex flex-col gap-3">
                {(selectedItem.status === "errors" || selectedItem.status === "unfinished") && (
                  <button
                    onClick={() => {
                      setSelectedId(null);
                      navigate("/check/upload");
                    }}
                    className="w-full py-3 bg-primary text-white rounded-2xl flex items-center justify-center gap-2"
                  >
                    <Pencil size={16} />
                    Исправить и проверить
                  </button>
                )}
                <button
                  onClick={() => {
                    setSelectedId(null);
                    navigate(selectedItem.mode === "help" ? "/help/upload" : "/check/scenario");
                  }}
                  className="w-full py-3 border border-primary text-primary rounded-2xl flex items-center justify-center gap-2"
                >
                  <RefreshCw size={16} />
                  Повторить
                </button>
                <button
                  onClick={() => setSelectedId(null)}
                  className="w-full py-3 text-muted-foreground"
                >
                  Закрыть
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}