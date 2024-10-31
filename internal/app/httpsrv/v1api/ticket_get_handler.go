package v1api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/SergeyParamoshkin/alerts/internal/app/domain"
	"github.com/SergeyParamoshkin/alerts/internal/app/httpresp"
)

type ticketGetResponseData struct {
	Total int           `json:"total"`
	Item  domain.Ticket `json:"item"`
}

type ticketGetResponse struct {
	Code string                `json:"code"`
	Data ticketGetResponseData `json:"data"`
}

func (resp *ticketGetResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	resp.Code = "OK"

	render.Status(r, http.StatusOK)

	return nil
}

// ticketGet godoc
//
//	@Summary		it's very to get all tickets
//	@Description	User list
//	@Tags			tickets
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"id"
//	@Param			search	query		string	false	"search"
//	@Success		200		{object}	ticketGetResponse
//	@Failure		400		{object}	httpresp.ErrResponse
//	@Failure		404		{object}	httpresp.ErrResponse
//	@Failure		500		{object}	httpresp.ErrResponse
//	@Router			/ticket/{id} [get]
func (a *API) ticketGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	search := r.URL.Query().Get("search")
	log.Println(search)
	// ticket, err := a.ticketService.Get(id)

	u, err := uuid.FromBytes([]byte(id))
	if err != nil {
		log.Printf("failed to parse UUID: %v", err)
		httpresp.Error(w, r, err)
		return
	}
	ticket := domain.Ticket{
		ID: u,
	}

	httpresp.Render(w, r, &ticketGetResponse{
		Code: "ok",
		Data: ticketGetResponseData{
			Total: 1,
			Item:  ticket,
		},
	})
}
