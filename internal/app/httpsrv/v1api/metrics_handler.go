package v1api

import (
	"net/http"
)

func (a *API) metricsHandler(w http.ResponseWriter, r *http.Request) {
	a.promHandler.ServeHTTP(w, r)
}
