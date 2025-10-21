package seeder

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/domain"
)

type User struct {
	UUID           string
	Login          string
	Email          string
	HashedPassword string
	CreatedAt      time.Time

	db *sql.DB
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		UUID:      uuid.MustParse(u.UUID),
		Login:     u.Login,
		Email:     u.Email,
		Password:  u.HashedPassword,
		CreatedAt: u.CreatedAt,
	}
}

func newUser(db *sql.DB) *User {
	return &User{
		UUID:           uuid.New().String(),
		Login:          "login",
		Email:          "email",
		HashedPassword: "hashedPassword",
		CreatedAt:      time.Now(),

		db: db,
	}
}

func (u *User) WithLogin(login string) *User {
	u.Login = login
	return u
}

func (u *User) WithEmail(email string) *User {
	u.Email = email
	return u
}

func (u *User) WithHashedPassword(hashedPassword string) *User {
	u.HashedPassword = hashedPassword
	return u
}

func (u *User) WithCreatedAt(createdAt time.Time) *User {
	u.CreatedAt = createdAt
	return u
}

func (u *User) WithUUID(uuid string) *User {
	u.UUID = uuid
	return u
}

func (u *User) Create() error {
	query := `
		INSERT INTO users (uuid, login, email, hashed_password, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uuid, login, email, hashed_password, created_at`

	err := u.db.QueryRow(query, u.UUID, u.Login, u.Email, u.HashedPassword, u.CreatedAt).Scan( //nolint:noctx
		&u.UUID,
		&u.Login,
		&u.Email,
		&u.HashedPassword,
		&u.CreatedAt,
	)
	if err != nil {
		return errors.Join(domain.ErrInternal, fmt.Errorf("user seeder: create: %w", err))
	}

	return nil
}

func getUserByLogin(db *sql.DB, login string) (*User, error) {
	var user User
	query := `
		SELECT uuid, login, email, hashed_password, created_at
		FROM users
		WHERE login = $1`

	err := db.QueryRow(query, login).Scan( //nolint:noctx
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user seeder: getByLogin: %w", err))
	}

	return &user, nil
}

func getUserByEmail(db *sql.DB, email string) (*User, error) {
	var user User
	query := `
		SELECT uuid, login, email, hashed_password, created_at
		FROM users
		WHERE email = $1`

	err := db.QueryRow(query, email).Scan( //nolint:noctx
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user seeder: getByEmail: %w", err))
	}

	return &user, nil
}

func getUserByUUID(db *sql.DB, uuid string) (*User, error) {
	var user User
	query := `
		SELECT uuid, login, email, hashed_password, created_at
		FROM users
		WHERE uuid = $1`

	err := db.QueryRow(query, uuid).Scan( //nolint:noctx
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, errors.Join(domain.ErrInternal, fmt.Errorf("user seeder: getByUUID: %w", err))
	}

	return &user, nil
}
