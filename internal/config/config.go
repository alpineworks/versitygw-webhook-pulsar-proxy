package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"error"`

	MetricsEnabled bool `env:"METRICS_ENABLED" envDefault:"true"`
	MetricsPort    int  `env:"METRICS_PORT" envDefault:"8081"`

	Local bool `env:"LOCAL" envDefault:"false"`

	TracingEnabled    bool    `env:"TRACING_ENABLED" envDefault:"false"`
	TracingSampleRate float64 `env:"TRACING_SAMPLERATE" envDefault:"0.01"`
	TracingService    string  `env:"TRACING_SERVICE" envDefault:"versitygw-webhook-pulsar-proxy"`
	TracingVersion    string  `env:"TRACING_VERSION"`

	ServerPort           int           `env:"SERVER_PORT" envDefault:"8080"`
	PulsarURL            string        `env:"PULSAR_URL" envDefault:"pulsar://localhost:6650"`
	PulsarTopic          string        `env:"PULSAR_TOPIC" envDefault:"s3-events"`
	PulsarProduceTimeout time.Duration `env:"PULSAR_PRODUCE_TIMEOUT_SECONDS" envDefault:"5s"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
