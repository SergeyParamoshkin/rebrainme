package postgres

const (
	defaultDatabaseURL = "localhost:5432"
)

type Config struct {
	DatabaseURL string `yaml:"database_url"`
}

func NewDefaultConfig() Config {
	return Config{
		DatabaseURL: defaultDatabaseURL,
	}
}
