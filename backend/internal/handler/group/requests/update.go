package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateGroupRequest struct {
	Name string `json:"name"`
}

func (r UpdateGroupRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
	)
}
