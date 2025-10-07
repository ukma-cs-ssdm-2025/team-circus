package user

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *UserRepository) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	query := `
		SELECT uuid, login, email, created_at 
		FROM users 
		WHERE uuid = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT uuid, login, email, created_at 
		FROM users 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.UUID,
			&user.Login,
			&user.Email,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
