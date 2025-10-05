package responses

import (
	"time"

	"github.com/google/uuid"
)

type GetDocumentResponse struct {
	UUID      uuid.UUID `json:"uuid"`
	GroupUUID uuid.UUID `json:"group_uuid"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type GetAllDocumentsResponse struct {
	Documents []GetDocumentResponse `json:"documents"`
}

type GetDocumentsByGroupResponse struct {
	Documents []GetDocumentResponse `json:"documents"`
}
