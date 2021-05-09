package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
	"github.com/gorilla/securecookie"
	"github.com/kieron-pivotal/menu-planner-app/db"
	"github.com/kieron-pivotal/menu-planner-app/handlers"
	"github.com/kieron-pivotal/menu-planner-app/jwt"
	"github.com/kieron-pivotal/menu-planner-app/routing"
	"github.com/kieron-pivotal/menu-planner-app/session"
	_ "github.com/lib/pq"
)

// TODO: get all these from env
const (
	aud    = "176462381984-bfq3v9mc00v0ipvpebiaiide4l22dmoh.apps.googleusercontent.com"
	webURI = "http://localhost:3000"
	port   = 8080
)

var (
	sign    = securecookie.GenerateRandomKey(32)
	encrypt = securecookie.GenerateRandomKey(32)
)

func main() {
	googleVerifier := new(googleAuthIDTokenVerifier.Verifier)
	jwtDecoder := jwt.NewJWT()

	connStr := mustGetEnv("DB_CONN_STR")
	pg, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	userStore := db.NewUserStore(pg)
	recipeStore := db.NewRecipeStore(pg)

	sessionManager := session.NewManager([][]byte{sign, encrypt})
	authHandler := handlers.NewAuthHandler(aud, googleVerifier, jwtDecoder, userStore, sessionManager)
	recipeHandler := handlers.NewRecipeHandler(sessionManager, recipeStore)
	routes := routing.New(webURI, sessionManager, authHandler, recipeHandler)
	r := routes.SetupRoutes()

	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(port), r))
}

func mustGetEnv(v string) string {
	s := os.Getenv(v)
	if s != "" {
		return s
	}
	panic(fmt.Sprintf("env var %q not set", v))
}
