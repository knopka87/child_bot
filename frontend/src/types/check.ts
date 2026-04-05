// src/types/check.ts

export type CheckScenario = 'single_photo' | 'two_photo';

export type ImageRole = 'task' | 'answer';

export interface CreateCheckAttemptResponse {
  id: string;
  status: string;
}

export interface CheckAttempt {
  id: string;
  childProfileId: string;
  mode: 'check';
  scenario: CheckScenario;
  status: CheckStatus;
  taskImageUrl?: string;
  answerImageUrl?: string;
  createdAt: string;
  updatedAt: string;
}

export type CheckStatus =
  | 'created'
  | 'uploading_task'
  | 'uploading_answer'
  | 'uploaded'
  | 'processing'
  | 'long_wait'
  | 'completed'
  | 'failed';

export interface CheckResult {
  attemptId: string;
  status: 'success' | 'error';
  errors?: CheckError[];
  coinsEarned: number;
  damageDealt: number;
}

export interface CheckError {
  id: string;
  stepNumber: number | null;
  lineReference: string | null;
  description: string;
  locationType: 'step' | 'line' | 'general';
  severity: 'error' | 'warning';
}

export interface ErrorLocalization {
  type: 'step' | 'line' | 'general';
  reference: string | null;
  description: string;
}
