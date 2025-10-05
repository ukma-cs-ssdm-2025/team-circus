package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	UUID      uuid.UUID
	Name      string
	CreatedAt time.Time
}

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrInternal      = errors.New("internal error")
)
