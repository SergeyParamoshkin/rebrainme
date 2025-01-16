package config

type Config struct{}

func NewConfig() Config {
	return Config{}
}

func (c *Config) Validate() error {
	// TODO: add validation logic here
	return nil
}
