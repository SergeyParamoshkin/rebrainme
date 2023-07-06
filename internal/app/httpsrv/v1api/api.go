package v1api

import (
	"fmt"
	"net/http"

	"dumper/internal/app/httpsrv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/swaggo/swag/example/basic/docs"
	"go.uber.org/zap"
)

const (
	MaxSizeBytes = 128 * 1024 // Kb
)

var _ httpsrv.API = &API{}

// @title           fw API
// @version         1.0
// @description     fw server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.apikey APIKeyAuth
// @in header
// @name Authorization
// @description Bearer Token Authorization with token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

type API struct {
	router *chi.Mux
	logger *zap.Logger

	promHandler http.Handler
}

func New(params Params) (Result, error) {
	logger := params.Logger.Named("api")

	api := &API{
		router: chi.NewRouter(),
		logger: logger,
	}

	api.router.Route(fmt.Sprintf("/%s", api.Version()), func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST"},
		}))
	})

	return Result{
		API: *api,
	}, nil
}

func (a *API) GenSwaggerJSON(host, basePath string, schemes []string) string {
	swaggerInfo := *docs.SwaggerInfo
	swaggerInfo.Host = host
	swaggerInfo.BasePath = basePath
	swaggerInfo.Schemes = schemes

	return swaggerInfo.ReadDoc()
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *API) Version() string {
	return "v1"
}
