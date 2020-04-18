package auth

import "fmt"

//counterfeiter:generate . User

type User interface {
	Email() string
	Name() string
	ID() string
}

//counterfeiter:generate . UserStore

type UserStore interface {
	FindByEmail(email string) (User, error)
	Create(email, name string) (User, error)
}

//counterfeiter:generate . JWTGenerator

type JWTGenerator interface {
	GenerateToken(id, name string) (string, error)
}

type LocalAuth struct {
	userStore    UserStore
	jwtGenerator JWTGenerator
}

func NewLocalAuth(userStore UserStore, jwtGen JWTGenerator) *LocalAuth {
	return &LocalAuth{
		userStore:    userStore,
		jwtGenerator: jwtGen,
	}
}

// LocalAuth takes email and name from a trusted 3rd-party authenticator, e.g.
// Google Sign-in, and produces a JWT for this system.
//
// It will return an ID from the database for the email, if a matching entry
// exists, or will create a new entry in the database for the user.
func (a *LocalAuth) LocalAuth(email, name string) (string, error) {
	user, err := a.userStore.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("local-auth failed looking up user %w", err)
	}
	if user == nil {
		user, err = a.userStore.Create(email, name)
		if err != nil {
			return "", fmt.Errorf("local-auth failed creating new user %w", err)
		}
	}
	return a.jwtGenerator.GenerateToken(user.ID(), name)
}
