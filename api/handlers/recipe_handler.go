package handlers

import (
	"encoding/json"
	"net/http"
)

type RecipeHandler struct {
	sessionManager SessionManager
}

func NewRecipeHandler(sessionManager SessionManager) *RecipeHandler {
	return &RecipeHandler{
		sessionManager: sessionManager,
	}
}

func (h *RecipeHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	sess, err := h.sessionManager.Get(r.Context())
	if err != nil || sess == nil || !sess.IsLoggedIn {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)

		return
	}

	w.Header().Add("Content-Type", "application/json")

	type Recipe struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	recipes := []Recipe{
		{Name: "Bangers and Mash", ID: "1"},
		{Name: "Fish and Chips", ID: "2"},
		{Name: "Creamy Salmon Pasta", ID: "3"},
	}

	if err = json.NewEncoder(w).Encode(recipes); err != nil {
		http.Error(w, "json encoding failure", http.StatusInternalServerError)

		return
	}
}
