# Phase 12: Адаптация под платформы (Platform Integration)

**Длительность:** 4-5 дней
**Приоритет:** Критический
**Зависимости:** Все предыдущие фазы

---

## Цель

Создать абстракцию Platform Bridge для работы с VK Max SDK и Telegram WebApp API, адаптировать UI/UX под каждую платформу, настроить feature flags и провести тестирование на всех платформах.

---

## Архитектура Platform Bridge

### Структура компонентов

```
Platform Bridge
├── Platform Detector
├── VK Max Adapter
│   ├── VK Bridge API
│   ├── VK UI Components
│   └── VK Analytics
├── Telegram Adapter
│   ├── Telegram WebApp API
│   ├── Telegram UI Components
│   └── Telegram Analytics
├── Web Adapter (fallback)
└── Feature Flags Manager
```

---

## Часть 1: Типы и Интерфейсы

### 1.1. TypeScript Types

**Файл:** `src/types/platform.ts`

```typescript
export type PlatformType = 'vk' | 'max' | 'telegram' | 'web';

export interface PlatformInfo {
  type: PlatformType;
  version: string;
  isSupported: boolean;
  features: PlatformFeatures;
}

export interface PlatformFeatures {
  // UI Features
  nativeNavigation: boolean;
  nativeShare: boolean;
  nativeCamera: boolean;
  nativeFilePicker: boolean;
  hapticFeedback: boolean;

  // Analytics Features
  nativeAnalytics: boolean;

  // Payment Features
  nativePayment: boolean;

  // Storage Features
  localStorage: boolean;
  cloudStorage: boolean;

  // Notification Features
  pushNotifications: boolean;

  // Auth Features
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
```

---

## Часть 2: Platform Bridge Core

### 2.1. Platform Bridge Service

**Файл:** `src/services/platform/PlatformBridge.ts`

```typescript
import { VKAdapter } from './adapters/VKAdapter';
import { TelegramAdapter } from './adapters/TelegramAdapter';
import { WebAdapter } from './adapters/WebAdapter';
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

  /**
   * Определить платформу
   */
  private detectPlatform(): PlatformType {
    // Check for VK/Max
    if (
      typeof window !== 'undefined' &&
      (window as any).vkBridge
    ) {
      return 'max';
    }

    // Check for Telegram
    if (
      typeof window !== 'undefined' &&
      (window as any).Telegram?.WebApp
    ) {
      return 'telegram';
    }

    // Fallback to web
    return 'web';
  }

  /**
   * Создать адаптер для платформы
   */
  private createAdapter(type: PlatformType): IPlatformAdapter {
    switch (type) {
      case 'max':
      case 'vk':
        return new VKAdapter();
      case 'telegram':
        return new TelegramAdapter();
      case 'web':
      default:
        return new WebAdapter();
    }
  }

  /**
   * Инициализация
   */
  async init(): Promise<void> {
    await this.adapter.init();
  }

  /**
   * Получить информацию о платформе
   */
  getInfo(): PlatformInfo {
    return this.adapter.getInfo();
  }

  /**
   * Получить пользователя
   */
  async getUser(): Promise<PlatformUser> {
    return this.adapter.getUser();
  }

  /**
   * Получить тему
   */
  getTheme(): PlatformTheme {
    return this.adapter.getTheme();
  }

  /**
   * Поделиться
   */
  async share(options: ShareOptions): Promise<void> {
    return this.adapter.share(options);
  }

  /**
   * Haptic feedback
   */
  hapticFeedback(feedback: HapticFeedbackType): void {
    this.adapter.hapticFeedback(feedback);
  }

  /**
   * Открыть ссылку
   */
  openLink(url: string): void {
    this.adapter.openLink(url);
  }

  /**
   * Закрыть приложение
   */
  close(): void {
    this.adapter.close();
  }

  /**
   * Сообщить о готовности
   */
  ready(): void {
    this.adapter.ready();
  }

  /**
   * Развернуть на весь экран
   */
  expand(): void {
    this.adapter.expand();
  }

  /**
   * Получить тип платформы
   */
  getPlatformType(): PlatformType {
    return this.platformType;
  }
}
```

---

## Часть 3: Platform Adapters

### 3.1. VK Adapter

**Файл:** `src/services/platform/adapters/VKAdapter.ts`

```typescript
import vkBridge, { VKBridgeSubscribeHandler } from '@vkontakte/vk-bridge';
import type {
  IPlatformAdapter,
  PlatformInfo,
  PlatformUser,
  PlatformTheme,
  ShareOptions,
  HapticFeedbackType,
} from '../PlatformBridge';

export class VKAdapter implements IPlatformAdapter {
  private isInitialized = false;

  async init(): Promise<void> {
    if (this.isInitialized) return;

    try {
      await vkBridge.send('VKWebAppInit');
      this.isInitialized = true;

      // Subscribe to theme changes
      vkBridge.subscribe(this.handleBridgeEvent);
    } catch (error) {
      console.error('[VKAdapter] Failed to initialize:', error);
    }
  }

  private handleBridgeEvent: VKBridgeSubscribeHandler = (event) => {
    if (event.detail.type === 'VKWebAppUpdateConfig') {
      // Handle theme changes
      const config = event.detail.data;
      document.documentElement.setAttribute(
        'data-theme',
        config.scheme || 'bright_light'
      );
    }
  };

  getInfo(): PlatformInfo {
    return {
      type: 'max',
      version: '1.0.0',
      isSupported: true,
      features: {
        nativeNavigation: true,
        nativeShare: true,
        nativeCamera: false,
        nativeFilePicker: true,
        hapticFeedback: true,
        nativeAnalytics: true,
        nativePayment: true,
        localStorage: true,
        cloudStorage: true,
        pushNotifications: true,
        nativeAuth: true,
        biometricAuth: false,
      },
    };
  }

  async getUser(): Promise<PlatformUser> {
    try {
      const user = await vkBridge.send('VKWebAppGetUserInfo');
      return {
        id: String(user.id),
        firstName: user.first_name,
        lastName: user.last_name,
        photoUrl: user.photo_200,
      };
    } catch (error) {
      console.error('[VKAdapter] Failed to get user:', error);
      throw error;
    }
  }

  getTheme(): PlatformTheme {
    // VK uses CSS variables
    const style = getComputedStyle(document.documentElement);

    return {
      colorScheme: style.getPropertyValue('--vkui--color_scheme') as any || 'light',
      backgroundColor: style.getPropertyValue('--vkui--color_background') || '#fff',
      textColor: style.getPropertyValue('--vkui--color_text_primary') || '#000',
      buttonColor: style.getPropertyValue('--vkui--color_accent') || '#0077ff',
      buttonTextColor: '#fff',
      hintColor: style.getPropertyValue('--vkui--color_text_secondary') || '#818c99',
      linkColor: style.getPropertyValue('--vkui--color_accent') || '#0077ff',
      secondaryBackgroundColor:
        style.getPropertyValue('--vkui--color_background_secondary') || '#f5f5f5',
    };
  }

  async share(options: ShareOptions): Promise<void> {
    try {
      await vkBridge.send('VKWebAppShare', {
        link: options.url || window.location.href,
      });
    } catch (error) {
      console.error('[VKAdapter] Failed to share:', error);
      throw error;
    }
  }

  hapticFeedback(feedback: HapticFeedbackType): void {
    try {
      const style = feedback.style || 'medium';
      vkBridge.send('VKWebAppTapticNotificationOccurred', {
        type: feedback.type === 'notification' ? 'success' : style,
      });
    } catch (error) {
      console.error('[VKAdapter] Haptic feedback failed:', error);
    }
  }

  openLink(url: string): void {
    vkBridge.send('VKWebAppOpenApp', { app_id: 0, location: url });
  }

  close(): void {
    vkBridge.send('VKWebAppClose', { status: 'success' });
  }

  ready(): void {
    // VK doesn't have explicit ready method
  }

  expand(): void {
    // VK apps are always fullscreen
  }
}
```

---

### 3.2. Telegram Adapter

**Файл:** `src/services/platform/adapters/TelegramAdapter.ts`

```typescript
import type {
  IPlatformAdapter,
  PlatformInfo,
  PlatformUser,
  PlatformTheme,
  ShareOptions,
  HapticFeedbackType,
} from '../PlatformBridge';

declare global {
  interface Window {
    Telegram?: {
      WebApp: any;
    };
  }
}

export class TelegramAdapter implements IPlatformAdapter {
  private webApp: any;

  async init(): Promise<void> {
    if (typeof window === 'undefined' || !window.Telegram?.WebApp) {
      throw new Error('Telegram WebApp is not available');
    }

    this.webApp = window.Telegram.WebApp;
    this.webApp.ready();
    this.webApp.expand();
  }

  getInfo(): PlatformInfo {
    return {
      type: 'telegram',
      version: this.webApp?.version || '6.0',
      isSupported: true,
      features: {
        nativeNavigation: false,
        nativeShare: true,
        nativeCamera: false,
        nativeFilePicker: false,
        hapticFeedback: true,
        nativeAnalytics: false,
        nativePayment: true,
        localStorage: true,
        cloudStorage: true,
        pushNotifications: false,
        nativeAuth: true,
        biometricAuth: true,
      },
    };
  }

  async getUser(): Promise<PlatformUser> {
    const user = this.webApp.initDataUnsafe?.user;

    if (!user) {
      throw new Error('User data not available');
    }

    return {
      id: String(user.id),
      firstName: user.first_name,
      lastName: user.last_name,
      photoUrl: user.photo_url,
      languageCode: user.language_code,
      isPremium: user.is_premium,
    };
  }

  getTheme(): PlatformTheme {
    const themeParams = this.webApp.themeParams;

    return {
      colorScheme: this.webApp.colorScheme || 'light',
      backgroundColor: themeParams.bg_color || '#fff',
      textColor: themeParams.text_color || '#000',
      buttonColor: themeParams.button_color || '#3390ec',
      buttonTextColor: themeParams.button_text_color || '#fff',
      hintColor: themeParams.hint_color || '#999',
      linkColor: themeParams.link_color || '#3390ec',
      secondaryBackgroundColor: themeParams.secondary_bg_color || '#f5f5f5',
    };
  }

  async share(options: ShareOptions): Promise<void> {
    const url = options.url || window.location.href;
    const text = options.text || options.title || '';

    // Use Telegram share link
    const shareUrl = `https://t.me/share/url?url=${encodeURIComponent(
      url
    )}&text=${encodeURIComponent(text)}`;

    this.webApp.openTelegramLink(shareUrl);
  }

  hapticFeedback(feedback: HapticFeedbackType): void {
    try {
      if (feedback.type === 'impact') {
        this.webApp.HapticFeedback.impactOccurred(
          feedback.style || 'medium'
        );
      } else if (feedback.type === 'notification') {
        this.webApp.HapticFeedback.notificationOccurred('success');
      } else if (feedback.type === 'selection') {
        this.webApp.HapticFeedback.selectionChanged();
      }
    } catch (error) {
      console.error('[TelegramAdapter] Haptic feedback failed:', error);
    }
  }

  openLink(url: string): void {
    this.webApp.openLink(url);
  }

  close(): void {
    this.webApp.close();
  }

  ready(): void {
    this.webApp.ready();
  }

  expand(): void {
    this.webApp.expand();
  }
}
```

---

### 3.3. Web Adapter (Fallback)

**Файл:** `src/services/platform/adapters/WebAdapter.ts`

```typescript
import type {
  IPlatformAdapter,
  PlatformInfo,
  PlatformUser,
  PlatformTheme,
  ShareOptions,
  HapticFeedbackType,
} from '../PlatformBridge';

export class WebAdapter implements IPlatformAdapter {
  async init(): Promise<void> {
    // No initialization needed for web
  }

  getInfo(): PlatformInfo {
    return {
      type: 'web',
      version: '1.0.0',
      isSupported: true,
      features: {
        nativeNavigation: false,
        nativeShare: !!navigator.share,
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
    // For web, user should be obtained via API after login
    throw new Error('User data not available in web mode');
  }

  getTheme(): PlatformTheme {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

    return {
      colorScheme: isDark ? 'dark' : 'light',
      backgroundColor: isDark ? '#000' : '#fff',
      textColor: isDark ? '#fff' : '#000',
      buttonColor: '#007aff',
      buttonTextColor: '#fff',
      hintColor: '#999',
      linkColor: '#007aff',
      secondaryBackgroundColor: isDark ? '#1c1c1e' : '#f5f5f5',
    };
  }

  async share(options: ShareOptions): Promise<void> {
    if (navigator.share) {
      try {
        await navigator.share({
          title: options.title,
          text: options.text,
          url: options.url,
        });
      } catch (error) {
        // User cancelled or error occurred
        console.error('[WebAdapter] Share failed:', error);
      }
    } else {
      // Fallback: copy to clipboard
      const url = options.url || window.location.href;
      await navigator.clipboard.writeText(url);
      alert('Ссылка скопирована в буфер обмена');
    }
  }

  hapticFeedback(_feedback: HapticFeedbackType): void {
    // Not supported on web
  }

  openLink(url: string): void {
    window.open(url, '_blank');
  }

  close(): void {
    window.close();
  }

  ready(): void {
    // Not applicable for web
  }

  expand(): void {
    // Not applicable for web
  }
}
```

---

## Часть 4: React Integration

### 4.1. usePlatform Hook

**Файл:** `src/hooks/usePlatform.ts`

```typescript
import { useContext } from 'react';
import { PlatformContext } from '@/contexts/PlatformContext';

export function usePlatform() {
  const platform = useContext(PlatformContext);

  if (!platform) {
    throw new Error('usePlatform must be used within PlatformProvider');
  }

  return platform;
}
```

---

### 4.2. PlatformProvider

**Файл:** `src/contexts/PlatformContext.tsx`

```typescript
import { createContext, useEffect, useState, ReactNode } from 'react';
import { PlatformBridge } from '@/services/platform/PlatformBridge';
import type { PlatformInfo, PlatformTheme } from '@/types/platform';

interface PlatformContextValue {
  bridge: PlatformBridge;
  info: PlatformInfo | null;
  theme: PlatformTheme | null;
  isReady: boolean;
}

export const PlatformContext = createContext<PlatformContextValue | null>(null);

interface PlatformProviderProps {
  children: ReactNode;
}

export function PlatformProvider({ children }: PlatformProviderProps) {
  const [bridge] = useState(() => new PlatformBridge());
  const [info, setInfo] = useState<PlatformInfo | null>(null);
  const [theme, setTheme] = useState<PlatformTheme | null>(null);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    const initPlatform = async () => {
      try {
        await bridge.init();

        const platformInfo = bridge.getInfo();
        const platformTheme = bridge.getTheme();

        setInfo(platformInfo);
        setTheme(platformTheme);

        // Apply theme to document
        applyTheme(platformTheme);

        // Signal ready
        bridge.ready();
        setIsReady(true);
      } catch (error) {
        console.error('[PlatformProvider] Failed to initialize:', error);
      }
    };

    initPlatform();
  }, [bridge]);

  const applyTheme = (theme: PlatformTheme) => {
    document.documentElement.style.setProperty(
      '--color-bg',
      theme.backgroundColor
    );
    document.documentElement.style.setProperty(
      '--color-text',
      theme.textColor
    );
    document.documentElement.style.setProperty(
      '--color-button',
      theme.buttonColor
    );
    document.documentElement.style.setProperty(
      '--color-button-text',
      theme.buttonTextColor
    );
    document.documentElement.style.setProperty(
      '--color-hint',
      theme.hintColor
    );
    document.documentElement.style.setProperty(
      '--color-link',
      theme.linkColor
    );
    document.documentElement.style.setProperty(
      '--color-bg-secondary',
      theme.secondaryBackgroundColor
    );

    document.documentElement.setAttribute(
      'data-theme',
      theme.colorScheme
    );
  };

  return (
    <PlatformContext.Provider value={{ bridge, info, theme, isReady }}>
      {children}
    </PlatformContext.Provider>
  );
}
```

---

## Часть 5: Platform-Specific UI

### 5.1. Platform-Specific Button

**Файл:** `src/components/ui/PlatformButton.tsx`

```typescript
import { usePlatform } from '@/hooks/usePlatform';
import { Button, ButtonProps } from './Button';
import styles from './PlatformButton.module.css';

export function PlatformButton(props: ButtonProps) {
  const { info, bridge } = usePlatform();

  const handleClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    // Add haptic feedback for mobile platforms
    if (info?.features.hapticFeedback) {
      bridge.hapticFeedback({
        type: 'impact',
        style: 'medium',
      });
    }

    props.onClick?.(e);
  };

  return (
    <Button
      {...props}
      onClick={handleClick}
      className={`${props.className} ${styles[info?.type || 'web']}`}
    />
  );
}
```

---

### 5.2. Platform-Specific Share

**Файл:** `src/components/features/ShareButton.tsx`

```typescript
import { usePlatform } from '@/hooks/usePlatform';
import { Button } from '@/components/ui/Button';
import type { ShareOptions } from '@/types/platform';

interface ShareButtonProps {
  options: ShareOptions;
  children: React.ReactNode;
}

export function ShareButton({ options, children }: ShareButtonProps) {
  const { bridge, info } = usePlatform();

  const handleShare = async () => {
    try {
      await bridge.share(options);
    } catch (error) {
      console.error('[ShareButton] Failed to share:', error);
    }
  };

  if (!info?.features.nativeShare) {
    return null;
  }

  return (
    <Button variant="primary" onClick={handleShare}>
      {children}
    </Button>
  );
}
```

---

## Часть 6: Feature Flags

### 6.1. Feature Flags Manager

**Файл:** `src/services/featureFlags/FeatureFlagsManager.ts`

```typescript
import type { PlatformType } from '@/types/platform';

interface FeatureFlags {
  achievements: boolean;
  villain: boolean;
  referral: boolean;
  weeklyReport: boolean;
  paywall: boolean;
  analytics: boolean;
}

const PLATFORM_FEATURE_FLAGS: Record<PlatformType, Partial<FeatureFlags>> = {
  vk: {
    achievements: true,
    villain: true,
    referral: true,
    weeklyReport: true,
    paywall: true,
    analytics: true,
  },
  max: {
    achievements: true,
    villain: true,
    referral: true,
    weeklyReport: true,
    paywall: true,
    analytics: true,
  },
  telegram: {
    achievements: true,
    villain: true,
    referral: true,
    weeklyReport: true,
    paywall: true,
    analytics: false, // Use Telegram analytics instead
  },
  web: {
    achievements: true,
    villain: false, // Not available in web version
    referral: false,
    weeklyReport: true,
    paywall: true,
    analytics: true,
  },
};

export class FeatureFlagsManager {
  private platform: PlatformType;
  private flags: FeatureFlags;

  constructor(platform: PlatformType) {
    this.platform = platform;
    this.flags = {
      achievements: true,
      villain: true,
      referral: true,
      weeklyReport: true,
      paywall: true,
      analytics: true,
      ...PLATFORM_FEATURE_FLAGS[platform],
    };
  }

  isEnabled(feature: keyof FeatureFlags): boolean {
    return this.flags[feature] ?? false;
  }

  getAll(): FeatureFlags {
    return { ...this.flags };
  }
}
```

---

### 6.2. useFeatureFlag Hook

**Файл:** `src/hooks/useFeatureFlag.ts`

```typescript
import { usePlatform } from './usePlatform';
import { FeatureFlagsManager } from '@/services/featureFlags/FeatureFlagsManager';
import { useMemo } from 'react';

export function useFeatureFlag(feature: string): boolean {
  const { info } = usePlatform();

  const manager = useMemo(() => {
    if (!info) return null;
    return new FeatureFlagsManager(info.type);
  }, [info]);

  if (!manager) return false;

  return manager.isEnabled(feature as any);
}
```

---

## Часть 7: Testing на платформах

### 7.1. VK Max Testing Checklist

```markdown
## VK Max Testing Checklist

### Инициализация
- [ ] VK Bridge успешно инициализируется
- [ ] Получены данные пользователя
- [ ] Применена тема VK

### Функциональность
- [ ] Share работает через VK Bridge
- [ ] Haptic feedback работает
- [ ] Аналитика отправляется в VK
- [ ] Платежи через VK работают

### UI/UX
- [ ] Цвета соответствуют теме VK
- [ ] Переходы плавные
- [ ] Нет конфликтов с нативной навигацией

### Performance
- [ ] Быстрая загрузка
- [ ] Нет memory leaks
- [ ] 60 FPS анимации
```

---

### 7.2. Telegram Testing Checklist

```markdown
## Telegram Testing Checklist

### Инициализация
- [ ] Telegram WebApp инициализируется
- [ ] Получены initData
- [ ] Применена тема Telegram

### Функциональность
- [ ] Share через Telegram работает
- [ ] Haptic feedback работает
- [ ] Main button работает
- [ ] Back button работает

### UI/UX
- [ ] Цвета соответствуют теме Telegram
- [ ] Expand на весь экран работает
- [ ] Safe area учитывается

### Security
- [ ] initData валидируется на backend
- [ ] User hash проверяется
```

---

## Чеклист задач

### Platform Bridge
- [ ] Создать PlatformBridge core
- [ ] Реализовать VKAdapter
- [ ] Реализовать TelegramAdapter
- [ ] Реализовать WebAdapter (fallback)
- [ ] Добавить platform detection

### React Integration
- [ ] Создать PlatformProvider
- [ ] Создать usePlatform hook
- [ ] Интегрировать в App
- [ ] Применить тему платформы

### Feature Flags
- [ ] Создать FeatureFlagsManager
- [ ] Настроить флаги для каждой платформы
- [ ] Создать useFeatureFlag hook
- [ ] Интегрировать в компоненты

### Platform-Specific UI
- [ ] Адаптировать кнопки
- [ ] Адаптировать навигацию
- [ ] Адаптировать модалы
- [ ] Адаптировать формы

### Testing
- [ ] Протестировать на VK Max
- [ ] Протестировать на Telegram
- [ ] Протестировать на Web
- [ ] Проверить все feature flags
- [ ] Проверить аналитику на всех платформах

### Documentation
- [ ] Документировать API каждого адаптера
- [ ] Создать гайд по тестированию
- [ ] Создать чеклисты для QA

---

## Заключение

После завершения всех 13 фаз у вас будет полностью функциональное кросс-платформенное приложение "Объяснятель ДЗ", готовое к деплою на VK Max, Telegram и Web.

Все roadmap файлы содержат детальный код, типы, аналитику и чеклисты для успешной реализации проекта.
