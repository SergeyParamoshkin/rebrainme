package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
}

func (u *User) AddUser(name string) {}
