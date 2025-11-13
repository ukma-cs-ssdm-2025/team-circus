package responses

import (
	"time"

	"github.com/google/uuid"
)

type GroupMemberResponse struct {
	GroupUUID uuid.UUID `json:"group_uuid"`
	UserUUID  uuid.UUID `json:"user_uuid"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UserLogin string    `json:"user_login"`
	UserEmail string    `json:"user_email"`
}

type GroupMembersResponse struct {
	Members []GroupMemberResponse `json:"members"`
}
