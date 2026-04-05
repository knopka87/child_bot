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
  reward_milestones: Array<{
    friends_count: number;
    reward: number;
    is_claimed: boolean;
    description: string;
  }>;
}

export const referralAPI = {
  /**
   * Получить реферальные данные
   */
  async getReferralData(_childProfileId: string): Promise<ReferralData> {
    console.log('[referralAPI] getReferralData called');
    const response = await apiClient.get<BackendReferralResponse>('/friends/referrals');
    console.log('[referralAPI] Response received:', response);

    // Находим первый незавершённый milestone
    const nextMilestone = response.reward_milestones.find(m => !m.is_claimed);
    const targetCount = nextMilestone?.friends_count || 5;
    const progressPercent = Math.min((response.active_invited / targetCount) * 100, 100);

    // Адаптируем данные к формату фронтенда
    return {
      referralCode: response.referral_code,
      referralLink: response.referral_link,
      invitedCount: response.active_invited,
      targetCount,
      progressPercent,
      currentGoal: {
        id: String(targetCount),
        targetCount,
        reward: {
          type: 'sticker',
          id: 'friendship_sticker',
          name: 'Редкий стикер «Дружба»',
          description: nextMilestone?.description || '',
          amount: nextMilestone?.reward,
        },
        isCompleted: response.active_invited >= targetCount,
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
