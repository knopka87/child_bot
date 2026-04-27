// src/api/achievements.ts
import { apiClient } from './client';
import type { Achievement, AchievementsStats } from '@/types/achievements';

export const achievementsAPI = {
  /**
   * Получить все достижения
   */
  async getAchievements(): Promise<Achievement[]> {
    return apiClient.get<Achievement[]>('/api/achievements');
  },

  /**
   * Получить статистику достижений
   */
  async getAchievementsStats(): Promise<AchievementsStats> {
    return apiClient.get<AchievementsStats>('/api/achievements/stats');
  },

  /**
   * Получить детали достижения
   */
  async getAchievementDetail(achievementId: string): Promise<Achievement> {
    return apiClient.get<Achievement>(`/api/achievements/${achievementId}`);
  },

  /**
   * Получить награду за достижение
   */
  async claimAchievementReward(achievementId: string): Promise<void> {
    return apiClient.post<void>(`/api/achievements/${achievementId}/claim`);
  },

  /**
   * Проверить есть ли новые достижения
   */
  async hasNewAchievements(): Promise<{ has_new: boolean }> {
    return apiClient.get<{ has_new: boolean }>('/api/achievements/has-new');
  },

  /**
   * Отметить что достижения просмотрены
   */
  async markAchievementsViewed(): Promise<void> {
    return apiClient.post<void>('/api/achievements/mark-viewed');
  },
};
