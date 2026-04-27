// src/api/home.ts
import { apiClient } from './client';
import type { HomeData, RecentAttempt } from '@/types/home';
import type { Attempt } from '@/types/domain';

export const homeAPI = {
  /**
   * Получить все данные для главного экрана
   */
  async getHomeData(childProfileId: string): Promise<HomeData> {
    return apiClient.get<HomeData>(`/api/home/${childProfileId}`);
  },

  /**
   * Получить незавершенную попытку
   */
  async getUnfinishedAttempt(childProfileId: string): Promise<Attempt | null> {
    return apiClient.get<Attempt | null>(`/api/attempts/unfinished`, {
      params: { childProfileId },
    });
  },

  /**
   * Получить последние попытки
   */
  async getRecentAttempts(
    childProfileId: string,
    limit: number = 3
  ): Promise<RecentAttempt[]> {
    return apiClient.get<RecentAttempt[]>(`/api/attempts/recent`, {
      params: { childProfileId, limit },
    });
  },

  /**
   * Удалить незавершенную попытку
   */
  async deleteAttempt(attemptId: string): Promise<void> {
    return apiClient.delete<void>(`/api/attempts/${attemptId}`);
  },
};
