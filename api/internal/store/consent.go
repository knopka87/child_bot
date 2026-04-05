package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// ParentConsent представляет согласие родителя
type ParentConsent struct {
	ID                       string
	ParentUserID             string
	PlatformID               string
	PrivacyPolicyVersion     string
	PrivacyPolicyAccepted    bool
	PrivacyPolicyAcceptedAt  time.Time
	TermsVersion             string
	TermsAccepted            bool
	TermsAcceptedAt          time.Time
	AdultConsent             bool
	AdultConsentAt           sql.NullTime
	IPAddress                sql.NullString
	UserAgent                sql.NullString
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

// SaveParentConsent сохраняет или обновляет согласие родителя (UPSERT)
func (s *Store) SaveParentConsent(
	ctx context.Context,
	platformID string,
	parentUserID string,
	privacyPolicyVersion string,
	termsVersion string,
	adultConsent bool,
	ipAddress, userAgent string,
) error {
	query := `
		INSERT INTO parent_consents (
			platform_id,
			parent_user_id,
			privacy_policy_version,
			privacy_policy_accepted,
			privacy_policy_accepted_at,
			terms_version,
			terms_accepted,
			terms_accepted_at,
			adult_consent,
			adult_consent_at,
			ip_address,
			user_agent
		) VALUES ($1, $2, $3, TRUE, NOW(), $4, TRUE, NOW(), $5, CASE WHEN $5 THEN NOW() ELSE NULL END, $6, $7)
		ON CONFLICT (platform_id, parent_user_id)
		DO UPDATE SET
			privacy_policy_version = EXCLUDED.privacy_policy_version,
			privacy_policy_accepted = TRUE,
			privacy_policy_accepted_at = NOW(),
			terms_version = EXCLUDED.terms_version,
			terms_accepted = TRUE,
			terms_accepted_at = NOW(),
			adult_consent = EXCLUDED.adult_consent,
			adult_consent_at = CASE WHEN EXCLUDED.adult_consent THEN NOW() ELSE parent_consents.adult_consent_at END,
			ip_address = EXCLUDED.ip_address,
			user_agent = EXCLUDED.user_agent,
			updated_at = NOW()
	`

	_, err := s.DB.ExecContext(
		ctx,
		query,
		platformID,
		parentUserID,
		privacyPolicyVersion,
		termsVersion,
		adultConsent,
		sql.NullString{String: ipAddress, Valid: ipAddress != ""},
		sql.NullString{String: userAgent, Valid: userAgent != ""},
	)
	if err != nil {
		return fmt.Errorf("save parent consent: %w", err)
	}

	return nil
}

// GetParentConsent получает согласие родителя
func (s *Store) GetParentConsent(ctx context.Context, platformID, parentUserID string) (*ParentConsent, error) {
	query := `
		SELECT
			id,
			parent_user_id,
			platform_id,
			privacy_policy_version,
			privacy_policy_accepted,
			privacy_policy_accepted_at,
			terms_version,
			terms_accepted,
			terms_accepted_at,
			adult_consent,
			adult_consent_at,
			ip_address,
			user_agent,
			created_at,
			updated_at
		FROM parent_consents
		WHERE platform_id = $1 AND parent_user_id = $2
	`

	var consent ParentConsent
	err := s.DB.QueryRowContext(ctx, query, platformID, parentUserID).Scan(
		&consent.ID,
		&consent.ParentUserID,
		&consent.PlatformID,
		&consent.PrivacyPolicyVersion,
		&consent.PrivacyPolicyAccepted,
		&consent.PrivacyPolicyAcceptedAt,
		&consent.TermsVersion,
		&consent.TermsAccepted,
		&consent.TermsAcceptedAt,
		&consent.AdultConsent,
		&consent.AdultConsentAt,
		&consent.IPAddress,
		&consent.UserAgent,
		&consent.CreatedAt,
		&consent.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get parent consent: %w", err)
	}

	return &consent, nil
}

// HasValidConsent проверяет, есть ли у родителя действительное согласие
func (s *Store) HasValidConsent(ctx context.Context, platformID, parentUserID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM parent_consents
			WHERE platform_id = $1
			  AND parent_user_id = $2
			  AND privacy_policy_accepted = TRUE
			  AND terms_accepted = TRUE
			  AND adult_consent = TRUE
		)
	`

	var exists bool
	err := s.DB.QueryRowContext(ctx, query, platformID, parentUserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check valid consent: %w", err)
	}

	return exists, nil
}
