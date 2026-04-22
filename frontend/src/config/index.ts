// src/config/index.ts
/**
 * Централизованная конфигурация приложения
 * Все настройки берутся из environment variables
 */

export type Environment = 'development' | 'staging' | 'production' | 'test';

export interface AppConfig {
  environment: Environment;
  isDevelopment: boolean;
  isProduction: boolean;
  isTest: boolean;

  api: {
    baseURL: string;
    timeout: number;
  };

  analytics: {
    enabled: boolean;
    debug: boolean;
    batchSize: number;
    batchInterval: number;
    retryAttempts: number;
    retryDelay: number;
  };

  features: {
    villainMode: boolean;
    achievements: boolean;
    referrals: boolean;
    offlineSupport: boolean;
  };

  platforms: {
    vk: {
      appId: string;
    };
    max: {
      appId: string;
    };
    telegram: {
      botUsername: string;
    };
  };

  app: {
    version: string;
    name: string;
  };
}

const getEnvironment = (): Environment => {
  const mode = import.meta.env.MODE;
  if (mode === 'test') return 'test';
  if (mode === 'production') return 'production';
  if (mode === 'staging') return 'staging';
  return 'development';
};

const environment = getEnvironment();

const config: AppConfig = {
  // Environment
  environment,
  isDevelopment: environment === 'development',
  isProduction: environment === 'production',
  isTest: environment === 'test',

  // API Configuration
  api: {
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
    timeout: parseInt(import.meta.env.VITE_API_TIMEOUT || '30000', 10),
  },

  // Analytics Configuration
  analytics: {
    enabled: import.meta.env.VITE_ANALYTICS_ENABLED === 'true' || import.meta.env.PROD,
    debug: import.meta.env.VITE_ANALYTICS_DEBUG === 'true' || import.meta.env.DEV,
    batchSize: parseInt(import.meta.env.VITE_ANALYTICS_BATCH_SIZE || '10', 10),
    batchInterval: parseInt(import.meta.env.VITE_ANALYTICS_BATCH_INTERVAL || '10000', 10),
    retryAttempts: parseInt(import.meta.env.VITE_ANALYTICS_RETRY_ATTEMPTS || '3', 10),
    retryDelay: parseInt(import.meta.env.VITE_ANALYTICS_RETRY_DELAY || '2000', 10),
  },

  // Feature Flags
  features: {
    villainMode: import.meta.env.VITE_FEATURE_VILLAIN !== 'false', // enabled by default
    achievements: import.meta.env.VITE_FEATURE_ACHIEVEMENTS !== 'false',
    referrals: import.meta.env.VITE_FEATURE_REFERRALS !== 'false',
    offlineSupport: import.meta.env.VITE_FEATURE_OFFLINE === 'true',
  },

  // Platform IDs
  platforms: {
    vk: {
      appId: import.meta.env.VITE_VK_APP_ID || '',
    },
    max: {
      appId: import.meta.env.VITE_MAX_APP_ID || '',
    },
    telegram: {
      botUsername: import.meta.env.VITE_TELEGRAM_BOT_USERNAME || '',
    },
  },

  // App Info
  app: {
    version: import.meta.env.VITE_APP_VERSION || '0.1.0',
    name: import.meta.env.VITE_APP_NAME || 'Homework Helper',
  },
};

// Validate critical config in non-production environments
if (!config.isProduction) {
  if (!config.api.baseURL) {
    console.warn('[Config] API base URL is not set');
  }
}

export default config;
