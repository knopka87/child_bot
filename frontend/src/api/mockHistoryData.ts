// src/api/mockHistoryData.ts
// Тестовые данные для истории попыток

import type { HistoryAttempt } from '@/types/profile';

export const MOCK_HISTORY_DATA: HistoryAttempt[] = [
  {
    id: '1',
    mode: 'help',
    status: 'success',
    scenarioType: 'single_photo',
    createdAt: '2026-04-03T14:30:00Z',
    completedAt: '2026-04-03T14:35:00Z',
    images: [
      {
        id: 'img1',
        role: 'single',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
    ],
    result: {
      status: 'correct',
      errorCount: 0,
      summary: 'Задача решена правильно!',
    },
    hintsUsed: 0,
  },
  {
    id: '2',
    mode: 'check',
    status: 'error',
    scenarioType: 'two_photo',
    createdAt: '2026-04-02T16:20:00Z',
    completedAt: '2026-04-02T16:25:00Z',
    images: [
      {
        id: 'img2a',
        role: 'task',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
      {
        id: 'img2b',
        role: 'answer',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
    ],
    result: {
      status: 'has_errors',
      errorCount: 2,
      feedback: [
        {
          id: 'err1',
          stepNumber: 2,
          lineReference: 'Строка 3',
          description: 'Неправильное вычисление: 5 + 3 = 7 (должно быть 8)',
          locationType: 'line',
        },
        {
          id: 'err2',
          stepNumber: 4,
          description: 'Пропущен шаг умножения',
          locationType: 'step',
        },
      ],
      summary: 'Найдено 2 ошибки в решении',
    },
    hintsUsed: 0,
  },
  {
    id: '3',
    mode: 'help',
    status: 'in_progress',
    scenarioType: 'single_photo',
    createdAt: '2026-04-02T10:15:00Z',
    images: [
      {
        id: 'img3',
        role: 'single',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
    ],
    result: {
      status: 'processing',
    },
    hintsUsed: 0,
  },
  {
    id: '4',
    mode: 'help',
    status: 'success',
    scenarioType: 'single_photo',
    createdAt: '2026-04-01T13:45:00Z',
    completedAt: '2026-04-01T13:52:00Z',
    images: [
      {
        id: 'img4',
        role: 'single',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
    ],
    result: {
      status: 'correct',
      errorCount: 0,
      summary: 'Отлично! Все правильно',
    },
    hintsUsed: 2,
  },
  {
    id: '5',
    mode: 'check',
    status: 'in_progress',
    scenarioType: 'single_photo',
    createdAt: '2026-03-31T11:30:00Z',
    images: [
      {
        id: 'img5',
        role: 'single',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
    ],
    result: {
      status: 'processing',
    },
    hintsUsed: 0,
  },
  {
    id: '6',
    mode: 'help',
    status: 'success',
    scenarioType: 'single_photo',
    createdAt: '2026-03-30T15:20:00Z',
    completedAt: '2026-03-30T15:28:00Z',
    images: [
      {
        id: 'img6',
        role: 'single',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
    ],
    result: {
      status: 'correct',
      errorCount: 0,
      summary: 'Правильно решено',
    },
    hintsUsed: 1,
  },
  {
    id: '7',
    mode: 'check',
    status: 'error',
    scenarioType: 'two_photo',
    createdAt: '2026-03-29T14:10:00Z',
    completedAt: '2026-03-29T14:15:00Z',
    images: [
      {
        id: 'img7a',
        role: 'task',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
      {
        id: 'img7b',
        role: 'answer',
        url: 'https://placehold.co/400x300',
        thumbnailUrl: 'https://placehold.co/200x150',
      },
    ],
    result: {
      status: 'has_errors',
      errorCount: 1,
      feedback: [
        {
          id: 'err3',
          stepNumber: 1,
          description: 'Ошибка в первом действии',
          locationType: 'step',
        },
      ],
      summary: 'Найдена 1 ошибка',
    },
    hintsUsed: 0,
  },
];

// Названия для карточек (mock)
export const MOCK_TASK_TITLES: Record<string, string> = {
  '1': 'Математика — задача про яблоки',
  '2': 'Русский — упражнение 45',
  '3': 'Математика — примеры',
  '4': 'Окружающий мир — вопросы',
  '5': 'Математика — задание 12',
  '6': 'Литература — чтение',
  '7': 'Математика — уравнения',
};
