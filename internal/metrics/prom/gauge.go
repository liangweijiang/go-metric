package prom

import (
	"context"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"go.opentelemetry.io/otel/metric"
)

// _ is a blank identifier used for type assertion to ensure that the Gauge struct implements the interfaces.Gauge interface.
var _ interfaces.Gauge = (*Gauge)(nil)

// Gauge is a struct representing a metric gauge which measures non-cumulative values like memory usage or CPU utilization.
// It embeds a Base for common attributes and a Float64Gauge for gauge operations.
type Gauge struct {
	base  Base
	gauge metric.Float64Gauge
}

// NewGauge creates a new Gauge interface instance wrapping a metric.Float64Gauge with a given name and initial gauge.
// It initializes the Gauge with a Base that includes the name and no initial tags.
func NewGauge(name string, gauge metric.Float64Gauge) interfaces.Gauge {
	return &Gauge{
		base: Base{
			name: name,
		},
		gauge: gauge,
	}
}

// Update records the given value to the gauge metric if the gauge is ready.
// It uses the provided context and adds the gauge's associated tags to the recorded metric.
// Parameters:
//
//	ctx: The context for recording the metric.
//	v: The value to update the gauge with.
//
// It returns nothing and does not indicate whether the update was successful.
func (g *Gauge) Update(ctx context.Context, v float64) {
	if !g.base.ready() {
		return
	}
	g.gauge.Record(ctx, v, metric.WithAttributes(g.base.tags...))
}

// AddTag adds a tag with the specified key and value to the Gauge's tags.
// It modifies the Gauge in place and returns the same instance for chaining calls.
// Key must adhere to the regex pattern `^[a-zA-Z_][a-zA-Z0-9_]*$`, avoiding double underscores at the start.
// Parameters:
//
//	key:   The tag key, which should be a valid identifier.
//	value: The value associated with the tag key.
func (g *Gauge) AddTag(key string, value string) interfaces.Gauge {
	g.base.AddTag(key, value)
	return g
}

// WithTags sets the provided tags on the Gauge, appending them to existing tags.
// It modifies the Gauge in place and returns the same instance for chaining calls.
// If the input map is nil or empty, no action is taken.
// Parameters:
//
//	tags: A map of tags to associate with the Gauge.
//
// Returns:
// The Gauge instance with updated tags.
func (g *Gauge) WithTags(tags map[string]string) interfaces.Gauge {
	g.base.WithTags(tags)
	return g
}
