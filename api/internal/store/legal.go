package store

import (
	"context"
	"database/sql"
	"time"
)

// LegalDocument представляет юридический документ
type LegalDocument struct {
	ID            string
	DocumentType  string
	Version       string
	Title         string
	Content       string
	Language      string
	IsActive      bool
	EffectiveDate time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// GetActiveLegalDocument получает активную версию юридического документа
func (s *Store) GetActiveLegalDocument(ctx context.Context, documentType, language string) (*LegalDocument, error) {
	query := `
		SELECT id, document_type, version, title, content, language,
		       is_active, effective_date, created_at, updated_at
		FROM legal_documents
		WHERE document_type = $1 AND language = $2 AND is_active = TRUE
		ORDER BY effective_date DESC
		LIMIT 1
	`

	var doc LegalDocument
	err := s.DB.QueryRowContext(ctx, query, documentType, language).Scan(
		&doc.ID,
		&doc.DocumentType,
		&doc.Version,
		&doc.Title,
		&doc.Content,
		&doc.Language,
		&doc.IsActive,
		&doc.EffectiveDate,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

// GetLegalDocumentByVersion получает конкретную версию документа
func (s *Store) GetLegalDocumentByVersion(ctx context.Context, documentType, version, language string) (*LegalDocument, error) {
	query := `
		SELECT id, document_type, version, title, content, language,
		       is_active, effective_date, created_at, updated_at
		FROM legal_documents
		WHERE document_type = $1 AND version = $2 AND language = $3
	`

	var doc LegalDocument
	err := s.DB.QueryRowContext(ctx, query, documentType, version, language).Scan(
		&doc.ID,
		&doc.DocumentType,
		&doc.Version,
		&doc.Title,
		&doc.Content,
		&doc.Language,
		&doc.IsActive,
		&doc.EffectiveDate,
		&doc.CreatedAt,
		&doc.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &doc, nil
}
