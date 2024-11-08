package sdk

import (
	"github.com/liangweijiang/go-metric/internal/meter/nop"
	"github.com/liangweijiang/go-metric/internal/meter/prom"
	"github.com/liangweijiang/go-metric/pkg/config"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
)

// NewMeter creates a new meter instance based on the provided options and configuration.
// It allows customization through options which modify the configuration before deciding the meter provider.
// In a development environment, it returns a no-op meter. For Prometheus configuration, it initializes a Prometheus meter.
// Otherwise, it defaults to a no-op meter.
// Returns a meter implementation and an error if one occurs during initialization.
func NewMeter(options ...interfaces.Option) (interfaces.Meter, error) {
	cfg := config.GetConfig()
	for _, option := range options {
		option.ApplyConfig(cfg)
	}

	if cfg.IsDev() {
		cfg.WriteInfoOrNot("under test environment, using NopMeter")
		return nop.NewNopMeter(), nil
	}

	switch cfg.MeterProvider {
	case config.MeterProviderTypePrometheus:
		meter, err := prom.NewPrometheusMeter(cfg)
		if err != nil {
			cfg.WriteErrorOrNot("set prometheus meter provider error: " + err.Error())
			return nil, err
		}
		return meter, err
	default:
		return nop.NewNopMeter(), nil
	}
}
