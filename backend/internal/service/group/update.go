package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Update(ctx context.Context, userUUID, groupUUID uuid.UUID, name string) (*domain.Group, error) {
	existing, err := s.repo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("group service: update get group: %w", err)
	}

	if existing == nil {
		return nil, domain.ErrGroupNotFound
	}

	if existing.AuthorUUID != userUUID {
		return nil, domain.ErrForbidden
	}

	group, err := s.repo.Update(ctx, groupUUID, name)
	if err != nil {
		return nil, fmt.Errorf("group service: update: %w", err)
	}

	if group == nil {
		return nil, domain.ErrGroupNotFound
	}

	group.AuthorUUID = existing.AuthorUUID
	group.Role = domain.GroupRoleAuthor

	return group, nil
}
