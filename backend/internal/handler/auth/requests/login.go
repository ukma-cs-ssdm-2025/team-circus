package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type LogInRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r LogInRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Login, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Password, validation.Required, validation.Length(1, 255)),
	)
}
