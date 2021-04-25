package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kieron-pivotal/menu-planner-app/models"
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

//counterfeiter:generate . UserStore

type UserStore interface {
	IsNotFoundErr(error) bool
	FindByEmail(email string) (models.User, error)
	Create(email, name string) (models.User, error)
}

//counterfeiter:generate . SessionManager

type SessionManager interface {
	Get(ctx context.Context) (*session.AuthInfo, error)
	Set(r *http.Request, w http.ResponseWriter, s *session.AuthInfo) error
}

type AuthHandler struct {
	audience       string
	tokenVerifier  TokenVerifier
	jwtDecoder     JWTDecoder
	userStore      UserStore
	sessionManager SessionManager
}

func NewAuthHandler(
	audience string,
	tokenVerifier TokenVerifier,
	jwtDecoder JWTDecoder,
	userStore UserStore,
	sessionSetter SessionManager,
) *AuthHandler {
	return &AuthHandler{
		audience:       audience,
		tokenVerifier:  tokenVerifier,
		jwtDecoder:     jwtDecoder,
		userStore:      userStore,
		sessionManager: sessionSetter,
	}
}

func (h *AuthHandler) AuthGoogle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read-body: %v\n", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var authReq struct {
		IDToken string `json:"idToken"`
	}

	err = json.Unmarshal(body, &authReq)
	if err != nil {
		log.Printf("json unmarshal: %v\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err = h.tokenVerifier.VerifyIDToken(authReq.IDToken, []string{h.audience}); err != nil {
		log.Printf("token-verifier: %v\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	claimSet, err := h.jwtDecoder.ClaimSet(authReq.IDToken)
	if err != nil {
		log.Printf("jwt-decoder: %v\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	var email, name string
	if email, err = extractString(claimSet, "email"); err != nil {
		log.Printf("extract-string: %v", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	user, err := h.userStore.FindByEmail(email)
	if err != nil {
		if !h.userStore.IsNotFoundErr(err) {
			log.Printf("user-store-find-by-email: %v\n", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if name, err = extractString(claimSet, "name"); err != nil {
			log.Printf("extract-string: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err = h.userStore.Create(email, name)
		if err != nil {
			log.Printf("user-store-create: %v\n", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	sess := session.AuthInfo{
		IsLoggedIn: true,
		ID:         user.ID(),
		Name:       user.Name(),
	}

	if err := h.sessionManager.Set(r, w, &sess); err != nil {
		log.Printf("failed-to-set-session: %v\n", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"name": "%s"}`, user.Name())
}

func (h *AuthHandler) WhoAmI(w http.ResponseWriter, r *http.Request) {
	sess, err := h.sessionManager.Get(r.Context())
	if err != nil || sess == nil || !sess.IsLoggedIn {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, `{"name": "%s"}`, sess.Name)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := h.sessionManager.Get(r.Context())
	if err != nil || sess == nil || !sess.IsLoggedIn {
		return
	}
	sess.IsLoggedIn = false
	if err = h.sessionManager.Set(r, w, sess); err != nil {
		log.Printf("failed-to-set-session: %v\n", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
	fmt.Fprint(w, "logged out")
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
