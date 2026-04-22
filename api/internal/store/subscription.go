package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SubscriptionPlan план подписки
type SubscriptionPlan struct {
	ID              string
	Name            string
	Description     string
	PriceCents      int
	Currency        string
	DurationDays    int
	TrialDays       int
	DiscountPercent int
	IsPopular       bool
	DisplayOrder    int
	IsActive        bool
	CreatedAt       time.Time
}

// Subscription подписка пользователя
type Subscription struct {
	ID                int64
	ChildProfileID    string
	PlanID            string
	Status            string
	StartedAt         time.Time
	TrialEndsAt       *time.Time
	ExpiresAt         time.Time
	CancelledAt       *time.Time
	AutoRenew         bool
	PaymentProvider   *string
	PaymentExternalID *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// GetSubscriptionPlan получает план подписки по ID
func (s *Store) GetSubscriptionPlan(ctx context.Context, planID string) (*SubscriptionPlan, error) {
	query := `
		SELECT id, name, description, price_cents, currency, duration_days,
		       trial_days, discount_percent, is_popular, display_order, is_active, created_at
		FROM subscription_plans
		WHERE id = $1
	`

	plan := &SubscriptionPlan{}
	err := s.DB.QueryRowContext(ctx, query, planID).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.PriceCents, &plan.Currency,
		&plan.DurationDays, &plan.TrialDays, &plan.DiscountPercent,
		&plan.IsPopular, &plan.DisplayOrder, &plan.IsActive, &plan.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("query subscription plan: %w", err)
	}

	return plan, nil
}

// GetActivePlans получает список активных планов подписки
func (s *Store) GetActivePlans(ctx context.Context) ([]*SubscriptionPlan, error) {
	query := `
		SELECT id, name, description, price_cents, currency, duration_days,
		       trial_days, discount_percent, is_popular, display_order, is_active, created_at
		FROM subscription_plans
		WHERE is_active = TRUE
		ORDER BY display_order ASC, created_at ASC
	`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query active plans: %w", err)
	}
	defer rows.Close()

	var plans []*SubscriptionPlan
	for rows.Next() {
		plan := &SubscriptionPlan{}
		err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Description, &plan.PriceCents, &plan.Currency,
			&plan.DurationDays, &plan.TrialDays, &plan.DiscountPercent,
			&plan.IsPopular, &plan.DisplayOrder, &plan.IsActive, &plan.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan plan: %w", err)
		}
		plans = append(plans, plan)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plans: %w", err)
	}

	return plans, nil
}

// GetActiveSubscription получает активную подписку пользователя
func (s *Store) GetActiveSubscription(ctx context.Context, childProfileID string) (*Subscription, error) {
	query := `
		SELECT id, child_profile_id, plan_id, status, started_at, trial_ends_at,
		       expires_at, cancelled_at, auto_renew, payment_provider, payment_external_id,
		       created_at, updated_at
		FROM subscriptions
		WHERE child_profile_id = $1
		  AND status IN ('trial', 'active')
		ORDER BY created_at DESC
		LIMIT 1
	`

	sub := &Subscription{}
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(
		&sub.ID, &sub.ChildProfileID, &sub.PlanID, &sub.Status,
		&sub.StartedAt, &sub.TrialEndsAt, &sub.ExpiresAt, &sub.CancelledAt,
		&sub.AutoRenew, &sub.PaymentProvider, &sub.PaymentExternalID,
		&sub.CreatedAt, &sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}
	if err != nil {
		return nil, fmt.Errorf("query subscription: %w", err)
	}

	return sub, nil
}

// GetSubscription получает подписку по ID
func (s *Store) GetSubscription(ctx context.Context, subscriptionID int64) (*Subscription, error) {
	query := `
		SELECT id, child_profile_id, plan_id, status, started_at, trial_ends_at,
		       expires_at, cancelled_at, auto_renew, payment_provider, payment_external_id,
		       created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	sub := &Subscription{}
	err := s.DB.QueryRowContext(ctx, query, subscriptionID).Scan(
		&sub.ID, &sub.ChildProfileID, &sub.PlanID, &sub.Status,
		&sub.StartedAt, &sub.TrialEndsAt, &sub.ExpiresAt, &sub.CancelledAt,
		&sub.AutoRenew, &sub.PaymentProvider, &sub.PaymentExternalID,
		&sub.CreatedAt, &sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}
	if err != nil {
		return nil, fmt.Errorf("query subscription: %w", err)
	}

	return sub, nil
}

// CancelSubscription отменяет подписку (но оставляет доступ до конца оплаченного периода)
func (s *Store) CancelSubscription(ctx context.Context, childProfileID string) error {
	query := `
		UPDATE subscriptions
		SET auto_renew = FALSE,
		    cancelled_at = NOW(),
		    updated_at = NOW()
		WHERE child_profile_id = $1
		  AND status IN ('trial', 'active')
		  AND cancelled_at IS NULL
	`

	result, err := s.DB.ExecContext(ctx, query, childProfileID)
	if err != nil {
		return fmt.Errorf("cancel subscription: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no active subscription to cancel")
	}

	return nil
}

// ResumeSubscription возобновляет отмененную подписку
func (s *Store) ResumeSubscription(ctx context.Context, childProfileID string) error {
	query := `
		UPDATE subscriptions
		SET auto_renew = TRUE,
		    cancelled_at = NULL,
		    updated_at = NOW()
		WHERE child_profile_id = $1
		  AND status IN ('trial', 'active')
		  AND cancelled_at IS NOT NULL
	`

	result, err := s.DB.ExecContext(ctx, query, childProfileID)
	if err != nil {
		return fmt.Errorf("resume subscription: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no cancelled subscription to resume")
	}

	return nil
}

// ExpireSubscriptions переводит истекшие подписки в статус expired
func (s *Store) ExpireSubscriptions(ctx context.Context) (int64, error) {
	query := `
		UPDATE subscriptions
		SET status = 'expired',
		    updated_at = NOW()
		WHERE status IN ('trial', 'active')
		  AND expires_at < NOW()
	`

	result, err := s.DB.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("expire subscriptions: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("get rows affected: %w", err)
	}

	return rows, nil
}
