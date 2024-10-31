package v1api

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
	"github.com/SergeyParamoshkin/alerts/internal/app/httpresp"
	fwerr "github.com/SergeyParamoshkin/alerts/internal/err"
)

type ticketListRequestFilter struct {
	Key *string `json:"key"`
}

type ticketListRequest struct {
	Size   int `json:"size"`
	Offset int `json:"offset"`

	Filter ticketListRequestFilter `json:"filter"`
}

func (req *ticketListRequest) Bind(_ *http.Request) error {
	return nil
}

type ticketListResponseData struct {
	Total int             `json:"total"`
	Items []domain.Ticket `json:"items"`
}

type ticketListResponse struct {
	Code string                 `json:"code"`
	Data ticketListResponseData `json:"data"`
}

func (resp *ticketListResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	resp.Code = "OK"

	render.Status(r, http.StatusOK)

	return nil
}

// ticketList godoc
//
//	@Summary		it's very to get all tickets
//	@Description	User list
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ticketListRequest	true	"request body"
//	@Success		200		{object}	ticketListResponse
//	@Failure		400		{object}	httpresp.ErrResponse
//	@Failure		404		{object}	httpresp.ErrResponse
//	@Failure		500		{object}	httpresp.ErrResponse
//	@Router			/ticket/list [post]
//
// .
func (a *API) ticketList(w http.ResponseWriter, r *http.Request) {
	var request ticketListRequest

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		httpresp.Error(w, r, fwerr.APIErrorWrap("failed to parse request", err))

		return
	}

	var total int

	ticketList, total, err := a.ticketService.TicketList(r.Context(), &domain.TicketFilter{
		Limit:      request.Size,
		Offset:     request.Offset,
		TicketLike: request.Filter.Key,
		Total:      true,
	})
	if err != nil {
		httpresp.Error(w, r, err)

		return
	}

	httpresp.Render(w, r, &ticketListResponse{
		Data: ticketListResponseData{
			Total: total,
			Items: ticketList,
		},
	})
}
