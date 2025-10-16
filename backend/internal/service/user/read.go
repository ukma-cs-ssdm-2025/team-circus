package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *UserService) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	user, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("user service: getByUUID: %w", err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("user service: getAll: %w", err)
	}

	return users, nil
}
