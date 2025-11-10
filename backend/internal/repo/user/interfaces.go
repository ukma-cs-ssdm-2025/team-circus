package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

const (
	// DefaultPageLimit is applied when the caller does not request a limit explicitly.
	DefaultPageLimit = 50
	// MaxPageLimit guards against unbounded queries.
	MaxPageLimit = 100
)

// PageParams describes limit/offset pagination options for user listings.
// Limit defaults to DefaultPageLimit and is capped at MaxPageLimit.
// Offset is coerced to zero when a negative value is supplied.
type PageParams struct {
	Limit  int
	Offset int
}

// Normalize applies default and max values to the pagination parameters.
func (p PageParams) Normalize() PageParams {
	if p.Limit <= 0 {
		p.Limit = DefaultPageLimit
	}
	if p.Limit > MaxPageLimit {
		p.Limit = MaxPageLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}

// UpdateUserParams describes optional fields that can be updated for a user.
// Nil pointers mean "leave the existing value". Password must contain the
// already hashed (e.g., bcrypt) value by the time it reaches the repository.
type UpdateUserParams struct {
	Login    *string
	Email    *string
	Password *string
}

// Repository describes the storage contract for users.
//
//go:generate go run github.com/golang/mock/mockgen -destination=../../mocks/mock_user_repository.go -package=mocks . Repository
type Repository interface {
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	GetByLogin(ctx context.Context, login string) (*domain.User, error)
	GetAll(ctx context.Context, params PageParams) ([]*domain.User, error)
	Update(ctx context.Context, uuid uuid.UUID, params UpdateUserParams) (*domain.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}
