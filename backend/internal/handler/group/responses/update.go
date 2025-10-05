package responses

import (
	"time"

	"github.com/google/uuid"
)

type UpdateGroupResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
