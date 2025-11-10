package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	userrepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) Update(ctx context.Context, uuid uuid.UUID, params userrepo.UpdateUserParams) (*domain.User, error) {
	if params.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*params.Password), s.hashingCost)
		if err != nil {
			return nil, fmt.Errorf("user service: update: hash password: %w", err)
		}
		hashedStr := string(hashed)
		params.Password = &hashedStr
	}

	user, err := s.repo.Update(ctx, uuid, params)
	if err != nil {
		return nil, fmt.Errorf("user service: update: %w", err)
	}

	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}
