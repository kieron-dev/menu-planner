package models

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . User

type User interface {
	Email() string
	Name() string
	ID() int
}
