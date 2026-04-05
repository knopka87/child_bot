// src/types/villain.ts
export interface Villain {
  id: string;
  name: string;
  description: string;
  image_url: string; // Backend использует snake_case
  hp: number; // Текущее здоровье
  max_hp: number; // Максимальное здоровье
  level: number;
  is_active: boolean;
  is_defeated: boolean;
  unlocked_at?: string;
  defeated_at?: string;
  taunt?: string; // Опциональное поле для реплики
}

export interface VillainBattle {
  villain_id: string;
  battle_stats: BattleStats;
  recent_damage: DamageEvent[];
  next_damage_at?: string;
  can_damage_now: boolean;
}

export interface BattleStats {
  total_damage_dealt: number;
  correct_tasks_count: number;
  damage_per_task: number;
  progress_percent: number;
}

export interface DamageEvent {
  id: string;
  damage: number;
  task_type: string;
  created_at: string;
}

export interface VillainVictory {
  villain_id: string;
  villain_name: string;
  defeated_at: string;
  total_damage: number;
  tasks_completed: number;
  rewards: VillainReward[];
  next_villain?: Villain;
}

export interface VillainReward {
  type: 'sticker' | 'achievement' | 'coins' | 'avatar';
  id: string;
  name: string;
  image_url?: string;
  amount?: number;
}

export type VillainTaunt =
  | 'Ха-ха! Попробуй-ка реши задачки!'
  | 'Думаешь, справишься?'
  | 'Я непобедим!'
  | 'Ещё немного, и ты сдашься!'
  | 'Ну давай, удиви меня!';
