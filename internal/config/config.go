package config

type Config struct {
	Http HTTPConfig
}

type HTTPConfig struct {
	Port string
}
