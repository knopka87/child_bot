// src/services/platform/adapters/MaxAdapter.ts
/**
 * MAX Bridge Adapter для интеграции с мессенджером Max (Mail.ru Group)
 * Документация: https://dev.max.ru/docs/webapps/bridge
 */

import { createLogger } from '@/lib/logger';
import type {
  PlatformInfo,
  PlatformUser,
  PlatformTheme,
  ShareOptions,
  HapticFeedbackType,
} from '@/types/platform';

const logger = createLogger('MaxAdapter');

// Типизация MAX Bridge API
interface MaxBridgeAPI {
  init(): Promise<void>;
  getUserInfo(): Promise<MaxUserInfo>;
  share(options: MaxShareOptions): Promise<void>;
  hapticFeedback(type: string, style?: string): void;
  openLink(url: string): void;
  close(): void;
  ready(): void;
  expand(): void;
  getTheme(): MaxThemeInfo;
  subscribe(event: string, callback: (data: any) => void): void;
}

interface MaxUserInfo {
  id: string;
  first_name?: string;
  last_name?: string;
  photo_url?: string;
  language_code?: string;
}

interface MaxShareOptions {
  url?: string;
  text?: string;
  title?: string;
}

interface MaxThemeInfo {
  theme: 'light' | 'dark';
  backgroundColor: string;
  textColor: string;
  buttonColor: string;
  buttonTextColor: string;
  hintColor: string;
  linkColor: string;
}

declare global {
  interface Window {
    MaxBridge?: MaxBridgeAPI;
  }
}

export class MaxAdapter {
  private bridge: MaxBridgeAPI | null = null;
  private isInitialized = false;

  async init(): Promise<void> {
    if (this.isInitialized) return;

    try {
      // Проверяем доступность MAX Bridge
      if (!window.MaxBridge) {
        throw new Error('MAX Bridge is not available');
      }

      this.bridge = window.MaxBridge;
      await this.bridge.init();
      this.isInitialized = true;

      // Подписываемся на изменения темы
      this.bridge.subscribe('themeChanged', (data) => {
        logger.debug('Theme changed', { data });
      });

      logger.info('Initialized successfully');
    } catch (error) {
      logger.error('Failed to initialize', { error });
      throw error;
    }
  }

  getInfo(): PlatformInfo {
    return {
      type: 'max',
      version: '1.0.0',
      isSupported: true,
      features: {
        nativeNavigation: true,
        nativeShare: true,
        nativeCamera: true, // Max поддерживает камеру
        nativeFilePicker: true,
        hapticFeedback: true,
        nativeAnalytics: true,
        nativePayment: true, // Max Pay
        localStorage: true,
        cloudStorage: true,
        pushNotifications: true,
        nativeAuth: true,
        biometricAuth: false,
      },
    };
  }

  async getUser(): Promise<PlatformUser> {
    if (!this.bridge) {
      throw new Error('MAX Bridge not initialized');
    }

    try {
      const user = await this.bridge.getUserInfo();
      return {
        id: user.id,
        firstName: user.first_name,
        lastName: user.last_name,
        photoUrl: user.photo_url,
        languageCode: user.language_code,
      };
    } catch (error) {
      logger.error('Failed to get user', { error });
      throw error;
    }
  }

  getTheme(): PlatformTheme {
    if (!this.bridge) {
      return this.getDefaultTheme();
    }

    try {
      const maxTheme = this.bridge.getTheme();

      return {
        colorScheme: maxTheme.theme,
        backgroundColor: maxTheme.backgroundColor || '#fff',
        textColor: maxTheme.textColor || '#000',
        buttonColor: maxTheme.buttonColor || '#0077ff',
        buttonTextColor: maxTheme.buttonTextColor || '#fff',
        hintColor: maxTheme.hintColor || '#818c99',
        linkColor: maxTheme.linkColor || '#0077ff',
        secondaryBackgroundColor:
          maxTheme.backgroundColor === '#000' ? '#1a1a1a' : '#f5f5f5',
      };
    } catch (error) {
      logger.error('Failed to get theme', { error });
      return this.getDefaultTheme();
    }
  }

  private getDefaultTheme(): PlatformTheme {
    return {
      colorScheme: 'light',
      backgroundColor: '#fff',
      textColor: '#000',
      buttonColor: '#0077ff',
      buttonTextColor: '#fff',
      hintColor: '#818c99',
      linkColor: '#0077ff',
      secondaryBackgroundColor: '#f5f5f5',
    };
  }

  async share(options: ShareOptions): Promise<void> {
    if (!this.bridge) {
      throw new Error('MAX Bridge not initialized');
    }

    try {
      await this.bridge.share({
        url: options.url || window.location.href,
        text: options.text,
        title: options.title,
      });
      logger.debug('Share completed', { options });
    } catch (error) {
      logger.error('Failed to share', { error });
      throw error;
    }
  }

  hapticFeedback(feedback: HapticFeedbackType): void {
    if (!this.bridge) return;

    try {
      const type = feedback.type;
      const style = feedback.style || 'medium';

      this.bridge.hapticFeedback(type, style);
      logger.debug('Haptic feedback triggered', { type, style });
    } catch (error) {
      logger.error('Haptic feedback failed', { error });
    }
  }

  openLink(url: string): void {
    if (!this.bridge) {
      window.open(url, '_blank');
      return;
    }

    try {
      this.bridge.openLink(url);
      logger.debug('Link opened', { url });
    } catch (error) {
      logger.error('Failed to open link', { url, error });
      window.open(url, '_blank');
    }
  }

  close(): void {
    if (!this.bridge) return;

    try {
      this.bridge.close();
      logger.debug('App closed');
    } catch (error) {
      logger.error('Failed to close', { error });
    }
  }

  ready(): void {
    if (!this.bridge) return;

    try {
      this.bridge.ready();
      logger.debug('Ready signal sent');
    } catch (error) {
      logger.error('Failed to send ready', { error });
    }
  }

  expand(): void {
    if (!this.bridge) return;

    try {
      this.bridge.expand();
      logger.debug('App expanded');
    } catch (error) {
      logger.error('Failed to expand', { error });
    }
  }
}
