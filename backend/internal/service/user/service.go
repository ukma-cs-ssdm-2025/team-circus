package user

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
)

type UserService struct {
	repo        *user.UserRepository
	hashingCost int
}

func NewUserService(repo *user.UserRepository, hashingCost int) *UserService {
	return &UserService{
		repo:        repo,
		hashingCost: hashingCost,
	}
}
