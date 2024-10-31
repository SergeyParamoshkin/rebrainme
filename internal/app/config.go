package app

import (
	"os"

	"github.com/SergeyParamoshkin/alerts/internal/app/httpsrv"
	"github.com/SergeyParamoshkin/alerts/internal/app/httpsrv/v1api"
	"github.com/SergeyParamoshkin/alerts/internal/app/service/ticketsvc"
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
	App  AppConfig      `yaml:"app"`
	HTTP httpsrv.Config `yaml:"http"`

	Postgres  postgres.Config `yaml:"postgres"`
	Telemetry tel.Config      `yaml:"telemetry"`
}

type ConfigOut struct {
	fx.Out

	App       *AppConfig
	HTTP      *httpsrv.Config
	V1API     *v1api.Config
	Ticket    *ticketsvc.Config
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

	return ConfigOut{
		App:       &c.App,
		HTTP:      &c.HTTP,
		V1API:     &v1API,
		Postgres:  &c.Postgres,
		Telemetry: &c.Telemetry,
	}, nil
}
