// src/types/home.ts
import type { Attempt, Villain, MascotState } from './domain';

export interface HomeData {
  profile: {
    id: string;
    displayName: string;
    level: number;
    xpTotal: number;
    xpForNextLevel: number;
    levelProgress: number; // 0-100 процентов
    coinsBalance: number;
    tasksSolvedCorrectCount: number;
  };
  mascot: {
    id: string;
    state: MascotState;
    imageUrl: string;
    message: string;
  };
  villain: Villain | null;
  unfinishedAttempt: Attempt | null;
  recentAttempts: RecentAttempt[];
  achievements: {
    unlockedCount: number;
    totalCount: number;
  };
}

export interface RecentAttempt {
  id: string;
  mode: 'help' | 'check';
  status: 'success' | 'error' | 'in_progress';
  createdAt: string;
  thumbnail?: string;
  resultSummary?: string;
}
