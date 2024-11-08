package config

import (
	"go.opentelemetry.io/otel/attribute"
	"os"
	"time"
)

type PushGatewayCfg struct {
	GatewayAddress string
	PushPeriod     time.Duration
}

type Config struct {
	PrometheusPort      int
	LocalIP             string
	PushGateway         *PushGatewayCfg
	HistogramBoundaries []float64
	BaseTags            map[string]string
	InfoLogWrite        func(s string)
	ErrorLogWrite       func(s string)
}

// WriteErrorOrNot logs an error message either to a custom error log function defined in Config or to stdout if not set.
// It prefixes the message with "[go-metrics][error]:" when writing to stdout.
//
// Parameters:
// s (string): The error message to be logged.
//
// Returns:
// None
func (cfg *Config) WriteErrorOrNot(s string) {
	if cfg.ErrorLogWrite == nil {
		_, _ = os.Stdout.WriteString("[go-metrics][error]: " + s + "\n")
	} else {
		cfg.ErrorLogWrite("[go-metrics] " + s)
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
func (cfg *Config) WriteInfoOrNot(s string) {
	if cfg.InfoLogWrite == nil {
		_, _ = os.Stdout.WriteString("[go-metrics][error]: " + s + "\n")
	} else {
		cfg.InfoLogWrite("[go-metrics] " + s)
	}
}

func (cfg *Config) WithAttributes() []attribute.KeyValue {
	var attributes []attribute.KeyValue
	for key, value := range cfg.BaseTags {
		attributes = append(attributes, attribute.String(key, value))
	}
	return attributes
}
