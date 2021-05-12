package models

type Recipe struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	UserID int    `json:"-"`
}
