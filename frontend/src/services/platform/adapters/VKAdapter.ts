// src/services/platform/adapters/VKAdapter.ts
/**
 * VK Bridge Adapter для интеграции с VK Mini Apps
 * Документация: https://dev.vk.com/bridge/getting-started
 */

import bridge from '@/lib/platform/bridge';
import type {
  VKBridgeEvent,
  UserInfo,
} from '@vkontakte/vk-bridge';

import { createLogger } from '@/lib/logger';
import type {
  PlatformInfo,
  PlatformUser,
  PlatformTheme,
  ShareOptions,
  HapticFeedbackType,
} from '@/types/platform';

const logger = createLogger('VKAdapter');

export class VKAdapter {
  private isInitialized = false;
  private userInfo: UserInfo | null = null;
  private theme: PlatformTheme | null = null;

  async init(): Promise<void> {
    if (this.isInitialized) return;

    try {
      // Инициализируем VK Bridge
      await bridge.send('VKWebAppInit');
      this.isInitialized = true;

      // Подписываемся на изменения темы
      bridge.subscribe((event: VKBridgeEvent<any>) => {
        if (event.detail.type === 'VKWebAppUpdateConfig') {
          this.theme = this.parseVKTheme(event.detail.data);
          logger.debug('Theme updated', { theme: this.theme });
        }
      });

      // Получаем информацию о пользователе
      try {
        this.userInfo = await bridge.send('VKWebAppGetUserInfo');
        logger.info('User info retrieved', {
          userId: this.userInfo.id,
        });
      } catch (error) {
        logger.warn('Failed to get user info', { error });
      }

      logger.info('Initialized successfully');
    } catch (error) {
      logger.error('Failed to initialize', { error });
      throw error;
    }
  }

  getInfo(): PlatformInfo {
    return {
      type: 'vk',
      version: '1.0.0',
      isSupported: true,
      features: {
        nativeNavigation: true,
        nativeShare: true,
        nativeCamera: true, // VK поддерживает камеру через VKWebAppOpenCamera
        nativeFilePicker: true,
        hapticFeedback: true,
        nativeAnalytics: true,
        nativePayment: true, // VK Pay
        localStorage: true,
        cloudStorage: true, // VK Storage
        pushNotifications: true,
        nativeAuth: true,
        biometricAuth: false,
      },
    };
  }

  async getUser(): Promise<PlatformUser> {
    if (!this.isInitialized) {
      throw new Error('VK Bridge not initialized');
    }

    try {
      // Если уже получили userInfo при инициализации
      if (this.userInfo) {
        return this.convertVKUser(this.userInfo);
      }

      // Иначе запрашиваем заново
      const user = await bridge.send('VKWebAppGetUserInfo');
      this.userInfo = user;
      return this.convertVKUser(user);
    } catch (error) {
      logger.error('Failed to get user', { error });
      throw error;
    }
  }

  private convertVKUser(vkUser: UserInfo): PlatformUser {
    return {
      id: vkUser.id.toString(),
      firstName: vkUser.first_name,
      lastName: vkUser.last_name,
      photoUrl: vkUser.photo_200 || vkUser.photo_100,
      languageCode: undefined, // VK UserInfo не содержит language_code
      isPremium: false, // VK не предоставляет информацию о премиуме в UserInfo
    };
  }

  getTheme(): PlatformTheme {
    if (this.theme) {
      return this.theme;
    }

    // Fallback тема
    return this.getDefaultTheme();
  }

  private parseVKTheme(config: any): PlatformTheme {
    const scheme = config.scheme || 'bright_light';
    const isDark = scheme.includes('space_gray') || scheme === 'vkcom_dark';

    return {
      colorScheme: isDark ? 'dark' : 'light',
      backgroundColor: config.background_color || (isDark ? '#000' : '#fff'),
      textColor: config.text_color || (isDark ? '#fff' : '#000'),
      buttonColor: config.button_color || '#0077ff',
      buttonTextColor: config.button_text_color || '#fff',
      hintColor: config.hint_color || '#818c99',
      linkColor: config.link_color || '#0077ff',
      secondaryBackgroundColor: config.header_background_color || (isDark ? '#1a1a1a' : '#f5f5f5'),
    };
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
    if (!this.isInitialized) {
      throw new Error('VK Bridge not initialized');
    }

    try {
      await bridge.send('VKWebAppShare', {
        link: options.url || window.location.href,
      });
    } catch (error) {
      console.error('[VKAdapter] Failed to share:', error);
      throw error;
    }
  }

  hapticFeedback(feedback: HapticFeedbackType): void {
    if (!this.isInitialized) return;

    try {
      // VK Bridge поддерживает impact, notification, selection
      const type = feedback.type;
      const style = feedback.style || 'medium';

      bridge.send('VKWebAppTapticNotificationOccurred', {
        type: this.mapHapticType(type, style),
      }).catch((error) => {
        console.error('[VKAdapter] Haptic feedback failed:', error);
      });
    } catch (error) {
      console.error('[VKAdapter] Haptic feedback failed:', error);
    }
  }

  private mapHapticType(type: string, style: string): 'error' | 'success' | 'warning' {
    // Маппим наши типы на VK типы
    if (type === 'notification') {
      if (style === 'light') return 'warning';
      if (style === 'heavy') return 'error';
      return 'success';
    }

    // Для impact и selection используем success как дефолт
    return 'success';
  }

  openLink(url: string): void {
    // VK Bridge не имеет прямого метода для открытия ссылок
    // Используем стандартный window.open
    window.open(url, '_blank');
  }

  close(): void {
    if (!this.isInitialized) return;

    try {
      // VKWebAppClose требует специальных параметров
      // Просто закрываем окно стандартным способом
      window.close();
    } catch (error) {
      console.error('[VKAdapter] Failed to close:', error);
    }
  }

  ready(): void {
    // VK Bridge не требует явного ready сигнала
    // Приложение становится ready после VKWebAppInit
  }

  expand(): void {
    // VK Mini Apps не требуют expand - они уже fullscreen
    // Метод оставляем пустым для совместимости с интерфейсом
  }
}
