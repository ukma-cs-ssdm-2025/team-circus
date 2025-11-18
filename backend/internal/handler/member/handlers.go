package member

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/member/responses"
)

func mapMemberToResponse(member *domain.Member) responses.GetMemberResponse {
	return responses.GetMemberResponse{
		GroupUUID: member.GroupUUID,
		UserUUID:  member.UserUUID,
		Role:      member.Role,
		CreatedAt: member.CreatedAt,
	}
}

func mapMembersToResponse(members []*domain.Member) []responses.GetMemberResponse {
	result := make([]responses.GetMemberResponse, len(members))
	for i, member := range members {
		result[i] = mapMemberToResponse(member)
	}
	return result
}
