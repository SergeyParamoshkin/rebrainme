package httpresp

import (
	"net/http"

	"github.com/go-chi/render"
)

type SuccessResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (resp SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)

	return nil
}

func NewSuccessResponse(message string) SuccessResponse {
	return SuccessResponse{
		Code:    "OK",
		Message: message,
	}
}
