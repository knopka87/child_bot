// src/api/monetization.ts
import { apiClient } from './client';
import type {
  SubscriptionPlan,
  Subscription,
  PaymentMethod,
  PaymentIntent,
  InAppPurchase,
} from '@/types/monetization';

export const monetizationAPI = {
  /**
   * Получить доступные планы подписки
   */
  async getPlans(): Promise<SubscriptionPlan[]> {
    return apiClient.get<SubscriptionPlan[]>('/subscriptions/plans');
  },

  /**
   * Получить текущую подписку пользователя
   */
  async getCurrentSubscription(userId: string): Promise<Subscription | null> {
    return apiClient.get<Subscription | null>(
      `/subscriptions/users/${userId}/current`
    );
  },

  /**
   * Получить доступные методы оплаты
   */
  async getPaymentMethods(): Promise<PaymentMethod[]> {
    return apiClient.get<PaymentMethod[]>('/payments/methods');
  },

  /**
   * Создать платёжный интент
   */
  async createPaymentIntent(
    planId: string,
    method: string
  ): Promise<PaymentIntent> {
    return apiClient.post<PaymentIntent>('/payments/intents', {
      planId,
      method,
    });
  },

  /**
   * Проверить статус платежа
   */
  async checkPaymentStatus(intentId: string): Promise<PaymentIntent> {
    return apiClient.get<PaymentIntent>(`/payments/intents/${intentId}`);
  },

  /**
   * Отменить подписку
   */
  async cancelSubscription(
    userId: string,
    subscriptionId: string
  ): Promise<void> {
    return apiClient.post<void>(
      `/subscriptions/users/${userId}/cancel`,
      { subscriptionId }
    );
  },

  /**
   * Возобновить подписку
   */
  async resumeSubscription(
    userId: string,
    subscriptionId: string
  ): Promise<void> {
    return apiClient.post<void>(
      `/subscriptions/users/${userId}/resume`,
      { subscriptionId }
    );
  },

  /**
   * Получить список доступных покупок
   */
  async getInAppPurchases(): Promise<InAppPurchase[]> {
    return apiClient.get<InAppPurchase[]>('/purchases/products');
  },

  /**
   * Совершить внутриигровую покупку
   */
  async purchaseItem(
    userId: string,
    itemId: string
  ): Promise<PaymentIntent> {
    return apiClient.post<PaymentIntent>(`/purchases/users/${userId}/buy`, {
      itemId,
    });
  },
};
