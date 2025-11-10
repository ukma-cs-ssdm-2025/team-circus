package requests

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateUserRequest struct {
	Login    *string `json:"login,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (r UpdateUserRequest) Validate() error {
	if r.Login == nil && r.Email == nil && r.Password == nil {
		return errors.New("at least one field must be provided")
	}

	return validation.ValidateStruct(&r,
		validation.Field(&r.Login, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Email, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Password, validation.NilOrNotEmpty, validation.Length(1, 255)),
	)
}
