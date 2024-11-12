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

// Config holds the configuration parameters for setting up metrics reporting, including port details, environment settings, meter provider types, push gateway configurations, histogram boundaries, base tags for metrics, and optional log output functions.
type Config struct {
	PrometheusPort        int
	LocalIP               string
	Env                   MeterEnv
	MeterProvider         MeterProviderType
	PushGateway           *PushGatewayCfg
	RuntimeMetricsCollect bool
	HistogramBoundaries   []float64
	BaseTags              map[string]string
	InfoLogWrite          func(s string)
	ErrorLogWrite         func(s string)
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
		_, _ = os.Stdout.WriteString("[go-metrics][info]: " + s + "\n")
	} else {
		c.InfoLogWrite("[go-metrics] " + s)
	}
}

// WithBaseTags creates a slice of attribute.KeyValue from the BaseTags map in the Config.
// Each key-value pair in the BaseTags map is converted into an attribute.KeyValue.
// This function is useful for populating common tags across metrics or traces.
func (c *Config) WithBaseTags() []attribute.KeyValue {
	var attributes []attribute.KeyValue
	for key, value := range c.BaseTags {
		attributes = append(attributes, attribute.String(key, value))
	}
	return attributes
}

// IsDev returns true if the configuration's environment is set to development (`MeterEnvDev`).
func (c *Config) IsDev() bool {
	return c.Env == MeterEnvDev
}
