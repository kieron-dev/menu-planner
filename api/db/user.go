package db

import (
	"database/sql"
	"errors"
	"fmt"
)

var notFoundErr = errors.New("no matching row found")

func NotFoundErr() error {
	return notFoundErr
}

func IsNotFoundErr(err error) bool {
	return err == notFoundErr
}

type UserStore struct {
	sqlDB *sql.DB
}

func NewUserStore(sqlDB *sql.DB) *UserStore {
	return &UserStore{
		sqlDB: sqlDB,
	}
}

type User struct {
	id    int
	email string
	name  string
	lid   []uint8
}

func (u User) Email() string {
	return u.email
}

func (u User) Name() string {
	return u.name
}

func (u User) Id() []uint8 {
	return u.lid
}

func (s *UserStore) FindByEmail(email string) (User, error) {
	var e, name string
	var id []uint8
	err := s.sqlDB.QueryRow(`
SELECT email, name, lid
FROM local_user
WHERE email = $1
`, email).Scan(&e, &name, &id)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, notFoundErr
		}
		return User{}, fmt.Errorf("find-by-email failed %w", err)
	}
	return User{
		email: e,
		name:  name,
		lid:   id,
	}, nil
}

func (s *UserStore) Create(email, name string) (User, error) {
	var uuid []uint8
	err := s.sqlDB.QueryRow(`
INSERT INTO local_user (email, name)
VALUES ($1, $2)
RETURNING lid`, email, name).Scan(&uuid)
	if err != nil {
		return User{}, fmt.Errorf("create-user failed %w", err)
	}

	return User{
		email: email,
		name:  name,
		lid:   uuid,
	}, nil
}
