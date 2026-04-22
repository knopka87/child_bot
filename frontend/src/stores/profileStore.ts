// src/stores/profileStore.ts
import { create } from 'zustand';

export interface ChildProfile {
  child_profile_id: string;
  parent_user_id: string;
  display_name: string;
  avatar_id: string;
  grade: number; // 1-11
  level: number;
  level_progress_percent: number; // 0-100
  coins_balance: number;
  tasks_solved_correct_count: number;
  wins_count: number;
  checks_correct_count: number;
  current_streak_days: number;
  has_unfinished_attempt: boolean;
  mascot_id: string;
  mascot_state: 'idle' | 'happy' | 'thinking' | 'celebrating';
  active_villain_id: string | null;
  invited_count_total: number;
  achievements_unlocked_count: number;
  created_at: string;
  updated_at: string;
}

interface ProfileState {
  profile: ChildProfile | null;
  isLoading: boolean;
  error: string | null;

  // Actions
  setProfile: (profile: ChildProfile) => void;
  updateProfile: (updates: Partial<ChildProfile>) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  clearProfile: () => void;
}

export const useProfileStore = create<ProfileState>((set) => ({
  profile: null,
  isLoading: false,
  error: null,

  setProfile: (profile) => set({ profile, error: null }),

  updateProfile: (updates) =>
    set((state) => ({
      profile: state.profile ? { ...state.profile, ...updates } : null,
    })),

  setLoading: (isLoading) => set({ isLoading }),

  setError: (error) => set({ error, isLoading: false }),

  clearProfile: () => set({ profile: null, error: null, isLoading: false }),
}));
