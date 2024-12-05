package v1api

import (
	"net/http"
)

// oauth2Login godoc
//
//	@Summary		login
//	@Description	Login
//	@Tags			oauth2
//	@Param			redirect_url	query		string	true	"redirect url"	""
//	@Success		200				{object}	oauth2CallbackResponse
//	@Failure		400				{object}	httpresp.ErrResponse
//	@Failure		500				{object}	httpresp.ErrResponse
//	@Router			/auth/oauth2/login [get]
//
// .
func (a *API) oauth2Login(w http.ResponseWriter, r *http.Request) {
	a.oauth2Controller.LoginHandler(w, r)
}
