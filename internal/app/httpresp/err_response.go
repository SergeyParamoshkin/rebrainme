package httpresp

import (
	"net/http"

	"github.com/go-chi/render"

	fwerr "github.com/SergeyParamoshkin/alerts/internal/err"
)

type ErrResponse struct {
	Err fwerr.Error `json:"-"`

	Code    fwerr.Code `json:"code"`
	Message string     `json:"message,omitempty"`
}

func (resp ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, ErrCodeToHTTPStatus(resp.Err.Code))

	return nil
}

func NewErrResponse(err error) ErrResponse {
	fwErr := fwerr.FromError(err)

	return ErrResponse{
		Err:     fwErr,
		Code:    fwErr.Code,
		Message: fwErr.String(),
	}
}

func ErrCodeToHTTPStatus(c fwerr.Code) int {
	switch c {
	case fwerr.CodeAuthError:
		return http.StatusUnauthorized
	case fwerr.CodeForbidden:
		return http.StatusForbidden
	case fwerr.CodeBadRequest:
		return http.StatusBadRequest
	case fwerr.CodeNotFound:
		return http.StatusNotFound
	case
		fwerr.CodeDBError,
		fwerr.CodeInternalError:
		fallthrough
	default:
		return http.StatusInternalServerError
	}
}
