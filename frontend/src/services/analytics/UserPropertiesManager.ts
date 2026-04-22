// src/services/analytics/UserPropertiesManager.ts
import type { UserProperties } from '@/types/analytics';

export class UserPropertiesManager {
  private properties: UserProperties = {};

  set(properties: UserProperties): void {
    this.properties = { ...properties };
  }

  update(properties: Partial<UserProperties>): void {
    this.properties = { ...this.properties, ...properties };
  }

  get(key: keyof UserProperties): any {
    return this.properties[key];
  }

  getAll(): UserProperties {
    return { ...this.properties };
  }

  clear(): void {
    this.properties = {};
  }
}
