package responses

import (
	"time"

	"github.com/google/uuid"
)

type GetMemberResponse struct {
	GroupUUID uuid.UUID `json:"group_uuid"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type GetAllMembersResponse struct {
	Members []GetMemberResponse `json:"members"`
}
