// src/lib/validation/schemas.ts
/**
 * Runtime validation schemas с использованием zod
 * Обеспечивает type safety для данных из внешних источников (API, localStorage, platform bridges)
 */

import { z } from 'zod';

// ============================================================================
// Villain Schemas
// ============================================================================

export const VillainSchema = z.object({
  id: z.string().min(1),
  name: z.string().min(1),
  description: z.string(),
  imageUrl: z.string().url().optional().or(z.literal('')),
  healthPercent: z.number().min(0).max(100),
  currentHealth: z.number().nonnegative(),
  maxHealth: z.number().positive(),
  taunt: z.string(),
  isActive: z.boolean(),
  isDefeated: z.boolean(),
});

export const VillainRewardSchema = z.object({
  type: z.enum(['sticker', 'achievement', 'coins', 'avatar']),
  id: z.string(),
  name: z.string(),
  description: z.string(),
  imageUrl: z.string().url().optional(),
  amount: z.number().optional(),
  rarity: z.enum(['common', 'rare', 'epic', 'legendary']).optional(),
});

export const VillainVictorySchema = z.object({
  villainId: z.string(),
  totalDamage: z.number().nonnegative(),
  totalAttempts: z.number().int().positive(),
  rewards: z.array(VillainRewardSchema),
  completedAt: z.number().optional(),
});

// ============================================================================
// Profile Schemas
// ============================================================================

export const UserProfileSchema = z.object({
  child_profile_id: z.string().min(1),
  display_name: z.string().min(1).max(50),
  grade: z.number().int().min(1).max(11),
  level: z.number().int().nonnegative(),
  coins_balance: z.number().int().nonnegative(),
  avatar_url: z.string().url().optional(),
  tasks_solved_correct_count: z.number().int().nonnegative().optional(),
  wins_count: z.number().int().nonnegative().optional(),
  checks_correct_count: z.number().int().nonnegative().optional(),
  current_streak_days: z.number().int().nonnegative().optional(),
  has_unfinished_attempt: z.boolean().optional(),
  active_villain_id: z.string().optional(),
  active_villain_health_percent: z.number().min(0).max(100).optional(),
  invited_count_total: z.number().int().nonnegative().optional(),
  achievements_unlocked_count: z.number().int().nonnegative().optional(),
});

// ============================================================================
// Achievement Schemas
// ============================================================================

export const AchievementSchema = z.object({
  id: z.string().min(1),
  title: z.string().min(1),
  description: z.string(),
  icon_url: z.string().url().optional(),
  is_unlocked: z.boolean(),
  unlocked_at: z.number().optional(),
  progress: z.number().min(0).max(100).optional(),
  target: z.number().positive().optional(),
  current: z.number().nonnegative().optional(),
  reward_type: z.enum(['coins', 'sticker', 'avatar', 'badge']).optional(),
  reward_amount: z.number().optional(),
});

// ============================================================================
// Task/Attempt Schemas
// ============================================================================

export const TaskSchema = z.object({
  id: z.string().min(1),
  grade: z.number().int().min(1).max(11),
  subject: z.string().min(1),
  difficulty: z.enum(['easy', 'medium', 'hard']),
  created_at: z.number(),
  updated_at: z.number().optional(),
});

export const AttemptSchema = z.object({
  id: z.string().min(1),
  task_id: z.string().min(1),
  child_profile_id: z.string().min(1),
  status: z.enum(['pending', 'processing', 'answered', 'checked']),
  image_url: z.string().url().optional(),
  answer: z.string().optional(),
  is_correct: z.boolean().optional(),
  damage: z.number().nonnegative().optional(),
  hints_used_count: z.number().int().nonnegative().optional(),
  created_at: z.number(),
  updated_at: z.number().optional(),
});

// ============================================================================
// Friends/Referral Schemas
// ============================================================================

export const FriendSchema = z.object({
  id: z.string().min(1),
  display_name: z.string().min(1),
  avatar_url: z.string().url().optional(),
  level: z.number().int().nonnegative(),
  is_online: z.boolean().optional(),
  last_active_at: z.number().optional(),
});

export const ReferralSchema = z.object({
  id: z.string().min(1),
  referred_user_id: z.string(),
  referred_user_name: z.string(),
  status: z.enum(['pending', 'active', 'completed']),
  reward_coins: z.number().int().nonnegative(),
  reward_claimed: z.boolean(),
  created_at: z.number(),
});

// ============================================================================
// Subscription Schemas
// ============================================================================

export const SubscriptionPlanSchema = z.object({
  id: z.string().min(1),
  name: z.string().min(1),
  description: z.string(),
  price: z.number().nonnegative(),
  currency: z.string().length(3),
  duration_days: z.number().int().positive(),
  features: z.array(z.string()),
});

export const SubscriptionStatusSchema = z.object({
  is_active: z.boolean(),
  plan_id: z.string().optional(),
  started_at: z.number().optional(),
  expires_at: z.number().optional(),
  auto_renew: z.boolean().optional(),
  is_trial: z.boolean().optional(),
  trial_ends_at: z.number().optional(),
});

// ============================================================================
// API Response Wrappers
// ============================================================================

export const ApiSuccessResponseSchema = <T extends z.ZodType>(dataSchema: T) =>
  z.object({
    success: z.literal(true),
    data: dataSchema,
  });

export const ApiErrorResponseSchema = z.object({
  success: z.literal(false),
  error: z.string(),
  message: z.string().optional(),
  details: z.any().optional(),
});

export const ApiResponseSchema = <T extends z.ZodType>(dataSchema: T) =>
  z.union([ApiSuccessResponseSchema(dataSchema), ApiErrorResponseSchema]);

// ============================================================================
// Type Inference
// ============================================================================

export type Villain = z.infer<typeof VillainSchema>;
export type VillainReward = z.infer<typeof VillainRewardSchema>;
export type VillainVictory = z.infer<typeof VillainVictorySchema>;
export type UserProfile = z.infer<typeof UserProfileSchema>;
export type Achievement = z.infer<typeof AchievementSchema>;
export type Task = z.infer<typeof TaskSchema>;
export type Attempt = z.infer<typeof AttemptSchema>;
export type Friend = z.infer<typeof FriendSchema>;
export type Referral = z.infer<typeof ReferralSchema>;
export type SubscriptionPlan = z.infer<typeof SubscriptionPlanSchema>;
export type SubscriptionStatus = z.infer<typeof SubscriptionStatusSchema>;

// ============================================================================
// Validation Helper Functions
// ============================================================================

/**
 * Безопасная валидация с логированием ошибок
 */
export function validateData<T>(
  schema: z.ZodType<T>,
  data: unknown,
  context?: string
): T {
  const result = schema.safeParse(data);

  if (!result.success) {
    const errorMessage = `Validation failed${context ? ` for ${context}` : ''}`;
    console.error(errorMessage, {
      errors: result.error.issues,
      data,
    });

    throw new Error(`${errorMessage}: ${result.error.issues[0].message}`);
  }

  return result.data;
}

/**
 * Валидация с fallback значением
 */
export function validateDataWithFallback<T>(
  schema: z.ZodType<T>,
  data: unknown,
  fallback: T
): T {
  const result = schema.safeParse(data);

  if (!result.success) {
    console.warn('Validation failed, using fallback', {
      errors: result.error.issues,
      fallback,
    });
    return fallback;
  }

  return result.data;
}

/**
 * Проверка является ли данные валидными без exception
 */
export function isValidData<T>(
  schema: z.ZodType<T>,
  data: unknown
): data is T {
  return schema.safeParse(data).success;
}
