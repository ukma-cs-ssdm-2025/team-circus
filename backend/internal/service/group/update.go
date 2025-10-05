package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Update(ctx context.Context, uuid uuid.UUID, name string) (*domain.Group, error) {
	group, err := s.repo.Update(ctx, uuid, name)
	if err != nil {
		return nil, fmt.Errorf("group service: update: %w", err)
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	return group, nil
}
