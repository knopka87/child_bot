// src/types/analytics.ts

// ============================================================================
// Typed Event Parameters - используем discriminated unions для type safety
// ============================================================================

// Базовые параметры, присутствующие во всех событиях
interface BaseEventParams {
  app_version?: string;
  platform_type?: string;
  child_profile_id?: string;
}

// Onboarding Events
interface OnboardingOpenedParams extends BaseEventParams {}
interface RegistrationOpenedParams extends BaseEventParams {}
interface ConsentScreenOpenedParams extends BaseEventParams {}
interface GradeSelectedParams extends BaseEventParams {
  grade: number;
}
interface AvatarSelectedParams extends BaseEventParams {
  avatar_id: string;
}
interface DisplayNameEnteredParams extends BaseEventParams {
  name_length: number;
}
interface AdultConsentCheckedParams extends BaseEventParams {}
interface PrivacyPolicyOpenedParams extends BaseEventParams {}
interface PrivacyPolicyAcceptedParams extends BaseEventParams {}
interface TermsOpenedParams extends BaseEventParams {}
interface TermsAcceptedParams extends BaseEventParams {}
interface EmailEnteredParams extends BaseEventParams {
  email_domain?: string;
}
interface EmailVerificationSentParams extends BaseEventParams {}
interface EmailVerificationSuccessParams extends BaseEventParams {}
interface OnboardingCompletedParams extends BaseEventParams {}

// Home Events
interface HomeOpenedParams extends BaseEventParams {}
interface LevelBarViewedParams extends BaseEventParams {
  level?: number;
  progress_percent?: number;
}
interface CoinsBalanceViewedParams extends BaseEventParams {
  coins?: number;
}
interface TasksCorrectCountViewedParams extends BaseEventParams {
  count?: number;
}
interface HomeHelpClickedParams extends BaseEventParams {}
interface HomeCheckClickedParams extends BaseEventParams {}
interface UnfinishedAttemptModalShownParams extends BaseEventParams {
  attempt_id?: string;
}
interface UnfinishedAttemptContinueClickedParams extends BaseEventParams {
  attempt_id?: string;
}
interface UnfinishedAttemptNewTaskClickedParams extends BaseEventParams {}
interface MascotClickedParams extends BaseEventParams {}
interface VillainClickedParams extends BaseEventParams {
  villain_id?: string;
}
interface RecentAttemptClickedParams extends BaseEventParams {
  attempt_id?: string;
}
interface RecentAttemptsViewAllClickedParams extends BaseEventParams {}

// Help Flow Events
interface HelpSourcePickerOpenedParams extends BaseEventParams {}
interface HelpImageSelectedParams extends BaseEventParams {
  source?: 'camera' | 'file';
}
interface HelpChooseFileClickedParams extends BaseEventParams {}
interface HelpCameraClickedParams extends BaseEventParams {}
interface HelpImageUploadStartedParams extends BaseEventParams {
  file_size?: number;
  upload_source?: string;
  attempt_id?: string;
}
interface HelpImageUploadCompletedParams extends BaseEventParams {
  attempt_id?: string;
}
interface HelpImageUploadFailedParams extends BaseEventParams {
  error_message?: string;
}
interface HelpUploadRetryClickedParams extends BaseEventParams {}
interface HelpUploadBackClickedParams extends BaseEventParams {}
interface HelpProcessingStartedParams extends BaseEventParams {
  attempt_id?: string;
}
interface HelpResultOpenedParams extends BaseEventParams {
  attempt_id?: string;
}
interface HelpHintRequestedParams extends BaseEventParams {
  attempt_id?: string;
  hint_index?: number;
}
interface HelpAnswerSubmittedParams extends BaseEventParams {
  attempt_id?: string;
  answer_length?: number;
}
interface HelpAnswerResultParams extends BaseEventParams {
  attempt_id?: string;
  is_correct?: boolean;
  damage?: number;
}

// Check Flow Events
interface CheckScenarioSelectionOpenedParams extends BaseEventParams {}
interface CheckScenarioSelectedParams extends BaseEventParams {
  scenario?: string;
}
interface CheckResultOpenedParams extends BaseEventParams {
  is_correct?: boolean;
}

// Achievements Events
interface AchievementsOpenedParams extends BaseEventParams {}
interface AchievementClickedParams extends BaseEventParams {
  achievement_id?: string;
}
interface AchievementUnlockedParams extends BaseEventParams {
  achievement_id?: string;
  achievement_type?: string;
}
interface AchievementRewardClaimedParams extends BaseEventParams {
  achievement_id?: string;
  reward_type?: string;
}

// Friends Events
interface FriendsOpenedParams extends BaseEventParams {}
interface ReferralLinkCopiedParams extends BaseEventParams {}
interface ReferralLinkSharedParams extends BaseEventParams {
  platform?: string;
}
interface ReferralShareSentParams extends BaseEventParams {}
interface ReferralInviteSentParams extends BaseEventParams {}

// Profile Events
interface ProfileOpenedParams extends BaseEventParams {}
interface ProfileHistoryClickedParams extends BaseEventParams {}
interface ProfileReportClickedParams extends BaseEventParams {}
interface ProfileSupportClickedParams extends BaseEventParams {}
interface ProfileSubscriptionClickedParams extends BaseEventParams {}
interface HistoryOpenedParams extends BaseEventParams {}
interface HistoryAttemptClickedParams extends BaseEventParams {
  attempt_id?: string;
}
interface HistoryItemClickedParams extends BaseEventParams {
  attempt_id?: string;
  mode?: 'help' | 'check';
  status?: string;
}
interface HistoryBackClickedParams extends BaseEventParams {}
interface HistoryFilterChangedParams extends BaseEventParams {
  filter_type?: string;
  filter_value?: string;
}
interface HistoryRetryClickedParams extends BaseEventParams {
  attempt_id?: string;
  mode?: 'help' | 'check';
}
interface HistoryFixErrorsClickedParams extends BaseEventParams {
  attempt_id?: string;
}

// Villain Events
interface VillainScreenOpenedParams extends BaseEventParams {
  villain_id?: string;
}
interface VillainTauntViewedParams extends BaseEventParams {
  villain_id?: string;
}
interface VillainHealthChangedParams extends BaseEventParams {
  villain_id?: string;
  health_percent?: number;
}
interface VillainVictoryTriggeredParams extends BaseEventParams {
  villain_id?: string;
}
interface VictoryScreenOpenedParams extends BaseEventParams {
  villain_id?: string;
  attempt_id?: string;
}
interface VictoryRewardViewedParams extends BaseEventParams {
  villain_id?: string;
  reward_type?: string;
  reward_id?: string;
}
interface VictoryContinueClickedParams extends BaseEventParams {
  villain_id?: string;
}

// Support Events
interface SupportOpenedParams extends BaseEventParams {}
interface SupportMessageSentParams extends BaseEventParams {
  message_length?: number;
}

// Paywall Events
interface PaywallOpenedParams extends BaseEventParams {
  source?: string;
}
interface PaymentStartedParams extends BaseEventParams {
  plan?: string;
  amount?: number;
}
interface PaymentSuccessParams extends BaseEventParams {
  plan?: string;
  amount?: number;
}
interface PaymentFailedParams extends BaseEventParams {
  plan?: string;
  error?: string;
}

// ============================================================================
// Discriminated Union для всех событий
// ============================================================================

export type AnalyticsEvent =
  // Onboarding
  | { name: 'onboarding_opened'; params: OnboardingOpenedParams }
  | { name: 'registration_opened'; params: RegistrationOpenedParams }
  | { name: 'consent_screen_opened'; params: ConsentScreenOpenedParams }
  | { name: 'grade_selected'; params: GradeSelectedParams }
  | { name: 'avatar_selected'; params: AvatarSelectedParams }
  | { name: 'display_name_entered'; params: DisplayNameEnteredParams }
  | { name: 'adult_consent_checked'; params: AdultConsentCheckedParams }
  | { name: 'privacy_policy_opened'; params: PrivacyPolicyOpenedParams }
  | { name: 'privacy_policy_accepted'; params: PrivacyPolicyAcceptedParams }
  | { name: 'terms_opened'; params: TermsOpenedParams }
  | { name: 'terms_accepted'; params: TermsAcceptedParams }
  | { name: 'email_entered'; params: EmailEnteredParams }
  | { name: 'email_verification_sent'; params: EmailVerificationSentParams }
  | { name: 'email_verification_success'; params: EmailVerificationSuccessParams }
  | { name: 'onboarding_completed'; params: OnboardingCompletedParams }
  // Home
  | { name: 'home_opened'; params: HomeOpenedParams }
  | { name: 'level_bar_viewed'; params: LevelBarViewedParams }
  | { name: 'coins_balance_viewed'; params: CoinsBalanceViewedParams }
  | { name: 'tasks_correct_count_viewed'; params: TasksCorrectCountViewedParams }
  | { name: 'home_help_clicked'; params: HomeHelpClickedParams }
  | { name: 'home_check_clicked'; params: HomeCheckClickedParams }
  | { name: 'unfinished_attempt_modal_shown'; params: UnfinishedAttemptModalShownParams }
  | { name: 'unfinished_attempt_continue_clicked'; params: UnfinishedAttemptContinueClickedParams }
  | { name: 'unfinished_attempt_new_task_clicked'; params: UnfinishedAttemptNewTaskClickedParams }
  | { name: 'mascot_clicked'; params: MascotClickedParams }
  | { name: 'villain_clicked'; params: VillainClickedParams }
  | { name: 'recent_attempt_clicked'; params: RecentAttemptClickedParams }
  | { name: 'recent_attempts_view_all_clicked'; params: RecentAttemptsViewAllClickedParams }
  // Help Flow
  | { name: 'help_source_picker_opened'; params: HelpSourcePickerOpenedParams }
  | { name: 'help_image_selected'; params: HelpImageSelectedParams }
  | { name: 'help_choose_file_clicked'; params: HelpChooseFileClickedParams }
  | { name: 'help_camera_clicked'; params: HelpCameraClickedParams }
  | { name: 'help_image_upload_started'; params: HelpImageUploadStartedParams }
  | { name: 'help_image_upload_completed'; params: HelpImageUploadCompletedParams }
  | { name: 'help_image_upload_failed'; params: HelpImageUploadFailedParams }
  | { name: 'help_upload_retry_clicked'; params: HelpUploadRetryClickedParams }
  | { name: 'help_upload_back_clicked'; params: HelpUploadBackClickedParams }
  | { name: 'help_processing_started'; params: HelpProcessingStartedParams }
  | { name: 'help_result_opened'; params: HelpResultOpenedParams }
  | { name: 'help_hint_requested'; params: HelpHintRequestedParams }
  | { name: 'help_answer_submitted'; params: HelpAnswerSubmittedParams }
  | { name: 'help_answer_result'; params: HelpAnswerResultParams }
  // Check Flow
  | { name: 'check_scenario_selection_opened'; params: CheckScenarioSelectionOpenedParams }
  | { name: 'check_scenario_selected'; params: CheckScenarioSelectedParams }
  | { name: 'check_result_opened'; params: CheckResultOpenedParams }
  // Achievements
  | { name: 'achievements_opened'; params: AchievementsOpenedParams }
  | { name: 'achievement_clicked'; params: AchievementClickedParams }
  | { name: 'achievement_unlocked'; params: AchievementUnlockedParams }
  | { name: 'achievement_reward_claimed'; params: AchievementRewardClaimedParams }
  // Friends
  | { name: 'friends_opened'; params: FriendsOpenedParams }
  | { name: 'referral_link_copied'; params: ReferralLinkCopiedParams }
  | { name: 'referral_link_shared'; params: ReferralLinkSharedParams }
  | { name: 'referral_share_sent'; params: ReferralShareSentParams }
  | { name: 'referral_invite_sent'; params: ReferralInviteSentParams }
  // Profile
  | { name: 'profile_opened'; params: ProfileOpenedParams }
  | { name: 'profile_history_clicked'; params: ProfileHistoryClickedParams }
  | { name: 'profile_report_clicked'; params: ProfileReportClickedParams }
  | { name: 'profile_support_clicked'; params: ProfileSupportClickedParams }
  | { name: 'profile_subscription_clicked'; params: ProfileSubscriptionClickedParams }
  | { name: 'history_opened'; params: HistoryOpenedParams }
  | { name: 'history_attempt_clicked'; params: HistoryAttemptClickedParams }
  | { name: 'history_item_clicked'; params: HistoryItemClickedParams }
  | { name: 'history_back_clicked'; params: HistoryBackClickedParams }
  | { name: 'history_filter_changed'; params: HistoryFilterChangedParams }
  | { name: 'history_retry_clicked'; params: HistoryRetryClickedParams }
  | { name: 'history_fix_errors_clicked'; params: HistoryFixErrorsClickedParams }
  // Villain
  | { name: 'villain_screen_opened'; params: VillainScreenOpenedParams }
  | { name: 'villain_taunt_viewed'; params: VillainTauntViewedParams }
  | { name: 'villain_health_changed'; params: VillainHealthChangedParams }
  | { name: 'villain_victory_triggered'; params: VillainVictoryTriggeredParams }
  | { name: 'victory_screen_opened'; params: VictoryScreenOpenedParams }
  | { name: 'victory_reward_viewed'; params: VictoryRewardViewedParams }
  | { name: 'victory_continue_clicked'; params: VictoryContinueClickedParams }
  // Support
  | { name: 'support_opened'; params: SupportOpenedParams }
  | { name: 'support_message_sent'; params: SupportMessageSentParams }
  // Paywall
  | { name: 'paywall_opened'; params: PaywallOpenedParams }
  | { name: 'payment_started'; params: PaymentStartedParams }
  | { name: 'payment_success'; params: PaymentSuccessParams }
  | { name: 'payment_failed'; params: PaymentFailedParams };

// Helper type для получения типа события по имени
export type AnalyticsEventName = AnalyticsEvent['name'];
export type AnalyticsEventParams<T extends AnalyticsEventName> = Extract<
  AnalyticsEvent,
  { name: T }
>['params'];

// Структура события для хранения
export interface StoredAnalyticsEvent {
  name: AnalyticsEventName;
  timestamp: number;
  sessionId: string;
  params: Record<string, any>; // Здесь остается any для сериализации
}

export interface AnalyticsConfig {
  enabled: boolean;
  debug: boolean;
  batchSize: number;
  batchInterval: number;
  retryAttempts: number;
  retryDelay: number;
  platforms: AnalyticsPlatform[];
}

export type AnalyticsPlatform = 'vk' | 'max' | 'backend';

export interface UserProperties {
  platform_type?: 'vk' | 'max' | 'web';
  subscription_status?: 'trial' | 'active' | 'expired' | 'cancelled';
  trial_status?: string;
  email_verified?: boolean;
  weekly_report_enabled?: boolean;
  report_archive_enabled?: boolean;
  grade?: number;
  level?: number;
  coins_balance?: number;
  tasks_solved_correct_count?: number;
  wins_count?: number;
  checks_correct_count?: number;
  current_streak_days?: number;
  has_unfinished_attempt?: boolean;
  active_villain_id?: string;
  active_villain_health_percent?: number;
  invited_count_total?: number;
  achievements_unlocked_count?: number;
}

export interface AnalyticsSession {
  id: string;
  startedAt: number;
  platform: string;
  appVersion: string;
  userId?: string;
  childProfileId?: string;
}
