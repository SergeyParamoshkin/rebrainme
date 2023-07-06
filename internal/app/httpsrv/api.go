package httpsrv

import (
	"net/http"
)

type API interface {
	http.Handler

	Version() string
}
