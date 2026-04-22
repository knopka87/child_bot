// src/types/monetization.ts
export type SubscriptionStatus = 'trial' | 'active' | 'expired' | 'cancelled';

export type SubscriptionPlanId = 'trial' | 'monthly' | 'yearly';

export interface SubscriptionPlan {
  id: SubscriptionPlanId;
  name: string;
  description: string;
  price: number;
  currency: string;
  period: 'trial' | 'month' | 'year';
  features: string[];
  isPopular?: boolean;
}

export interface Subscription {
  id: string;
  planId: SubscriptionPlanId;
  status: SubscriptionStatus;
  startedAt: string;
  expiresAt?: string;
  cancelledAt?: string;
  autoRenew: boolean;
  trialDaysRemaining?: number;
}

export interface PaymentMethod {
  type: 'vk_pay' | 'card' | 'sbp';
  title: string;
  icon: string;
  isAvailable: boolean;
}

export interface PaymentIntent {
  id: string;
  amount: number;
  currency: string;
  description: string;
  status: 'pending' | 'processing' | 'succeeded' | 'failed';
  createdAt: string;
}

export interface AdConfig {
  enabled: boolean;
  provider: 'vk_ads' | 'yandex_ads';
  placementId: string;
  frequency: number;
}

export interface InAppPurchase {
  id: string;
  type: 'consumable' | 'non_consumable';
  title: string;
  description: string;
  price: number;
  currency: string;
  icon?: string;
}
