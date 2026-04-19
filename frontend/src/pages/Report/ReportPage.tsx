// src/pages/Report/ReportPage.tsx
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Download, Mail, Send } from 'lucide-react';
import { motion } from 'framer-motion';
import { useAnalytics } from '@/hooks/useAnalytics';
import { useProfileStore } from '@/stores/profileStore';
import { profileAPI } from '@/api/profile';
import { ROUTES } from '@/config/routes';
import { BottomNav } from '@/components/layout/BottomNav';

interface PastReport {
  id: string;
  date: string;
  status: string;
}

const MOCK_PAST_REPORTS: PastReport[] = [
  { id: '1', date: '10 марта 2026', status: 'Отправлен' },
  { id: '2', date: '3 марта 2026', status: 'Отправлен' },
  { id: '3', date: '24 февраля 2026', status: 'Отправлен' },
];

export default function ReportPage() {
  const navigate = useNavigate();
  const analytics = useAnalytics();
  const profile = useProfileStore((state) => state.profile);

  const [email, setEmail] = useState('');
  const [weeklyEnabled, setWeeklyEnabled] = useState(true);
  const [pastReports] = useState<PastReport[]>(MOCK_PAST_REPORTS);
  const [isLoading, setIsLoading] = useState(false);

  // Загружаем настройки отчётов при монтировании
  useEffect(() => {
    const loadSettings = async () => {
      if (!profile?.child_profile_id) return;

      try {
        const settings = await profileAPI.getReportSettings(profile.child_profile_id);
        setEmail(settings.email);
        setWeeklyEnabled(settings.weeklyReportEnabled);
      } catch {
        setEmail('');
        setWeeklyEnabled(true);
      }
    };

    loadSettings();

    analytics.trackEvent('profile_report_clicked', {
      child_profile_id: profile?.child_profile_id,
    });
  }, [profile?.child_profile_id, analytics]);

  const handleEmailChange = async (newEmail: string) => {
    setEmail(newEmail);

    if (profile?.child_profile_id && newEmail.includes('@')) {
      try {
        await profileAPI.updateReportSettings(profile.child_profile_id, { email: newEmail });
        analytics.trackEvent('report_email_changed' as any, {
          child_profile_id: profile.child_profile_id,
          email_domain: newEmail.split('@')[1],
        });
      } catch (error) {
        console.error('[ReportPage] Failed to update email:', error);
      }
    }
  };

  const handleWeeklyToggle = async () => {
    const newValue = !weeklyEnabled;
    setWeeklyEnabled(newValue);

    if (profile?.child_profile_id) {
      try {
        await profileAPI.updateReportSettings(profile.child_profile_id, {
          weeklyReportEnabled: newValue,
        });
        analytics.trackEvent('weekly_report_toggled' as any, {
          child_profile_id: profile.child_profile_id,
          enabled: newValue,
        });
      } catch (error) {
        console.error('[ReportPage] Failed to toggle weekly report:', error);
        setWeeklyEnabled(!newValue);
      }
    }
  };

  const handleSendTestReport = async () => {
    if (!profile?.child_profile_id) return;

    setIsLoading(true);
    try {
      await profileAPI.sendTestReport(profile.child_profile_id);
      analytics.trackEvent('test_report_sent' as any, {
        child_profile_id: profile.child_profile_id,
      });
      alert('✅ Тестовый отчёт отправлен!');
    } catch (error) {
      console.error('[ReportPage] Failed to send test report:', error);
      alert('⚠️ Не удалось отправить тестовый отчёт');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDownloadReport = async (reportId: string) => {
    if (!profile?.child_profile_id) return;

    try {
      analytics.trackEvent('report_download_clicked' as any, {
        child_profile_id: profile.child_profile_id,
        report_id: reportId,
      });
      alert('📄 Отчёт скоро будет доступен для скачивания');
    } catch (error) {
      console.error('[ReportPage] Failed to download report:', error);
    }
  };

  return (
    <div className="flex flex-col min-h-full px-5 pt-4 pb-6 bg-gradient-to-b from-[#F0F4FF] to-background">
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
              weeklyEnabled ? 'bg-primary' : 'bg-switch-background'
            }`}
          >
            <div
              className={`w-5 h-5 bg-white rounded-full absolute top-1 transition-all shadow-sm ${
                weeklyEnabled ? 'right-1' : 'left-1'
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
          (!email || isLoading) ? 'opacity-50 cursor-not-allowed' : ''
        }`}
      >
        <Send size={18} />
        {isLoading ? 'Отправляем...' : 'Отправить тестовый отчёт'}
      </motion.button>

      {/* Archive section */}
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.2 }}
      >
        <h3 className="text-foreground mb-3">Архив отчётов</h3>
        <div className="flex flex-col gap-2">
          {pastReports.map((report) => (
            <div key={report.id} className="bg-white rounded-2xl p-3 flex items-center justify-between shadow-sm">
              <div>
                <p className="text-[14px] text-foreground">{report.date}</p>
                <p className="text-[12px] text-[#00B894]">{report.status}</p>
              </div>
              <button
                onClick={() => handleDownloadReport(report.id)}
                className="w-9 h-9 bg-primary/10 rounded-xl flex items-center justify-center"
              >
                <Download size={18} className="text-primary" />
              </button>
            </div>
          ))}
        </div>
      </motion.div>

      {/* Bottom Nav */}
      <BottomNav />
    </div>
  );
}
