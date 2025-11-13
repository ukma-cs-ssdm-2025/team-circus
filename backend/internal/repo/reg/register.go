package reg

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *RegRepository) Register(ctx context.Context, login string, email string, password string) (*domain.User, error) {
	query := `
		INSERT INTO users (login, email, hashed_password) 
		VALUES ($1, $2, $3) 
		RETURNING uuid, login, email, hashed_password, created_at`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, login, email, password).Scan(
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, domain.ErrAlreadyExists
		}
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("reg repository: register: %w", err))
	}

	return &user, nil
}
