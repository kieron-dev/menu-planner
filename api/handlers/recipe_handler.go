package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kieron-pivotal/menu-planner-app/models"
)

//counterfeiter:generate . RecipeStore

type RecipeStore interface {
	List(userID int) ([]models.Recipe, error)
}

type RecipeHandler struct {
	sessionManager SessionManager
	recipeStore    RecipeStore
}

func NewRecipeHandler(sessionManager SessionManager, recipeStore RecipeStore) *RecipeHandler {
	return &RecipeHandler{
		sessionManager: sessionManager,
		recipeStore:    recipeStore,
	}
}

func (h *RecipeHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	sess, err := h.sessionManager.Get(r.Context())
	if err != nil || sess == nil || !sess.IsLoggedIn {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)

		return
	}

	recipes, err := h.recipeStore.List(sess.ID)
	if err != nil {
		// TODO: handle err
	}

	w.Header().Add("Content-Type", "application/json")

	type Recipe struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}

	list := []Recipe{}

	for _, r := range recipes {
		list = append(list, Recipe{Name: r.Name(), ID: r.ID()})
	}

	if err = json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, "json encoding failure", http.StatusInternalServerError)

		return
	}
}
