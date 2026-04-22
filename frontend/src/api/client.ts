// src/api/client.ts
import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig, AxiosRequestConfig } from 'axios';
import { vkStorage, storageKeys } from '@/lib/platform/vk-storage';
import { createLogger } from '@/lib/logger';
import config from '@/config';

const logger = createLogger('APIClient');

interface ApiErrorResponse {
  message?: string;
  error?: string;
  details?: any;
}

/**
 * Базовый API клиент с interceptors, request deduplication и error handling
 */
class APIClient {
  private client: AxiosInstance;
  private requestCache = new Map<string, Promise<any>>();

  constructor() {
    this.client = axios.create({
      baseURL: config.api.baseURL,
      timeout: config.api.timeout,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor - добавляем Platform ID и Child Profile ID
    this.client.interceptors.request.use(
        async (requestConfig: InternalAxiosRequestConfig) => {
        console.log('[APIClient] Request interceptor START', { url: requestConfig.url });

        // Добавляем X-Platform-ID (определяем платформу)
        console.log('[APIClient] Getting platform_id from storage...');
        let platformID = localStorage.getItem('platform_id');
        console.log('[APIClient] platformID from storage:', platformID);

        if (!platformID) {
          // Fallback: определяем из URL или используем web
          const urlParams = new URLSearchParams(window.location.search);
          platformID = urlParams.get('vk_platform') ? 'vk' : 'web';
          console.log('[APIClient] Using fallback platformID:', platformID);
          localStorage.setItem('platform_id', platformID);
        }

        // Dev режим локально - всегда устанавливаем платформу принудительно
        if (import.meta.env.DEV) {
          platformID = 'web';
        }

        if (requestConfig.headers) {
          requestConfig.headers['X-Platform-ID'] = platformID;
          console.log('[APIClient] Set X-Platform-ID header:', platformID);

          // Добавляем X-Child-Profile-ID если есть
          console.log('[APIClient] Getting profile_id from storage...');
          const childProfileId = await vkStorage.getItem(storageKeys.PROFILE_ID);
          console.log('[APIClient] childProfileId from storage:', childProfileId);
          if (childProfileId) {
            requestConfig.headers['X-Child-Profile-ID'] = childProfileId;
            console.log('[APIClient] Set X-Child-Profile-ID header:', childProfileId);
          }
        }

        console.log('[APIClient] Request interceptor COMPLETE', {
          url: requestConfig.url,
          headers: requestConfig.headers,
        });

        logger.debug('API Request', {
          method: requestConfig.method?.toUpperCase(),
          url: requestConfig.url,
          platformID,
          hasProfileID: !!requestConfig.headers?.['X-Child-Profile-ID'],
        });

        return requestConfig;
      },
      (error) => {
        logger.error('Request interceptor error', { error: error.message });
        return Promise.reject(error);
      }
    );

    // Response interceptor - обработка ошибок
    this.client.interceptors.response.use(
      (response) => {
        logger.debug('API Response', {
          status: response.status,
          url: response.config.url,
        });
        return response;
      },
      async (error: AxiosError<ApiErrorResponse>) => {
        if (error.response) {
          // Сервер вернул ошибку
          const status = error.response.status;
          const data = error.response.data;

          logger.error('API Error Response', {
            status,
            url: error.config?.url,
            message: data?.message || data?.error,
            details: data?.details,
          });

          if (status === 401) {
            // Unauthorized - очищаем данные профиля
            logger.warn('Unauthorized - clearing profile data');
            await vkStorage.removeItem(storageKeys.PROFILE_ID);
            await vkStorage.removeItem(storageKeys.ONBOARDING_COMPLETED);
            // Редирект обрабатывается в AppInitializer
          }

          if (status === 429) {
            // Rate limit
            logger.warn('Rate limit exceeded');
          }

          if (status >= 500) {
            // Server error
            logger.error('Server error', {
              status,
              url: error.config?.url,
            });
          }
        } else if (error.request) {
          // Запрос был отправлен, но ответа не получено
          logger.error('No response received', {
            url: error.config?.url,
            timeout: error.code === 'ECONNABORTED',
          });
        } else {
          // Ошибка при настройке запроса
          logger.error('Request setup error', {
            message: error.message,
          });
        }

        return Promise.reject(error);
      }
    );
  }

  /**
   * Генерирует ключ для кеша запросов
   */
  private getCacheKey(method: string, url: string, data?: any): string {
    const dataKey = data ? JSON.stringify(data) : '';
    return `${method}:${url}:${dataKey}`;
  }

  /**
   * Проверяет нужно ли retry для данной ошибки
   */
  private shouldRetry(error: AxiosError): boolean {
    if (!error.response) {
      // Network error - можно retry
      return true;
    }

    const status = error.response.status;
    // Retry для 5xx и 429
    return status >= 500 || status === 429;
  }

  /**
   * Retry с exponential backoff
   */
  private async retryRequest<T>(
    requestFn: () => Promise<T>,
    maxRetries: number = 2,
    baseDelay: number = 1000
  ): Promise<T> {
    let lastError: Error | null = null;

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        return await requestFn();
      } catch (error) {
        lastError = error as Error;

        if (
          attempt < maxRetries &&
          axios.isAxiosError(error) &&
          this.shouldRetry(error)
        ) {
          const delay = baseDelay * Math.pow(2, attempt);
          logger.warn('Retrying request', {
            attempt: attempt + 1,
            maxRetries,
            delay,
          });

          await new Promise((resolve) => setTimeout(resolve, delay));
          continue;
        }

        throw error;
      }
    }

    throw lastError;
  }

  /**
   * GET request с deduplication и retry
   */
  async get<T>(url: string, requestConfig: AxiosRequestConfig = {}): Promise<T> {
    const cacheKey = this.getCacheKey('GET', url);

    // Если запрос уже выполняется - возвращаем существующий promise
    if (this.requestCache.has(cacheKey)) {
      logger.debug('Using cached request', { url });
      return this.requestCache.get(cacheKey);
    }

    const requestPromise = this.retryRequest(async () => {
      const response = await this.client.get<T>(url, requestConfig);
      return response.data;
    })
      .finally(() => {
        this.requestCache.delete(cacheKey);
      });

    this.requestCache.set(cacheKey, requestPromise);
    return requestPromise;
  }

  /**
   * POST request
   */
  async post<T>(url: string, data?: unknown, requestConfig: AxiosRequestConfig = {}): Promise<T> {
    const response = await this.client.post<T>(url, data, requestConfig);
    return response.data;
  }

  /**
   * PATCH request
   */
  async patch<T>(url: string, data?: unknown, requestConfig: AxiosRequestConfig = {}): Promise<T> {
    const response = await this.client.patch<T>(url, data, requestConfig);
    return response.data;
  }

  /**
   * DELETE request
   */
  async delete<T>(url: string, requestConfig: AxiosRequestConfig = {}): Promise<T> {
    const response = await this.client.delete<T>(url, requestConfig);
    return response.data;
  }

  /**
   * Multipart form data upload (для фото)
   */
  async upload<T>(
    url: string,
    file: Blob | FormData,
    onProgress?: (progress: number) => void
  ): Promise<T> {
    const formData = file instanceof Blob
      ? (() => {
          const fd = new FormData();
          fd.append('file', file);
          return fd;
        })()
      : file;

    logger.info('Starting file upload', {
      url,
      fileSize: file instanceof Blob ? file.size : undefined,
    });

    const response = await this.client.post<T>(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      timeout: 60000, // 60 seconds для загрузки файлов
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percentage = Math.round(
            (progressEvent.loaded * 100) / progressEvent.total
          );
          onProgress(percentage);

          logger.debug('Upload progress', {
            url,
            percentage,
            loaded: progressEvent.loaded,
            total: progressEvent.total,
          });
        }
      },
    });

    logger.info('File upload completed', { url });
    return response.data;
  }

  /**
   * Очищает кеш запросов
   */
  clearCache(): void {
    this.requestCache.clear();
    logger.debug('Request cache cleared');
  }
}

// Singleton instance
export const apiClient = new APIClient();
