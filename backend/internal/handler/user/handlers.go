package user

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/user/responses"
)

func mapUserToCreateResponse(user *domain.User) responses.CreateUserResponse {
	return responses.CreateUserResponse{
		UUID:      user.UUID,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func mapUserToGetResponse(user *domain.User) responses.GetUserResponse {
	return responses.GetUserResponse{
		UUID:      user.UUID,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func mapUserToUpdateResponse(user *domain.User) responses.UpdateUserResponse {
	return responses.UpdateUserResponse{
		UUID:      user.UUID,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func mapUsersToGetAllResponse(users []*domain.User) []responses.GetUserResponse {
	result := make([]responses.GetUserResponse, len(users))
	for i, user := range users {
		result[i] = mapUserToGetResponse(user)
	}
	return result
}
