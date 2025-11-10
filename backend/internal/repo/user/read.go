package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *UserRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	query := `
		SELECT uuid, login, email, hashed_password, created_at
		FROM users
		WHERE uuid = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user repository: getByUUID: %w", err))
	}

	return &user, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `
    SELECT uuid, login, email, hashed_password, created_at
    FROM users
    WHERE login = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, login).Scan(
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user repository: getByLogin: %w", err))
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, params PageParams) ([]*domain.User, error) {
	params = params.Normalize()

	query := `
		SELECT uuid, login, email, hashed_password, created_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, params.Limit, params.Offset)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user repository: getAll query: %w", err))
	}
	defer rows.Close() //nolint:errcheck

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.UUID,
			&user.Login,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user repository: getAll scan: %w", err))
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user repository: getAll rows err: %w", err))
	}

	return users, nil
}
