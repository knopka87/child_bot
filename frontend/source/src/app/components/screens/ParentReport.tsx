import { useState } from "react";
import { useNavigate } from "react-router";
import { ArrowLeft, Download, Mail, Send } from "lucide-react";
import { toast } from "sonner";

const pastReports = [
  { id: 1, date: "10 марта 2026", status: "Отправлен" },
  { id: 2, date: "3 марта 2026", status: "Отправлен" },
  { id: 3, date: "24 февраля 2026", status: "Отправлен" },
];

export function ParentReport() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("parent@example.com");
  const [weeklyEnabled, setWeeklyEnabled] = useState(true);

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
      <button
        onClick={() => navigate("/profile")}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Профиль</span>
      </button>

      <h2 className="text-primary mb-4">Отчёт родителю</h2>

      <div className="bg-white rounded-2xl p-4 shadow-sm mb-4">
        {/* Email */}
        <div className="mb-4">
          <label className="text-[13px] text-muted-foreground mb-1 block">E-mail</label>
          <div className="flex items-center gap-2">
            <Mail size={18} className="text-muted-foreground" />
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="flex-1 bg-muted/30 rounded-xl px-3 py-2 text-[14px] border border-border outline-none focus:border-primary"
            />
          </div>
        </div>

        {/* Weekly toggle */}
        <div className="flex items-center justify-between mb-3">
          <span className="text-[14px] text-foreground">Еженедельный отчёт</span>
          <button
            onClick={() => setWeeklyEnabled(!weeklyEnabled)}
            className={`w-12 h-7 rounded-full transition-all relative ${
              weeklyEnabled ? "bg-primary" : "bg-switch-background"
            }`}
          >
            <div
              className={`w-5 h-5 bg-white rounded-full absolute top-1 transition-all shadow-sm ${
                weeklyEnabled ? "right-1" : "left-1"
              }`}
            />
          </button>
        </div>

        {weeklyEnabled && (
          <p className="text-[12px] text-muted-foreground">
            Отчёт отправляется каждый понедельник в 10:00
          </p>
        )}
      </div>

      <button
        onClick={() => toast.success("Тестовый отчёт отправлен!")}
        className="w-full py-3 bg-primary text-white rounded-2xl flex items-center justify-center gap-2 mb-6 shadow-lg shadow-primary/20"
      >
        <Send size={18} />
        Отправить тестовый отчёт
      </button>

      <h3 className="text-foreground mb-3">Архив отчётов</h3>
      <div className="flex flex-col gap-2">
        {pastReports.map((r) => (
          <div key={r.id} className="bg-white rounded-2xl p-3 flex items-center justify-between shadow-sm">
            <div>
              <p className="text-[14px] text-foreground">{r.date}</p>
              <p className="text-[12px] text-[#00B894]">{r.status}</p>
            </div>
            <button className="w-9 h-9 bg-primary/10 rounded-xl flex items-center justify-center">
              <Download size={18} className="text-primary" />
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
