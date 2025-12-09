package responses

import (
	"time"

	"github.com/google/uuid"
)

type ShareDocumentResponse struct {
	DocumentUUID uuid.UUID `json:"document_uuid"`
	URL          string    `json:"url"`
	ExpiresAt    time.Time `json:"expires_at"`
}
