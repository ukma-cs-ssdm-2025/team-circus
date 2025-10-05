package domain

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	UUID      uuid.UUID
	Name      string
	CreatedAt time.Time
}
