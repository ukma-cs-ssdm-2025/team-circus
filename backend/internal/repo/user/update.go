package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (r *UserRepository) Update(ctx context.Context, uuid uuid.UUID, params UpdateUserParams) (*domain.User, error) {
	setClauses := make([]string, 0, 3)
	args := make([]interface{}, 0, 4)
	argIdx := 1

	if params.Login != nil {
		setClauses = append(setClauses, fmt.Sprintf("login = $%d", argIdx))
		args = append(args, *params.Login)
		argIdx++
	}
	if params.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, *params.Email)
		argIdx++
	}
	if params.Password != nil {
		setClauses = append(setClauses, fmt.Sprintf("hashed_password = $%d", argIdx))
		args = append(args, *params.Password)
		argIdx++
	}

	if len(setClauses) == 0 {
		return r.GetByUUID(ctx, uuid)
	}

	var builder strings.Builder
	builder.WriteString("UPDATE users SET ")
	builder.WriteString(strings.Join(setClauses, ", "))
	builder.WriteString(" WHERE uuid = $")
	builder.WriteString(strconv.Itoa(argIdx))
	builder.WriteString(" RETURNING uuid, login, email, hashed_password, created_at")

	query := builder.String()

	args = append(args, uuid)

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
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
