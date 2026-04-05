// src/stores/onboardingStore.ts
import { create } from 'zustand';
import type { OnboardingStep, OnboardingState } from '@/types/onboarding';

interface OnboardingStore extends OnboardingState {
  // Actions
  setStep: (step: OnboardingStep) => void;
  setGrade: (grade: number) => void;
  setAvatar: (avatarId: string) => void;
  setDisplayName: (displayName: string) => void;
  setAdultConsent: (adultConsent: boolean) => void;
  setPrivacyAccepted: (privacyAccepted: boolean) => void;
  setTermsAccepted: (termsAccepted: boolean) => void;
  setEmail: (email: string) => void;
  setEmailVerified: (emailVerified: boolean) => void;
  reset: () => void;

  // Computed
  canProceed: () => boolean;
  progressPercent: () => number;
}

const initialState: OnboardingState = {
  currentStep: 'welcome',
  grade: null,
  avatarId: null,
  displayName: '',
  adultConsent: false,
  privacyAccepted: false,
  termsAccepted: false,
  email: '',
  emailVerified: false,
};

export const useOnboardingStore = create<OnboardingStore>((set, get) => ({
  ...initialState,

  setStep: (currentStep) => set({ currentStep }),
  setGrade: (grade) => set({ grade }),
  setAvatar: (avatarId) => set({ avatarId }),
  setDisplayName: (displayName) => set({ displayName }),
  setAdultConsent: (adultConsent) => set({ adultConsent }),
  setPrivacyAccepted: (privacyAccepted) => set({ privacyAccepted }),
  setTermsAccepted: (termsAccepted) => set({ termsAccepted }),
  setEmail: (email) => set({ email }),
  setEmailVerified: (emailVerified) => set({ emailVerified }),

  reset: () => set(initialState),

  canProceed: () => {
    const state = get();
    switch (state.currentStep) {
      case 'welcome':
        return true;
      case 'grade':
        return state.grade !== null;
      case 'avatar':
        return state.avatarId !== null;
      case 'display_name':
        return state.displayName.trim().length >= 2;
      case 'consent':
        return (
          state.adultConsent &&
          state.privacyAccepted &&
          state.termsAccepted
        );
      case 'email':
        return state.email.includes('@');
      case 'email_verification':
        return state.emailVerified;
      default:
        return true;
    }
  },

  progressPercent: () => {
    const state = get();
    const steps: OnboardingStep[] = [
      'welcome',
      'grade',
      'avatar',
      'display_name',
      'consent',
      'email',
      'email_verification',
      'completed',
    ];
    const currentIndex = steps.indexOf(state.currentStep);
    return Math.round(((currentIndex + 1) / steps.length) * 100);
  },
}));
