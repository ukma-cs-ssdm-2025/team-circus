package group

import (
	"context"
	"database/sql"

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
		return nil, err
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
		return nil, err
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
			return nil, err
		}
		groups = append(groups, &group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}
