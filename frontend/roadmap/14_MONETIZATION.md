# Phase 13: Monetization

**Длительность:** 4-5 дней
**Приоритет:** Высокий
**Зависимости:** 09_PROFILE.md (Paywall)

---

## Цель

Интегрировать VK Pay для подписок, Ads SDK для монетизации через рекламу, и настроить backend webhooks для обработки платежей.

---

## Часть 1: VK Pay Integration

### 1.1. Frontend: Инициация платежа

**Файл:** `src/api/payments.ts`

```typescript
import bridge from '@vkontakte/vk-bridge';
import { apiClient } from './client';
import { haptics } from '@/lib/platform/haptics';

export interface VKPayOrderRequest {
  billing_plan_id: string;
}

export interface VKPayOrder {
  order_id: string;
  merchant_id: string;
  amount: number; // в копейках
  description: string;
  sign: string; // подпись с backend
}

export interface VKPayResult {
  status: 'success' | 'cancel' | 'fail';
  transaction_id?: string;
  amount?: number;
}

/**
 * Создание заказа на backend
 */
export async function createVKPayOrder(
  planId: string
): Promise<VKPayOrder> {
  return apiClient.post<VKPayOrder>('/subscription/vk-pay/order', {
    billing_plan_id: planId,
  });
}

/**
 * Открытие VK Pay формы
 */
export async function openVKPay(order: VKPayOrder): Promise<VKPayResult> {
  try {
    haptics.impact('medium');

    const result = await bridge.send('VKWebAppOpenPayForm', {
      app_id: Number(import.meta.env.VITE_VK_APP_ID),
      action: 'pay-to-service',
      params: {
        merchant_id: order.merchant_id,
        amount: order.amount,
        description: order.description,
        order_id: order.order_id,
        sign: order.sign,
      },
    });

    if (result.status === 'success') {
      haptics.notification('success');
    }

    return result as VKPayResult;
  } catch (error) {
    console.error('[VKPay] Payment failed:', error);
    haptics.notification('error');
    throw error;
  }
}

/**
 * Проверка статуса заказа (polling)
 */
export async function checkOrderStatus(
  orderId: string,
  maxAttempts = 10,
  intervalMs = 2000
): Promise<'paid' | 'pending' | 'failed'> {
  for (let i = 0; i < maxAttempts; i++) {
    const status = await apiClient.get<{ status: string }>(
      `/subscription/vk-pay/order/${orderId}/status`
    );

    if (status.status === 'paid') {
      return 'paid';
    }

    if (status.status === 'failed') {
      return 'failed';
    }

    await new Promise((resolve) => setTimeout(resolve, intervalMs));
  }

  return 'pending';
}
```

---

### 1.2. Компонент выбора подписки

**Файл:** `src/pages/Profile/PaywallPage.tsx`

```typescript
import { useState } from 'react';
import { Group, Card, Button, Div, Title, Text, Snackbar, Spinner } from '@vkontakte/vkui';
import { Icon24Done } from '@vkontakte/icons';
import {
  createVKPayOrder,
  openVKPay,
  checkOrderStatus,
} from '@/api/payments';
import { analytics } from '@/lib/analytics/tracker';

interface BillingPlan {
  id: string;
  name: string;
  price: number; // в рублях
  duration_days: number;
  features: string[];
}

const BILLING_PLANS: BillingPlan[] = [
  {
    id: 'basic_monthly',
    name: 'Базовый',
    price: 199,
    duration_days: 30,
    features: ['Безлимитные подсказки', 'Проверка работ', 'Без рекламы'],
  },
  {
    id: 'premium_monthly',
    name: 'Премиум',
    price: 349,
    duration_days: 30,
    features: [
      'Все из Базового',
      'Приоритетная обработка',
      'Эксклюзивные стикеры',
      'VIP поддержка',
    ],
  },
  {
    id: 'premium_yearly',
    name: 'Премиум (год)',
    price: 2990,
    duration_days: 365,
    features: [
      'Все из Премиум',
      'Скидка 30%',
      'Бонусные монеты',
      'Ранний доступ к функциям',
    ],
  },
];

export function PaywallPage() {
  const [selectedPlan, setSelectedPlan] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(false);
  const [snackbar, setSnackbar] = useState<React.ReactNode>(null);

  const handlePurchase = async (planId: string) => {
    setSelectedPlan(planId);
    setIsProcessing(true);

    try {
      // Отправляем событие начала покупки
      analytics.sendEvent({
        event_name: 'payment_initiated',
        billing_plan_id: planId,
        payment_method: 'vk_pay',
      });

      // 1. Создаем заказ на backend
      const order = await createVKPayOrder(planId);

      // 2. Открываем VK Pay
      const result = await openVKPay(order);

      if (result.status === 'success') {
        // 3. Ждем подтверждения через webhook (polling)
        const status = await checkOrderStatus(order.order_id);

        if (status === 'paid') {
          // Успех!
          analytics.sendEvent({
            event_name: 'payment_success',
            order_id: order.order_id,
            amount: order.amount / 100, // в рублях
            billing_plan_id: planId,
          });

          setSnackbar(
            <Snackbar
              onClose={() => setSnackbar(null)}
              before={<Icon24Done fill="var(--vkui--color_icon_positive)" />}
            >
              Подписка активирована!
            </Snackbar>
          );

          // Обновляем профиль
          setTimeout(() => {
            window.location.reload();
          }, 2000);
        } else {
          throw new Error('Payment verification failed');
        }
      } else if (result.status === 'cancel') {
        analytics.sendEvent({
          event_name: 'payment_cancelled',
          billing_plan_id: planId,
        });
      }
    } catch (error) {
      console.error('[Paywall] Purchase failed:', error);

      analytics.sendEvent({
        event_name: 'payment_failed',
        billing_plan_id: planId,
        error_message: (error as Error).message,
      });

      setSnackbar(
        <Snackbar onClose={() => setSnackbar(null)}>
          Ошибка оплаты. Попробуйте снова.
        </Snackbar>
      );
    } finally {
      setIsProcessing(false);
      setSelectedPlan(null);
    }
  };

  return (
    <Group header={<Title level="2">Выберите подписку</Title>}>
      {BILLING_PLANS.map((plan) => (
        <Card key={plan.id} mode="outline" style={{ margin: 12 }}>
          <Div>
            <Title level="3" weight="2">
              {plan.name}
            </Title>
            <Text weight="1" style={{ fontSize: 32, marginTop: 8 }}>
              {plan.price} ₽
            </Text>
            <Text style={{ color: 'var(--vkui--color_text_secondary)' }}>
              на {plan.duration_days} дней
            </Text>

            <div style={{ marginTop: 16 }}>
              {plan.features.map((feature, i) => (
                <Text key={i} style={{ marginTop: 4 }}>
                  ✓ {feature}
                </Text>
              ))}
            </div>

            <Button
              size="l"
              stretched
              mode="primary"
              onClick={() => handlePurchase(plan.id)}
              loading={isProcessing && selectedPlan === plan.id}
              disabled={isProcessing}
              style={{ marginTop: 16 }}
            >
              {isProcessing && selectedPlan === plan.id ? (
                <Spinner size="small" />
              ) : (
                'Оформить подписку'
              )}
            </Button>
          </Div>
        </Card>
      ))}

      {snackbar}
    </Group>
  );
}
```

---

### 1.3. Backend: VK Pay Webhook Handler (Go)

**Файл:** `internal/api/v1/payments/vk_pay_webhook.go`

```go
package payments

import (
    "crypto/hmac"
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "sort"
    "strings"
)

type VKPayNotification struct {
    NotificationType string `json:"notification_type"`
    AppID            int    `json:"app_id"`
    UserID           int    `json:"user_id"`
    ReceiverID       int    `json:"receiver_id"`
    OrderID          string `json:"order_id"`
    Date             int64  `json:"date"`
    Status           string `json:"status"` // "chargeable"
    Item             string `json:"item"`
    ItemID           string `json:"item_id"`
    ItemTitle        string `json:"item_title"`
    ItemPhotoURL     string `json:"item_photo_url"`
    ItemPrice        int    `json:"item_price"`
    Sig              string `json:"sig"`
}

// HandleVKPayWebhook обрабатывает webhook от VK Pay
func (h *PaymentHandler) HandleVKPayWebhook(w http.ResponseWriter, r *http.Request) {
    // 1. Читаем тело запроса
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    var notification VKPayNotification
    if err := json.Unmarshal(body, &notification); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // 2. Проверяем подпись
    if !h.validateVKPaySignature(notification) {
        log.Warn("Invalid VK Pay signature", "order_id", notification.OrderID)
        http.Error(w, "Invalid signature", http.StatusForbidden)
        return
    }

    // 3. Обрабатываем уведомление
    switch notification.NotificationType {
    case "order_status_change", "order_status_change_test":
        if notification.Status == "chargeable" {
            // Активируем подписку
            err := h.activateSubscription(r.Context(), notification)
            if err != nil {
                log.Error("Failed to activate subscription", "error", err)
                http.Error(w, "Internal error", http.StatusInternalServerError)
                return
            }

            // Отправляем аналитику
            h.analytics.SendEvent("payment_webhook_received", map[string]interface{}{
                "order_id":   notification.OrderID,
                "user_id":    notification.UserID,
                "amount":     notification.ItemPrice,
                "item_id":    notification.ItemID,
            })
        }
    }

    // 4. Важно: возвращаем "ok"
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("ok"))
}

// validateVKPaySignature проверяет подпись VK Pay
func (h *PaymentHandler) validateVKPaySignature(n VKPayNotification) bool {
    // Собираем параметры для подписи
    params := []string{
        fmt.Sprintf("notification_type=%s", n.NotificationType),
        fmt.Sprintf("app_id=%d", n.AppID),
        fmt.Sprintf("user_id=%d", n.UserID),
        fmt.Sprintf("receiver_id=%d", n.ReceiverID),
        fmt.Sprintf("order_id=%s", n.OrderID),
        fmt.Sprintf("date=%d", n.Date),
    }

    // Сортируем
    sort.Strings(params)
    queryString := strings.Join(params, "&")

    // MD5 hash
    hash := md5.New()
    hash.Write([]byte(queryString + h.vkPaySecret))
    expectedSig := hex.EncodeToString(hash.Sum(nil))

    return n.Sig == expectedSig
}

// activateSubscription активирует подписку для пользователя
func (h *PaymentHandler) activateSubscription(ctx context.Context, n VKPayNotification) error {
    // Находим заказ в БД
    order, err := h.orderRepo.GetByOrderID(ctx, n.OrderID)
    if err != nil {
        return fmt.Errorf("order not found: %w", err)
    }

    // Проверяем, не активирован ли уже
    if order.Status == "paid" {
        log.Info("Order already paid", "order_id", n.OrderID)
        return nil
    }

    // Обновляем статус заказа
    order.Status = "paid"
    order.TransactionID = fmt.Sprintf("%d_%s", n.Date, n.OrderID)
    if err := h.orderRepo.Update(ctx, order); err != nil {
        return fmt.Errorf("failed to update order: %w", err)
    }

    // Активируем подписку
    subscription := &domain.Subscription{
        ChildProfileID:  order.ChildProfileID,
        BillingPlanID:   order.BillingPlanID,
        Status:          "active",
        StartedAt:       time.Now(),
        ExpiresAt:       time.Now().AddDate(0, 0, order.DurationDays),
    }

    if err := h.subscriptionRepo.Create(ctx, subscription); err != nil {
        return fmt.Errorf("failed to create subscription: %w", err)
    }

    log.Info("Subscription activated", "order_id", n.OrderID, "profile_id", order.ChildProfileID)
    return nil
}
```

---

## Часть 2: VK Ads SDK Integration

### 2.1. Frontend: Ads Service

**Файл:** `src/lib/ads/vk-ads.ts`

```typescript
import bridge from '@vkontakte/vk-bridge';
import { analytics } from '@/lib/analytics/tracker';

export type AdFormat = 'reward' | 'interstitial';

export interface RewardedAdResult {
  watched: boolean;
  reward?: {
    type: 'extra_hints' | 'coins';
    amount: number;
  };
}

/**
 * Проверка доступности рекламы
 */
async function checkAdAvailable(format: AdFormat): Promise<boolean> {
  try {
    const result = await bridge.send('VKWebAppCheckNativeAds', {
      ad_format: format,
    });
    return result.result === true;
  } catch (error) {
    console.warn('[Ads] Check failed:', error);
    return false;
  }
}

/**
 * Показать rewarded video (за награду)
 */
export async function showRewardedAd(): Promise<RewardedAdResult> {
  try {
    analytics.sendEvent({
      event_name: 'ad_rewarded_initiated',
    });

    const isAvailable = await checkAdAvailable('reward');
    if (!isAvailable) {
      analytics.sendEvent({
        event_name: 'ad_rewarded_unavailable',
      });
      return { watched: false };
    }

    await bridge.send('VKWebAppShowNativeAds', { ad_format: 'reward' });

    analytics.sendEvent({
      event_name: 'ad_rewarded_watched',
    });

    return {
      watched: true,
      reward: {
        type: 'extra_hints',
        amount: 1,
      },
    };
  } catch (error) {
    console.error('[Ads] Rewarded ad failed:', error);

    analytics.sendEvent({
      event_name: 'ad_rewarded_failed',
      error_message: (error as Error).message,
    });

    return { watched: false };
  }
}

/**
 * Показать interstitial (полноэкранная реклама)
 */
export async function showInterstitialAd(): Promise<boolean> {
  try {
    analytics.sendEvent({
      event_name: 'ad_interstitial_initiated',
    });

    const isAvailable = await checkAdAvailable('interstitial');
    if (!isAvailable) {
      analytics.sendEvent({
        event_name: 'ad_interstitial_unavailable',
      });
      return false;
    }

    await bridge.send('VKWebAppShowNativeAds', { ad_format: 'interstitial' });

    analytics.sendEvent({
      event_name: 'ad_interstitial_shown',
    });

    return true;
  } catch (error) {
    console.error('[Ads] Interstitial ad failed:', error);

    analytics.sendEvent({
      event_name: 'ad_interstitial_failed',
      error_message: (error as Error).message,
    });

    return false;
  }
}

/**
 * Ads Strategy Manager
 */
export class AdsStrategy {
  private attemptCount = 0;
  private interstitialFrequency = 5; // каждые 5 попыток

  /**
   * Увеличить счетчик попыток и показать рекламу если нужно
   */
  async onAttemptComplete(): Promise<void> {
    this.attemptCount++;

    if (this.attemptCount % this.interstitialFrequency === 0) {
      await showInterstitialAd();
    }
  }

  /**
   * Сбросить счетчик (при покупке подписки)
   */
  reset(): void {
    this.attemptCount = 0;
  }
}

export const adsStrategy = new AdsStrategy();
```

---

### 2.2. Компонент для разблокировки подсказки через рекламу

**Файл:** `src/components/features/hints/UnlockHintWithAd.tsx`

```typescript
import { Button, Snackbar, Div, Text } from '@vkontakte/vkui';
import { Icon24VideoOutline } from '@vkontakte/icons';
import { useState } from 'react';
import { showRewardedAd } from '@/lib/ads/vk-ads';
import { haptics } from '@/lib/platform/haptics';

interface UnlockHintWithAdProps {
  onUnlocked: () => void;
}

export function UnlockHintWithAd({ onUnlocked }: UnlockHintWithAdProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [snackbar, setSnackbar] = useState<React.ReactNode>(null);

  const handleWatchAd = async () => {
    setIsLoading(true);

    try {
      const result = await showRewardedAd();

      if (result.watched) {
        haptics.notification('success');

        setSnackbar(
          <Snackbar onClose={() => setSnackbar(null)}>
            Подсказка разблокирована!
          </Snackbar>
        );

        onUnlocked();
      } else {
        setSnackbar(
          <Snackbar onClose={() => setSnackbar(null)}>
            Реклама недоступна. Попробуйте позже.
          </Snackbar>
        );
      }
    } catch (error) {
      setSnackbar(
        <Snackbar onClose={() => setSnackbar(null)}>
          Ошибка при загрузке рекламы
        </Snackbar>
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Div>
      <Text style={{ textAlign: 'center', marginBottom: 12 }}>
        Подсказка уровня 3 заблокирована
      </Text>

      <Button
        size="l"
        mode="primary"
        stretched
        before={<Icon24VideoOutline />}
        onClick={handleWatchAd}
        loading={isLoading}
      >
        Посмотреть рекламу
      </Button>

      <Text
        style={{
          textAlign: 'center',
          marginTop: 8,
          color: 'var(--vkui--color_text_secondary)',
        }}
      >
        После просмотра вы получите доступ к подсказке
      </Text>

      {snackbar}
    </Div>
  );
}
```

---

## Часть 3: Subscription Management

### 3.1. API для проверки подписки

**Файл:** `src/api/subscription.ts`

```typescript
import { apiClient } from './client';

export interface Subscription {
  id: string;
  billing_plan_id: string;
  status: 'active' | 'expired' | 'cancelled';
  started_at: string;
  expires_at: string;
  auto_renew: boolean;
}

export interface BillingPlan {
  id: string;
  name: string;
  price: number;
  duration_days: number;
  features: string[];
}

/**
 * Получить текущую подписку
 */
export async function getCurrentSubscription(): Promise<Subscription | null> {
  try {
    return await apiClient.get<Subscription>('/subscription/current');
  } catch (error) {
    if ((error as any).response?.status === 404) {
      return null;
    }
    throw error;
  }
}

/**
 * Получить доступные планы
 */
export async function getBillingPlans(): Promise<BillingPlan[]> {
  return apiClient.get<BillingPlan[]>('/subscription/plans');
}

/**
 * Отменить подписку
 */
export async function cancelSubscription(subscriptionId: string): Promise<void> {
  return apiClient.post(`/subscription/${subscriptionId}/cancel`);
}
```

---

## Чеклист задач

### VK Pay
- [ ] Создать API для создания VK Pay заказа
- [ ] Реализовать openVKPay функцию
- [ ] Реализовать polling для проверки статуса
- [ ] Создать PaywallPage компонент
- [ ] Настроить backend webhook handler (Go)
- [ ] Протестировать VK Pay flow (песочница)
- [ ] Добавить аналитику для платежей

### VK Ads SDK
- [ ] Реализовать showRewardedAd функцию
- [ ] Реализовать showInterstitialAd функцию
- [ ] Создать AdsStrategy для управления показом
- [ ] Создать UnlockHintWithAd компонент
- [ ] Интегрировать interstitial ads после попыток
- [ ] Добавить аналитику для рекламы
- [ ] Протестировать ads на реальном устройстве

### Subscription Management
- [ ] Создать API для текущей подписки
- [ ] Создать UI для отмены подписки
- [ ] Реализовать проверку подписки при входе
- [ ] Добавить badge "PRO" для подписчиков
- [ ] Скрыть рекламу для подписчиков
- [ ] Протестировать весь subscription flow

### Backend (Go)
- [ ] Реализовать VK Pay webhook handler
- [ ] Валидация подписи VK Pay
- [ ] Создание/обновление subscription в БД
- [ ] Проверка статуса подписки
- [ ] Отправка webhook событий в аналитику

---

## Следующий этап

После завершения Monetization переходи к **15_VK_MODERATION.md** для подготовки к модерации VK.
