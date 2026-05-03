// src/services/platform/adapters/WebAdapter.ts
import type {
  PlatformInfo,
  PlatformUser,
  PlatformTheme,
  ShareOptions,
  HapticFeedbackType,
} from '@/types/platform';

export class WebAdapter {
  async init(): Promise<void> {
    // Web doesn't need initialization
  }

  getInfo(): PlatformInfo {
    return {
      type: 'web',
      version: '1.0.0',
      isSupported: true,
      features: {
        nativeNavigation: false,
        nativeShare: typeof navigator.share !== 'undefined',
        nativeCamera: false,
        nativeFilePicker: true,
        hapticFeedback: false,
        nativeAnalytics: false,
        nativePayment: false,
        localStorage: true,
        cloudStorage: false,
        pushNotifications: false,
        nativeAuth: false,
        biometricAuth: false,
      },
    };
  }

  async getUser(): Promise<PlatformUser> {
    // Generate or retrieve unique user ID for web platform
    // ВАЖНО: НЕ использовать константу 'web-user' - это создаёт коллизии!
    const STORAGE_KEY = 'web_platform_user_id';

    let userId = localStorage.getItem(STORAGE_KEY);
    if (!userId) {
      // Generate new unique ID
      userId = `web-${crypto.randomUUID()}`;
      localStorage.setItem(STORAGE_KEY, userId);
      console.log('[WebAdapter] Generated new user ID:', userId);
    }

    return {
      id: userId,
      firstName: 'Web',
      lastName: 'User',
    };
  }

  getTheme(): PlatformTheme {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

    return {
      colorScheme: isDark ? 'dark' : 'light',
      backgroundColor: isDark ? '#000' : '#fff',
      textColor: isDark ? '#fff' : '#000',
      buttonColor: '#0077ff',
      buttonTextColor: '#fff',
      hintColor: '#818c99',
      linkColor: '#0077ff',
      secondaryBackgroundColor: isDark ? '#1a1a1a' : '#f5f5f5',
    };
  }

  async share(options: ShareOptions): Promise<void> {
    if (navigator.share) {
      await navigator.share({
        title: options.title,
        text: options.text,
        url: options.url || window.location.href,
      });
    } else {
      const url = options.url || window.location.href;
      await navigator.clipboard.writeText(url);
    }
  }

  hapticFeedback(_feedback: HapticFeedbackType): void {
    // Web doesn't support haptic feedback
  }

  openLink(url: string): void {
    window.open(url, '_blank');
  }

  close(): void {
    window.close();
  }

  ready(): void {
    // Web doesn't need ready signal
  }

  expand(): void {
    // Web doesn't need expand
  }
}
