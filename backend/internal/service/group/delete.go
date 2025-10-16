package group

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Delete(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.Delete(ctx, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrGroupNotFound
		}
		return fmt.Errorf("group service: delete: %w", err)
	}

	return nil
}
