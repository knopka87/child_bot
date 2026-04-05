// src/services/platform/adapters/TelegramAdapter.ts
/**
 * Telegram WebApp Adapter для интеграции с Telegram Mini Apps
 * Документация: https://core.telegram.org/bots/webapps
 */

import { createLogger } from '@/lib/logger';
import type {
  PlatformInfo,
  PlatformUser,
  PlatformTheme,
  ShareOptions,
  HapticFeedbackType,
} from '@/types/platform';

const logger = createLogger('TelegramAdapter');

// Типизация Telegram WebApp API
interface TelegramWebApp {
  initData: string;
  initDataUnsafe: TelegramWebAppInitData;
  version: string;
  platform: string;
  colorScheme: 'light' | 'dark';
  themeParams: TelegramThemeParams;
  isExpanded: boolean;
  viewportHeight: number;
  viewportStableHeight: number;
  headerColor: string;
  backgroundColor: string;
  isClosingConfirmationEnabled: boolean;
  BackButton: TelegramBackButton;
  MainButton: TelegramMainButton;
  HapticFeedback: TelegramHapticFeedback;

  ready(): void;
  expand(): void;
  close(): void;
  openLink(url: string, options?: { try_instant_view?: boolean }): void;
  openTelegramLink(url: string): void;
  showPopup(params: TelegramPopupParams, callback?: (button_id: string) => void): void;
  showAlert(message: string, callback?: () => void): void;
  showConfirm(message: string, callback?: (confirmed: boolean) => void): void;
  sendData(data: string): void;
  switchInlineQuery(query: string, choose_chat_types?: string[]): void;
  onEvent(eventType: string, callback: () => void): void;
  offEvent(eventType: string, callback: () => void): void;
}

interface TelegramWebAppInitData {
  query_id?: string;
  user?: TelegramUser;
  receiver?: TelegramUser;
  chat?: TelegramChat;
  chat_type?: string;
  chat_instance?: string;
  start_param?: string;
  can_send_after?: number;
  auth_date: number;
  hash: string;
}

interface TelegramUser {
  id: number;
  is_bot?: boolean;
  first_name: string;
  last_name?: string;
  username?: string;
  language_code?: string;
  is_premium?: boolean;
  photo_url?: string;
}

interface TelegramChat {
  id: number;
  type: string;
  title: string;
  username?: string;
  photo_url?: string;
}

interface TelegramThemeParams {
  bg_color?: string;
  text_color?: string;
  hint_color?: string;
  link_color?: string;
  button_color?: string;
  button_text_color?: string;
  secondary_bg_color?: string;
}

interface TelegramBackButton {
  isVisible: boolean;
  show(): void;
  hide(): void;
  onClick(callback: () => void): void;
  offClick(callback: () => void): void;
}

interface TelegramMainButton {
  text: string;
  color: string;
  textColor: string;
  isVisible: boolean;
  isActive: boolean;
  isProgressVisible: boolean;
  setText(text: string): void;
  onClick(callback: () => void): void;
  offClick(callback: () => void): void;
  show(): void;
  hide(): void;
  enable(): void;
  disable(): void;
  showProgress(leaveActive?: boolean): void;
  hideProgress(): void;
  setParams(params: { text?: string; color?: string; text_color?: string; is_active?: boolean; is_visible?: boolean }): void;
}

interface TelegramHapticFeedback {
  impactOccurred(style: 'light' | 'medium' | 'heavy' | 'rigid' | 'soft'): void;
  notificationOccurred(type: 'error' | 'success' | 'warning'): void;
  selectionChanged(): void;
}

interface TelegramPopupParams {
  title?: string;
  message: string;
  buttons?: Array<{ id?: string; type?: string; text?: string }>;
}

declare global {
  interface Window {
    Telegram?: {
      WebApp: TelegramWebApp;
    };
  }
}

export class TelegramAdapter {
  private webApp: TelegramWebApp | null = null;
  private isInitialized = false;

  async init(): Promise<void> {
    if (this.isInitialized) return;

    try {
      // Проверяем доступность Telegram WebApp
      if (!window.Telegram?.WebApp) {
        throw new Error('Telegram WebApp is not available');
      }

      this.webApp = window.Telegram.WebApp;

      // Сигнализируем о готовности приложения
      this.webApp.ready();

      // Разворачиваем приложение на весь экран
      if (!this.webApp.isExpanded) {
        this.webApp.expand();
      }

      this.isInitialized = true;

      logger.info('Initialized successfully', {
        version: this.webApp.version,
        platform: this.webApp.platform,
        colorScheme: this.webApp.colorScheme,
      });
    } catch (error) {
      logger.error('Failed to initialize', { error });
      throw error;
    }
  }

  getInfo(): PlatformInfo {
    return {
      type: 'telegram',
      version: this.webApp?.version || '1.0',
      isSupported: true,
      features: {
        nativeNavigation: true, // Back Button
        nativeShare: true, // switchInlineQuery
        nativeCamera: false, // Telegram не предоставляет прямой доступ к камере
        nativeFilePicker: false,
        hapticFeedback: true,
        nativeAnalytics: false,
        nativePayment: true, // Telegram Payments
        localStorage: true,
        cloudStorage: true, // Cloud Storage API
        pushNotifications: true,
        nativeAuth: true,
        biometricAuth: false,
      },
    };
  }

  async getUser(): Promise<PlatformUser> {
    if (!this.webApp) {
      throw new Error('Telegram WebApp not initialized');
    }

    try {
      const user = this.webApp.initDataUnsafe.user;

      if (!user) {
        throw new Error('User data not available');
      }

      return {
        id: user.id.toString(),
        firstName: user.first_name,
        lastName: user.last_name,
        photoUrl: user.photo_url,
        languageCode: user.language_code,
        isPremium: user.is_premium,
      };
    } catch (error) {
      logger.error('Failed to get user', { error });
      throw error;
    }
  }

  getTheme(): PlatformTheme {
    if (!this.webApp) {
      return this.getDefaultTheme();
    }

    try {
      const theme = this.webApp.themeParams;
      const isDark = this.webApp.colorScheme === 'dark';

      return {
        colorScheme: this.webApp.colorScheme,
        backgroundColor: theme.bg_color || (isDark ? '#000' : '#fff'),
        textColor: theme.text_color || (isDark ? '#fff' : '#000'),
        buttonColor: theme.button_color || '#0077ff',
        buttonTextColor: theme.button_text_color || '#fff',
        hintColor: theme.hint_color || '#818c99',
        linkColor: theme.link_color || '#0077ff',
        secondaryBackgroundColor: theme.secondary_bg_color || (isDark ? '#1a1a1a' : '#f5f5f5'),
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
    if (!this.webApp) {
      throw new Error('Telegram WebApp not initialized');
    }

    try {
      const text = options.text || options.title || '';
      const url = options.url || window.location.href;
      const shareText = text ? `${text} ${url}` : url;

      // Используем switchInlineQuery для шаринга
      // Пользователь сможет выбрать чат для отправки
      this.webApp.switchInlineQuery(shareText, ['users', 'groups', 'channels']);
      logger.debug('Share initiated', { text, url });
    } catch (error) {
      logger.error('Failed to share', { error });
      throw error;
    }
  }

  hapticFeedback(feedback: HapticFeedbackType): void {
    if (!this.webApp) return;

    try {
      const haptic = this.webApp.HapticFeedback;
      const type = feedback.type;
      const style = feedback.style || 'medium';

      switch (type) {
        case 'impact':
          haptic.impactOccurred(style as 'light' | 'medium' | 'heavy' | 'rigid' | 'soft');
          break;
        case 'notification':
          // Маппим style на notification type
          const notificationType = style === 'light' ? 'warning' : style === 'heavy' ? 'error' : 'success';
          haptic.notificationOccurred(notificationType);
          break;
        case 'selection':
          haptic.selectionChanged();
          break;
      }

      logger.debug('Haptic feedback triggered', { type, style });
    } catch (error) {
      logger.error('Haptic feedback failed', { error });
    }
  }

  openLink(url: string): void {
    if (!this.webApp) {
      window.open(url, '_blank');
      return;
    }

    try {
      // Telegram ссылки открываем через openTelegramLink
      if (url.startsWith('https://t.me/')) {
        this.webApp.openTelegramLink(url);
      } else {
        // Остальные ссылки через openLink
        this.webApp.openLink(url);
      }
      logger.debug('Link opened', { url });
    } catch (error) {
      logger.error('Failed to open link', { url, error });
      window.open(url, '_blank');
    }
  }

  close(): void {
    if (!this.webApp) return;

    try {
      this.webApp.close();
      logger.debug('App closed');
    } catch (error) {
      logger.error('Failed to close', { error });
    }
  }

  ready(): void {
    if (!this.webApp) return;

    try {
      this.webApp.ready();
      logger.debug('Ready signal sent');
    } catch (error) {
      logger.error('Failed to send ready', { error });
    }
  }

  expand(): void {
    if (!this.webApp) return;

    try {
      this.webApp.expand();
      logger.debug('App expanded');
    } catch (error) {
      logger.error('Failed to expand', { error });
    }
  }
}
