package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *GroupRepository) Update(ctx context.Context, uuid uuid.UUID, name string) (*domain.Group, error) {
	query := `
		UPDATE groups 
		SET name = $1 
		WHERE uuid = $2 
		RETURNING uuid, name, created_at`

	var group domain.Group
	err := r.db.QueryRowContext(ctx, query, name, uuid).Scan(
		&group.UUID,
		&group.Name,
		&group.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: update: %w", err))
	}

	return &group, nil
}
