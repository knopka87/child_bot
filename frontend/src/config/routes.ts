// src/config/routes.ts

/**
 * Маршруты приложения
 */
export const ROUTES = {
  HOME: '/',
  ONBOARDING: '/onboarding',
  LEGAL_PRIVACY: '/legal/privacy',
  LEGAL_TERMS: '/legal/terms',
  HELP: '/help',
  HELP_UPLOAD: '/help/upload',
  HELP_QUALITY: '/help/quality',
  HELP_PROCESSING: '/help/processing',
  HELP_RESULT: '/help/result',
  CHECK: '/check',
  CHECK_SCENARIO: '/check/scenario',
  CHECK_PROCESSING: '/check/processing',
  ACHIEVEMENTS: '/achievements',
  FRIENDS: '/friends',
  PROFILE: '/profile',
  PROFILE_HISTORY: '/profile/history',
  PROFILE_REPORT: '/profile/report',
  VILLAIN: '/villain',
  VILLAIN_VICTORY: '/villain/victory',
} as const;

export type RouteKey = keyof typeof ROUTES;
export type RoutePath = (typeof ROUTES)[RouteKey];
