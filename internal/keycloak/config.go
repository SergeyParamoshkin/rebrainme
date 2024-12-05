package keycloak

const (
	defaultBaseURL          = "http://localhost:8080"
	defaultRealm            = "your_realm_name"
	defaultClientID         = "your_client_id"
	defaultClientSecret     = "your_client_secret"
	defaultRedirectURI      = "http://localhost:8080/v1/oauth2/callback"
	defaultLoginRedirectURL = ""
)

type Config struct {
	BaseURL                 string `yaml:"base_url"`
	Realm                   string `yaml:"realm"`
	ClientID                string `yaml:"client_id"`
	ClientSecret            string `yaml:"client_secret"`
	RedirectURI             string `yaml:"redirect_uri"`
	DefaultLoginRedirectURL string `yaml:"default_login_redirect_url"`
}

func NewDefaultConfig() Config {
	return Config{
		BaseURL:                 defaultBaseURL,
		Realm:                   defaultRealm,
		ClientID:                defaultClientID,
		ClientSecret:            defaultClientSecret,
		RedirectURI:             defaultRedirectURI,
		DefaultLoginRedirectURL: defaultLoginRedirectURL,
	}
}
