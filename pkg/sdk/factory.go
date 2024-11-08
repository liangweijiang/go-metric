package sdk

import (
	"go-mertric/internal/meter/nop"
	"go-mertric/internal/meter/prom"
	"go-mertric/pkg/config"
	"go-mertric/pkg/interfaces"
)

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
