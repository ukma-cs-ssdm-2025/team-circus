package group

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Delete(ctx context.Context, userUUID, groupUUID uuid.UUID) error {
	member, err := s.memberRepo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return fmt.Errorf("group service: delete: %w", err)
	}
	if member == nil {
		return domain.ErrForbidden
	}
	if member.Role != domain.RoleAuthor {
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
