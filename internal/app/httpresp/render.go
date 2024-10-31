package httpresp

import (
	"net/http"

	"github.com/go-chi/render"

	fwctx "github.com/SergeyParamoshkin/alerts/internal/ctx"
)

func Error(w http.ResponseWriter, r *http.Request, err error) {
	fwctx.RecordError(r.Context(), err)

	//nolint: errcheck // last resort, must not be error
	_ = render.Render(w, r, NewErrResponse(err))
}

func Render(w http.ResponseWriter, r *http.Request, resp render.Renderer) {
	if err := render.Render(w, r, resp); err != nil {
		Error(w, r, err)

		return
	}
}
