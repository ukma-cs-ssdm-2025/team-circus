package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *GroupRepository) Create(ctx context.Context, name string) (*domain.Group, error) {
	query := `
		INSERT INTO groups (name) 
		VALUES ($1) 
		RETURNING uuid, name, created_at`

	var group domain.Group
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&group.UUID,
		&group.Name,
		&group.CreatedAt,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: create: %w", err))
	}

	return &group, nil
}
