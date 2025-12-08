package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *MemberService) CreateMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID,
	role string) (*domain.Member, error) {
	group, err := s.groupRepo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: create member get group: %w", err)
	}
	if group == nil {
		return nil, domain.ErrGroupNotFound
	}

	if role == domain.RoleAuthor {
		return nil, domain.ErrOnlyAuthor
	}

	member, err := s.repo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: create member get actor: %w", err)
	}
	if member == nil {
		return nil, domain.ErrForbidden
	}
	if member.Role != domain.RoleAuthor {
		return nil, domain.ErrForbidden
	}

	user, err := s.userRepo.GetByUUID(ctx, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: create member get user: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	existingMember, err := s.repo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: create member get existing: %w", err)
	}
	if existingMember != nil {
		return nil, domain.ErrAlreadyExists
	}

	createdMember, err := s.repo.CreateMember(ctx, groupUUID, memberUUID, role)
	if err != nil {
		return nil, fmt.Errorf("member service: create: %w", err)
	}

	return createdMember, nil
}
