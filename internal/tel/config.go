package tel

const (
	defaultServiceName     = "fw"
	defaultSamplerFraction = 1.0 // always sample
	defaultJaegerAgentHost = "localhost"
	defaultJaegerAgentPort = "6831"
)

type JaegerConfig struct {
	AgentHost string `yaml:"agent_host"`
	AgentPort string `yaml:"agent_port"`
}

type Config struct {
	ServiceName     string  `yaml:"service_name"`
	ServiceVersion  string  `yaml:"-"` // set manually
	SamplerFraction float64 `yaml:"sampler_fraction"`

	Jaeger JaegerConfig
}

func NewDefaultConfig() Config {
	return Config{
		ServiceName:     defaultServiceName,
		SamplerFraction: defaultSamplerFraction,
		Jaeger: JaegerConfig{
			AgentHost: defaultJaegerAgentHost,
			AgentPort: defaultJaegerAgentPort,
		},
	}
}
