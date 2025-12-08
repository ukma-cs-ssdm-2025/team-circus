package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

type UpdateMemberRequest struct {
	Role string `json:"role"`
}

func (r UpdateMemberRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Role, validation.Required, validation.Length(1, 50),
			validation.In(
				domain.RoleAuthor,
				domain.RoleEditor,
				domain.RoleViewer,
			)),
	)
}
