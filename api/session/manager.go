/* Package session handles everything to do with web sessions */
package session

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

type sessionKey int

const (
	ctxSessionKey     sessionKey = 0
	sessionCookieName            = "_id"
	authInfoKey                  = "authInfo"
)

type Manager struct {
	sessionStore sessions.Store
}

type AuthInfo struct {
	Name       string
	ID         int
	IsLoggedIn bool
}

func init() {
	gob.Register(&AuthInfo{})
}

func NewManager(sessionStoreKeys [][]byte) *Manager {
	store := createSessionStore(sessionStoreKeys)
	return &Manager{
		sessionStore: store,
	}
}

func (m *Manager) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := m.sessionStore.Get(r, sessionCookieName)
		if err != nil {
			delCookie := http.Cookie{
				Name:   sessionCookieName,
				MaxAge: -1,
			}
			http.SetCookie(w, &delCookie)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		val := session.Values[authInfoKey]
		if val != nil {
			var ourSession *AuthInfo
			var ok bool
			if ourSession, ok = val.(*AuthInfo); !ok {
				log.Printf("not a *AuthInfo, actually a %T", val)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r = r.Clone(context.WithValue(r.Context(), ctxSessionKey, ourSession))

			session.Save(r, w)
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Manager) Set(r *http.Request, w http.ResponseWriter, authInfo *AuthInfo) error {
	session, err := m.sessionStore.Get(r, sessionCookieName)
	if err != nil {
		return fmt.Errorf("session-set: failed to get session %w", err)
	}
	session.Values[authInfoKey] = authInfo

	if err = session.Save(r, w); err != nil {
		return fmt.Errorf("session-set: failed to save session %w", err)
	}

	return nil
}

func (m *Manager) Get(ctx context.Context) (*AuthInfo, error) {
	val := ctx.Value(ctxSessionKey)
	if val == nil {
		return nil, nil
	}

	authInfo, ok := val.(*AuthInfo)
	if !ok {
		return nil, fmt.Errorf("wrong type %T for session, expected *session.AuthInfo", val)
	}
	return authInfo, nil
}

func createSessionStore(sessionStoreKeys [][]byte) sessions.Store {
	store := sessions.NewFilesystemStore("", sessionStoreKeys...)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.MaxAge = 60 * 15
	return store
}
