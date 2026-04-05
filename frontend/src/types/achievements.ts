// src/types/achievements.ts

// Achievement ID - это string (uuid или slug), НЕ enum
export type AchievementID = string;

// Категории приходят с бекенда, могут добавляться новые
export type AchievementCategory =
  | 'streak'
  | 'tasks'
  | 'fixes'
  | 'milestones'
  | 'mastery'
  | 'speed_solver'
  | 'villain_defeater'
  | 'wise_owl'
  | 'collector'
  | 'knowledge_rocket'
  | 'superstar'
  | 'marathoner'
  | 'genius'
  | string; // Backend может добавить новые

export type AchievementType = AchievementCategory;

export interface Achievement {
  id: string;
  type: AchievementType;
  title: string;
  description: string;
  icon: string; // emoji or image URL
  is_unlocked: boolean;
  unlocked_at?: string;
  progress: {
    current: number;
    total: number;
    percent: number;
  };
  reward: AchievementReward;
  shelf_order: number; // 1, 2, 3
  position_in_shelf: number; // 0-3
}

export interface AchievementReward {
  type: 'sticker' | 'coins' | 'avatar' | 'badge';
  id: string;
  name: string;
  image_url?: string;
  amount?: number; // для coins
}

export interface AchievementsStats {
  totalAchievements: number;
  unlockedAchievements: number;
  totalCoinsEarned: number;
  // Backend может вернуть в snake_case
  total_count?: number;
  unlocked_count?: number;
  totalCount?: number;
  unlockedCount?: number;
}
