package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Create(ctx context.Context, userUUID uuid.UUID, name string) (*domain.Group, error) {
	group, err := s.repo.Create(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("group service: create: %w", err)
	}

	if err = s.repo.AddMember(ctx, group.UUID, userUUID, domain.GroupRoleOwner); err != nil {
		return nil, fmt.Errorf("group service: add creator as member: %w", err)
	}

	return group, nil
}
