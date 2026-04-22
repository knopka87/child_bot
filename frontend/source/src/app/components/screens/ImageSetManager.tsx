import { useState } from "react";
import { useNavigate, useLocation } from "react-router";
import { ArrowLeft, Plus, Trash2, ArrowUpDown, RefreshCw } from "lucide-react";
import { motion } from "motion/react";

interface ImageItem {
  id: number;
  label: string;
  uploaded: boolean;
}

export function ImageSetManager() {
  const navigate = useNavigate();
  const location = useLocation();
  const mode = location.pathname.startsWith("/help") ? "help" : "check";
  const params = new URLSearchParams(location.search);
  const scenario = params.get("scenario") || "two-separate";

  const getInitialImages = (): ImageItem[] => {
    if (scenario === "two-separate") {
      return [
        { id: 1, label: "Задание", uploaded: true },
        { id: 2, label: "Ответ", uploaded: false },
      ];
    }
    return [
      { id: 1, label: "Страница 1", uploaded: true },
      { id: 2, label: "Страница 2", uploaded: false },
    ];
  };

  const [images, setImages] = useState<ImageItem[]>(getInitialImages);

  const handleUpload = (id: number) => {
    setImages((prev) => prev.map((img) => (img.id === id ? { ...img, uploaded: true } : img)));
  };

  const handleDelete = (id: number) => {
    setImages((prev) => prev.map((img) => (img.id === id ? { ...img, uploaded: false } : img)));
  };

  const allUploaded = images.every((img) => img.uploaded);

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Назад</span>
      </button>

      <h2 className="text-primary mb-1">Изображения</h2>
      <p className="text-muted-foreground text-[14px] mb-6">
        Загрузи все необходимые фото
      </p>

      <div className="flex flex-col gap-4 mb-6 flex-1">
        {images.map((img, i) => (
          <motion.div
            key={img.id}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.1 }}
            className="bg-white rounded-2xl p-4 shadow-sm border border-border"
          >
            <div className="flex items-center justify-between mb-3">
              <span className="text-[14px] text-foreground">{img.label}</span>
              <div className="flex items-center gap-2">
                {img.uploaded && (
                  <>
                    <button
                      onClick={() => handleUpload(img.id)}
                      className="w-8 h-8 bg-primary/10 rounded-lg flex items-center justify-center"
                    >
                      <RefreshCw size={14} className="text-primary" />
                    </button>
                    <button
                      onClick={() => handleDelete(img.id)}
                      className="w-8 h-8 bg-red-50 rounded-lg flex items-center justify-center"
                    >
                      <Trash2 size={14} className="text-destructive" />
                    </button>
                  </>
                )}
              </div>
            </div>

            {img.uploaded ? (
              <div className="w-full aspect-[16/9] bg-muted/30 rounded-xl flex items-center justify-center">
                <div className="text-center">
                  <div className="text-[32px]">📸</div>
                  <p className="text-[11px] text-muted-foreground">Загружено</p>
                </div>
              </div>
            ) : (
              <button
                onClick={() => handleUpload(img.id)}
                className="w-full aspect-[16/9] bg-muted/20 rounded-xl border-2 border-dashed border-border flex items-center justify-center gap-2 text-muted-foreground active:scale-[0.98] transition-transform"
              >
                <Plus size={20} />
                <span className="text-[13px]">Загрузить фото</span>
              </button>
            )}
          </motion.div>
        ))}
      </div>

      {images.length > 1 && (
        <button className="w-full py-3 bg-white border border-border text-foreground rounded-2xl flex items-center justify-center gap-2 mb-3 active:scale-[0.98] transition-transform">
          <ArrowUpDown size={16} className="text-primary" />
          Поменять порядок
        </button>
      )}

      <button
        onClick={() => navigate(`/${mode}/quality`)}
        disabled={!allUploaded}
        className={`w-full py-4 rounded-2xl transition-all ${
          allUploaded
            ? "bg-primary text-white shadow-lg shadow-primary/20 active:scale-[0.98]"
            : "bg-muted text-muted-foreground cursor-not-allowed"
        }`}
      >
        Продолжить
      </button>
    </div>
  );
}
