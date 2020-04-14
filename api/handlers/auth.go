package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . TokenVerifier

type TokenVerifier interface {
	VerifyIDToken(token string, audience []string) error
}

//counterfeiter:generate . JWTDecoder

type JWTDecoder interface {
	ClaimSet(token string) (map[string]string, error)
}

//counterfeiter:generate . LocalAuther

type LocalAuther interface {
	LocalAuth(email, name string) (token string, err error)
}

type Handlers struct {
	audience      string
	tokenVerifier TokenVerifier
	jwtDecoder    JWTDecoder
	localAuther   LocalAuther
}

func New(audience string, tokenVerifier TokenVerifier, jwtDecoder JWTDecoder, localAuther LocalAuther) *Handlers {
	return &Handlers{
		audience:      audience,
		tokenVerifier: tokenVerifier,
		jwtDecoder:    jwtDecoder,
		localAuther:   localAuther,
	}
}

func (h *Handlers) AuthGoogle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read-body: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var authReq struct {
		TokenID string `json:"tokenID"`
	}

	err = json.Unmarshal(body, &authReq)
	if err != nil {
		log.Printf("json unmarshal: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("done writing the header\n")
		return
	}

	if err = h.tokenVerifier.VerifyIDToken(authReq.TokenID, []string{h.audience}); err != nil {
		log.Printf("token-verifier: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claimSet, err := h.jwtDecoder.ClaimSet(authReq.TokenID)
	if err != nil {
		log.Printf("jwt-decoder: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email, ok := claimSet["email"]
	if !ok {
		log.Printf("email-missing: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name, ok := claimSet["name"]
	if !ok {
		log.Printf("name-missing: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString, err := h.localAuther.LocalAuth(email, name)
	if err != nil {
		log.Printf("local-auth: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var respObj struct {
		Token string `json:"token"`
	}
	respObj.Token = tokenString
	respBytes, err := json.Marshal(respObj)

	if err != nil {
		log.Printf("marshal-response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
}
