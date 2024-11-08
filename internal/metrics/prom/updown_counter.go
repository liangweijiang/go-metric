package prom

import (
	"context"
	"go-mertric/pkg/interfaces"
	"go.opentelemetry.io/otel/metric"
)

// _ is a blank identifier used for type assertion to ensure that *UpDownCounter implements the interfaces.UpDownCounter interface.
var _ interfaces.UpDownCounter = (*UpDownCounter)(nil)

type UpDownCounter struct {
	base    Base
	counter metric.Float64UpDownCounter
}

// NewUpDownCounter creates a new UpDownCounter instance wrapping the provided metric.Float64UpDownCounter with a given name and optional tags management.
// It returns an implementation of interfaces.UpDownCounter that delegates to the underlying counter for Update, IncrOne, DecrOne, AddTag, and WithTags operations.
func NewUpDownCounter(name string, counter metric.Float64UpDownCounter) interfaces.UpDownCounter {
	return &UpDownCounter{
		base: Base{
			name: name,
		},
		counter: counter,
	}

}

// Update adjusts the counter by the given delta.
// It requires a context and a float64 value representing the change.
// If the counter is not ready, the update is ignored.
func (c *UpDownCounter) Update(ctx context.Context, delta float64) {
	if !c.base.ready() {
		return
	}
	c.counter.Add(ctx, delta, metric.WithAttributes(c.base.tags...))
}

// IncrOne increments the UpDownCounter by one, given a context. This is a convenience method wrapping around Update with a delta of 1.
func (c *UpDownCounter) IncrOne(ctx context.Context) {
	c.Update(ctx, 1)
}

// DecrOne decreases the counter by one.
// It uses the provided context and updates the counter with a delta of -1.
// ctx context.Context: The context for the operation.
func (c *UpDownCounter) DecrOne(ctx context.Context) {
	c.Update(ctx, -1)
}

// AddTag adds a tag with the specified key and value to the UpDownCounter's base tags.
// It returns the UpDownCounter itself for chaining calls.
// Key must match the regular expression pattern "^[a-zA-Z_][a-zA-Z0-9_]*$" and cannot start with "__".
func (c *UpDownCounter) AddTag(key string, value string) interfaces.UpDownCounter {
	c.base.AddTag(key, value)
	return c
}

func (c *UpDownCounter) WithTags(tags map[string]string) interfaces.UpDownCounter {
	c.base.WithTags(tags)
	return c
}
