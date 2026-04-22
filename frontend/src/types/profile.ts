// src/types/profile.ts

export interface ProfileData {
  id: string;
  displayName: string;
  avatarId: string;
  avatarUrl: string;
  grade: number;
  subscription: SubscriptionData;
}

export interface SubscriptionData {
  status: 'trial' | 'active' | 'expired' | 'cancelled';
  planId?: string;
  planName?: string;
  trialDaysRemaining?: number;
  expiresAt?: string;
}

export interface HistoryAttempt {
  id: string;
  mode: 'help' | 'check';
  status: 'success' | 'error' | 'in_progress' | 'completed' | 'failed';
  scenarioType?: 'single_photo' | 'two_photo';
  createdAt: string;
  completedAt?: string;
  images: HistoryImage[];
  result?: HistoryResult;
  hintsUsed?: number;
}

export interface HistoryImage {
  id: string;
  role: 'task' | 'answer' | 'single';
  url: string;
  thumbnailUrl: string;
}

export interface HistoryResult {
  status: 'correct' | 'has_errors' | 'processing';
  errorCount?: number;
  feedback?: ErrorFeedback[];
  summary?: string;
}

export interface ErrorFeedback {
  id: string;
  stepNumber?: number;
  lineReference?: string;
  description: string;
  locationType: 'line' | 'step' | 'general';
}

export interface HistoryFilters {
  mode?: 'help' | 'check' | 'all';
  status?: 'success' | 'error' | 'in_progress' | 'all';
  dateFrom?: string;
  dateTo?: string;
}

export interface ReportSettings {
  email: string;
  emailVerified: boolean;
  weeklyReportEnabled: boolean;
  archiveEnabled: boolean;
}

export interface WeeklyReport {
  id: string;
  periodStart: string;
  periodEnd: string;
  generatedAt: string;
  downloadUrl?: string;
  stats: {
    totalAttempts: number;
    successfulAttempts: number;
    errorsFixed: number;
    streakDays: number;
  };
}

export interface ReportSettings {
  email: string;
  emailVerified: boolean;
  weeklyReportEnabled: boolean;
  archiveEnabled: boolean;
}

export interface WeeklyReport {
  id: string;
  periodStart: string;
  periodEnd: string;
  generatedAt: string;
  downloadUrl?: string;
  stats: {
    totalAttempts: number;
    successfulAttempts: number;
    errorsFixed: number;
    streakDays: number;
  };
}
