package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *UserService) Update(ctx context.Context, uuid uuid.UUID, login string, email string, password string) (*domain.User, error) {
	user, err := s.repo.Update(ctx, uuid, login, email, password)
	if err != nil {
		return nil, fmt.Errorf("user service: update: %w", err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}
