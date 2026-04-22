// src/types/help.ts

export type UploadSource = 'camera' | 'file' | 'clipboard' | 'dragdrop';

export interface CreateAttemptResponse {
  id: string;
  status: string;
}

export interface HelpAttempt {
  id: string;
  childProfileId: string;
  mode: 'help';
  status: HelpStatus;
  imageUrl?: string;
  thumbnailUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export type HelpStatus =
  | 'created'
  | 'uploading'
  | 'uploaded'
  | 'quality_check'
  | 'processing'
  | 'long_wait'
  | 'completed'
  | 'failed';

export interface Hint {
  id: string;
  level: 1 | 2 | 3; // Уровень подсказки
  title: string;
  content: string;
  order: number;
}

export interface HelpResult {
  attemptId: string;
  hints: Hint[];
  coinsEarned: number;
  damageDealt: number; // Урон злодею
  taskImage?: string; // Изображение задания для передачи на страницу проверки
}

export interface CropArea {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface UploadProgress {
  loaded: number;
  total: number;
  percentage: number;
}
