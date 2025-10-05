package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type CreateDocumentRequest struct {
	GroupUUID uuid.UUID `json:"group_uuid"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
}

func (r CreateDocumentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.GroupUUID, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Content, validation.Required),
	)
}
