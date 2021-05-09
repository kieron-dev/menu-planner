package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kieron-pivotal/menu-planner-app/models"
)

var errNotFound = errors.New("no matching row found")

func NotFoundErr() error {
	return errNotFound
}

type UserStore struct {
	sqlDB DB
}

func NewUserStore(sqlDB DB) *UserStore {
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

func (u User) ID() int {
	return u.id
}

func (s *UserStore) IsNotFoundErr(err error) bool {
	return err == errNotFound
}

func (s *UserStore) FindByEmail(email string) (models.User, error) {
	var e, name string
	var id int
	err := s.sqlDB.QueryRow(`
SELECT id, email, name
FROM local_user
WHERE email = $1
`, email).Scan(&id, &e, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, errNotFound
		}
		return User{}, fmt.Errorf("find-by-email failed %w", err)
	}
	return User{
		email: e,
		name:  name,
		id:    id,
	}, nil
}

func (s *UserStore) Create(email, name string) (models.User, error) {
	var id int
	err := s.sqlDB.QueryRow(`
INSERT INTO local_user (email, name)
VALUES ($1, $2)
RETURNING id`, email, name).Scan(&id)
	if err != nil {
		return User{}, fmt.Errorf("create-user failed %w", err)
	}

	return User{
		id:    id,
		email: email,
		name:  name,
	}, nil
}
