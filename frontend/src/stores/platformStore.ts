// src/stores/platformStore.ts
import { create } from 'zustand';
import { type Platform } from '@/lib/platform/platform-detection';

interface PlatformState {
  platform: Platform | null;
  isIOS: boolean;
  isAndroid: boolean;
  isDesktop: boolean;
  isMobile: boolean;

  // Actions
  setPlatform: (platformInfo: {
    platform: Platform;
    isIOS: boolean;
    isAndroid: boolean;
    isDesktop: boolean;
    isMobile: boolean;
  }) => void;
}

export const usePlatformStore = create<PlatformState>((set) => ({
  platform: null,
  isIOS: false,
  isAndroid: false,
  isDesktop: false,
  isMobile: false,

  setPlatform: (platformInfo) => set(platformInfo),
}));
