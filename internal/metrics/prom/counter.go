package prom

import (
	"context"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"go.opentelemetry.io/otel/metric"
)

// _ is a blank identifier used for type assertion to ensure that *Counter implements the interfaces.Counter interface.
var _ interfaces.Counter = (*Counter)(nil)

// Counter combines a Base structure for metric identification and tagging with a metric.Float64Counter to track incremental values.
// It provides methods to increment the counter, add tags, and manage context-specific metadata.
type Counter struct {
	base    Base
	counter metric.Float64Counter
}

// NewCounter creates and returns a new Counter instance wrapping a metric.Float64Counter with a given name and initial counter.
// It initializes the Counter with a base structure containing the provided name and prepares it for metric increments.
// Parameters:
//
//	name: The name of the counter metric.
//	counter: The underlying Float64Counter to wrap with the Counter interface.
//
// Returns an implementation of interfaces.Counter.
func NewCounter(name string, counter metric.Float64Counter) interfaces.Counter {
	return &Counter{
		base: Base{
			name: name,
		},
		counter: counter,
	}
}

// Incr increments the counter by the given delta, provided the context and ensuring the counter is ready for operations.
func (c *Counter) Incr(ctx context.Context, delta float64) {
	if !c.base.ready() {
		return
	}
	c.counter.Add(ctx, delta, metric.WithAttributes(c.base.tags...))
}

// IncrOne increments the counter by one, given a context. It is a convenience method wrapping around Incr with a fixed delta of 1.
func (c *Counter) IncrOne(ctx context.Context) {
	c.Incr(ctx, 1)
}

// AddTag adds a tag with the specified key and value to the Counter's base tags.
// It returns the Counter instance to allow for method chaining.
// Key must adhere to the pattern ^[a-zA-Z_][a-zA-Z0-9_]*$, avoiding __ prefix.
// Parameters:
//
//	key: The tag key to add.
//	value: The value associated with the added tag key.
//
// Returns:
//
//	The updated Counter instance.
func (c *Counter) AddTag(key string, value string) interfaces.Counter {
	c.base.AddTag(key, value)
	return c
}

// WithTags sets the provided tags on the Counter's base instance, appending them to existing tags.
// It allows for adding contextual metadata to the Counter in the form of a tag map.
// Parameters:
//
//	tags: A map of tags to set on the Counter.
//
// Returns:
//
//	The Counter instance with updated tags.
func (c *Counter) WithTags(tags map[string]string) interfaces.Counter {
	c.base.WithTags(tags)
	return c
}
