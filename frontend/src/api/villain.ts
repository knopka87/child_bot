// src/api/villain.ts
import { apiClient } from './client';
import type {
  Villain,
  VillainBattle,
  VillainVictory,
} from '@/types/villain';

export const villainAPI = {
  /**
   * Получить активного злодея
   */
  async getActiveVillain(): Promise<Villain | null> {
    return apiClient.get<Villain | null>('/api/villains/active');
  },

  /**
   * Получить детали битвы
   */
  async getBattleDetails(villainId: string): Promise<VillainBattle> {
    return apiClient.get<VillainBattle>(`/api/villains/${villainId}/battle`);
  },

  /**
   * Получить детали победы
   */
  async getVictoryDetails(villainId: string): Promise<VillainVictory> {
    return apiClient.get<VillainVictory>(`/api/villains/${villainId}/victory`);
  },

  /**
   * Нанести урон злодею
   */
  async dealDamage(
    villainId: string,
    attemptId: string,
    damage: number
  ): Promise<{
    damage_dealt: number;
    villain_hp: number;
    is_defeated: boolean;
    message: string;
  }> {
    return apiClient.post(`/villains/${villainId}/damage`, {
      attempt_id: attemptId,
      damage,
    });
  },
};
