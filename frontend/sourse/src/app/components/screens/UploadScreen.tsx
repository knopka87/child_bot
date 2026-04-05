import { useState } from "react";
import { useNavigate, useLocation } from "react-router";
import { ArrowLeft, Camera, Image, AlertCircle } from "lucide-react";
import { motion } from "motion/react";

export function UploadScreen() {
  const navigate = useNavigate();
  const location = useLocation();
  const [error, setError] = useState<string | null>(null);

  const mode = location.pathname.startsWith("/help") ? "help" : "check";
  const isHelp = mode === "help";
  const title = isHelp ? "Помоги разобраться" : "Проверка ДЗ";
  const params = new URLSearchParams(location.search);
  const scenario = params.get("scenario");
  const needsMultiple = scenario === "two-separate" || scenario === "two-pages";

  const handleFileSelect = () => {
    setError(null);
    if (needsMultiple) {
      navigate(`/${mode}/images?scenario=${scenario}`);
    } else {
      navigate(`/${mode}/quality`);
    }
  };

  const handleCamera = () => {
    setError(null);
    if (needsMultiple) {
      navigate(`/${mode}/images?scenario=${scenario}`);
    } else {
      navigate(`/${mode}/quality`);
    }
  };

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Назад</span>
      </button>

      <div className="text-center mb-8">
        <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-3">
          <Image size={30} className="text-primary" />
        </div>
        <h2 className="text-primary">{title}</h2>
        <p className="text-muted-foreground text-[14px] mt-1">
          {needsMultiple
            ? "Загрузи первое фото"
            : "Загрузи фото задания из учебника"}
        </p>
      </div>

      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="flex-1 flex flex-col gap-4"
      >
        <button
          onClick={handleFileSelect}
          className="bg-white rounded-3xl p-6 flex items-center gap-4 shadow-sm border border-border active:scale-[0.98] transition-transform"
        >
          <div className="w-12 h-12 bg-primary/10 rounded-2xl flex items-center justify-center">
            <Image size={24} className="text-primary" />
          </div>
          <div className="text-left">
            <p className="text-foreground">Выбрать изображение</p>
            <p className="text-muted-foreground text-[12px]">JPG, PNG</p>
          </div>
        </button>

        <button
          onClick={handleCamera}
          className="bg-white rounded-3xl p-6 flex items-center gap-4 shadow-sm border border-border active:scale-[0.98] transition-transform"
        >
          <div className="w-12 h-12 bg-[#FDCB6E]/20 rounded-2xl flex items-center justify-center">
            <Camera size={24} className="text-[#E17055]" />
          </div>
          <div className="text-left">
            <p className="text-foreground">Сфотографировать</p>
            <p className="text-muted-foreground text-[12px]">Открыть камеру</p>
          </div>
        </button>

        {error && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className="bg-red-50 rounded-2xl p-4 flex gap-3 items-start"
          >
            <AlertCircle size={20} className="text-destructive flex-shrink-0 mt-0.5" />
            <div>
              <p className="text-destructive text-[14px]">{error}</p>
              <button
                onClick={() => setError(null)}
                className="text-primary text-[13px] mt-2 underline"
              >
                Попробовать снова
              </button>
            </div>
          </motion.div>
        )}
      </motion.div>
    </div>
  );
}
