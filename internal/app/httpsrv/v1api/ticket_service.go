package v1api

import (
	"context"

	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
)

type TicketService interface {
	TicketList(
		ctx context.Context, filter *domain.TicketFilter,
	) (tickets []domain.Ticket, total int, err error)
}
