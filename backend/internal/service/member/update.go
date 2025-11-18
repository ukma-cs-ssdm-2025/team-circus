package member

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *MemberService) UpdateMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID,
	role string) (*domain.Member, error) {
	group, err := s.groupRepo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: update member get group: %w", err)
	}
	if group == nil {
		return nil, domain.ErrGroupNotFound
	}

	actor, err := s.repo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: update member get actor: %w", err)
	}
	if actor == nil {
		return nil, domain.ErrForbidden
	}
	if actor.Role != domain.RoleAuthor {
		return nil, domain.ErrForbidden
	}

	targetMember, err := s.repo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return nil, fmt.Errorf("member service: update member get target: %w", err)
	}
	if targetMember == nil {
		return nil, domain.ErrUserNotFound
	}
	if targetMember.Role == role {
		return targetMember, nil
	}

	if targetMember.Role == domain.RoleAuthor && role != domain.RoleAuthor {
		return nil, domain.ErrOnlyAuthor
	}

	if role == domain.RoleAuthor && targetMember.Role != domain.RoleAuthor {
		member, err := s.repo.UpdateMember(ctx, groupUUID, memberUUID, role)
		if err != nil {
			return nil, fmt.Errorf("member service: update member promote: %w", err)
		}
		if member == nil {
			return nil, domain.ErrUserNotFound
		}

		_, err = s.repo.UpdateMember(ctx, groupUUID, actor.UserUUID, domain.RoleEditor)
		if err != nil {
			// very unlucky
		}

		return member, nil
	}

	member, err := s.repo.UpdateMember(ctx, groupUUID, memberUUID, role)
	if err != nil {
		return nil, fmt.Errorf("member service: update member: %w", err)
	}
	if member == nil {
		return nil, domain.ErrUserNotFound
	}

	return member, nil
}
