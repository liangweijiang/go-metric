package config

import (
	"go.opentelemetry.io/otel/attribute"
	"os"
	"time"
)

// MeterEnv represents an enumeration of environments for metering purposes, such as "production", "test", or "dev".
type MeterEnv string

const (

	// MeterEnvProduct represents the production environment for metering.
	MeterEnvProduct MeterEnv = "production"

	// MeterEnvTest represents the test environment for metering.
	MeterEnvTest MeterEnv = "test"

	// MeterEnvDev represents the development environment for metering. Create nop meter.
	MeterEnvDev MeterEnv = "dev"
)

type MeterProviderType int

const (
	MeterProviderTypePrometheus MeterProviderType = iota + 1
)

type PushGatewayCfg struct {
	GatewayAddress string
	PushPeriod     time.Duration
}

type Config struct {
	PrometheusPort      int
	Env                 MeterEnv
	MeterProvider       MeterProviderType
	LocalIP             string
	PushGateway         *PushGatewayCfg
	HistogramBoundaries []float64
	BaseTags            map[string]string
	InfoLogWrite        func(s string)
	ErrorLogWrite       func(s string)
}

func GetConfig() *Config {
	return new(Config)
}

// WriteErrorOrNot logs an error message either to a custom error log function defined in Config or to stdout if not set.
// It prefixes the message with "[go-metrics][error]:" when writing to stdout.
//
// Parameters:
// s (string): The error message to be logged.
//
// Returns:
// None
func (c *Config) WriteErrorOrNot(s string) {
	if c.ErrorLogWrite == nil {
		_, _ = os.Stdout.WriteString("[go-metrics][error]: " + s + "\n")
	} else {
		c.ErrorLogWrite("[go-metrics] " + s)
	}
}

// WriteInfoOrNot logs an informational message to either stdout or a custom info log function based on the configuration.
// If the InfoLogWrite function is not set in Config, it defaults to writing to stdout with a prefixed label.
//
// Parameters:
// s (string): The informational message to log.
//
// Returns:
// None
func (c *Config) WriteInfoOrNot(s string) {
	if c.InfoLogWrite == nil {
		_, _ = os.Stdout.WriteString("[go-metrics][error]: " + s + "\n")
	} else {
		c.InfoLogWrite("[go-metrics] " + s)
	}
}

func (c *Config) WithAttributes() []attribute.KeyValue {
	var attributes []attribute.KeyValue
	for key, value := range c.BaseTags {
		attributes = append(attributes, attribute.String(key, value))
	}
	return attributes
}

func (c *Config) IsDev() bool {
	return c.Env == MeterEnvDev
}
