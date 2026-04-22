import { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router";
import { ArrowLeft, Save, RefreshCw, X } from "lucide-react";
import { motion } from "motion/react";
import { Mascot } from "../Mascot";

export function ProcessingScreen() {
  const navigate = useNavigate();
  const location = useLocation();
  const mode = location.pathname.startsWith("/help") ? "help" : "check";
  const [seconds, setSeconds] = useState(0);
  const [dots, setDots] = useState("");
  const [isLongWait, setIsLongWait] = useState(false);

  useEffect(() => {
    const timer = setInterval(() => setSeconds((s) => s + 1), 1000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    const dotTimer = setInterval(() => {
      setDots((d) => (d.length >= 3 ? "" : d + "."));
    }, 500);
    return () => clearInterval(dotTimer);
  }, []);

  // Show long wait actions after 10s
  useEffect(() => {
    if (seconds >= 10) setIsLongWait(true);
  }, [seconds]);

  // Auto-navigate after 5s (demo)
  useEffect(() => {
    const timeout = setTimeout(() => {
      navigate(mode === "check" ? "/check/result" : "/help/result");
    }, 5000);
    return () => clearTimeout(timeout);
  }, [navigate, mode]);

  const messages = [
    { time: 0, text: mode === "help" ? "Думаю" : "Проверяю решение" },
    { time: 8, text: "Нужно чуть больше времени" },
    { time: 15, text: "Почти готово, подожди ещё немного" },
  ];

  const currentMessage = [...messages].reverse().find((m) => seconds >= m.time)?.text || messages[0].text;

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF]">
      <button
        onClick={() => navigate(-1)}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Назад</span>
      </button>

      <div className="flex-1 flex flex-col items-center justify-center">
        <Mascot size="lg" message={currentMessage + dots} />

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.3 }}
          className="mt-8 text-center"
        >
          <p className="text-primary text-[18px]">{currentMessage}{dots}</p>

          <div className="w-48 h-2 bg-white rounded-full mt-6 overflow-hidden mx-auto">
            <motion.div
              animate={{ x: ["-100%", "100%"] }}
              transition={{ repeat: Infinity, duration: 1.5, ease: "easeInOut" }}
              className="w-full h-full bg-gradient-to-r from-primary to-accent rounded-full"
            />
          </div>

          <p className="text-muted-foreground text-[12px] mt-4">
            Обработка продолжится, даже если ты выйдешь
          </p>
        </motion.div>

        {/* Long wait actions */}
        {isLongWait && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="mt-8 w-full flex flex-col gap-3"
          >
            <button
              onClick={() => navigate("/")}
              className="w-full py-3 bg-white border border-border text-foreground rounded-2xl flex items-center justify-center gap-2"
            >
              <Save size={16} className="text-primary" />
              Сохранить и подождать
            </button>
            <button
              onClick={() => {
                setSeconds(0);
                setIsLongWait(false);
              }}
              className="w-full py-3 bg-white border border-border text-foreground rounded-2xl flex items-center justify-center gap-2"
            >
              <RefreshCw size={16} className="text-primary" />
              Повторить
            </button>
            <button
              onClick={() => navigate("/")}
              className="w-full py-3 text-muted-foreground rounded-2xl flex items-center justify-center gap-2"
            >
              <X size={16} />
              Отменить
            </button>
          </motion.div>
        )}
      </div>
    </div>
  );
}
