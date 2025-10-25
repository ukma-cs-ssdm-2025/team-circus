package testdb

import (
	"database/sql"
	"errors"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=mcd sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ResetDB(db *sql.DB) error {
	var err error
	var errs []error

	_, err = db.Exec("DELETE FROM user_groups") //nolint:noctx
	errs = append(errs, err)
	_, err = db.Exec("DELETE FROM documents") //nolint:noctx
	errs = append(errs, err)
	_, err = db.Exec("DELETE FROM groups") //nolint:noctx
	errs = append(errs, err)
	_, err = db.Exec("DELETE FROM users") //nolint:noctx
	errs = append(errs, err)

	resetErr := errors.Join(errs...)
	return resetErr
}
