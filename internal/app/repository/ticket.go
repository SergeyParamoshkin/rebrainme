package repository

import (
	"time"

	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
	"github.com/google/uuid"
)

type Ticket struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Queue     string
	Key       string
	Summary   string
	Type      string
}

func (t *Ticket) ToDomain() *domain.Ticket {
	return &domain.Ticket{
		ID:        t.ID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Queue:     t.Queue,
		Key:       t.Key,
		Summary:   t.Summary,
		Type:      t.Type,
	}
}
