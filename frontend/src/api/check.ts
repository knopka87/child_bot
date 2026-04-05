// src/api/check.ts
import { apiClient } from './client';
import type {
  CreateCheckAttemptResponse,
  CheckResult,
  CheckScenario,
  ImageRole,
} from '@/types/check';
import type { CropArea } from '@/types/help';

export const checkAPI = {
  /**
   * Создать новую попытку проверки
   */
  async createAttempt(
    childProfileId: string,
    scenario: CheckScenario
  ): Promise<CreateCheckAttemptResponse> {
    const response = await apiClient.post<{ attempt_id: string; status: string }>(
      '/attempts',
      {
        child_profile_id: childProfileId,
        type: 'check',
        scenario_type: scenario,
      }
    );

    // Преобразуем snake_case в camelCase
    return {
      id: response.attempt_id,
      status: response.status,
    };
  },

  /**
   * Загрузить изображение (задание или ответ)
   */
  async uploadImage(
    attemptId: string,
    imageRole: ImageRole,
    file: Blob,
    onProgress?: (progress: number) => void
  ): Promise<{ imageUrl: string; thumbnailUrl: string }> {
    // Конвертируем Blob в base64 data URI (с префиксом data:image/...;base64,)
    const base64DataUri = await new Promise<string>((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        const result = reader.result as string;
        // Оставляем полный data URI (бэкенд ожидает префикс data:image/)
        resolve(result);
      };
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });

    // Имитация прогресса (т.к. base64 конвертация мгновенная)
    if (onProgress) {
      onProgress(50);
    }

    const response = await apiClient.post<{ image_url: string; message: string }>(
      `/attempts/${attemptId}/images`,
      {
        image_type: imageRole, // 'task' или 'answer'
        image_data: base64DataUri,
      }
    );

    if (onProgress) {
      onProgress(100);
    }

    // Преобразуем snake_case в camelCase
    return {
      imageUrl: response.image_url,
      thumbnailUrl: response.image_url, // Бэкенд пока возвращает только один URL
    };
  },

  /**
   * Подтвердить качество изображения
   */
  async confirmQuality(
    attemptId: string,
    imageRole: ImageRole
  ): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/confirm-quality`, {
      imageRole,
    });
  },

  /**
   * Обрезать изображение
   */
  async cropImage(
    attemptId: string,
    imageRole: ImageRole,
    cropArea: CropArea
  ): Promise<{ imageUrl: string }> {
    return apiClient.post<{ imageUrl: string }>(
      `/attempts/${attemptId}/crop`,
      { imageRole, cropArea }
    );
  },

  /**
   * Начать обработку проверки
   */
  async processAttempt(attemptId: string): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/process`);
  },

  /**
   * Получить результат проверки (polling)
   */
  async getResult(attemptId: string): Promise<CheckResult> {
    return apiClient.get<CheckResult>(`/attempts/${attemptId}/result`);
  },

  /**
   * Исправить и переотправить
   */
  async resubmit(
    attemptId: string,
    file: Blob,
    onProgress?: (progress: number) => void
  ): Promise<{ success: boolean }> {
    await checkAPI.uploadImage(attemptId, 'answer', file, onProgress);
    await checkAPI.processAttempt(attemptId);
    return { success: true };
  },

  /**
   * Сохранить и подождать
   */
  async saveAndWait(attemptId: string): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/save-and-wait`);
  },

  /**
   * Отменить попытку
   */
  async cancelAttempt(attemptId: string): Promise<void> {
    return apiClient.delete<void>(`/attempts/${attemptId}`);
  },
};
