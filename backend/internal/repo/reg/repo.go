package reg

import (
	"database/sql"
)

type RegRepository struct {
	db *sql.DB
}

func NewRegRepository(db *sql.DB) *RegRepository {
	return &RegRepository{
		db: db,
	}
}
