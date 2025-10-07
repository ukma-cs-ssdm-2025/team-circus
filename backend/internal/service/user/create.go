package user

import (
	"context"
	"fmt"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *UserService) Create(ctx context.Context, login string, email string, password string) (*domain.User, error) {
	user, err := s.repo.Create(ctx, login, email, password)
	if err != nil {
		return nil, fmt.Errorf("user service: create: %w", err)
	}

	return user, nil
}
