# Phase 0: Настройка проекта

**Длительность:** 1-2 дня
**Приоритет:** Критический
**Статус:** To Do

---

## Цель

Инициализировать проект React + TypeScript + Vite с полной настройкой инфраструктуры для разработки VK MiniApp.

---

## Шаг 1: Инициализация Vite проекта

### 1.1. Создание проекта

```bash
cd /Users/a.yanover/Xsolla/child_bot/frontend
npm create vite@latest . -- --template react-ts
```

### 1.2. Установка зависимостей

```bash
# Core dependencies
npm install react@^18.2.0 react-dom@^18.2.0
npm install react-router-dom@^6.22.0
npm install zustand@^4.5.0
npm install axios@^1.6.7

# VK Bridge и VKUI (обязательно для VK miniapps!)
npm install @vkontakte/vk-bridge@^2.12.0
npm install @vkontakte/vkui@^6.0.0
npm install @vkontakte/icons@^2.0.0

# Forms and validation
npm install react-hook-form@^7.51.0
npm install zod@^3.22.4

# Image processing
npm install browser-image-compression@^2.0.2
npm install react-image-crop@^11.0.5

# Utilities
npm install clsx@^2.1.0
npm install date-fns@^3.3.1

# Dev dependencies
npm install -D @types/react@^18.2.0
npm install -D @types/react-dom@^18.2.0
npm install -D @typescript-eslint/eslint-plugin@^7.1.0
npm install -D @typescript-eslint/parser@^7.1.0
npm install -D eslint@^8.57.0
npm install -D eslint-plugin-react-hooks@^4.6.0
npm install -D eslint-plugin-react-refresh@^0.4.5
npm install -D prettier@^3.2.5
npm install -D vitest@^1.3.1
npm install -D @testing-library/react@^14.2.1
npm install -D @testing-library/jest-dom@^6.4.2
npm install -D rollup-plugin-visualizer@^5.12.0
```

---

## Шаг 2: Структура папок проекта

### 2.1. Создание структуры

```bash
mkdir -p src/{api,assets,components,config,hooks,lib,pages,stores,types,utils}
mkdir -p src/components/{ui,layout,features}
mkdir -p src/assets/{images,icons}
mkdir -p public
```

### 2.2. Финальная структура

```
frontend/
├── public/
│   ├── vk-app-icon.png           # Иконка приложения VK
│   └── manifest.json              # Манифест VK MiniApp
├── src/
│   ├── api/                       # API клиенты
│   │   ├── client.ts              # Базовый HTTP клиент
│   │   ├── endpoints.ts           # API эндпоинты
│   │   ├── profile.ts             # Профиль API
│   │   ├── attempts.ts            # Попытки API
│   │   ├── achievements.ts        # Достижения API
│   │   └── analytics.ts           # Аналитика API
│   ├── assets/                    # Статические файлы
│   │   ├── images/                # Изображения (маскоты, злодеи)
│   │   └── icons/                 # SVG иконки
│   ├── components/                # React компоненты
│   │   ├── ui/                    # UI Kit компоненты
│   │   │   ├── Button/
│   │   │   ├── Input/
│   │   │   ├── Card/
│   │   │   ├── Modal/
│   │   │   ├── Spinner/
│   │   │   └── ProgressBar/
│   │   ├── layout/                # Layout компоненты
│   │   │   ├── Header/
│   │   │   ├── BottomNav/
│   │   │   └── Container/
│   │   └── features/              # Feature-специфичные компоненты
│   │       ├── MascotCard/
│   │       ├── VillainCard/
│   │       ├── HintCard/
│   │       └── AchievementCard/
│   ├── config/                    # Конфигурация
│   │   ├── constants.ts           # Константы приложения
│   │   ├── routes.ts              # Роутинг
│   │   └── theme.ts               # Тема и стили
│   ├── hooks/                     # Custom hooks
│   │   ├── usePlatform.ts         # Платформенная абстракция
│   │   ├── useAnalytics.ts        # Аналитика
│   │   ├── useImageUpload.ts      # Загрузка изображений
│   │   └── useProfile.ts          # Профиль пользователя
│   ├── lib/                       # Библиотеки и утилиты
│   │   ├── platform/              # Платформенная абстракция
│   │   │   ├── bridge.ts          # Абстракция VK Bridge
│   │   │   ├── vk.ts              # VK специфичный код
│   │   │   └── types.ts           # Типы платформы
│   │   └── analytics/             # Аналитика
│   │       ├── tracker.ts         # Трекер событий
│   │       └── events.ts          # Определения событий
│   ├── pages/                     # Страницы приложения
│   │   ├── Home/
│   │   ├── Help/
│   │   ├── Check/
│   │   ├── Achievements/
│   │   ├── Friends/
│   │   ├── Profile/
│   │   └── Onboarding/
│   ├── stores/                    # Zustand stores
│   │   ├── profileStore.ts        # Профиль пользователя
│   │   ├── attemptStore.ts        # Текущая попытка
│   │   └── platformStore.ts       # Платформенные данные
│   ├── types/                     # TypeScript типы
│   │   ├── api.ts                 # API типы
│   │   ├── domain.ts              # Доменные модели
│   │   └── analytics.ts           # Аналитические события
│   ├── utils/                     # Утилиты
│   │   ├── format.ts              # Форматирование данных
│   │   ├── validation.ts          # Валидация
│   │   └── image.ts               # Обработка изображений
│   ├── App.tsx                    # Главный компонент
│   ├── main.tsx                   # Entry point
│   └── index.css                  # Глобальные стили
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
├── .eslintrc.json
├── .prettierrc
└── .env.example
```

---

## Шаг 3: Конфигурация TypeScript

### 3.1. `tsconfig.json`

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "useDefineForClassFields": true,
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,

    /* Bundler mode */
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",

    /* Linting */
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,

    /* Path mapping */
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@/api/*": ["./src/api/*"],
      "@/components/*": ["./src/components/*"],
      "@/hooks/*": ["./src/hooks/*"],
      "@/lib/*": ["./src/lib/*"],
      "@/pages/*": ["./src/pages/*"],
      "@/stores/*": ["./src/stores/*"],
      "@/types/*": ["./src/types/*"],
      "@/utils/*": ["./src/utils/*"]
    }
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
```

---

## Шаг 4: Конфигурация ESLint и Prettier

### 4.1. `.eslintrc.json`

```json
{
  "root": true,
  "env": { "browser": true, "es2020": true },
  "extends": [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:react-hooks/recommended"
  ],
  "ignorePatterns": ["dist", ".eslintrc.cjs"],
  "parser": "@typescript-eslint/parser",
  "plugins": ["react-refresh"],
  "rules": {
    "react-refresh/only-export-components": [
      "warn",
      { "allowConstantExport": true }
    ],
    "@typescript-eslint/no-unused-vars": [
      "error",
      { "argsIgnorePattern": "^_" }
    ],
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/no-explicit-any": "warn"
  }
}
```

### 4.2. `.prettierrc`

```json
{
  "semi": true,
  "trailingComma": "es5",
  "singleQuote": true,
  "printWidth": 80,
  "tabWidth": 2,
  "useTabs": false,
  "arrowParens": "always"
}
```

---

## Шаг 5: Конфигурация Vite

### 5.1. `vite.config.ts`

```typescript
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import { visualizer } from 'rollup-plugin-visualizer';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    // Bundle analyzer - анализ размера bundle
    visualizer({
      open: false,
      filename: 'dist/stats.html',
      gzipSize: true,
      brotliSize: true,
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@/api': path.resolve(__dirname, './src/api'),
      '@/components': path.resolve(__dirname, './src/components'),
      '@/hooks': path.resolve(__dirname, './src/hooks'),
      '@/lib': path.resolve(__dirname, './src/lib'),
      '@/pages': path.resolve(__dirname, './src/pages'),
      '@/stores': path.resolve(__dirname, './src/stores'),
      '@/types': path.resolve(__dirname, './src/types'),
      '@/utils': path.resolve(__dirname, './src/utils'),
    },
  },
  server: {
    port: 3000,
    host: true, // Необходимо для доступа с мобильных устройств
    https: true, // VK MiniApps требует HTTPS
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    // Оптимизация для production
    minify: 'esbuild',
    target: 'es2015',
    rollupOptions: {
      output: {
        // Code splitting для оптимизации bundle size
        manualChunks: {
          // React ecosystem
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          // VK SDK
          'vk-vendor': ['@vkontakte/vk-bridge', '@vkontakte/vkui', '@vkontakte/icons'],
          // State management & data fetching
          'ui-vendor': ['zustand', 'axios'],
        },
      },
    },
    // Warning при превышении 1MB на chunk
    chunkSizeWarningLimit: 1000,
  },
  // Оптимизация для production
  esbuild: {
    drop: process.env.NODE_ENV === 'production' ? ['console', 'debugger'] : [],
  },
});
```

---

## Шаг 6: Environment Variables

### 6.1. `.env.example`

```bash
# API Configuration
VITE_API_BASE_URL=https://api.example.com/api/v1
VITE_API_TIMEOUT=30000

# Platform Configuration
VITE_VK_APP_ID=YOUR_VK_APP_ID
VITE_PLATFORM=vk

# Analytics
VITE_ANALYTICS_ENDPOINT=https://api.example.com/analytics
VITE_ANALYTICS_DEBUG=false

# Feature Flags
VITE_ENABLE_ANALYTICS=true
VITE_ENABLE_MOCK_API=false

# Upload Configuration
VITE_MAX_IMAGE_SIZE_MB=10
VITE_IMAGE_COMPRESSION_QUALITY=0.8
```

### 6.2. `.env.development`

```bash
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_ANALYTICS_DEBUG=true
VITE_ENABLE_MOCK_API=true
```

---

## Шаг 7: VK Bridge Integration и Platform Detection

### 7.1. `src/lib/platform/vk-storage.ts`

**VK Storage API** - надежная альтернатива localStorage для VK miniapps:

```typescript
// src/lib/platform/vk-storage.ts
import bridge from '@vkontakte/vk-bridge';

export const vkStorage = {
  /**
   * Получить значение по ключу
   */
  async getItem(key: string): Promise<string | null> {
    try {
      const data = await bridge.send('VKWebAppStorageGet', { keys: [key] });
      return data.keys[0]?.value || null;
    } catch (error) {
      console.warn('[VKStorage] Fallback to localStorage:', error);
      return localStorage.getItem(key);
    }
  },

  /**
   * Сохранить значение
   */
  async setItem(key: string, value: string): Promise<void> {
    try {
      await bridge.send('VKWebAppStorageSet', { key, value });
    } catch (error) {
      console.warn('[VKStorage] Fallback to localStorage:', error);
      localStorage.setItem(key, value);
    }
  },

  /**
   * Удалить значение
   */
  async removeItem(key: string): Promise<void> {
    try {
      await bridge.send('VKWebAppStorageSet', { key, value: '' });
    } catch (error) {
      console.warn('[VKStorage] Fallback to localStorage:', error);
      localStorage.removeItem(key);
    }
  },

  /**
   * Получить несколько значений сразу
   */
  async getItems(keys: string[]): Promise<Record<string, string>> {
    try {
      const data = await bridge.send('VKWebAppStorageGet', { keys });
      const result: Record<string, string> = {};
      for (const item of data.keys) {
        if (item.value) {
          result[item.key] = item.value;
        }
      }
      return result;
    } catch (error) {
      console.warn('[VKStorage] Fallback to localStorage:', error);
      const result: Record<string, string> = {};
      for (const key of keys) {
        const value = localStorage.getItem(key);
        if (value) result[key] = value;
      }
      return result;
    }
  },
};

// Typed wrappers для типобезопасности
export const storageKeys = {
  TOKEN: 'auth_token',
  USER_ID: 'user_id',
  PROFILE_ID: 'profile_id',
  ONBOARDING_COMPLETED: 'onboarding_completed',
} as const;

/**
 * Type-safe wrapper для работы с storage
 */
export async function getStorageValue<T>(
  key: string,
  defaultValue: T
): Promise<T> {
  const value = await vkStorage.getItem(key);
  if (!value) return defaultValue;

  try {
    return JSON.parse(value) as T;
  } catch {
    return value as T;
  }
}

export async function setStorageValue<T>(key: string, value: T): Promise<void> {
  const stringValue = typeof value === 'string' ? value : JSON.stringify(value);
  await vkStorage.setItem(key, stringValue);
}
```

---

### 7.2. `src/lib/platform/haptics.ts`

**Haptic Feedback** - тактильная обратная связь для лучшего UX:

```typescript
// src/lib/platform/haptics.ts
import bridge from '@vkontakte/vk-bridge';

export const haptics = {
  /**
   * Impact vibration - для кликов и жестов
   */
  impact(style: 'light' | 'medium' | 'heavy' = 'medium'): void {
    bridge
      .send('VKWebAppTapticImpactOccurred', { style })
      .catch((error) => {
        console.debug('[Haptics] Impact not available:', error);
      });
  },

  /**
   * Notification vibration - для результатов операций
   */
  notification(type: 'error' | 'success' | 'warning'): void {
    bridge
      .send('VKWebAppTapticNotificationOccurred', { type })
      .catch((error) => {
        console.debug('[Haptics] Notification not available:', error);
      });
  },

  /**
   * Selection vibration - для переключения между элементами
   */
  selection(): void {
    bridge
      .send('VKWebAppTapticSelectionChanged', {})
      .catch((error) => {
        console.debug('[Haptics] Selection not available:', error);
      });
  },
};

// Использование в компонентах
export function useHaptics() {
  return {
    /**
     * Легкая вибрация для кнопок
     */
    onButtonClick: () => haptics.impact('light'),

    /**
     * Средняя вибрация для важных действий
     */
    onImportantAction: () => haptics.impact('medium'),

    /**
     * Сильная вибрация для критичных действий
     */
    onCriticalAction: () => haptics.impact('heavy'),

    /**
     * Вибрация успеха
     */
    onSuccess: () => haptics.notification('success'),

    /**
     * Вибрация ошибки
     */
    onError: () => haptics.notification('error'),

    /**
     * Вибрация предупреждения
     */
    onWarning: () => haptics.notification('warning'),

    /**
     * Вибрация при переключении
     */
    onSelect: () => haptics.selection(),
  };
}
```

---

### 7.3. `src/lib/platform/platform-detection.ts`

**Platform Detection** - определение платформы для адаптивного UI:

```typescript
// src/lib/platform/platform-detection.ts
import bridge from '@vkontakte/vk-bridge';
import { Platform } from '@vkontakte/vkui';

export type VKPlatform =
  | 'mobile_iphone'
  | 'mobile_ipad'
  | 'mobile_android'
  | 'desktop_web'
  | 'mobile_web';

export interface PlatformInfo {
  platform: Platform; // для VKUI
  vkPlatform: VKPlatform;
  isIOS: boolean;
  isAndroid: boolean;
  isDesktop: boolean;
  isMobile: boolean;
}

let cachedPlatformInfo: PlatformInfo | null = null;

/**
 * Определение платформы через VK Bridge
 */
export async function detectPlatform(): Promise<PlatformInfo> {
  if (cachedPlatformInfo) {
    return cachedPlatformInfo;
  }

  try {
    const launchParams = await bridge.send('VKWebAppGetLaunchParams');
    const vkPlatform = launchParams.vk_platform as VKPlatform;

    const isIOS = vkPlatform.includes('iphone') || vkPlatform.includes('ipad');
    const isAndroid = vkPlatform.includes('android');
    const isDesktop = vkPlatform.includes('desktop');
    const isMobile = !isDesktop;

    let platform: Platform;
    if (isIOS) {
      platform = 'ios';
    } else if (isAndroid) {
      platform = 'android';
    } else {
      platform = 'vkcom'; // desktop
    }

    cachedPlatformInfo = {
      platform,
      vkPlatform,
      isIOS,
      isAndroid,
      isDesktop,
      isMobile,
    };

    console.log('[Platform] Detected:', cachedPlatformInfo);
    return cachedPlatformInfo;
  } catch (error) {
    console.error('[Platform] Detection failed, using fallback:', error);

    // Fallback based on user agent
    const ua = navigator.userAgent.toLowerCase();
    const isIOS = /iphone|ipad|ipod/.test(ua);
    const isAndroid = /android/.test(ua);
    const isMobile = /mobile/.test(ua);

    cachedPlatformInfo = {
      platform: isIOS ? 'ios' : isAndroid ? 'android' : 'vkcom',
      vkPlatform: 'desktop_web',
      isIOS,
      isAndroid,
      isDesktop: !isMobile,
      isMobile,
    };

    return cachedPlatformInfo;
  }
}

/**
 * React hook для использования в компонентах
 */
export function usePlatformDetection() {
  const [platformInfo, setPlatformInfo] = React.useState<PlatformInfo | null>(
    cachedPlatformInfo
  );

  React.useEffect(() => {
    detectPlatform().then(setPlatformInfo);
  }, []);

  return platformInfo;
}
```

---

### 7.4. `src/lib/platform/bridge.ts`

```typescript
import bridge, {
  VKBridgeEvent,
  UserInfo,
} from '@vkontakte/vk-bridge';

export interface PlatformBridge {
  init(): Promise<void>;
  getUser(): Promise<PlatformUser>;
  shareLink(url: string): Promise<void>;
  copyToClipboard(text: string): Promise<void>;
  openURL(url: string): Promise<void>;
  hapticFeedback(type: 'light' | 'medium' | 'heavy'): void;
  requestPhotoAccess(): Promise<boolean>;
  showPopup(params: PopupParams): Promise<void>;
}

export interface PlatformUser {
  id: number;
  firstName: string;
  lastName: string;
  photo?: string;
  platform: 'vk' | 'max' | 'telegram' | 'web';
}

export interface PopupParams {
  title: string;
  message: string;
  buttons: Array<{ text: string; type?: 'primary' | 'destructive' }>;
}

class VKBridge implements PlatformBridge {
  private isInitialized = false;

  async init(): Promise<void> {
    if (this.isInitialized) return;

    try {
      await bridge.send('VKWebAppInit');
      this.isInitialized = true;
      console.log('[VKBridge] Initialized successfully');
    } catch (error) {
      console.error('[VKBridge] Initialization failed:', error);
      throw error;
    }
  }

  async getUser(): Promise<PlatformUser> {
    try {
      const userInfo = await bridge.send('VKWebAppGetUserInfo');

      return {
        id: userInfo.id,
        firstName: userInfo.first_name,
        lastName: userInfo.last_name,
        photo: userInfo.photo_200,
        platform: 'vk',
      };
    } catch (error) {
      console.error('[VKBridge] Failed to get user info:', error);
      throw error;
    }
  }

  async shareLink(url: string): Promise<void> {
    try {
      await bridge.send('VKWebAppShare', { link: url });
    } catch (error) {
      console.error('[VKBridge] Failed to share link:', error);
      throw error;
    }
  }

  async copyToClipboard(text: string): Promise<void> {
    try {
      await bridge.send('VKWebAppCopyText', { text });
    } catch (error) {
      console.error('[VKBridge] Failed to copy to clipboard:', error);
      throw error;
    }
  }

  async openURL(url: string): Promise<void> {
    try {
      await bridge.send('VKWebAppOpenLink', { url });
    } catch (error) {
      console.error('[VKBridge] Failed to open URL:', error);
      throw error;
    }
  }

  hapticFeedback(type: 'light' | 'medium' | 'heavy'): void {
    const impactMap = {
      light: 'light',
      medium: 'medium',
      heavy: 'heavy',
    } as const;

    bridge
      .send('VKWebAppTapticImpactOccurred', {
        style: impactMap[type],
      })
      .catch((error) => {
        console.warn('[VKBridge] Haptic feedback not supported:', error);
      });
  }

  async requestPhotoAccess(): Promise<boolean> {
    try {
      const result = await bridge.send('VKWebAppGetAuthToken', {
        app_id: Number(import.meta.env.VITE_VK_APP_ID),
        scope: 'photos',
      });
      return !!result.access_token;
    } catch (error) {
      console.error('[VKBridge] Failed to request photo access:', error);
      return false;
    }
  }

  async showPopup(params: PopupParams): Promise<void> {
    try {
      await bridge.send('VKWebAppShowAlert', {
        title: params.title,
        message: params.message,
        actions: params.buttons.map((btn) => ({
          title: btn.text,
          mode: btn.type || 'default',
        })),
      });
    } catch (error) {
      console.error('[VKBridge] Failed to show popup:', error);
      throw error;
    }
  }
}

// Singleton instance
export const platformBridge = new VKBridge();
```

### 7.2. `src/hooks/usePlatform.ts`

```typescript
import { useEffect, useState } from 'react';
import { platformBridge, PlatformUser } from '@/lib/platform/bridge';

export function usePlatform() {
  const [user, setUser] = useState<PlatformUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const initPlatform = async () => {
      try {
        setIsLoading(true);
        await platformBridge.init();
        const userData = await platformBridge.getUser();
        setUser(userData);
      } catch (err) {
        setError(err as Error);
        console.error('[usePlatform] Initialization failed:', err);
      } finally {
        setIsLoading(false);
      }
    };

    initPlatform();
  }, []);

  return {
    user,
    isLoading,
    error,
    platform: platformBridge,
  };
}
```

---

## Шаг 8: Базовая тема и стили

### 8.1. `src/config/theme.ts`

```typescript
export const theme = {
  colors: {
    primary: '#5181B8', // VK Blue
    secondary: '#6C7A89',
    success: '#4BB34B',
    error: '#E64646',
    warning: '#FFA000',

    background: {
      primary: '#FFFFFF',
      secondary: '#F5F5F5',
      tertiary: '#EDEEF0',
    },

    text: {
      primary: '#000000',
      secondary: '#818C99',
      tertiary: '#99A2AD',
      inverse: '#FFFFFF',
    },

    border: {
      light: '#E1E3E6',
      medium: '#D3D9DE',
    },
  },

  spacing: {
    xs: '4px',
    sm: '8px',
    md: '16px',
    lg: '24px',
    xl: '32px',
  },

  borderRadius: {
    sm: '4px',
    md: '8px',
    lg: '12px',
    xl: '16px',
    full: '9999px',
  },

  fontSize: {
    xs: '12px',
    sm: '14px',
    md: '16px',
    lg: '18px',
    xl: '24px',
    xxl: '32px',
  },

  fontWeight: {
    regular: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
  },

  shadows: {
    sm: '0 1px 2px rgba(0, 0, 0, 0.05)',
    md: '0 4px 8px rgba(0, 0, 0, 0.1)',
    lg: '0 8px 16px rgba(0, 0, 0, 0.15)',
  },

  zIndex: {
    base: 0,
    dropdown: 1000,
    sticky: 1100,
    modal: 1300,
    popover: 1400,
    tooltip: 1500,
  },
} as const;

export type Theme = typeof theme;
```

### 8.2. `src/index.css`

```css
/* Reset and base styles */
*,
*::before,
*::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

html,
body,
#root {
  height: 100%;
  width: 100%;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #f5f5f5;
  color: #000000;
  overflow-x: hidden;
}

/* VK MiniApp specific */
.vkui {
  height: 100%;
}

/* Disable text selection on interactive elements */
button,
a {
  -webkit-tap-highlight-color: transparent;
  -webkit-touch-callout: none;
  user-select: none;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: #f5f5f5;
}

::-webkit-scrollbar-thumb {
  background: #d3d9de;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #99a2ad;
}
```

---

## Шаг 9: Базовый роутинг

### 9.1. `src/config/routes.ts`

```typescript
export const ROUTES = {
  HOME: '/',
  ONBOARDING: '/onboarding',
  HELP: '/help',
  HELP_UPLOAD: '/help/upload',
  HELP_CROP: '/help/crop',
  HELP_PROCESSING: '/help/processing',
  HELP_RESULT: '/help/result',
  CHECK: '/check',
  CHECK_SCENARIO: '/check/scenario',
  CHECK_UPLOAD: '/check/upload',
  CHECK_PROCESSING: '/check/processing',
  CHECK_RESULT: '/check/result',
  ACHIEVEMENTS: '/achievements',
  FRIENDS: '/friends',
  PROFILE: '/profile',
  PROFILE_HISTORY: '/profile/history',
  PROFILE_SETTINGS: '/profile/settings',
  PROFILE_SUPPORT: '/profile/support',
  VILLAIN: '/villain',
  PAYWALL: '/paywall',
} as const;

export type RoutePath = (typeof ROUTES)[keyof typeof ROUTES];
```

### 9.2. `src/App.tsx`

```typescript
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ROUTES } from '@/config/routes';
import { usePlatform } from '@/hooks/usePlatform';

// Pages (будут созданы позже)
import HomePage from '@/pages/Home';
import OnboardingPage from '@/pages/Onboarding';

function App() {
  const { user, isLoading } = usePlatform();

  if (isLoading) {
    return <div>Loading...</div>; // TODO: Replace with Spinner component
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route path={ROUTES.HOME} element={<HomePage />} />
        <Route path={ROUTES.ONBOARDING} element={<OnboardingPage />} />
        {/* TODO: Add other routes */}
        <Route path="*" element={<Navigate to={ROUTES.HOME} replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
```

---

## Шаг 10: VK MiniApp Manifest

### 10.1. `public/manifest.json`

```json
{
  "name": "Объяснятель ДЗ",
  "short_name": "ДЗ",
  "description": "Помощь с домашними заданиями для школьников",
  "icons": [
    {
      "src": "/vk-app-icon.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ],
  "display": "standalone",
  "start_url": "/",
  "background_color": "#ffffff",
  "theme_color": "#5181B8"
}
```

---

## Шаг 11: VK Sign Validation (Backend Setup)

**КРИТИЧНО:** VK требует проверку `sign` параметра на backend для безопасности!

### 11.1. Backend: Проверка sign (Go)

```go
// Backend: internal/auth/vk_sign.go
package auth

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "sort"
    "strings"
)

// ValidateVKSign проверяет подпись launch params от VK
func ValidateVKSign(params map[string]string, secretKey string) bool {
    // Извлекаем sign
    sign := params["sign"]
    if sign == "" {
        return false
    }

    // Удаляем sign из параметров
    delete(params, "sign")

    // Собираем все vk_* параметры
    var keys []string
    for k := range params {
        if strings.HasPrefix(k, "vk_") {
            keys = append(keys, k)
        }
    }

    // Сортируем ключи
    sort.Strings(keys)

    // Строим query string
    var queryParts []string
    for _, k := range keys {
        queryParts = append(queryParts, k+"="+params[k])
    }
    queryString := strings.Join(queryParts, "&")

    // Вычисляем HMAC-SHA256
    h := hmac.New(sha256.New, []byte(secretKey))
    h.Write([]byte(queryString))
    expectedSign := base64.URLEncoding.EncodeToString(h.Sum(nil))

    // Сравниваем
    return sign == strings.TrimRight(expectedSign, "=")
}
```

### 11.2. Frontend: Отправка launch params на backend

```typescript
// src/api/auth.ts
import bridge from '@vkontakte/vk-bridge';
import { apiClient } from './client';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';

export interface VKAuthResponse {
  access_token: string;
  child_profile_id: string;
  parent_user_id: string;
}

/**
 * Аутентификация через VK Bridge
 */
export async function authenticateVK(): Promise<VKAuthResponse> {
  try {
    // 1. Получаем launch params от VK
    const launchParams = await bridge.send('VKWebAppGetLaunchParams');

    // 2. Отправляем ВСЕ параметры на backend (включая sign!)
    const response = await apiClient.post<VKAuthResponse>('/auth/vk', {
      vk_user_id: launchParams.vk_user_id,
      vk_app_id: launchParams.vk_app_id,
      vk_platform: launchParams.vk_platform,
      vk_language: launchParams.vk_language,
      vk_is_app_user: launchParams.vk_is_app_user,
      vk_are_notifications_enabled: launchParams.vk_are_notifications_enabled,
      vk_ts: launchParams.vk_ts,
      sign: launchParams.sign, // КРИТИЧНО!
    });

    // 3. Сохраняем токен в VK Storage (не localStorage!)
    await vkStorage.setItem(storageKeys.TOKEN, response.access_token);
    await vkStorage.setItem(storageKeys.USER_ID, response.parent_user_id);
    await vkStorage.setItem(storageKeys.PROFILE_ID, response.child_profile_id);

    console.log('[Auth] VK authentication successful');
    return response;
  } catch (error) {
    console.error('[Auth] VK authentication failed:', error);
    throw error;
  }
}

/**
 * Проверка, авторизован ли пользователь
 */
export async function isAuthenticated(): Promise<boolean> {
  const token = await vkStorage.getItem(storageKeys.TOKEN);
  return !!token;
}

/**
 * Выход (очистка токена)
 */
export async function logout(): Promise<void> {
  await vkStorage.removeItem(storageKeys.TOKEN);
  await vkStorage.removeItem(storageKeys.USER_ID);
  await vkStorage.removeItem(storageKeys.PROFILE_ID);
}
```

### 11.3. API Client с JWT токеном

Обновите `src/api/client.ts`:

```typescript
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';

// В setupInterceptors добавьте:
this.client.interceptors.request.use(
  async (config) => {
    // Получаем токен из VK Storage (асинхронно!)
    const token = await vkStorage.getItem(storageKeys.TOKEN);
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    console.log('[API] Request:', config.method?.toUpperCase(), config.url);
    return config;
  },
  (error) => {
    console.error('[API] Request error:', error);
    return Promise.reject(error);
  }
);
```

---

## Чеклист задач

### Базовая настройка
- [ ] Установить Node.js и npm
- [ ] Создать Vite проект с React + TypeScript
- [ ] Установить все зависимости (включая VKUI и icons)
- [ ] Создать структуру папок проекта
- [ ] Настроить TypeScript (`tsconfig.json`)
- [ ] Настроить ESLint и Prettier
- [ ] Настроить Vite с bundle optimization (`vite.config.ts`)
- [ ] Создать `.env.example` и `.env.development`

### VK Bridge Integration
- [ ] Реализовать VK Storage wrapper (`vk-storage.ts`)
- [ ] Реализовать Haptics service (`haptics.ts`)
- [ ] Реализовать Platform Detection (`platform-detection.ts`)
- [ ] Реализовать VK Bridge абстракцию (`bridge.ts`)
- [ ] Создать hook `usePlatform`
- [ ] Создать hook `useHaptics`

### Security & Auth
- [ ] Настроить VK Sign Validation на backend (Go)
- [ ] Реализовать VK аутентификацию на frontend (`auth.ts`)
- [ ] Обновить API client для использования VK Storage
- [ ] Протестировать sign validation flow

### UI & Styling
- [ ] Настроить тему и глобальные стили
- [ ] Импортировать VKUI styles
- [ ] Настроить базовый роутинг
- [ ] Создать VK MiniApp manifest

### Testing
- [ ] Протестировать запуск приложения (`npm run dev`)
- [ ] Проверить bundle size (`npm run build`)
- [ ] Протестировать VK Bridge в VK Dev окружении
- [ ] Протестировать haptics на реальном устройстве
- [ ] Протестировать VK Storage

### Final
- [ ] Добавить `.gitignore` (node_modules, dist, .env)
- [ ] Проверить все import paths работают
- [ ] Запустить линтер и форматтер
- [ ] Создать первый commit

---

## Команды для разработки

```bash
# Разработка
npm run dev

# Сборка
npm run build

# Превью сборки
npm run preview

# Линтинг
npm run lint

# Форматирование
npm run format
```

---

## Следующий этап

После завершения настройки переходи к **02_CORE.md** для создания UI Kit и базовой инфраструктуры.
