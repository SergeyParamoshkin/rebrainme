package httpsrv

import "time"

const (
	defaultAddr              = ":8080"
	defaultReadHeaderTimeout = 3 * time.Second
)

type Config struct {
	Addr string `yaml:"addr"`

	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
}

func NewDefaultConfig() Config {
	return Config{
		Addr:              defaultAddr,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}
}
