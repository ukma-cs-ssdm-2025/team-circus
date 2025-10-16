package document

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *DocumentRepository) Create(ctx context.Context, groupUUID uuid.UUID, name, content string) (*domain.Document, error) {
	query := `
		INSERT INTO documents (group_uuid, name, content) 
		VALUES ($1, $2, $3) 
		RETURNING uuid, group_uuid, name, content, created_at`

	var document domain.Document
	err := r.db.QueryRowContext(ctx, query, groupUUID, name, content).Scan(
		&document.UUID,
		&document.GroupUUID,
		&document.Name,
		&document.Content,
		&document.CreatedAt,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("document repository: create: %w", err))
	}

	return &document, nil
}
