package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

// Repository describes the storage contract for users.
//
//go:generate go run github.com/golang/mock/mockgen -destination=../../mocks/mock_user_repository.go -package=mocks . Repository
type Repository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Update(ctx context.Context, uuid uuid.UUID, login string, email string, password string) (*domain.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}
