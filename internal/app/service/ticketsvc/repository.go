package ticketsvc

import (
	"context"

	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
	"github.com/SergeyParamoshkin/alerts/internal/app/repository"
)

type Repo interface {
	WithConn(ctx context.Context, callback func(ctx context.Context, q repository.Queryable) error) error
	WithTx(ctx context.Context, callback func(ctx context.Context, q repository.Queryable) error) error

	TicketList(
		ctx context.Context,
		q repository.Queryable,
		filter *domain.TicketFilter,
	) (tickets []domain.Ticket, total int, err error)
}
