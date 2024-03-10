package entities

import "github.com/google/uuid"

type Company struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Phone    string
	Password string
}
type Address struct {
	ID       uuid.UUID
	Country  string
	State    string
	District string
	City     string
}
type Link struct {
	ID    uuid.UUID
	Title string
	URL   string
}
