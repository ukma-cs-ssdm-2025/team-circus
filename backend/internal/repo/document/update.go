package document

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *DocumentRepository) Update(ctx context.Context, uuid uuid.UUID, name, content string) (*domain.Document, error) {
	query := `
		UPDATE documents 
		SET name = $1, content = $2 
		WHERE uuid = $3 
		RETURNING uuid, group_uuid, name, content, created_at`

	var document domain.Document
	err := r.db.QueryRowContext(ctx, query, name, content, uuid).Scan(
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
