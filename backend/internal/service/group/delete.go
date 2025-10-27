package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Delete(ctx context.Context, userUUID, groupUUID uuid.UUID) error {
	group, err := s.repo.GetByUUID(ctx, groupUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrGroupNotFound
		}
		return fmt.Errorf("group service: delete get group: %w", err)
	}

	if group == nil {
		return domain.ErrGroupNotFound
	}

	if group.AuthorUUID != userUUID {
		return domain.ErrForbidden
	}

	err = s.repo.Delete(ctx, groupUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.ErrGroupNotFound
		}
		return fmt.Errorf("group service: delete: %w", err)
	}

	return nil
}
