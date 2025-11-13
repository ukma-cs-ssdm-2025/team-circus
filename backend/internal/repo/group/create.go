package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
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

func (r *GroupRepository) AddMember(ctx context.Context, groupUUID, userUUID uuid.UUID, role string) error {
	const query = `
		INSERT INTO user_groups (user_uuid, group_uuid, role)
		VALUES ($1, $2, $3)`

	if _, err := r.db.ExecContext(ctx, query, userUUID, groupUUID, role); err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("group repository: addMember: %w", err))
	}

	return nil
}
