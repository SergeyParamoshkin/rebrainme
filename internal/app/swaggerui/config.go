package swaggerui

const (
	defaultScheme = "http"
)

type Config struct {
	Host     string   `yaml:"host"`
	BasePath string   `yaml:"base_path"`
	Schemes  []string `yaml:"schemes"`
}

func NewDefaultConfig(host, basePath string) Config {
	return Config{
		Host:     host,
		BasePath: basePath,
		Schemes:  []string{defaultScheme},
	}
}
