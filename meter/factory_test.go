package meter

import (
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"testing"

	"github.com/liangweijiang/go-metric/internal/meter/nop"
	"github.com/liangweijiang/go-metric/internal/meter/prom"
	"github.com/liangweijiang/go-metric/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestNewMeter(t *testing.T) {
	tests := []struct {
		name       string
		options    []interfaces.Option
		wantMeter  interfaces.Meter
		wantErr    bool
		errMessage string
	}{
		{
			name:      "DevEnvironment",
			wantMeter: &nop.Meter{},
			wantErr:   false,
		},
		{
			name:      "PrometheusEnabled",
			wantMeter: &prom.PrometheusMeter{},
			wantErr:   false,
		},
		{
			name:       "UnknownMeterProvider",
			wantMeter:  &nop.Meter{},
			wantErr:    false,
			errMessage: "set prometheus meter provider error: unsupported meter provider type: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.GetConfig()
			for _, opt := range tt.options {
				opt.ApplyConfig(cfg)
			}

			meter, err := NewMeter(tt.options...)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.wantMeter, meter)
			}
		})
	}
}
