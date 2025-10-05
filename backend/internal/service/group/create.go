package group

import (
	"context"
	"fmt"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Create(ctx context.Context, name string) (*domain.Group, error) {
	group, err := s.repo.Create(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("group service: create: %w", err)
	}

	return group, nil
}
