package v1api

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/SergeyParamoshkin/alerts/internal/app/httpresp"
	"github.com/SergeyParamoshkin/alerts/pkg/info"
)

type infoResponse struct {
	CommitSHA string `json:"commit_sha"`
	Version   string `json:"version"`
}

func (resp *infoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	resp.CommitSHA = info.CommitSHA
	resp.Version = info.Version

	render.Status(r, http.StatusOK)

	return nil
}

// info godoc
//
//	@Summary		info
//	@Description	info
//	@Tags			info
//	@Produce		json
//	@Success		200	{object}	infoResponse
//	@Router			/info [get]
func (a *API) info(w http.ResponseWriter, r *http.Request) {
	httpresp.Render(w, r, &infoResponse{
		CommitSHA: info.CommitSHA,
		Version:   info.Version,
	})
}
