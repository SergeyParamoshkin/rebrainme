package v1api

import (
	"net/http"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-chi/render"
)

type oauth2CallbackResponse struct { //nolint:unused
	*gocloak.JWT
}

func (resp *oauth2CallbackResponse) Render(_ http.ResponseWriter, r *http.Request) error { //nolint:unused,lll
	render.Status(r, http.StatusOK)

	return nil
}

// oauth2Callback godoc
//
//	@Summary		callback
//	@Description	Callback
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string	true	"authorization code"	""
//	@Success		200		{object}	oauth2CallbackResponse
//	@Failure		400		{object}	httpresp.ErrResponse
//	@Failure		500		{object}	httpresp.ErrResponse
//	@Router			/auth/oauth2/callback [get]
//
// .
func (a *API) oauth2Callback(w http.ResponseWriter, r *http.Request) {
	a.oauth2Controller.CallbackHandler(w, r)
}
