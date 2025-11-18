package member

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/group"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/member"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/repo/user"
)

type MemberService struct {
	repo      *member.MemberRepository
	groupRepo *group.GroupRepository
	userRepo  *user.UserRepository
}

func NewMemberService(repo *member.MemberRepository, groupRepo *group.GroupRepository,
	userRepo *user.UserRepository) *MemberService {
	return &MemberService{
		repo:      repo,
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}
