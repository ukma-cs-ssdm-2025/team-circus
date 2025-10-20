package seeder

import "database/sql"

type Seeder struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{
		db: db,
	}
}

func (s *Seeder) NewUser() *User {
	return newUser(s.db)
}

func (s *Seeder) GetUserByLogin(login string) (*User, error) {
	return getUserByLogin(s.db, login)
}

func (s *Seeder) GetUserByEmail(email string) (*User, error) {
	return getUserByEmail(s.db, email)
}

func (s *Seeder) GetUserByUUID(uuid string) (*User, error) {
	return getUserByUUID(s.db, uuid)
}
