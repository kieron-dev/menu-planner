package routing

import (
	"net/http"

	"github.com/gorilla/mux"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . AuthHandler

type AuthHandler interface {
	AuthGoogle(w http.ResponseWriter, r *http.Request)
	WhoAmI(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

//counterfeiter:generate . RecipeHandler

type RecipeHandler interface {
	GetRecipes(w http.ResponseWriter, r *http.Request)
}

//counterfeiter:generate . SessionManager

type SessionManager interface {
	SessionMiddleware(next http.Handler) http.Handler
}

type Routes struct {
	frontendURI    string
	sessionManager SessionManager
	authHandler    AuthHandler
	recipeHandler  RecipeHandler
}

func New(
	frontendURI string, sessionManager SessionManager,
	authHandler AuthHandler, recipeHandler RecipeHandler) Routes {
	return Routes{
		frontendURI:    frontendURI,
		sessionManager: sessionManager,
		authHandler:    authHandler,
		recipeHandler:  recipeHandler,
	}
}

func (r Routes) SetupRoutes() *mux.Router {
	m := mux.NewRouter()

	m.HandleFunc("/authGoogle", r.authHandler.AuthGoogle).Methods("POST", "OPTIONS")
	m.HandleFunc("/whoami", r.authHandler.WhoAmI).Methods("GET", "OPTIONS")
	m.HandleFunc("/logout", r.authHandler.Logout).Methods("POST", "OPTIONS")
	m.HandleFunc("/recipes", r.recipeHandler.GetRecipes).Methods("GET", "OPTIONS")
	m.Use(mux.CORSMethodMiddleware(m))
	m.Use(r.CORSOriginMiddleware)
	m.Use(r.sessionManager.SessionMiddleware)

	return m
}

func (r Routes) CORSOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.frontendURI)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if req.Method != http.MethodOptions {
			next.ServeHTTP(w, req)
		}
	})
}
