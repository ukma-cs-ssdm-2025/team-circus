package reg

import (
	"context"
	"fmt"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

func (s *RegService) Register(ctx context.Context, login string, email string, password string) (*domain.User, error) {
	user, err := s.repo.Register(ctx, login, email, password)
	if err != nil {
		return nil, fmt.Errorf("registration service: %w", err)
	}

	return user, nil
}
