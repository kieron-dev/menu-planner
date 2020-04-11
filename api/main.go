package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	googleTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/authGoogle", authGoogle).Methods("POST", "OPTIONS")
	r.Use(mux.CORSMethodMiddleware(r))

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

type AuthResp struct {
	Token string
}

type GoogleAuthReq struct {
	TokenID string `json:"tokenID"`
}

type Claims struct {
	Name string
	jwt.StandardClaims
}

var jwtKey = []byte("a very secret string")

func authGoogle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read-body: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	var authReq GoogleAuthReq
	err = json.Unmarshal(body, &authReq)
	if err != nil {
		log.Printf("json unmarshal: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	v := googleTokenVerifier.Verifier{}
	aud := "176462381984-bfq3v9mc00v0ipvpebiaiide4l22dmoh.apps.googleusercontent.com"
	err = v.VerifyIDToken(authReq.TokenID, []string{
		aud,
	})
	if err != nil {
		log.Printf("token verifier: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	claimSet, err := googleTokenVerifier.Decode(authReq.TokenID)
	if err != nil {
		log.Printf("token-verifier-decode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Name: claimSet.Name,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("jwt-token-sign: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var resp AuthResp
	resp.Token = tokenString

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json-marshall: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
}
