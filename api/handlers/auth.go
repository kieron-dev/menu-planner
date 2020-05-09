package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kieron-pivotal/menu-planner-app/db"
	"github.com/kieron-pivotal/menu-planner-app/session"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . TokenVerifier

type TokenVerifier interface {
	VerifyIDToken(token string, audience []string) error
}

//counterfeiter:generate . JWTDecoder

type JWTDecoder interface {
	ClaimSet(token string) (map[string]interface{}, error)
}

//counterfeiter:generate . User

type User interface {
	Email() string
	Name() string
	ID() int
}

//counterfeiter:generate . UserStore

type UserStore interface {
	FindByEmail(email string) (User, error)
	Create(email, name string) (User, error)
}

//counterfeiter:generate . SessionSetter

type SessionSetter interface {
	Set(r *http.Request, w http.ResponseWriter, s *session.Session) error
}

type AuthHandler struct {
	audience      string
	tokenVerifier TokenVerifier
	jwtDecoder    JWTDecoder
	userStore     UserStore
	sessionSetter SessionSetter
}

func New(
	audience string,
	tokenVerifier TokenVerifier,
	jwtDecoder JWTDecoder,
	userStore UserStore,
	sessionSetter SessionSetter,
) *AuthHandler {

	return &AuthHandler{
		audience:      audience,
		tokenVerifier: tokenVerifier,
		jwtDecoder:    jwtDecoder,
		userStore:     userStore,
		sessionSetter: sessionSetter,
	}
}

func (h *AuthHandler) AuthGoogle(w http.ResponseWriter, r *http.Request) {
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

	var email, name string
	if email, err = extractString(claimSet, "email"); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userStore.FindByEmail(email)
	if err != nil {
		if !db.IsNotFoundErr(err) {
			log.Printf("user-store-find-by-email failed: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if name, err = extractString(claimSet, "name"); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err = h.userStore.Create(email, name)
		if err != nil {
			log.Printf("user-store-create failed: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	sess := session.Session{
		IsLoggedIn: true,
		ID:         user.ID(),
		Name:       user.Name(),
	}

	if err := h.sessionSetter.Set(r, w, &sess); err != nil {
		log.Printf("failed-to-set-session: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func extractString(claimSet map[string]interface{}, key string) (string, error) {
	val, ok := claimSet[key]
	if !ok {
		log.Printf("missing key %q\n", key)
		return "", fmt.Errorf("key %q not in claimSet", key)
	}

	valStr, ok := val.(string)
	if !ok {
		log.Printf("%q not-a-string\n", key)
		return "", fmt.Errorf("%q is a %t - expected string", val, val)
	}
	return valStr, nil
}
