import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, Download, Mail, Send, FileText } from "lucide-react";
import { motion } from "framer-motion";
import { getCurrentChildProfileId } from "@/lib/auth";
import { ROUTES } from "@/config/routes";
import { BottomNav } from "@/components/layout/BottomNav";
import styles from "@/pages/Profile/ProfilePage.module.css";

// API configuration
const API_BASE_URL = "http://localhost:8080";
const PLATFORM_ID = "web";

interface ReportInfo {
  id: string;
  reportDate: string;
  sentAt?: string;
  createdAt: string;
}

export default function ReportPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [weeklyEnabled, setWeeklyEnabled] = useState(true);
  const [reportsList, setReportsList] = useState<ReportInfo[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    const childProfileId = await getCurrentChildProfileId();
    if (!childProfileId) return;

    try {
      // Load email settings from backend API
      const settingsResponse = await fetch(
        `${API_BASE_URL}/reports/${childProfileId}/settings`,
        {
          headers: {
            "X-Platform-ID": PLATFORM_ID,
            "X-Child-Profile-ID": childProfileId,
          },
        }
      );

      if (settingsResponse.ok) {
        const settings = await settingsResponse.json();
        setEmail(settings.email || "");
        setWeeklyEnabled(settings.weeklyReportEnabled ?? true);
      }

      // Load reports list
      const listResponse = await fetch(
        `${API_BASE_URL}/reports/${childProfileId}/list`,
        {
          headers: {
            "X-Platform-ID": PLATFORM_ID,
            "X-Child-Profile-ID": childProfileId,
          },
        }
      );

      if (listResponse.ok) {
        const data = await listResponse.json();
        setReportsList(data || []);
      }
    } catch (error) {
      console.error("[ReportPage] Failed to load data:", error);
    }
  };

  const handleEmailChange = async (newEmail: string) => {
    setEmail(newEmail);

    // Save to backend API
    if (newEmail.includes("@")) {
      const childProfileId = await getCurrentChildProfileId();
      if (!childProfileId) return;

      try {
        await fetch(
          `${API_BASE_URL}/reports/${childProfileId}/settings`,
          {
            method: "PUT",
            headers: {
              "Content-Type": "application/json",
              "X-Platform-ID": PLATFORM_ID,
              "X-Child-Profile-ID": childProfileId,
            },
            body: JSON.stringify({ email: newEmail }),
          }
        );
      } catch (error) {
        console.error("[ReportPage] Failed to save email:", error);
      }
    }
  };

  const handleWeeklyToggle = async () => {
    const newValue = !weeklyEnabled;
    setWeeklyEnabled(newValue);

    const childProfileId = await getCurrentChildProfileId();
    if (!childProfileId) return;

    try {
      await fetch(
        `${API_BASE_URL}/reports/${childProfileId}/settings`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            "X-Platform-ID": PLATFORM_ID,
            "X-Child-Profile-ID": childProfileId,
          },
          body: JSON.stringify({ weeklyReportEnabled: newValue }),
        }
      );
    } catch (error) {
      console.error("[ReportPage] Failed to save weekly setting:", error);
      // Откатываем изменение в случае ошибки
      setWeeklyEnabled(!newValue);
    }
  };

  const handleSendTestReport = async () => {
    if (!email || !email.includes("@")) {
      alert("⚠️ Пожалуйста, введите корректный email");
      return;
    }

    const childProfileId = await getCurrentChildProfileId();
    if (!childProfileId) {
      alert("⚠️ Профиль не найден");
      return;
    }

    setIsLoading(true);
    try {
      const response = await fetch(
        `${API_BASE_URL}/reports/${childProfileId}/send-test`,
        {
          method: "POST",
          headers: {
            "X-Platform-ID": PLATFORM_ID,
            "X-Child-Profile-ID": childProfileId,
          },
        }
      );

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to send report");
      }

      const result = await response.json();
      alert(`✅ ${result.message || "Тестовый отчёт отправлен на " + email}`);
    } catch (error) {
      console.error("[ReportPage] Failed to send test report:", error);
      alert("⚠️ Не удалось отправить тестовый отчёт");
    } finally {
      setIsLoading(false);
    }
  };

  const handleDownloadReport = async (reportDate: string) => {
    const childProfileId = await getCurrentChildProfileId();
    if (!childProfileId) {
      alert("Профиль не найден");
      return;
    }

    try {
      const response = await fetch(
        `${API_BASE_URL}/reports/${childProfileId}/${reportDate}/download`,
        {
          headers: {
            "X-Platform-ID": PLATFORM_ID,
            "X-Child-Profile-ID": childProfileId,
          },
        }
      );

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }

      // Получаем blob и создаем ссылку для скачивания
      const blob = await response.blob();
      const downloadUrl = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = downloadUrl;
      link.download = `report_${reportDate}.pdf`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(downloadUrl);
    } catch (err) {
      console.error("Failed to download report:", err);
      alert("Не удалось скачать отчёт");
    }
  };

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString("ru-RU", {
      day: "2-digit",
      month: "long",
      year: "numeric",
    });
  };

  const formatShortDate = (dateStr: string) => {
    const date = new Date(dateStr);
    return date.toLocaleDateString("ru-RU", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    });
  };

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-20 bg-gradient-to-b from-[#F0F4FF] to-background" style={{ color: '#2C2D2E' }}>
      {/* Back button */}
      <button
        onClick={() => navigate(ROUTES.PROFILE)}
        className="flex items-center gap-2 text-primary mb-6"
      >
        <ArrowLeft size={20} />
        <span className="text-[14px]">Профиль</span>
      </button>

      {/* Title */}
      <h2 className="text-primary mb-4 text-[24px] font-semibold">Отчёт родителю</h2>

      {/* Settings card */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        className="bg-white rounded-2xl p-4 shadow-sm mb-4"
      >
        {/* Email */}
        <div className="mb-4">
          <label className="text-[13px] text-muted-foreground mb-1 block">E-mail</label>
          <div className="flex items-center gap-2">
            <Mail size={18} className="text-muted-foreground" />
            <input
              type="email"
              value={email}
              onChange={(e) => handleEmailChange(e.target.value)}
              placeholder="Введите email родителя"
              className="flex-1 bg-muted/30 rounded-xl px-3 py-2 text-[14px] border border-border outline-none focus:border-primary"
            />
          </div>
        </div>

        {/* Weekly toggle */}
        <div className="flex items-center justify-between mb-3">
          <span className="text-[14px] text-foreground">Еженедельный отчёт</span>
          <button
            onClick={handleWeeklyToggle}
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
      </motion.div>

      {/* Send test report button */}
      <motion.button
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.1 }}
        onClick={handleSendTestReport}
        disabled={isLoading || !email}
        className={`w-full py-3 bg-primary text-white rounded-2xl flex items-center justify-center gap-2 mb-6 shadow-lg shadow-primary/20 ${
          (!email || isLoading) ? "opacity-50 cursor-not-allowed" : ""
        }`}
      >
        <Send size={18} />
        {isLoading ? "Отправляем..." : "Отправить тестовый отчёт"}
      </motion.button>

      {/* Archive section */}
      {reportsList.length > 0 && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
        >
          <h3 className="text-foreground mb-3">📚 Архив отчётов</h3>
          <div className="flex flex-col gap-2">
            {reportsList.map((report) => (
              <div
                key={report.id}
                className="bg-white rounded-2xl p-4 shadow-sm border border-border"
              >
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 bg-primary/10 rounded-xl flex items-center justify-center">
                      <FileText size={20} className="text-primary" />
                    </div>
                    <div>
                      <p className="text-[14px] text-foreground font-medium">
                        {formatDate(report.reportDate)}
                      </p>
                      <p className="text-[12px] text-muted-foreground">
                        Создан {formatShortDate(report.createdAt)}
                      </p>
                    </div>
                  </div>
                </div>
                <button
                  onClick={() => handleDownloadReport(report.reportDate)}
                  className="w-full py-2.5 bg-primary/10 text-primary rounded-xl flex items-center justify-center gap-2 active:scale-95 transition-transform"
                >
                  <Download size={16} />
                  <span className="text-[14px] font-medium">Скачать PDF</span>
                </button>
              </div>
            ))}
          </div>
        </motion.div>
      )}

      {/* Bottom Nav */}
      <BottomNav />
    </div>
  );
}
