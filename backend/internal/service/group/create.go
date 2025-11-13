package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Create(ctx context.Context, authorUUID uuid.UUID, name string) (*domain.Group, error) {
	group, err := s.repo.Create(ctx, name, authorUUID)
	if err != nil {
		return nil, fmt.Errorf("group service: create: %w", err)
	}

	group.AuthorUUID = authorUUID
	group.Role = domain.GroupRoleAuthor

	return group, nil
}
