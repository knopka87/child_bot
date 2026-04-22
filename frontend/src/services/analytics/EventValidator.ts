// src/services/analytics/EventValidator.ts
import type { AnalyticsEventName } from '@/types/analytics';
import { ANALYTICS_SCHEMA } from './schema';

interface ValidationResult {
  isValid: boolean;
  errors: string[];
}

export class EventValidator {
  validate(
    name: AnalyticsEventName,
    params: Record<string, any>
  ): ValidationResult {
    const schema = ANALYTICS_SCHEMA[name];

    if (!schema) {
      return {
        isValid: true,
        errors: [],
      };
    }

    const errors: string[] = [];

    schema.required?.forEach((param) => {
      if (!(param in params) || params[param] === undefined || params[param] === null) {
        errors.push(`Missing required parameter: ${param}`);
      }
    });

    Object.entries(params).forEach(([key, value]) => {
      const expectedType = schema.params[key];
      if (!expectedType) return;

      const actualType = typeof value;
      if (actualType !== expectedType) {
        errors.push(
          `Invalid type for ${key}: expected ${expectedType}, got ${actualType}`
        );
      }
    });

    return {
      isValid: errors.length === 0,
      errors,
    };
  }
}
