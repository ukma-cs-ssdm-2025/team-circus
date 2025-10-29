package requests

import (
	"errors"

	"github.com/google/uuid"
)

type AddMemberRequest struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	Role     string    `json:"role"`
}

func (r AddMemberRequest) Validate() error {
	if r.UserUUID == uuid.Nil {
		return errors.New("user_uuid is required")
	}
	if r.Role == "" {
		return errors.New("role is required")
	}
	return nil
}

type UpdateMemberRequest struct {
	Role string `json:"role"`
}

func (r UpdateMemberRequest) Validate() error {
	if r.Role == "" {
		return errors.New("role is required")
	}
	return nil
}
