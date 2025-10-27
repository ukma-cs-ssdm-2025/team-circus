package group

import (
	"context"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	grouprepo "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
)

type Repository interface {
	Create(ctx context.Context, name string, authorUUID uuid.UUID) (*domain.Group, error)
	GetByUUID(ctx context.Context, groupUUID uuid.UUID) (*domain.Group, error)
	GetMember(ctx context.Context, groupUUID, userUUID uuid.UUID) (*domain.GroupMember, error)
	GetAll(ctx context.Context) ([]*domain.Group, error)
	GetAllForUser(ctx context.Context, userUUID uuid.UUID) ([]*domain.Group, error)
	Update(ctx context.Context, uuid uuid.UUID, name string) (*domain.Group, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}

// ensure the concrete repository satisfies the interface at compile time.
var _ Repository = (*grouprepo.GroupRepository)(nil)

type GroupService struct {
	repo Repository
}

func NewGroupService(repo Repository) *GroupService {
	return &GroupService{
		repo: repo,
	}
}
