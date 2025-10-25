package reg

import (
	"context"
	"fmt"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func (s *RegService) Register(ctx context.Context, login string, email string, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), s.hashingCost)
	if err != nil {
		return nil, fmt.Errorf("registration service: %w", err)
	}

	user, err := s.repo.Register(ctx, login, email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("registration service: %w", err)
	}

	return user, nil
}
