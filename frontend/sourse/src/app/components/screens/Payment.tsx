import { useState } from "react";
import { useNavigate } from "react-router";
import { ArrowLeft, Check, Star, Shield, Zap, Clock, BookOpen, CheckCircle, XCircle } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";

type PayState = "paywall" | "plans" | "success" | "error";

const benefits = [
  { icon: <Zap size={20} />, text: "Безлимитная проверка ДЗ" },
  { icon: <BookOpen size={20} />, text: "Подробные объяснения заданий" },
  { icon: <Star size={20} />, text: "Все достижения и стикеры" },
  { icon: <Shield size={20} />, text: "Еженедельные отчёты родителям" },
  { icon: <Clock size={20} />, text: "Приоритетная обработка" },
];

const plans = [
  { id: "monthly", label: "Месяц", price: "299 ₽", period: "/мес", popular: false },
  { id: "yearly", label: "Год", price: "1 990 ₽", period: "/год", popular: true, savings: "Экономия 44%" },
];

export function PaymentScreen() {
  const navigate = useNavigate();
  const [state, setState] = useState<PayState>("paywall");
  const [selectedPlan, setSelectedPlan] = useState("yearly");

  if (state === "success") {
    return (
      <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#E8FFF8] to-background">
        <div className="flex-1 flex flex-col items-center justify-center text-center">
          <motion.div initial={{ scale: 0 }} animate={{ scale: 1 }} transition={{ type: "spring" }}>
            <CheckCircle size={64} className="text-[#00B894] mx-auto mb-4" />
          </motion.div>
          <h1 className="text-[#00B894] mb-2">Оплата прошла успешно!</h1>
          <p className="text-muted-foreground text-[14px]">
            Все возможности разблокированы. Приятного обучения!
          </p>
        </div>
        <button
          onClick={() => navigate(-1)}
          className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20 active:scale-[0.98] transition-transform"
        >
          Продолжить
        </button>
      </div>
    );
  }

  if (state === "error") {
    return (
      <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#FFF0F0] to-background">
        <div className="flex-1 flex flex-col items-center justify-center text-center">
          <motion.div initial={{ scale: 0 }} animate={{ scale: 1 }} transition={{ type: "spring" }}>
            <XCircle size={64} className="text-destructive mx-auto mb-4" />
          </motion.div>
          <h1 className="text-destructive mb-2">Ошибка оплаты</h1>
          <p className="text-muted-foreground text-[14px]">
            Не удалось обработать платёж. Попробуйте ещё раз.
          </p>
        </div>
        <div className="flex flex-col gap-3">
          <button
            onClick={() => setState("plans")}
            className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/20"
          >
            Попробовать ещё раз
          </button>
          <button
            onClick={() => navigate(-1)}
            className="w-full py-3 text-muted-foreground"
          >
            Позже
          </button>
        </div>
      </div>
    );
  }

  if (state === "plans") {
    return (
      <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF]">
        <button onClick={() => setState("paywall")} className="flex items-center gap-2 text-primary mb-4">
          <ArrowLeft size={20} /><span className="text-[14px]">Назад</span>
        </button>

        <h2 className="text-primary mb-1">Выбери план</h2>
        <p className="text-muted-foreground text-[14px] mb-6">Выбери удобный вариант подписки</p>

        <div className="flex flex-col gap-4 mb-6">
          {plans.map((plan) => (
            <button
              key={plan.id}
              onClick={() => setSelectedPlan(plan.id)}
              className={`rounded-3xl p-5 text-left transition-all relative ${
                selectedPlan === plan.id
                  ? "bg-white border-2 border-primary shadow-lg"
                  : "bg-white border border-border shadow-sm"
              }`}
            >
              {plan.popular && (
                <span className="absolute -top-3 right-4 bg-accent text-white text-[11px] px-3 py-1 rounded-full">
                  Популярный
                </span>
              )}
              <div className="flex items-baseline gap-1">
                <span className="text-[28px] text-foreground">{plan.price.split(" ")[0]}</span>
                <span className="text-[14px] text-muted-foreground"> {plan.price.split(" ")[1]}{plan.period}</span>
              </div>
              <p className="text-[14px] text-foreground mt-1">{plan.label}</p>
              {plan.savings && <p className="text-[#00B894] text-[12px] mt-1">{plan.savings}</p>}
            </button>
          ))}
        </div>

        <div className="flex flex-col gap-3 mt-auto">
          <button
            onClick={() => {
              // Simulate: 80% success, 20% error
              setState(Math.random() > 0.2 ? "success" : "error");
            }}
            className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/30 active:scale-[0.98] transition-transform"
          >
            Оплатить
          </button>
          <button onClick={() => navigate(-1)} className="w-full py-3 text-muted-foreground text-[14px]">
            Позже
          </button>
        </div>
      </div>
    );
  }

  // Paywall
  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF]">
      <button onClick={() => navigate(-1)} className="flex items-center gap-2 text-primary mb-4">
        <ArrowLeft size={20} /><span className="text-[14px]">Назад</span>
      </button>

      <div className="text-center mb-6">
        <div className="w-16 h-16 bg-gradient-to-br from-primary to-accent rounded-full flex items-center justify-center mx-auto mb-3">
          <Star size={30} className="text-white" />
        </div>
        <h1 className="text-primary">Доступ к сервису</h1>
        <p className="text-muted-foreground text-[14px] mt-1">
          Разблокируй все возможности!
        </p>
      </div>

      <div className="bg-white rounded-3xl p-5 shadow-sm mb-6">
        <div className="flex flex-col gap-3.5">
          {benefits.map((b, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, x: -10 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.08 }}
              className="flex items-center gap-3"
            >
              <div className="w-9 h-9 bg-primary/10 rounded-xl flex items-center justify-center text-primary flex-shrink-0">
                {b.icon}
              </div>
              <span className="text-[14px] text-foreground">{b.text}</span>
              <Check size={16} className="text-[#00B894] ml-auto" />
            </motion.div>
          ))}
        </div>
      </div>

      <div className="flex flex-col gap-3 mt-auto">
        <button
          onClick={() => setState("plans")}
          className="w-full py-4 bg-primary text-white rounded-2xl shadow-lg shadow-primary/30 active:scale-[0.98] transition-transform"
        >
          Выбрать план
        </button>
        <button onClick={() => navigate(-1)} className="w-full py-3 text-muted-foreground text-[14px]">
          Позже
        </button>
      </div>
    </div>
  );
}
