// src/api/onboarding.ts
import { apiClient } from './client';
import type { Avatar } from '@/types/onboarding';

export const onboardingAPI = {
  /**
   * Получить список аватаров
   */
  async getAvatars(): Promise<Avatar[]> {
    return apiClient.get<Avatar[]>('/avatars');
  },

  /**
   * Создать профиль ребёнка
   */
  async createChildProfile(data: {
    parentUserId: string;
    grade: number;
    avatarId: string;
    displayName: string;
    referralCode?: string;
  }): Promise<{ childProfileId: string }> {
    return apiClient.post<{ childProfileId: string }>('/profiles/child', data);
  },

  /**
   * Отправить email для верификации
   */
  async sendEmailVerification(data: {
    email: string;
    parentUserId: string;
  }): Promise<{ message: string; expiresAt: string; devCode?: string }> {
    return apiClient.post<{ message: string; expiresAt: string; devCode?: string }>(
      '/email/verify/send',
      data
    );
  },

  /**
   * Проверить код верификации email
   */
  async verifyEmailCode(data: {
    email: string;
    code: string;
  }): Promise<{ verified: boolean; message: string }> {
    return apiClient.post<{ verified: boolean; message: string }>(
      '/email/verify/check',
      data
    );
  },

  /**
   * Проверить статус верификации email
   */
  async checkEmailVerification(email: string): Promise<{ email: string; verified: boolean }> {
    return apiClient.get<{ email: string; verified: boolean }>(
      `/email/verify/status?email=${encodeURIComponent(email)}`
    );
  },

  /**
   * Сохранить согласие на обработку данных
   */
  async saveConsent(data: {
    parentUserId: string;
    privacyPolicyVersion: string;
    termsVersion: string;
    adultConsent: boolean;
  }): Promise<void> {
    return apiClient.post<void>('/consent', data);
  },

  /**
   * Завершить онбординг
   */
  async completeOnboarding(data: {
    parentUserId: string;
    childProfileId: string;
  }): Promise<void> {
    return apiClient.post<void>('/onboarding/complete', data);
  },
};
