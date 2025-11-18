package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *GroupService) Create(ctx context.Context, ownerUUID uuid.UUID, name string) (*domain.Group, error) {
	group, err := s.repo.Create(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("group service: create: %w", err)
	}

	if _, err := s.memberRepo.CreateMember(ctx, group.UUID, ownerUUID, domain.RoleAuthor); err != nil {
		if deleteErr := s.repo.Delete(ctx, group.UUID); deleteErr != nil {
			return nil, fmt.Errorf("group service: create add owner: %w, cleanup: %v", err, deleteErr)
		}
		return nil, fmt.Errorf("group service: create add owner: %w", err)
	}

	return group, nil
}
