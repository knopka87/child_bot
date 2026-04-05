// src/types/referral.ts

export interface ReferralData {
  referralCode: string;
  referralLink: string;
  invitedCount: number;
  targetCount: number;
  progressPercent: number;
  currentGoal: ReferralGoal;
  invitedFriends: InvitedFriend[];
}

export interface ReferralGoal {
  id: string;
  targetCount: number;
  reward: ReferralReward;
  isCompleted: boolean;
  completedAt?: string;
}

export interface ReferralReward {
  type: 'sticker' | 'coins' | 'avatar' | 'premium_feature';
  id: string;
  name: string;
  description: string;
  imageUrl?: string;
  amount?: number;
}

export interface InvitedFriend {
  id: string;
  displayName: string;
  avatarUrl?: string;
  invitedAt: string;
  status: 'pending' | 'active' | 'completed_first_task';
}

export interface ShareChannel {
  type: 'telegram' | 'vk' | 'whatsapp' | 'copy' | 'native';
  name: string;
  icon: string;
}
