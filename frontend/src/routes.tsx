// src/routes.tsx
import { Routes, Route, Navigate } from 'react-router-dom';
import { ROUTES } from '@/config/routes';

// Pages
import HomePage from '@/pages/Home/HomePage';
import { OnboardingPageNew as OnboardingPage } from '@/pages/Onboarding/OnboardingPageNew';
import { HelpPage } from '@/pages/Help/HelpPage';
import { PrivacyPolicyPage } from '@/pages/Legal/PrivacyPolicyPage';
import { TermsPage } from '@/pages/Legal/TermsPage';
import HelpImageUploader from '@/pages/Help/HelpImageUploader';
import UploadPage from '@/pages/Help/UploadPage';
import ProcessingPage from '@/pages/Help/ProcessingPage';
import ResultPage from '@/pages/Help/ResultPage';
import { CheckPage } from '@/pages/Check/CheckPage';
import ScenarioSelection from '@/pages/Check/ScenarioSelection';
import CheckImageSetManager from '@/pages/Check/CheckImageSetManager';
import CheckQualityTwoPage from '@/pages/Check/CheckQualityTwoPage';
import CheckQualitySinglePage from '@/pages/Check/CheckQualitySinglePage';
import CheckProcessingPage from '@/pages/Check/CheckProcessingPage';
import CheckResultPage from '@/pages/Check/CheckResultPage';
import { AchievementsPage } from '@/pages/Achievements/AchievementsPage';
import { FriendsPage } from '@/pages/Friends/FriendsPage';
import { ProfilePage } from '@/pages/Profile/ProfilePage';
import { HistoryPage } from '@/pages/Profile/History/HistoryPage';
import ReportPage from '@/pages/Report/ReportPage';
import { VillainPage } from '@/pages/Villain/VillainPage';
import { VictoryPage } from '@/pages/Villain/VictoryPage';

export function AppRoutes() {
  return (
    <Routes>
      <Route path={ROUTES.HOME} element={<HomePage />} />
      <Route path={ROUTES.ONBOARDING} element={<OnboardingPage />} />

      {/* Help Flow */}
      <Route path={ROUTES.HELP} element={<HelpPage />} />
      <Route path={ROUTES.HELP_UPLOAD} element={<HelpImageUploader />} />
      <Route path="/help/upload-progress" element={<UploadPage />} />
      <Route path={ROUTES.HELP_PROCESSING} element={<ProcessingPage />} />
      <Route path="/help/result/:attemptId" element={<ResultPage />} />

      {/* Check Flow */}
      <Route path={ROUTES.CHECK} element={<CheckPage />} />
      <Route path={ROUTES.CHECK_SCENARIO} element={<ScenarioSelection />} />
      <Route path="/check/upload-images" element={<CheckImageSetManager />} />
      <Route path="/check/quality-two" element={<CheckQualityTwoPage />} />
      <Route path="/check/quality-single" element={<CheckQualitySinglePage />} />
      <Route path={ROUTES.CHECK_PROCESSING} element={<CheckProcessingPage />} />
      <Route path="/check/result" element={<CheckResultPage />} />

      <Route path={ROUTES.ACHIEVEMENTS} element={<AchievementsPage />} />
      <Route path={ROUTES.FRIENDS} element={<FriendsPage />} />
      <Route path={ROUTES.PROFILE} element={<ProfilePage />} />
      <Route path={ROUTES.PROFILE_HISTORY} element={<HistoryPage />} />
      <Route path={ROUTES.PROFILE_REPORT} element={<ReportPage />} />
      <Route path={ROUTES.VILLAIN} element={<VillainPage />} />
      <Route path={ROUTES.VILLAIN_VICTORY} element={<VictoryPage />} />

      {/* Legal */}
      <Route path={ROUTES.LEGAL_PRIVACY} element={<PrivacyPolicyPage />} />
      <Route path={ROUTES.LEGAL_TERMS} element={<TermsPage />} />

      {/* Fallback */}
      <Route path="*" element={<Navigate to={ROUTES.HOME} replace />} />
    </Routes>
  );
}
