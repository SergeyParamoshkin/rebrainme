package main

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Article struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Text  string    `json:"string"`

	UserID uuid.UUID `json:"user_id"`
}
