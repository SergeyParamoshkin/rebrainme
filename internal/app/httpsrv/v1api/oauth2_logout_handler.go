package v1api

import (
	"net/http"

	"github.com/SergeyParamoshkin/alerts/internal/app/httpresp"
)

// oauth2Logout godoc
//
//	@Security		KeycloakAuth
//	@Summary		logout
//	@Description	Logout
//	@Tags			oauth2
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	httpresp.SuccessResponse
//	@Failure		400	{object}	httpresp.ErrResponse
//	@Failure		500	{object}	httpresp.ErrResponse
//	@Router			/auth/oauth2/logout [post]
//
// .
func (a *API) oauth2Logout(w http.ResponseWriter, r *http.Request) {
	cookies, err := a.oauth2Controller.Logout(r)
	if err != nil {
		httpresp.Error(w, r, err)

		return
	}

	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}

	httpresp.Render(w, r, httpresp.NewSuccessResponse("OK"))
}
