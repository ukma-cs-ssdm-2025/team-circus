package group

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (s *GroupService) Delete(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.Delete(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrGroupNotFound
		}
		return fmt.Errorf("group service: delete: %w", err)
	}

	return nil
}
