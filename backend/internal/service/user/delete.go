package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (s *UserService) Delete(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.Delete(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return fmt.Errorf("user service: delete: %w", err)
	}

	return nil
}
