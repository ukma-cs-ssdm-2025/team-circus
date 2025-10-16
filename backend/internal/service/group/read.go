package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Group, error) {
	group, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("group service: getByUUID: %w", err)
	}

	if group == nil {
		return nil, domain.ErrGroupNotFound
	}

	return group, nil
}

func (s *GroupService) GetAll(ctx context.Context) ([]*domain.Group, error) {
	groups, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("group service: getAll: %w", err)
	}

	return groups, nil
}
