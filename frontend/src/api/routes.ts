// src/api/routes.ts
/**
 * Централизованные API routes с type-safety
 * Все API endpoints должны быть определены здесь
 */

export const API_ROUTES = {
  // Analytics
  analytics: {
    events: '/analytics/events',
    properties: '/analytics/properties',
  },

  // Tasks & Homework
  tasks: {
    list: '/tasks',
    byId: (id: string) => `/tasks/${id}`,
    submit: (id: string) => `/tasks/${id}/submit`,
    hints: (id: string) => `/tasks/${id}/hints`,
    check: '/tasks/check',
  },

  // Villains
  villains: {
    list: '/villains',
    active: '/villains/active',
    byId: (id: string) => `/villains/${id}`,
    victory: (id: string) => `/villains/${id}/victory`,
    damage: (id: string) => `/villains/${id}/damage`,
  },

  // Profile
  profile: {
    get: '/profile',
    update: '/profile',
    history: '/profile/history',
    stats: '/profile/stats',
  },

  // Achievements
  achievements: {
    list: '/achievements',
    unlocked: '/achievements/unlocked',
    byId: (id: string) => `/achievements/${id}`,
    claim: (id: string) => `/achievements/${id}/claim`,
  },

  // Friends & Referrals
  friends: {
    list: '/friends',
    invite: '/friends/invite',
    referrals: '/friends/referrals',
    leaderboard: '/friends/leaderboard',
  },

  // Subscription & Payments
  subscription: {
    status: '/subscription/status',
    plans: '/subscription/plans',
    subscribe: '/subscription/subscribe',
    cancel: '/subscription/cancel',
  },

  // Support
  support: {
    send: '/support/message',
    history: '/support/history',
  },
} as const;