package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

type CreateMemberRequest struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	Role     string    `json:"role"`
}

func (r CreateMemberRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserUUID, validation.Required),
		validation.Field(&r.Role, validation.Required, validation.Length(1, 50),
			validation.In(
				domain.RoleAuthor,
				domain.RoleEditor,
				domain.RoleViewer,
			),
		),
	)
}
