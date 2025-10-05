package group

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
)

type GroupService struct {
	repo *group.GroupRepository
}

func NewGroupService(repo *group.GroupRepository) *GroupService {
	return &GroupService{
		repo: repo,
	}
}
