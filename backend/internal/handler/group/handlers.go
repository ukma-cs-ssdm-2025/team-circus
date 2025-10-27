package group

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/group/responses"
)

func mapGroupToCreateResponse(group *domain.Group) responses.CreateGroupResponse {
	return responses.CreateGroupResponse{
		UUID:       group.UUID,
		Name:       group.Name,
		CreatedAt:  group.CreatedAt,
		AuthorUUID: group.AuthorUUID,
		Role:       group.Role,
	}
}

func mapGroupToGetResponse(group *domain.Group) responses.GetGroupResponse {
	return responses.GetGroupResponse{
		UUID:       group.UUID,
		Name:       group.Name,
		CreatedAt:  group.CreatedAt,
		AuthorUUID: group.AuthorUUID,
		Role:       group.Role,
	}
}

func mapGroupToUpdateResponse(group *domain.Group) responses.UpdateGroupResponse {
	return responses.UpdateGroupResponse{
		UUID:       group.UUID,
		Name:       group.Name,
		CreatedAt:  group.CreatedAt,
		AuthorUUID: group.AuthorUUID,
		Role:       group.Role,
	}
}

func mapGroupsToGetAllResponse(groups []*domain.Group) []responses.GetGroupResponse {
	result := make([]responses.GetGroupResponse, len(groups))
	for i, group := range groups {
		result[i] = mapGroupToGetResponse(group)
	}
	return result
}
