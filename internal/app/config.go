package app

import (
	"os"

	"github.com/SergeyParamoshkin/alerts/internal/app/auth"
	"github.com/SergeyParamoshkin/alerts/internal/app/httpsrv"
	"github.com/SergeyParamoshkin/alerts/internal/app/httpsrv/v1api"
	"github.com/SergeyParamoshkin/alerts/internal/app/service/ticketsvc"
	"github.com/SergeyParamoshkin/alerts/internal/keycloak"
	"github.com/SergeyParamoshkin/alerts/internal/postgres"
	"github.com/SergeyParamoshkin/alerts/internal/tel"
	"go.uber.org/config"
	"go.uber.org/fx"
)

const (
	appName = "fw"

	defaultHostname = "localhost"
)

type AppConfig struct {
	Name     string `yaml:"name"`
	Hostname string `yaml:"hostname"`
	Debug    bool   `yaml:"debug"`
}

type fileConfig struct {
	App       AppConfig       `yaml:"app"`
	Keycloak  keycloak.Config `yaml:"keycloak"`
	HTTP      httpsrv.Config  `yaml:"http"`
	Postgres  postgres.Config `yaml:"postgres"`
	Telemetry tel.Config      `yaml:"telemetry"`
}

type ConfigOut struct {
	fx.Out

	App       *AppConfig
	HTTP      *httpsrv.Config
	V1API     *v1api.Config
	Ticket    *ticketsvc.Config
	Keycloak  *keycloak.Config
	OAuth2    *auth.OAuth2Config
	Postgres  *postgres.Config
	Telemetry *tel.Config
}

func NewConfig(args *Args) (ConfigOut, error) {
	provider, err := config.NewYAML(
		config.Permissive(),
		config.Static(fileConfig{ // default config
			App: AppConfig{
				Name:     appName,
				Hostname: defaultHostname,
				Debug:    false,
			},
			Keycloak:  keycloak.NewDefaultConfig(),
			HTTP:      httpsrv.NewDefaultConfig(),
			Postgres:  postgres.NewDefaultConfig(),
			Telemetry: tel.NewDefaultConfig(),
		}),
		config.Expand(os.LookupEnv),
		config.File(args.ConfigFileName),
	)
	if err != nil {
		return ConfigOut{}, err
	}

	c := fileConfig{}

	err = provider.Get("").Populate(&c)
	if err != nil {
		return ConfigOut{}, err
	}

	v1API := v1api.NewDefaultConfig()
	v1API.Debug = c.App.Debug

	err = provider.Get("http").Populate(&v1API)
	if err != nil {
		return ConfigOut{}, err
	}

	oauth2 := auth.NewDefaultOAuth2Config()

	err = provider.Get("http.oauth2").Populate(&oauth2)
	if err != nil {
		return ConfigOut{}, err
	}

	return ConfigOut{
		App:       &c.App,
		HTTP:      &c.HTTP,
		Keycloak:  &c.Keycloak,
		OAuth2:    &oauth2,
		V1API:     &v1API,
		Postgres:  &c.Postgres,
		Telemetry: &c.Telemetry,
	}, nil
}
