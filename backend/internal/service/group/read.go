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

func (s *GroupService) GetByUUIDForUser(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.Group, error) {
	member, err := s.memberRepo.GetMember(ctx, groupUUID, userUUID)
	if err != nil {
		return nil, fmt.Errorf("group service: getByUUIDForUser: %w", err)
	}
	if member == nil {
		return nil, domain.ErrForbidden
	}

	group, err := s.repo.GetByUUID(ctx, groupUUID)
	if err != nil {
		return nil, fmt.Errorf("group service: getByUUIDForUser: %w", err)
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

func (s *GroupService) GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Group, error) {
	groups, err := s.repo.GetAllForUser(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("group service: getAllForUser: %w", err)
	}

	return groups, nil
}
