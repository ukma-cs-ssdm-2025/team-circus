package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	UUID       uuid.UUID
	Name       string
	CreatedAt  time.Time
	AuthorUUID uuid.UUID
	Role       string
}

type GroupMember struct {
	GroupUUID uuid.UUID
	UserUUID  uuid.UUID
	Role      string
	CreatedAt time.Time
	UserLogin string
	UserEmail string
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
	ErrForbidden        = errors.New("forbidden")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvalidRole      = errors.New("invalid role")
	ErrLastAuthor       = errors.New("cannot remove last author from group")
)

const (
	GroupRoleAuthor   = "author"
	GroupRoleCoAuthor = "coauthor"
	GroupRoleReviewer = "reviewer"
)
