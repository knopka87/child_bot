// src/services/platform/PlatformBridge.ts
import { WebAdapter } from './adapters/WebAdapter';
import { MaxAdapter } from './adapters/MaxAdapter';
import { VKAdapter } from './adapters/VKAdapter';
import { TelegramAdapter } from './adapters/TelegramAdapter';
import type {
  PlatformType,
  PlatformInfo,
  PlatformUser,
  PlatformTheme,
  ShareOptions,
  HapticFeedbackType,
} from '@/types/platform';

export interface IPlatformAdapter {
  init(): Promise<void>;
  getInfo(): PlatformInfo;
  getUser(): Promise<PlatformUser>;
  getTheme(): PlatformTheme;
  share(options: ShareOptions): Promise<void>;
  hapticFeedback(feedback: HapticFeedbackType): void;
  openLink(url: string): void;
  close(): void;
  ready(): void;
  expand(): void;
}

export class PlatformBridge {
  private adapter: IPlatformAdapter;
  private platformType: PlatformType;

  constructor() {
    this.platformType = this.detectPlatform();
    this.adapter = this.createAdapter(this.platformType);
  }

  private detectPlatform(): PlatformType {
    if (typeof window === 'undefined') return 'web';

    // Check for VK Bridge - проверяем несколько признаков
    const urlParams = new URLSearchParams(window.location.search);
    const vkPlatform = urlParams.get('vk_platform');
    const vkUserId = urlParams.get('vk_user_id');
    const vkAppId = urlParams.get('vk_app_id');
    const hasVKBridge = typeof (window as any).vkBridge !== 'undefined' || typeof (window as any).VK !== 'undefined';

    // Если есть ЛЮБОЙ признак VK - это VK платформа
    if (vkPlatform || vkUserId || vkAppId || hasVKBridge) {
      console.log('[PlatformBridge] Detected VK platform:', {
        vkPlatform,
        vkUserId: !!vkUserId,
        vkAppId: !!vkAppId,
        hasVKBridge,
      });
      return 'vk';
    }

    // КРИТИЧЕСКИ ВАЖНО: Если VK параметров нет в URL (например, после redirect),
    // но platform_id уже сохранён в localStorage - используем его
    // Это предотвращает потерю платформы при навигации внутри приложения
    const savedPlatform = localStorage.getItem('platform_id') as PlatformType;
    if (savedPlatform && (savedPlatform === 'vk' || savedPlatform === 'max' || savedPlatform === 'telegram')) {
      console.log('[PlatformBridge] Using saved platform from storage:', savedPlatform);
      return savedPlatform;
    }

    // Check for MAX Bridge
    if ((window as any).MaxBridge || document.referrer.includes('max.ru')) {
      return 'max';
    }

    // Check for Telegram WebApp
    if ((window as any).Telegram?.WebApp) {
      return 'telegram';
    }

    console.log('[PlatformBridge] Detected web platform (fallback)');
    return 'web';
  }

  private createAdapter(type: PlatformType): IPlatformAdapter {
    switch (type) {
      case 'vk':
        return new VKAdapter();
      case 'max':
        return new MaxAdapter();
      case 'telegram':
        return new TelegramAdapter();
      case 'web':
      default:
        return new WebAdapter();
    }
  }

  async init(): Promise<void> {
    await this.adapter.init();
  }

  getInfo(): PlatformInfo {
    return this.adapter.getInfo();
  }

  async getUser(): Promise<PlatformUser> {
    return this.adapter.getUser();
  }

  getTheme(): PlatformTheme {
    return this.adapter.getTheme();
  }

  async share(options: ShareOptions): Promise<void> {
    return this.adapter.share(options);
  }

  hapticFeedback(feedback: HapticFeedbackType): void {
    this.adapter.hapticFeedback(feedback);
  }

  openLink(url: string): void {
    this.adapter.openLink(url);
  }

  close(): void {
    this.adapter.close();
  }

  ready(): void {
    this.adapter.ready();
  }

  expand(): void {
    this.adapter.expand();
  }

  getPlatformType(): PlatformType {
    return this.platformType;
  }
}
