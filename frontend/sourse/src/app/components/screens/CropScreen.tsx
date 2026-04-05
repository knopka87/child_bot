import { useNavigate, useLocation } from "react-router";
import { ArrowLeft, Check, RefreshCw } from "lucide-react";
import { motion } from "motion/react";

export function CropScreen() {
  const navigate = useNavigate();
  const location = useLocation();
  const mode = location.pathname.startsWith("/help") ? "help" : "check";

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Назад</span>
      </button>

      <div className="text-center mb-6">
        <h2 className="text-primary">Обрежь до одного задания</h2>
        <p className="text-muted-foreground text-[14px] mt-1">
          Выдели область с одним заданием
        </p>
      </div>

      {/* Crop area */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="bg-white rounded-3xl p-3 shadow-sm border border-border mb-6 flex-1 relative"
      >
        <div className="w-full aspect-[3/4] bg-muted/30 rounded-2xl relative overflow-hidden">
          {/* Simulated crop overlay */}
          <div className="absolute inset-0 bg-black/30" />
          <div className="absolute top-[15%] left-[10%] right-[10%] bottom-[25%] border-2 border-white border-dashed rounded-lg bg-transparent">
            {/* Corner handles */}
            {["top-left", "top-right", "bottom-left", "bottom-right"].map((pos) => (
              <div
                key={pos}
                className={`absolute w-4 h-4 bg-white rounded-full shadow ${
                  pos.includes("top") ? "-top-2" : "-bottom-2"
                } ${pos.includes("left") ? "-left-2" : "-right-2"}`}
              />
            ))}
          </div>
          <div className="absolute inset-0 flex items-center justify-center">
            <p className="text-white text-[13px] bg-black/40 px-3 py-1 rounded-full">
              Перетащи границы
            </p>
          </div>
        </div>
      </motion.div>

      {/* Actions */}
      <div className="flex flex-col gap-3">
        <button
          onClick={() => navigate(`/${mode}/processing`)}
          className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform flex items-center justify-center gap-2"
        >
          <Check size={18} />
          Подтвердить
        </button>
        <button
          onClick={() => navigate(-1)}
          className="w-full py-3 border border-primary text-primary rounded-2xl flex items-center justify-center gap-2 active:scale-[0.98] transition-transform"
        >
          <RefreshCw size={18} />
          Переснять
        </button>
      </div>
    </div>
  );
}
