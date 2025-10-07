package document

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *DocumentRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Document, error) {
	query := `
		SELECT uuid, group_uuid, name, content, created_at 
		FROM documents 
		WHERE uuid = $1`

	var document domain.Document
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&document.UUID,
		&document.GroupUUID,
		&document.Name,
		&document.Content,
		&document.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &document, nil
}

func (r *DocumentRepository) GetByGroupUUID(ctx context.Context, groupUUID uuid.UUID) ([]*domain.Document, error) {
	query := `
		SELECT uuid, group_uuid, name, content, created_at 
		FROM documents 
		WHERE group_uuid = $1 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, groupUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var documents []*domain.Document
	for rows.Next() {
		var document domain.Document
		err := rows.Scan(
			&document.UUID,
			&document.GroupUUID,
			&document.Name,
			&document.Content,
			&document.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *DocumentRepository) GetAll(ctx context.Context) ([]*domain.Document, error) {
	query := `
		SELECT uuid, group_uuid, name, content, created_at 
		FROM documents 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck

	var documents []*domain.Document
	for rows.Next() {
		var document domain.Document
		err := rows.Scan(
			&document.UUID,
			&document.GroupUUID,
			&document.Name,
			&document.Content,
			&document.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}
