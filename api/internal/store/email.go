package store

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// EmailVerification представляет запись верификации email
type EmailVerification struct {
	ID               string
	Email            string
	VerificationCode string
	IsVerified       bool
	VerifiedAt       *time.Time
	ExpiresAt        time.Time
	SendAttempts     int
	VerifyAttempts   int
	ParentUserID     string
	PlatformID       string
	IPAddress        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// CreateEmailVerification создает новую запись верификации
func (s *Store) CreateEmailVerification(
	ctx context.Context,
	email, code, parentUserID, platformID, ipAddress string,
	expiresAt time.Time,
) error {
	// Сначала проверяем, есть ли уже неверифицированная запись
	var existingID string
	checkQuery := `
		SELECT id FROM email_verifications
		WHERE email = $1 AND is_verified = FALSE
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := s.DB.QueryRowContext(ctx, checkQuery, email).Scan(&existingID)

	if err == nil {
		// Запись существует, обновляем её
		updateQuery := `
			UPDATE email_verifications
			SET verification_code = $1,
			    expires_at = $2,
			    send_attempts = send_attempts + 1,
			    parent_user_id = $3,
			    platform_id = $4,
			    ip_address = $5,
			    updated_at = NOW()
			WHERE id = $6
		`
		_, err = s.DB.ExecContext(
			ctx,
			updateQuery,
			code,
			expiresAt,
			parentUserID,
			platformID,
			ipAddress,
			existingID,
		)
		return err
	} else if err != sql.ErrNoRows {
		return err
	}

	// Записи нет, создаём новую
	query := `
		INSERT INTO email_verifications (
			email, verification_code, expires_at,
			parent_user_id, platform_id, ip_address
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = s.DB.ExecContext(
		ctx,
		query,
		email,
		code,
		expiresAt,
		parentUserID,
		platformID,
		ipAddress,
	)

	if err != nil {
		log.Printf("[Store] Failed to create email verification: %v", err)
		return err
	}

	return nil
}

// VerifyEmailCode проверяет код верификации
func (s *Store) VerifyEmailCode(ctx context.Context, email, code string) (bool, error) {
	// Проверяем код и срок действия
	query := `
		SELECT id, expires_at, is_verified
		FROM email_verifications
		WHERE email = $1 AND verification_code = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var id string
	var expiresAt time.Time
	var isVerified bool

	err := s.DB.QueryRowContext(ctx, query, email, code).Scan(&id, &expiresAt, &isVerified)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// Уже верифицирован
	if isVerified {
		return true, nil
	}

	// Проверяем срок действия
	if time.Now().After(expiresAt) {
		// Инкрементируем счётчик попыток
		s.DB.ExecContext(
			ctx,
			"UPDATE email_verifications SET verify_attempts = verify_attempts + 1 WHERE id = $1",
			id,
		)
		return false, nil
	}

	// Код валиден, помечаем как верифицированный
	updateQuery := `
		UPDATE email_verifications
		SET is_verified = TRUE,
		    verified_at = NOW(),
		    verify_attempts = verify_attempts + 1,
		    updated_at = NOW()
		WHERE id = $1
	`

	_, err = s.DB.ExecContext(ctx, updateQuery, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// IsEmailVerified проверяет, верифицирован ли email
func (s *Store) IsEmailVerified(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT is_verified
		FROM email_verifications
		WHERE email = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var isVerified bool
	err := s.DB.QueryRowContext(ctx, query, email).Scan(&isVerified)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return isVerified, nil
}

// GetEmailVerification получает запись верификации
func (s *Store) GetEmailVerification(ctx context.Context, email string) (*EmailVerification, error) {
	query := `
		SELECT id, email, verification_code, is_verified, verified_at,
		       expires_at, send_attempts, verify_attempts,
		       parent_user_id, platform_id, ip_address,
		       created_at, updated_at
		FROM email_verifications
		WHERE email = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var v EmailVerification
	err := s.DB.QueryRowContext(ctx, query, email).Scan(
		&v.ID,
		&v.Email,
		&v.VerificationCode,
		&v.IsVerified,
		&v.VerifiedAt,
		&v.ExpiresAt,
		&v.SendAttempts,
		&v.VerifyAttempts,
		&v.ParentUserID,
		&v.PlatformID,
		&v.IPAddress,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &v, nil
}
