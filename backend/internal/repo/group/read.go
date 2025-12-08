package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *GroupRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Group, error) {
	query := `
		SELECT uuid, name, created_at 
		FROM groups 
		WHERE uuid = $1`

	var group domain.Group
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&group.UUID,
		&group.Name,
		&group.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getByUUID: %w", err))
	}

	return &group, nil
}

func (r *GroupRepository) GetAll(ctx context.Context) ([]*domain.Group, error) {
	query := `
		SELECT uuid, name, created_at 
		FROM groups 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAll query: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var groups []*domain.Group
	for rows.Next() {
		var group domain.Group
		err := rows.Scan(
			&group.UUID,
			&group.Name,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAll scan: %w", err))
		}
		groups = append(groups, &group)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAll rows err: %w", err))
	}

	return groups, nil
}

func (r *GroupRepository) GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Group, error) {
	query := `
		SELECT g.uuid, g.name, g.created_at
		FROM groups g
		INNER JOIN user_groups ug ON ug.group_uuid = g.uuid
		WHERE ug.user_uuid = $1
		ORDER BY g.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userUUID)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAllForUser query: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var groups []*domain.Group
	for rows.Next() {
		var group domain.Group
		err := rows.Scan(
			&group.UUID,
			&group.Name,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAllForUser scan: %w", err))
		}
		groups = append(groups, &group)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAllForUser rows err: %w", err))
	}

	return groups, nil
}
