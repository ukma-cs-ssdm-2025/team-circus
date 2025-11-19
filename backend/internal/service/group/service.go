package group

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/member"
)

type GroupService struct {
	repo       *group.GroupRepository
	memberRepo *member.MemberRepository
}

func NewGroupService(repo *group.GroupRepository, memberRepo *member.MemberRepository) *GroupService {
	return &GroupService{
		repo:       repo,
		memberRepo: memberRepo,
	}
}
