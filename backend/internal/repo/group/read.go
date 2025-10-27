package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *GroupRepository) GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*domain.Group, error) {
	const query = `
		SELECT g.uuid, g.name, g.created_at, owner.user_uuid
		FROM groups g
		LEFT JOIN user_groups owner ON owner.group_uuid = g.uuid AND owner.role = $2
		WHERE g.uuid = $1`

	var group domain.Group
	var authorUUID sql.NullString
	err := r.db.QueryRowContext(ctx, query, groupUUID, domain.GroupRoleAuthor).Scan(
		&group.UUID,
		&group.Name,
		&group.CreatedAt,
		&authorUUID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getByUUID: %w", err))
	}

	if authorUUID.Valid {
		parsedAuthorUUID, parseErr := uuid.Parse(authorUUID.String)
		if parseErr != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getByUUID parse author: %w", parseErr))
		}
		group.AuthorUUID = parsedAuthorUUID
	}

	return &group, nil
}

func (r *GroupRepository) IsMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM user_groups
			WHERE group_uuid = $1 AND user_uuid = $2
		)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, groupUUID, userUUID).Scan(&exists)
	if err != nil {
		return false, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: isMember query: %w", err))
	}
	return exists, nil
}

func (r *GroupRepository) GetAll(ctx context.Context) ([]*domain.Group, error) {
	const query = `
		SELECT g.uuid, g.name, g.created_at, owner.user_uuid
		FROM groups g
		LEFT JOIN user_groups owner ON owner.group_uuid = g.uuid AND owner.role = $1
		ORDER BY g.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, domain.GroupRoleAuthor)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAll query: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var groups []*domain.Group
	for rows.Next() {
		var group domain.Group
		var authorUUID sql.NullString
		err := rows.Scan(
			&group.UUID,
			&group.Name,
			&group.CreatedAt,
			&authorUUID,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAll scan: %w", err))
		}
		if authorUUID.Valid {
			parsedAuthorUUID, parseErr := uuid.Parse(authorUUID.String)
			if parseErr != nil {
				return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAll parse author: %w", parseErr))
			}
			group.AuthorUUID = parsedAuthorUUID
		}
		groups = append(groups, &group)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAll rows err: %w", err))
	}

	return groups, nil
}

func (r *GroupRepository) GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Group, error) {
	const query = `
		SELECT g.uuid, g.name, g.created_at, owner.user_uuid, ug.role
		FROM groups g
		INNER JOIN user_groups ug ON ug.group_uuid = g.uuid
		LEFT JOIN user_groups owner ON owner.group_uuid = g.uuid AND owner.role = $2
		WHERE ug.user_uuid = $1
		ORDER BY g.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userUUID, domain.GroupRoleAuthor)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAllForUser query: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var groups []*domain.Group
	for rows.Next() {
		var group domain.Group
		var authorUUID sql.NullString
		var role sql.NullString
		err := rows.Scan(
			&group.UUID,
			&group.Name,
			&group.CreatedAt,
			&authorUUID,
			&role,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAllForUser scan: %w", err))
		}
		if authorUUID.Valid {
			parsedAuthorUUID, parseErr := uuid.Parse(authorUUID.String)
			if parseErr != nil {
				return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAllForUser parse author: %w", parseErr))
			}
			group.AuthorUUID = parsedAuthorUUID
		}
		if role.Valid {
			group.Role = role.String
		}
		groups = append(groups, &group)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("group repository: getAllForUser rows err: %w", err))
	}

	return groups, nil
}
