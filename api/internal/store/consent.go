package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// ParentConsent представляет согласие родителя
type ParentConsent struct {
	ID                      string
	ParentUserID            string
	PlatformID              string
	PrivacyPolicyVersion    string
	PrivacyPolicyAccepted   bool
	PrivacyPolicyAcceptedAt time.Time
	TermsVersion            string
	TermsAccepted           bool
	TermsAcceptedAt         time.Time
	AdultConsent            bool
	AdultConsentAt          sql.NullTime
	IPAddress               sql.NullString
	UserAgent               sql.NullString
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// ParentConsentHistory представляет запись истории изменения согласия
type ParentConsentHistory struct {
	ID                    string
	ConsentID             string
	ParentUserID          string
	PlatformID            string
	Action                string
	PrivacyPolicyVersion  string
	PrivacyPolicyAccepted bool
	TermsVersion          string
	TermsAccepted         bool
	AdultConsent          bool
	ChangedFields         []string
	PreviousValues        map[string]interface{}
	IPAddress             sql.NullString
	UserAgent             sql.NullString
	CreatedAt             time.Time
}

// SaveParentConsent сохраняет или обновляет согласие родителя (UPSERT) с audit trail
func (s *Store) SaveParentConsent(
	ctx context.Context,
	platformID string,
	parentUserID string,
	privacyPolicyVersion string,
	termsVersion string,
	adultConsent bool,
	ipAddress, userAgent string,
) error {
	// Начинаем транзакцию
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Получаем текущее согласие (если существует)
	var existingConsent *ParentConsent
	existingConsent, err = s.GetParentConsentTx(ctx, tx, platformID, parentUserID)
	if err != nil {
		return fmt.Errorf("get existing consent: %w", err)
	}

	// Определяем тип действия и изменения
	action := "created"
	var changedFields []string
	previousValues := make(map[string]interface{})

	if existingConsent != nil {
		action = "updated"

		// Проверяем что изменилось
		if existingConsent.PrivacyPolicyVersion != privacyPolicyVersion {
			changedFields = append(changedFields, "privacy_policy_version")
			previousValues["privacy_policy_version"] = existingConsent.PrivacyPolicyVersion
		}
		if existingConsent.TermsVersion != termsVersion {
			changedFields = append(changedFields, "terms_version")
			previousValues["terms_version"] = existingConsent.TermsVersion
		}
		if existingConsent.AdultConsent != adultConsent {
			changedFields = append(changedFields, "adult_consent")
			previousValues["adult_consent"] = existingConsent.AdultConsent
		}
	}

	// UPSERT согласия
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
		RETURNING id
	`

	var consentID string
	err = tx.QueryRowContext(
		ctx,
		query,
		platformID,
		parentUserID,
		privacyPolicyVersion,
		termsVersion,
		adultConsent,
		sql.NullString{String: ipAddress, Valid: ipAddress != ""},
		sql.NullString{String: userAgent, Valid: userAgent != ""},
	).Scan(&consentID)
	if err != nil {
		return fmt.Errorf("save parent consent: %w", err)
	}

	// Сохраняем в историю
	err = s.SaveConsentHistoryTx(
		ctx,
		tx,
		consentID,
		platformID,
		parentUserID,
		action,
		privacyPolicyVersion,
		termsVersion,
		adultConsent,
		changedFields,
		previousValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		return fmt.Errorf("save consent history: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GetParentConsentTx получает согласие родителя внутри транзакции
func (s *Store) GetParentConsentTx(ctx context.Context, tx *sql.Tx, platformID, parentUserID string) (*ParentConsent, error) {
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
	err := tx.QueryRowContext(ctx, query, platformID, parentUserID).Scan(
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

// SaveConsentHistoryTx сохраняет запись истории изменения согласия
func (s *Store) SaveConsentHistoryTx(
	ctx context.Context,
	tx *sql.Tx,
	consentID string,
	platformID string,
	parentUserID string,
	action string,
	privacyPolicyVersion string,
	termsVersion string,
	adultConsent bool,
	changedFields []string,
	previousValues map[string]interface{},
	ipAddress, userAgent string,
) error {
	// Конвертируем previousValues в JSONB
	var previousValuesJSON []byte
	if len(previousValues) > 0 {
		var err error
		previousValuesJSON, err = json.Marshal(previousValues)
		if err != nil {
			return fmt.Errorf("marshal previous values: %w", err)
		}
	}

	query := `
		INSERT INTO parent_consent_history (
			consent_id,
			platform_id,
			parent_user_id,
			action,
			privacy_policy_version,
			privacy_policy_accepted,
			terms_version,
			terms_accepted,
			adult_consent,
			changed_fields,
			previous_values,
			ip_address,
			user_agent
		) VALUES ($1, $2, $3, $4, $5, TRUE, $6, TRUE, $7, $8, $9, $10, $11)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		consentID,
		platformID,
		parentUserID,
		action,
		privacyPolicyVersion,
		termsVersion,
		adultConsent,
		pq.Array(changedFields),
		sql.NullString{String: string(previousValuesJSON), Valid: len(previousValuesJSON) > 0},
		sql.NullString{String: ipAddress, Valid: ipAddress != ""},
		sql.NullString{String: userAgent, Valid: userAgent != ""},
	)
	if err != nil {
		return fmt.Errorf("insert consent history: %w", err)
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

// GetConsentHistory получает историю изменений согласия родителя
func (s *Store) GetConsentHistory(ctx context.Context, platformID, parentUserID string, limit int) ([]*ParentConsentHistory, error) {
	if limit <= 0 {
		limit = 50 // default
	}

	query := `
		SELECT
			id,
			consent_id,
			parent_user_id,
			platform_id,
			action,
			privacy_policy_version,
			privacy_policy_accepted,
			terms_version,
			terms_accepted,
			adult_consent,
			changed_fields,
			previous_values,
			ip_address,
			user_agent,
			created_at
		FROM parent_consent_history
		WHERE platform_id = $1 AND parent_user_id = $2
		ORDER BY created_at DESC
		LIMIT $3
	`

	rows, err := s.DB.QueryContext(ctx, query, platformID, parentUserID, limit)
	if err != nil {
		return nil, fmt.Errorf("query consent history: %w", err)
	}
	defer rows.Close()

	var history []*ParentConsentHistory
	for rows.Next() {
		var h ParentConsentHistory
		var previousValuesJSON sql.NullString

		err := rows.Scan(
			&h.ID,
			&h.ConsentID,
			&h.ParentUserID,
			&h.PlatformID,
			&h.Action,
			&h.PrivacyPolicyVersion,
			&h.PrivacyPolicyAccepted,
			&h.TermsVersion,
			&h.TermsAccepted,
			&h.AdultConsent,
			pq.Array(&h.ChangedFields),
			&previousValuesJSON,
			&h.IPAddress,
			&h.UserAgent,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan consent history: %w", err)
		}

		// Парсим previous_values JSON
		if previousValuesJSON.Valid && previousValuesJSON.String != "" {
			if err := json.Unmarshal([]byte(previousValuesJSON.String), &h.PreviousValues); err != nil {
				return nil, fmt.Errorf("unmarshal previous values: %w", err)
			}
		}

		history = append(history, &h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return history, nil
}
