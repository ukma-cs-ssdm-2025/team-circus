package group

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *GroupRepository) Create(ctx context.Context, name string, authorUUID uuid.UUID) (*domain.Group, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: create begin tx: %w", err))
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	const insertGroupQuery = `
		INSERT INTO groups (name)
		VALUES ($1)
		RETURNING uuid, name, created_at`

	var group domain.Group
	if err := tx.QueryRowContext(ctx, insertGroupQuery, name).Scan(
		&group.UUID,
		&group.Name,
		&group.CreatedAt,
	); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: create insert group: %w", err))
	}

	const insertMembershipQuery = `
		INSERT INTO user_groups (user_uuid, group_uuid, role)
		VALUES ($1, $2, $3)`

	if _, err := tx.ExecContext(ctx, insertMembershipQuery, authorUUID, group.UUID, domain.GroupRoleAuthor); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: create insert membership: %w", err))
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: create commit: %w", err))
	}
	committed = true

	group.AuthorUUID = authorUUID
	group.Role = domain.GroupRoleAuthor

	return &group, nil
}
