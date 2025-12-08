package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *MemberService) GetAllMembersForUser(ctx context.Context, userUUID, groupUUID uuid.UUID) ([]*domain.Member, error) {
	group, err := s.groupRepo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: get all members get group: %w", err)
	}
	if group == nil {
		return nil, domain.ErrGroupNotFound
	}

	member, err := s.repo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: get all members get actor: %w", err)
	}
	if member == nil {
		return nil, domain.ErrForbidden
	}

	members, err := s.repo.GetAllMembers(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: get all members: %w", err)
	}

	return members, nil
}
