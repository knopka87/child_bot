// src/services/analytics/schema.ts
import type { AnalyticsEventName } from '@/types/analytics';

interface EventSchema {
  required?: string[];
  optional?: string[];
  params: Record<string, 'string' | 'number' | 'boolean'>;
}

export const ANALYTICS_SCHEMA: Partial<Record<AnalyticsEventName, EventSchema>> =
  {
    onboarding_opened: {
      required: [], // platform_type и session_id добавляются автоматически
      optional: ['entry_point'],
      params: {
        platform_type: 'string',
        session_id: 'string',
        entry_point: 'string',
      },
    },
    grade_selected: {
      required: ['grade'],
      optional: ['child_profile_id'],
      params: {
        grade: 'number',
        child_profile_id: 'string',
      },
    },
    home_opened: {
      required: ['child_profile_id'],
      optional: ['entry_point'],
      params: {
        child_profile_id: 'string',
        entry_point: 'string',
      },
    },
    achievements_opened: {
      required: ['child_profile_id'],
      params: {
        child_profile_id: 'string',
      },
    },
    friends_opened: {
      required: ['child_profile_id'],
      params: {
        child_profile_id: 'string',
      },
    },
    profile_opened: {
      required: ['child_profile_id'],
      params: {
        child_profile_id: 'string',
      },
    },
    villain_screen_opened: {
      required: ['child_profile_id', 'villain_id'],
      params: {
        child_profile_id: 'string',
        villain_id: 'string',
      },
    },
    victory_screen_opened: {
      required: ['child_profile_id', 'villain_id'],
      params: {
        child_profile_id: 'string',
        villain_id: 'string',
        attempt_id: 'string',
      },
    },
  };
