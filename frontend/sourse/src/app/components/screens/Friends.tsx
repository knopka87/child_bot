import { useState } from "react";
import { Copy, Send, Users, CheckCircle, Gift } from "lucide-react";
import { motion } from "motion/react";

export function FriendsScreen() {
  const [copied, setCopied] = useState(false);
  const referralLink = "https://homework.app/invite/abc123";
  const invitedCount = 2;
  const targetCount = 5;

  const handleCopy = () => {
    navigator.clipboard?.writeText(referralLink).catch(() => {});
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleSend = () => {
    if (navigator.share) {
      navigator.share({ title: "Помощник ДЗ", text: "Присоединяйся!", url: referralLink });
    }
  };

  return (
    <div className="flex flex-col min-h-full px-5 pt-8 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <h1 className="text-primary mb-1">Друзья</h1>
      <p className="text-muted-foreground text-[14px] mb-6">
        Пригласи друзей и учитесь вместе!
      </p>

      {/* Invite reward progress */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-gradient-to-r from-[#FDCB6E] to-[#F9CA24] rounded-3xl p-5 mb-6 text-center"
      >
        <div className="flex items-center justify-center gap-2 mb-2">
          <Gift size={20} className="text-white" />
          <p className="text-white text-[14px]">Пригласи {targetCount} друзей — получи редкий стикер!</p>
        </div>
        <div className="flex gap-2 justify-center mb-2">
          {Array.from({ length: targetCount }).map((_, i) => (
            <div
              key={i}
              className={`w-8 h-8 rounded-full flex items-center justify-center text-[14px] ${
                i < invitedCount
                  ? "bg-white text-[#E17055]"
                  : "bg-white/30 text-white/80"
              }`}
            >
              {i < invitedCount ? "✓" : i + 1}
            </div>
          ))}
        </div>
        <p className="text-white/80 text-[12px]">{invitedCount} из {targetCount}</p>
        {/* Reward preview */}
        <div className="mt-3 bg-white/20 rounded-xl px-3 py-2 inline-flex items-center gap-2">
          <span className="text-[20px]">⭐</span>
          <span className="text-white text-[12px]">Редкий стикер «Дружба»</span>
        </div>
      </motion.div>

      {/* Invite card */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        className="bg-white rounded-3xl p-6 text-center shadow-sm mb-6"
      >
        <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4">
          <Users size={30} className="text-primary" />
        </div>
        <h3 className="text-foreground mb-2">Пригласи друга</h3>
        <p className="text-muted-foreground text-[13px] mb-5">
          Отправь ссылку другу и получите оба бонусные стикеры!
        </p>

        <div className="bg-muted rounded-xl px-4 py-3 mb-4 text-[13px] text-muted-foreground break-all">
          {referralLink}
        </div>

        <div className="flex gap-3">
          <button
            onClick={handleCopy}
            className={`flex-1 py-3 rounded-2xl flex items-center justify-center gap-2 transition-all text-[14px] ${
              copied ? "bg-[#00B894] text-white" : "bg-primary text-white"
            }`}
          >
            {copied ? <CheckCircle size={18} /> : <Copy size={18} />}
            {copied ? "Скопировано" : "Скопировать"}
          </button>
          <button
            onClick={handleSend}
            className="flex-1 py-3 rounded-2xl border border-primary text-primary flex items-center justify-center gap-2 text-[14px]"
          >
            <Send size={18} />
            Отправить
          </button>
        </div>
      </motion.div>

      {/* Stats */}
      <div className="bg-white rounded-2xl p-4 shadow-sm">
        <div className="flex justify-between items-center">
          <span className="text-[14px] text-foreground">Приглашено друзей</span>
          <span className="text-primary text-[20px]">{invitedCount}</span>
        </div>
      </div>
    </div>
  );
}
