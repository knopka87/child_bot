// src/api/referral.ts
import { apiClient } from './client';
import type {
  ReferralData,
  ReferralGoal,
  InvitedFriend,
} from '@/types/referral';

// Тип ответа от бэкенда
interface BackendReferralResponse {
  referral_code: string;
  referral_link: string;
  total_invited: number;
  active_invited: number;
  total_rewards: number;
  invited_friends: Array<{
    id: string;
    display_name: string;
    avatar_url: string;
    joined_at: string;
    is_active: boolean;
    reward_earned: number;
  }>;
  current_achievement?: {
    achievement_id: string;
    title: string;
    description: string;
    icon: string;
    target_count: number;
    current_count: number;
    previous_level: number;
    is_unlocked: boolean;
    next_level?: number;
  };
}

export const referralAPI = {
  /**
   * Получить реферальные данные
   */
  async getReferralData(_childProfileId: string): Promise<ReferralData> {
    console.log('[referralAPI] getReferralData called');
    const response = await apiClient.get<BackendReferralResponse>('/friends/referrals');
    console.log('[referralAPI] Response received:', response);

    // Используем данные о достижении "Дружба"
    let targetCount: number; // Относительное количество (сколько друзей нужно пригласить от предыдущего уровня)
    let invitedCount: number; // Относительное количество (сколько приглашено от предыдущего уровня)
    let goalDescription: string;
    let rewardName: string;

    if (response.current_achievement) {
      // Вычисляем относительные значения
      const previousLevel = response.current_achievement.previous_level;
      targetCount = response.current_achievement.target_count - previousLevel;
      invitedCount = Math.max(0, response.current_achievement.current_count - previousLevel);

      goalDescription = response.current_achievement.description;
      rewardName = `Стикер «${response.current_achievement.title}»`;
    } else {
      // Нет данных о достижении - используем значения по умолчанию
      targetCount = 5;
      invitedCount = 0;
      goalDescription = 'За 5 приглашённых друзей';
      rewardName = 'Стикер «Дружба»';
    }

    const progressPercent = targetCount > 0 ? Math.min((invitedCount / targetCount) * 100, 100) : 0;

    // Адаптируем данные к формату фронтенда
    return {
      referralCode: response.referral_code,
      referralLink: response.referral_link,
      invitedCount, // Относительное количество (от предыдущего уровня)
      targetCount, // Относительная цель (от предыдущего уровня)
      totalInvited: response.total_invited, // Абсолютное общее количество
      progressPercent,
      currentGoal: {
        id: response.current_achievement?.achievement_id || String(targetCount),
        targetCount,
        reward: {
          type: 'sticker',
          id: response.current_achievement?.achievement_id || 'friendship_sticker',
          name: rewardName,
          description: goalDescription,
        },
        isCompleted: response.current_achievement?.is_unlocked || false,
      },
      invitedFriends: response.invited_friends.map(f => ({
        id: f.id,
        displayName: f.display_name,
        avatarUrl: f.avatar_url,
        invitedAt: f.joined_at,
        status: f.is_active ? 'active' : 'pending',
      })),
    };
  },

  /**
   * Сгенерировать реферальный код
   */
  async generateReferralCode(
    childProfileId: string
  ): Promise<{ code: string; link: string }> {
    return apiClient.post<{ code: string; link: string }>(
      `/referrals/${childProfileId}/generate`
    );
  },

  /**
   * Получить список приглашённых друзей
   */
  async getInvitedFriends(childProfileId: string): Promise<InvitedFriend[]> {
    return apiClient.get<InvitedFriend[]>(
      `/referrals/${childProfileId}/friends`
    );
  },

  /**
   * Получить текущую цель
   */
  async getCurrentGoal(childProfileId: string): Promise<ReferralGoal> {
    return apiClient.get<ReferralGoal>(`/referrals/${childProfileId}/goal`);
  },

  /**
   * Получить награду за достижение цели
   */
  async claimReferralReward(
    childProfileId: string,
    goalId: string
  ): Promise<void> {
    return apiClient.post<void>(
      `/referrals/${childProfileId}/goal/${goalId}/claim`
    );
  },

  /**
   * Трекинг отправки приглашения
   */
  async trackInviteSent(
    childProfileId: string,
    channel: string
  ): Promise<void> {
    return apiClient.post<void>(`/referrals/${childProfileId}/track-invite`, {
      channel,
    });
  },
};
