package repository

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
)

func (r *Repo) userApplyFilter(
	sb squirrel.SelectBuilder,
	filter *domain.TicketFilter,
) squirrel.SelectBuilder {
	return sb
}

func (r *Repo) userApplyLimitOffset(
	sb squirrel.SelectBuilder,
	filter *domain.TicketFilter,
) squirrel.SelectBuilder {
	if filter.Limit != 0 {
		sb = sb.Limit(uint64(filter.Limit))
	}

	if filter.Offset != 0 {
		sb = sb.Offset(uint64(filter.Offset))
	}

	return sb
}

func (r *Repo) ticketListTotal(ctx context.Context, q Queryable, filter *domain.TicketFilter) (total int, err error) {
	sb := r.stmtBuilder.Select("COUNT (*)").From(TableTickets)

	query, args, err := sb.ToSql()
	if err != nil {
		return 0, err
	}

	err = q.QueryRow(ctx, query, args...).Scan(&total)

	return total, nil
}

func (r *Repo) ticketList(
	ctx context.Context,
	q Queryable,
	filter *domain.TicketFilter,
) (tickets []domain.Ticket, err error) {
	ticket := Ticket{}
	scanTargets := []interface{}{
		&ticket.ID,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.Queue,
		&ticket.Key,
		&ticket.Summary,
		&ticket.Type,
	}

	sb := r.stmtBuilder.Select(
		ColTicketsID,
		ColTicketsCreatedAt,
		ColTicketsUpdatedAt,
		ColTicketsQueue,
		ColTicketsKey,
		ColTicketsSummary,
		ColTicketsType,
	)

	sb = sb.From(TableTickets)

	sb = r.userApplyFilter(sb, filter)

	query, args, err := r.userApplyLimitOffset(sb, filter).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(scanTargets...)
		if err != nil {
			return nil, err
		}

		tickets = append(tickets, *ticket.ToDomain())
	}

	return tickets, nil
}

func (r *Repo) TicketList(
	ctx context.Context, q Queryable, filter *domain.TicketFilter,
) (tickets []domain.Ticket, total int, err error) {
	tickets, err = r.ticketList(ctx, q, filter)
	if err != nil {
		return nil, 0, err
	}

	total, err = r.ticketListTotal(ctx, q, filter)
	if err != nil {
		return nil, 0, err
	}

	// if filter.Total {
	// 	total, err = r.userListTotal(ctx, q, filter)
	// 	if err != nil {
	// 		return nil, 0, err
	// 	}
	// }

	return tickets, total, nil
}
