// src/api/routes.ts
/**
 * Централизованные API routes с type-safety
 * Все API endpoints должны быть определены здесь
 */

export const API_ROUTES = {
  // Analytics
  analytics: {
    events: '/api/analytics/events',
    properties: '/api/analytics/properties',
  },

  // Tasks & Homework
  tasks: {
    list: '/api/tasks',
    byId: (id: string) => `/api/tasks/${id}`,
    submit: (id: string) => `/api/tasks/${id}/submit`,
    hints: (id: string) => `/api/tasks/${id}/hints`,
    check: '/api/tasks/check',
  },

  // Villains
  villains: {
    list: '/api/villains',
    active: '/api/villains/active',
    byId: (id: string) => `/api/villains/${id}`,
    victory: (id: string) => `/api/villains/${id}/victory`,
    damage: (id: string) => `/api/villains/${id}/damage`,
  },

  // Profile
  profile: {
    get: '/api/profile',
    update: '/api/profile',
    history: '/api/profile/history',
    stats: '/api/profile/stats',
  },

  // Achievements
  achievements: {
    list: '/api/achievements',
    unlocked: '/api/achievements/unlocked',
    byId: (id: string) => `/api/achievements/${id}`,
    claim: (id: string) => `/api/achievements/${id}/claim`,
  },

  // Friends & Referrals
  friends: {
    list: '/api/friends',
    invite: '/api/friends/invite',
    referrals: '/api/friends/referrals',
    leaderboard: '/api/friends/leaderboard',
  },

  // Subscription & Payments
  subscription: {
    status: '/api/subscription/status',
    plans: '/api/subscription/plans',
    subscribe: '/api/subscription/subscribe',
    cancel: '/api/subscription/cancel',
  },

  // Support
  support: {
    send: '/api/support/message',
    history: '/api/support/history',
  },
} as const;