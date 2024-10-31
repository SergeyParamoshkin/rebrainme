package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Queue     string
	Key       string
	Summary   string
	Type      string
}
