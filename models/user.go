package models

type User struct {
	ID string
	Email string
	Username string
	Password []byte
	Admin bool
}