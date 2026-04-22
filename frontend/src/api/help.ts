// src/api/help.ts
import { apiClient } from './client';
import type {
  CreateAttemptResponse,
  HelpResult,
  Hint,
  CropArea,
} from '@/types/help';

export const helpAPI = {
  /**
   * Создать новую попытку помощи
   */
  async createAttempt(childProfileId: string): Promise<CreateAttemptResponse> {
    const response = await apiClient.post<{ attempt_id: string; status: string }>(
      '/attempts',
      {
        child_profile_id: childProfileId,
        type: 'help',
      }
    );

    // Преобразуем snake_case в camelCase
    return {
      id: response.attempt_id,
      status: response.status,
    };
  },

  /**
   * Загрузить изображение задания
   */
  async uploadImage(
    attemptId: string,
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
        image_type: 'task',
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
  async confirmQuality(attemptId: string): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/confirm-quality`);
  },

  /**
   * Обрезать изображение
   */
  async cropImage(attemptId: string, cropArea: CropArea): Promise<{ imageUrl: string }> {
    return apiClient.post<{ imageUrl: string }>(
      `/attempts/${attemptId}/crop`,
      { cropArea }
    );
  },

  /**
   * Начать обработку задания
   */
  async processAttempt(attemptId: string): Promise<void> {
    return apiClient.post<void>(`/attempts/${attemptId}/process`);
  },

  /**
   * Получить результат обработки (polling)
   */
  async getResult(attemptId: string): Promise<HelpResult> {
    const response = await apiClient.get<any>(`/attempts/${attemptId}/result`);

    return {
      attemptId: response.attempt_id || attemptId,
      hints: response.result?.hints || [],
      coinsEarned: response.result?.coins_earned || 0,
      damageDealt: response.result?.damage_dealt || 0,
      taskImage: response.result?.task_image, // Изображение задания
    };
  },

  /**
   * Получить следующую подсказку
   */
  async getNextHint(attemptId: string, currentLevel: number): Promise<Hint> {
    return apiClient.post<Hint>(`/attempts/${attemptId}/next-hint`, {
      currentLevel,
    });
  },

  /**
   * Отправить ответ пользователя
   */
  async submitAnswer(
    attemptId: string,
    answer: string
  ): Promise<{ success: boolean; coinsEarned: number }> {
    return apiClient.post<{ success: boolean; coinsEarned: number }>(
      `/attempts/${attemptId}/submit-answer`,
      { answer }
    );
  },

  /**
   * Сохранить и подождать (для long-wait)
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
