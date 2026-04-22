// src/types/domain.ts

export type MascotState = 'idle' | 'happy' | 'thinking' | 'celebrating' | 'encouraging';

export interface Villain {
  id: string;
  name: string;
  imageUrl: string;
  healthPercent: number; // 0-100
  isDefeated: boolean;
}

export interface Attempt {
  id: string;
  mode: 'help' | 'check';
  status: 'in_progress' | 'completed' | 'failed';
  createdAt: string;
  updatedAt: string;
  taskText?: string;
  thumbnail?: string;
}
