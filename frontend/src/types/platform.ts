// src/types/platform.ts
export type PlatformType = 'vk' | 'max' | 'telegram' | 'web';

export interface PlatformInfo {
  type: PlatformType;
  version: string;
  isSupported: boolean;
  features: PlatformFeatures;
}

export interface PlatformFeatures {
  nativeNavigation: boolean;
  nativeShare: boolean;
  nativeCamera: boolean;
  nativeFilePicker: boolean;
  hapticFeedback: boolean;
  nativeAnalytics: boolean;
  nativePayment: boolean;
  localStorage: boolean;
  cloudStorage: boolean;
  pushNotifications: boolean;
  nativeAuth: boolean;
  biometricAuth: boolean;
}

export interface PlatformUser {
  id: string;
  firstName?: string;
  lastName?: string;
  photoUrl?: string;
  languageCode?: string;
  isPremium?: boolean;
}

export interface PlatformTheme {
  colorScheme: 'light' | 'dark';
  backgroundColor: string;
  textColor: string;
  buttonColor: string;
  buttonTextColor: string;
  hintColor: string;
  linkColor: string;
  secondaryBackgroundColor: string;
}

export interface ShareOptions {
  title?: string;
  text?: string;
  url?: string;
  files?: File[];
}

export interface HapticFeedbackType {
  type: 'impact' | 'notification' | 'selection';
  style?: 'light' | 'medium' | 'heavy' | 'soft' | 'rigid';
}
