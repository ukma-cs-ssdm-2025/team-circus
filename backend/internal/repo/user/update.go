package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *UserRepository) Update(ctx context.Context, uuid uuid.UUID, login string, email string, password string) (*domain.User, error) {
	query := `
		UPDATE users 
		SET login = $1, email = $2, hashed_password = $3
		WHERE uuid = $4
		RETURNING uuid, login, email, hashed_password, created_at`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, login, email, password, uuid).Scan(
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
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user repository: update: %w", err))
	}

	return &user, nil
}
