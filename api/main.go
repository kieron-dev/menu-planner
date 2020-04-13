package main

import (
	"log"
	"net/http"

	googleTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
	"github.com/kieron-pivotal/menu-planner-app/auth"
	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/jwt"
	"github.com/kieron-pivotal/menu-planner-app/routing"
)

// TODO: get from env
const aud = "176462381984-bfq3v9mc00v0ipvpebiaiide4l22dmoh.apps.googleusercontent.com"
const webURI = "http://localhost:3000"

func main() {
	googleVerifier := new(googleTokenVerifier.Verifier)
	jwtDecoder := jwt.NewJWT()
	localAuth := auth.NewLocalAuth()
	handlers := handlers.New(aud, googleVerifier, jwtDecoder, localAuth)
	routes := routing.New(webURI, handlers)
	r := routes.SetupRoutes()

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
