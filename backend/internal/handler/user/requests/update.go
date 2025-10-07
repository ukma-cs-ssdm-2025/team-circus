package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateUserRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Login, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Email, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Password, validation.Required, validation.Length(1, 255)),
	)
}
