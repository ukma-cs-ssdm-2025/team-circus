package responses

import (
	"time"

	"github.com/google/uuid"
)

type GetUserResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type GetAllUsersResponse struct {
	Users []GetUserResponse `json:"users"`
}
