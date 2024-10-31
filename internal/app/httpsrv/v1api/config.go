package v1api

import "github.com/SergeyParamoshkin/alerts/internal/app/swaggerui"

const (
	defaultSwaggerUIHost     = "localhost:8080"
	defaultSwaggerUIBasePath = "/v1"

	defaultOauth2CookieDomain = "localhost:8080"
	defaultOauth2CookiePath   = "/"
	defaultOauth2CookieSecure = false
)

type Oauth2CookieConfig struct {
	Domain string `yaml:"domain"`
	Path   string `yaml:"path"`
	Secure bool   `yaml:"secure"`
}

type Config struct {
	SwaggerUI    swaggerui.Config   `yaml:"swaggerUi"`
	Oauth2Cookie Oauth2CookieConfig `yaml:"oauth2Cookie"`

	DisableAuth bool `yaml:"disableAuth"`
	Debug       bool `yaml:"debug"`
}

func NewDefaultConfig() Config {
	return Config{
		SwaggerUI: swaggerui.NewDefaultConfig(
			defaultSwaggerUIHost,
			defaultSwaggerUIBasePath,
		),
		Oauth2Cookie: Oauth2CookieConfig{
			Domain: defaultOauth2CookieDomain,
			Path:   defaultOauth2CookiePath,
			Secure: defaultOauth2CookieSecure,
		},
		DisableAuth: false,
		Debug:       false,
	}
}
