package reg

import (
	"context"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

type RegRepository interface {
	Register(ctx context.Context, login string, email string, password string) (*domain.User, error)
}

type RegService struct {
	repo        RegRepository
	hashingCost int
}

func NewRegService(repo RegRepository, cost int) *RegService {
	return &RegService{
		repo:        repo,
		hashingCost: cost,
	}
}
