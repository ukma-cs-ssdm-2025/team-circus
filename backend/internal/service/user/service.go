package user

import "github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"

type UserService struct {
	repo        user.Repository
	hashingCost int
}

func NewUserService(repo user.Repository, hashingCost int) *UserService {
	return &UserService{
		repo:        repo,
		hashingCost: hashingCost,
	}
}
