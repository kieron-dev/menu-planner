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
		recipe := models.Recipe{}
		rows.Scan(&recipe.ID, &recipe.Name)

		res = append(res, recipe)
	}

	return res, nil
}

func (s *RecipeStore) Insert(recipe models.Recipe) (models.Recipe, error) {
	row := s.sqlDB.QueryRow(`INSERT INTO recipe
    (name, user_id)
    VALUES ($1, $2)
    RETURNING (id)`, recipe.Name, recipe.UserID)

	if err := row.Scan(&recipe.ID); err != nil {
		return models.Recipe{}, fmt.Errorf("insert failed: %w", err)
	}

	return recipe, nil
}
