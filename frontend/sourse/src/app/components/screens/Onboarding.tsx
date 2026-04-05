import { useState } from "react";
import { useNavigate } from "react-router";
import { Check, ArrowRight, Shield } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";

const avatars = ["🦊", "🐱", "🐶", "🐰", "🦄", "🐼", "🐸", "🦉"];
const grades = ["1 класс", "2 класс", "3 класс", "4 класс"];

type Step = "consent" | "profile";

export function Onboarding() {
  const navigate = useNavigate();
  const [step, setStep] = useState<Step>("consent");
  const [adultConsent, setAdultConsent] = useState(false);
  const [privacyAccepted, setPrivacyAccepted] = useState(false);
  const [name, setName] = useState("");
  const [selectedAvatar, setSelectedAvatar] = useState<number | null>(null);
  const [selectedGrade, setSelectedGrade] = useState<string>("");

  const canProceedConsent = adultConsent && privacyAccepted;
  const canFinish = name.trim().length > 0 && selectedAvatar !== null;

  if (step === "consent") {
    return (
      <div className="flex flex-col min-h-screen px-6 py-8 bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF]">
        <div className="flex-1 flex flex-col items-center justify-center text-center">
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ type: "spring" }}
            className="w-20 h-20 bg-primary/10 rounded-full flex items-center justify-center mb-6"
          >
            <Shield size={36} className="text-primary" />
          </motion.div>

          <h1 className="text-primary mb-2">Добро пожаловать!</h1>
          <p className="text-muted-foreground text-[14px] mb-8">
            Перед началом работы, пожалуйста, подтвердите согласие
          </p>

          <div className="w-full flex flex-col gap-4 text-left">
            <button
              onClick={() => setAdultConsent(!adultConsent)}
              className="flex items-start gap-3 bg-white rounded-2xl p-4 shadow-sm"
            >
              <div
                className={`w-6 h-6 mt-0.5 rounded-md flex items-center justify-center flex-shrink-0 transition-all ${
                  adultConsent ? "bg-primary" : "bg-white border-2 border-border"
                }`}
              >
                {adultConsent && <Check size={14} className="text-white" />}
              </div>
              <span className="text-[14px] text-foreground">
                Я подтверждаю, что являюсь взрослым (родителем или опекуном) и даю согласие на использование сервиса ребёнком
              </span>
            </button>

            <button
              onClick={() => setPrivacyAccepted(!privacyAccepted)}
              className="flex items-start gap-3 bg-white rounded-2xl p-4 shadow-sm"
            >
              <div
                className={`w-6 h-6 mt-0.5 rounded-md flex items-center justify-center flex-shrink-0 transition-all ${
                  privacyAccepted ? "bg-primary" : "bg-white border-2 border-border"
                }`}
              >
                {privacyAccepted && <Check size={14} className="text-white" />}
              </div>
              <span className="text-[13px] text-foreground leading-relaxed">
                Я согласен с{" "}
                <span className="text-primary underline">Политикой конфиденциальности</span> и{" "}
                <span className="text-primary underline">Пользовательским соглашением</span>
              </span>
            </button>
          </div>
        </div>

        <button
          disabled={!canProceedConsent}
          onClick={() => setStep("profile")}
          className={`w-full py-4 rounded-2xl transition-all text-white flex items-center justify-center gap-2 ${
            canProceedConsent
              ? "bg-primary shadow-lg shadow-primary/30 active:scale-[0.98]"
              : "bg-muted text-muted-foreground cursor-not-allowed"
          }`}
        >
          Продолжить
          <ArrowRight size={18} />
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-screen px-6 py-8 bg-gradient-to-b from-[#F0F4FF] to-[#E8E4FF]">
      <div className="text-center mb-6">
        <h1 className="text-primary mb-1">Создание профиля</h1>
        <p className="text-muted-foreground text-[14px]">Расскажи о себе!</p>
      </div>

      {/* Avatar selection */}
      <div className="mb-5">
        <label className="text-[14px] text-muted-foreground mb-2 block">Выбери аватар</label>
        <div className="grid grid-cols-4 gap-3">
          {avatars.map((a, i) => (
            <button
              key={i}
              onClick={() => setSelectedAvatar(i)}
              className={`w-full aspect-square rounded-2xl flex items-center justify-center text-[32px] transition-all ${
                selectedAvatar === i
                  ? "bg-primary shadow-lg scale-105 ring-3 ring-primary/30"
                  : "bg-white shadow-sm"
              }`}
            >
              {a}
            </button>
          ))}
        </div>
      </div>

      {/* Name input */}
      <div className="mb-5">
        <label className="text-[14px] text-muted-foreground mb-1.5 block">Как тебя зовут?</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Введи имя"
          className="w-full bg-white rounded-xl px-4 py-3 border border-border focus:ring-2 focus:ring-primary/30 focus:border-primary outline-none transition-all"
        />
      </div>

      {/* Grade selector */}
      <div className="mb-5">
        <label className="text-[14px] text-muted-foreground mb-1.5 block">Класс (необязательно)</label>
        <div className="flex gap-2 flex-wrap">
          {grades.map((g) => (
            <button
              key={g}
              onClick={() => setSelectedGrade(selectedGrade === g ? "" : g)}
              className={`px-4 py-2 rounded-xl transition-all text-[14px] ${
                selectedGrade === g
                  ? "bg-primary text-white shadow"
                  : "bg-white text-foreground border border-border"
              }`}
            >
              {g}
            </button>
          ))}
        </div>
      </div>

      <div className="mt-auto">
        <button
          disabled={!canFinish}
          onClick={() => navigate("/")}
          className={`w-full py-4 rounded-2xl transition-all text-white ${
            canFinish
              ? "bg-primary shadow-lg shadow-primary/30 active:scale-[0.98]"
              : "bg-muted text-muted-foreground cursor-not-allowed"
          }`}
        >
          Начать!
        </button>
      </div>
    </div>
  );
}
