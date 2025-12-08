package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Update(ctx context.Context, userUUID, groupUUID uuid.UUID, name string) (*domain.Group, error) {
	member, err := s.memberRepo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("group service: update member: %w", err)
	}
	if member == nil {
		return nil, domain.ErrForbidden
	}
	if member.Role != domain.RoleAuthor {
		return nil, domain.ErrForbidden
	}

	group, err := s.repo.Update(ctx, groupUUID, name)
	if err != nil {
		return nil, fmt.Errorf("group service: update: %w", err)
	}

	if group == nil {
		return nil, domain.ErrGroupNotFound
	}

	return group, nil
}
