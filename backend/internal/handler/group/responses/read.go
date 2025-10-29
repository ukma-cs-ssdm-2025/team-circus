package responses

import (
	"time"

	"github.com/google/uuid"
)

type GetGroupResponse struct {
	UUID       uuid.UUID `json:"uuid"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	AuthorUUID uuid.UUID `json:"author_uuid"`
	Role       string    `json:"role,omitempty"`
}

type GetAllGroupsResponse struct {
	Groups []GetGroupResponse `json:"groups"`
}
