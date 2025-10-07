package responses

import (
	"time"

	"github.com/google/uuid"
)

type UpdateUserResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
