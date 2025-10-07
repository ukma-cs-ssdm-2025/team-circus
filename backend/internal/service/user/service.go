package user

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
)

type UserService struct {
	repo *user.UserRepository
}

func NewUserService(repo *user.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
