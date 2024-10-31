package domain

type TicketFilter struct {
	TicketLike *string

	Limit  int
	Offset int
	Total  bool
}
