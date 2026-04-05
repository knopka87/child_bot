// src/types/onboarding.ts

export type OnboardingStep =
  | 'welcome'
  | 'grade'
  | 'avatar'
  | 'display_name'
  | 'consent'
  | 'email'
  | 'email_verification'
  | 'completed';

export interface OnboardingState {
  currentStep: OnboardingStep;
  grade: number | null;
  avatarId: string | null;
  displayName: string;
  adultConsent: boolean;
  privacyAccepted: boolean;
  termsAccepted: boolean;
  email: string;
  emailVerified: boolean;
}

export interface Avatar {
  id: string;
  imageUrl: string;
  name: string;
  isPremium: boolean;
}

export interface ConsentDocument {
  type: 'privacy_policy' | 'terms_of_service';
  version: string;
  url: string;
  acceptedAt?: string;
}

export interface OnboardingProgress {
  step: OnboardingStep;
  completedSteps: OnboardingStep[];
  totalSteps: number;
  progressPercent: number;
}
