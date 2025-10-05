package responses

import (
	"time"

	"github.com/google/uuid"
)

type GetGroupResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type GetAllGroupsResponse struct {
	Groups []GetGroupResponse `json:"groups"`
}
