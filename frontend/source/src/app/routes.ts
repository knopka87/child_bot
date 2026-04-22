import { createBrowserRouter } from "react-router";
import { Layout } from "./components/Layout";
import { MobileFrame } from "./components/MobileFrame";
import { Onboarding } from "./components/screens/Onboarding";
import { HomeScreen } from "./components/screens/Home";
import { UploadScreen } from "./components/screens/UploadScreen";
import { ImageQualityCheck } from "./components/screens/ImageQualityCheck";
import { CropScreen } from "./components/screens/CropScreen";
import { ImageSetManager } from "./components/screens/ImageSetManager";
import { ProcessingScreen } from "./components/screens/ProcessingScreen";
import { HelpResult } from "./components/screens/HelpResult";
import { CheckScenario } from "./components/screens/CheckScenario";
import { CheckResult } from "./components/screens/CheckResult";
import { VillainScreen } from "./components/screens/VillainScreen";
import { VictoryScreen } from "./components/screens/VictoryScreen";
import { AchievementsScreen } from "./components/screens/Achievements";
import { FriendsScreen } from "./components/screens/Friends";
import { ProfileScreen } from "./components/screens/Profile";
import { HistoryScreen } from "./components/screens/HistoryScreen";
import { ParentReport } from "./components/screens/ParentReport";
import { PaymentScreen } from "./components/screens/Payment";
import { HelpScreen } from "./components/screens/HelpScreen";

export const router = createBrowserRouter([
  // Standalone screens (no tab bar)
  {
    Component: MobileFrame,
    children: [
      { path: "/onboarding", Component: Onboarding },
      { path: "/payment", Component: PaymentScreen },
      { path: "/villain", Component: VillainScreen },
      { path: "/victory", Component: VictoryScreen },
      // Help flow
      { path: "/help/upload", Component: UploadScreen },
      { path: "/help/images", Component: ImageSetManager },
      { path: "/help/quality", Component: ImageQualityCheck },
      { path: "/help/crop", Component: CropScreen },
      { path: "/help/processing", Component: ProcessingScreen },
      { path: "/help/result", Component: HelpResult },
      // Check flow
      { path: "/check/scenario", Component: CheckScenario },
      { path: "/check/upload", Component: UploadScreen },
      { path: "/check/images", Component: ImageSetManager },
      { path: "/check/quality", Component: ImageQualityCheck },
      { path: "/check/crop", Component: CropScreen },
      { path: "/check/processing", Component: ProcessingScreen },
      { path: "/check/result", Component: CheckResult },
    ],
  },
  // Main app with tab bar
  {
    path: "/",
    Component: Layout,
    children: [
      { index: true, Component: HomeScreen },
      { path: "achievements", Component: AchievementsScreen },
      { path: "friends", Component: FriendsScreen },
      { path: "profile", Component: ProfileScreen },
      { path: "profile/history", Component: HistoryScreen },
      { path: "profile/report", Component: ParentReport },
      { path: "profile/help", Component: HelpScreen },
    ],
  },
]);
