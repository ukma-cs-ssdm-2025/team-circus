package member

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *MemberService) DeleteMemberByUser(ctx context.Context, userUUID, groupUUID, memberUUID uuid.UUID) error {
	group, err := s.groupRepo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return fmt.Errorf("member service: delete member get group: %w", err)
	}
	if group == nil {
		return domain.ErrGroupNotFound
	}

	actor, err := s.repo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return fmt.Errorf("member service: delete member get actor: %w", err)
	}
	if actor == nil {
		return domain.ErrForbidden
	}
	if actor.Role != domain.RoleAuthor {
		return domain.ErrForbidden
	}

	targetMember, err := s.repo.GetMember(ctx, groupUUID, memberUUID)
	if err != nil {
		return fmt.Errorf("member service: delete member get target: %w", err)
	}
	if targetMember == nil {
		return domain.ErrUserNotFound
	}
	if targetMember.Role == domain.RoleAuthor {
		return domain.ErrOnlyAuthor
	}

	if err := s.repo.DeleteMember(ctx, groupUUID, memberUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("member service: delete member: %w", err)
	}

	return nil
}
