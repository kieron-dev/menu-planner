package db

import (
	"database/sql"
	"fmt"

	"github.com/kieron-pivotal/menu-planner-app/models"
)

type DB interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
}

type RecipeStore struct {
	sqlDB DB
}

func NewRecipeStore(sqlDB DB) *RecipeStore {
	return &RecipeStore{
		sqlDB: sqlDB,
	}
}

type Recipe struct {
	id   int
	name string
}

func (r Recipe) Name() string {
	return r.name
}

func (r Recipe) ID() int {
	return r.id
}

func (s *RecipeStore) IsNotFoundErr(err error) bool {
	return err == errNotFound
}

func (s *RecipeStore) List(userID int) ([]models.Recipe, error) {
	res := []models.Recipe{}

	rows, err := s.sqlDB.Query(`
SELECT id, name
FROM recipe
WHERE user_id = $1
`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, errNotFound
		}
		return res, fmt.Errorf("list-recipes failed %w", err)
	}

	for rows.Next() {
		recipe := Recipe{}
		rows.Scan(&recipe.id, &recipe.name)

		res = append(res, recipe)
	}

	return res, nil
}
