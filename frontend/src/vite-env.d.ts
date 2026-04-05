/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_API_TIMEOUT: string;
  readonly VITE_VK_APP_ID: string;
  readonly VITE_PLATFORM: string;
  readonly VITE_ANALYTICS_ENDPOINT: string;
  readonly VITE_ANALYTICS_DEBUG: string;
  readonly VITE_ENABLE_ANALYTICS: string;
  readonly VITE_ENABLE_MOCK_API: string;
  readonly VITE_MAX_IMAGE_SIZE_MB: string;
  readonly VITE_IMAGE_COMPRESSION_QUALITY: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
