package ticketsvc

import (
	"context"

	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
	"github.com/SergeyParamoshkin/alerts/internal/app/repository"
)

func (s *Service) TicketList(
	ctx context.Context, filter *domain.TicketFilter,
) (tickets []domain.Ticket, total int, err error) {
	err = s.repo.WithTx(ctx, func(ctx context.Context, q repository.Queryable) error {
		tickets, total, err = s.repo.TicketList(ctx, q, filter)

		return err
	})

	return
}
