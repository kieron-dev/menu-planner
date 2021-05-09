package models

//counterfeiter:generate . Recipe

type Recipe interface {
	Name() string
	ID() int
}
