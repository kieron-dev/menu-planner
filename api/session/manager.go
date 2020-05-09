package session

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

const sessionCookieName = "_id"

type sessionKey int

const ctxSessionKey sessionKey = 0
const valuesSessionKey = "session"

type Manager struct {
	sessionStore sessions.Store
}

type Session struct {
	Name       string
	ID         int
	IsLoggedIn bool
}

func init() {
	gob.Register(&Session{})
}

func createSessionStore(sessionStoreKeys [][]byte) sessions.Store {
	store := sessions.NewFilesystemStore("", sessionStoreKeys...)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.MaxAge = 60 * 15
	return store
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
		// session.Save(r, w)

		val := session.Values[valuesSessionKey]
		if val != nil {
			var ourSession *Session
			var ok bool
			if ourSession, ok = val.(*Session); !ok {
				log.Printf("not a *Session, actually a %T", val)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r = r.Clone(context.WithValue(r.Context(), ctxSessionKey, ourSession))
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Manager) Set(r *http.Request, w http.ResponseWriter, s *Session) error {
	session, err := m.sessionStore.Get(r, sessionCookieName)
	if err != nil {
		return fmt.Errorf("session-set: failed to get session %w", err)
	}
	session.Values[valuesSessionKey] = s

	if err = session.Save(r, w); err != nil {
		return fmt.Errorf("session-set: failed to save session %w", err)
	}

	return nil
}

func (m *Manager) Get(ctx context.Context) (*Session, error) {
	val := ctx.Value(ctxSessionKey)
	if val == nil {
		return nil, nil
	}

	sesh, ok := val.(*Session)
	if !ok {
		return nil, fmt.Errorf("wrong type %T for session, expected *session.Session",
			ctx.Value(ctxSessionKey))
	}
	return sesh, nil
}
