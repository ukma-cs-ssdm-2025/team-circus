package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateDocumentRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (r UpdateDocumentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Content, validation.Required),
	)
}
