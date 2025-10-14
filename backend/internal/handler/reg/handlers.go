package reg

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/handler/reg/responses"
)

func mapUserToRegResponse(user *domain.User) responses.RegResponse {
	return responses.RegResponse{
		UUID:      user.UUID,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
