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

    // Check for VK Bridge
    const urlParams = new URLSearchParams(window.location.search);
    const vkPlatform = urlParams.get('vk_platform');
    if (vkPlatform || (window as any).vkBridge) {
      return 'vk';
    }

    // Check for MAX Bridge
    if ((window as any).MaxBridge || document.referrer.includes('max.ru')) {
      return 'max';
    }

    // Check for Telegram WebApp
    if ((window as any).Telegram?.WebApp) {
      return 'telegram';
    }

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
