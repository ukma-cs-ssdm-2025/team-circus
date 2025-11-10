package user

import "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"

type UserService struct {
	repo user.Repository
}

func NewUserService(repo user.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}
