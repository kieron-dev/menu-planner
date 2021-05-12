package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/kieron-pivotal/menu-planner-app/models"
)

//counterfeiter:generate . RecipeStore

type RecipeStore interface {
	List(userID int) ([]models.Recipe, error)
	Insert(recipe models.Recipe) (models.Recipe, error)
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

	list := []models.Recipe{}

	for _, r := range recipes {
		list = append(list, models.Recipe{Name: r.Name, ID: r.ID})
	}

	if err = json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, "json encoding failure", http.StatusInternalServerError)

		return
	}
}

func (h *RecipeHandler) NewRecipe(w http.ResponseWriter, r *http.Request) {
	sess, err := h.sessionManager.Get(r.Context())
	if err != nil || sess == nil || !sess.IsLoggedIn {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)

		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)

		return
	}

	recipe := models.Recipe{}
	if err = json.Unmarshal(body, &recipe); err != nil {
		http.Error(w, "", http.StatusBadRequest)

		return
	}

	recipe.UserID = sess.ID

	recipe, err = h.recipeStore.Insert(recipe)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}
