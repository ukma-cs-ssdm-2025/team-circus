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

type Document struct {
	UUID      uuid.UUID
	GroupUUID uuid.UUID
	Name      string
	Content   string
	CreatedAt time.Time
}

type User struct {
	UUID      uuid.UUID
	Login     string
	Email     string
	Password  string
	CreatedAt time.Time
}

var (
	ErrGroupNotFound    = errors.New("group not found")
	ErrDocumentNotFound = errors.New("document not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrInternal         = errors.New("internal error")
	ErrLoginTaken       = errors.New("login already taken")
	ErrEmailTaken       = errors.New("email already taken")
	ErrUserExists       = errors.New("user already exists")
)
