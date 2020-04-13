package routing

import (
	"net/http"

	"github.com/gorilla/mux"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Handlers

type Handlers interface {
	AuthGoogle(w http.ResponseWriter, r *http.Request)
}

type Routes struct {
	frontendURI string
	handlers    Handlers
}

func New(frontendURI string, handlers Handlers) Routes {
	return Routes{
		frontendURI: frontendURI,
		handlers:    handlers,
	}
}

func (r Routes) SetupRoutes() *mux.Router {
	m := mux.NewRouter()

	m.HandleFunc("/authGoogle", r.handlers.AuthGoogle).Methods("POST", "OPTIONS")
	m.Use(mux.CORSMethodMiddleware(m))
	m.Use(r.CORSOriginMiddleware)

	return m
}

func (r Routes) CORSOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.frontendURI)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if req.Method != http.MethodOptions {
			next.ServeHTTP(w, req)
		}
	})
}
