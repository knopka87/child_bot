import { useNavigate, useLocation } from "react-router";
import { ArrowLeft, Check, RefreshCw, Crop } from "lucide-react";
import { motion } from "motion/react";

export function ImageQualityCheck() {
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
        <h2 className="text-primary">Всё ли видно?</h2>
        <p className="text-muted-foreground text-[14px] mt-1">
          Проверь, что задание хорошо видно на фото
        </p>
      </div>

      {/* Image preview */}
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="bg-white rounded-3xl p-3 shadow-sm border border-border mb-6 flex-1"
      >
        <div className="w-full aspect-[3/4] bg-muted/30 rounded-2xl flex items-center justify-center">
          <div className="text-center text-muted-foreground">
            <div className="text-[48px] mb-2">📸</div>
            <p className="text-[13px]">Предпросмотр изображения</p>
          </div>
        </div>
      </motion.div>

      {/* Actions */}
      <div className="flex flex-col gap-3">
        <button
          onClick={() => navigate(`/${mode}/crop`)}
          className="w-full py-3 bg-white border border-border text-foreground rounded-2xl flex items-center justify-center gap-2 active:scale-[0.98] transition-transform"
        >
          <Crop size={18} className="text-primary" />
          Обрезать
        </button>
        <button
          onClick={() => navigate(`/${mode}/processing`)}
          className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform flex items-center justify-center gap-2"
        >
          <Check size={18} />
          Всё видно, продолжить
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
