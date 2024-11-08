package prom

import (
	"context"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"go.opentelemetry.io/otel/metric"
	"time"
)

// _ is a blank identifier used for type assertion to ensure that (*Histogram) implements the interfaces.Histogram interface.
var _ interfaces.Histogram = (*Histogram)(nil)

// Histogram represents a distribution of values over time.
// It is used to measure value distributions and supports updating with different time units.
// Histogram also allows adding tags for context and provides a method to time functions and record their durations.
type Histogram struct {
	base      Base
	histogram metric.Float64Histogram
}

// NewHistogram creates and returns a new Histogram instance wrapping the provided float64 histogram.
// It assigns the given name to the histogram and associates it with the base metrics structure.
// Parameters:
//
//	name: The name of the histogram metric.
//	histogram: The underlying float64 histogram implementation to use.
//
// Returns:
//
//	An interfaces.Histogram instance for tracking value distributions over time.
func NewHistogram(name string, histogram metric.Float64Histogram) interfaces.Histogram {
	return &Histogram{
		base: Base{
			name: name,
		},
		histogram: histogram,
	}
}

// Update adjusts the histogram with the duration in seconds converted from the given time.Duration value.
// It uses the context to associate the update with a tracing span, if one exists.
// The actual update is performed by calling UpdateInSeconds.
func (h *Histogram) Update(ctx context.Context, d time.Duration) {
	h.UpdateInSeconds(ctx, d.Seconds())
}

// UpdateInSeconds records a value in seconds to the histogram.
// It requires a context to optionally associate the update with a tracing span.
// No operation is performed if the histogram's base is not ready.
func (h *Histogram) UpdateInSeconds(ctx context.Context, s float64) {
	if !h.base.ready() {
		return
	}
	h.histogram.Record(ctx, s, metric.WithAttributes(h.base.tags...))
}

// UpdateInMilliseconds updates the histogram with a value in milliseconds, converting it to seconds before recording.
// This method takes a context to optionally associate the update with a tracing span and a float64 value representing the measurement in milliseconds.
// It internally calls UpdateInSeconds after converting the input to seconds.
// ctx context.Context: The context for optional tracing.
// m float64: The value in milliseconds to record in the histogram.
func (h *Histogram) UpdateInMilliseconds(ctx context.Context, m float64) {
	h.UpdateInSeconds(ctx, m/1000)
}

// UpdateSine calculates the elapsed time since the given start time and updates the histogram using UpdateInSeconds.
// This method is useful for timing the execution of a function or process and recording its duration in seconds.
// The update is associated with the provided context, which can include tracing spans.
// Parameters:
//
//	ctx: The context carrying optional tracing information.
//	start: The start time from which to calculate elapsed time.
func (h *Histogram) UpdateSine(ctx context.Context, start time.Time) {
	elapsed := time.Now().Sub(start)
	h.UpdateInSeconds(ctx, elapsed.Seconds())
}

// Time executes the provided function f and records its duration in seconds to the histogram.
// It starts a timer before calling f, and upon completion, it calculates the elapsed time and updates the histogram using UpdateSine.
// The context.Background() is used for this operation, which can be useful for tracing purposes.
func (h *Histogram) Time(f func()) {
	start := time.Now()
	f()
	h.UpdateSine(context.Background(), start)
}

// AddTag adds a tag with the specified key and value to the Histogram's base tags.
// It returns the modified Histogram instance allowing for method chaining.
// Key must be a valid identifier matching the regex (^[a-zA-Z_][a-zA-Z0-9_]*$).
// Value is the string value associated with the tag key.
// Tags starting with double underscores will be automatically escaped.
func (h *Histogram) AddTag(key string, value string) interfaces.Histogram {
	h.base.AddTag(key, value)
	return h
}

// WithTags initializes all tags for the histogram using the provided map.
// It updates the histogram's base tags with the new set of tags.
// Tags starting with double underscores will be automatically escaped.
// Parameters:
//
//	tags: A map of string key-value pairs representing tags to be added.
//
// Returns:
//
//	The updated Histogram instance with the new tags.
func (h *Histogram) WithTags(tags map[string]string) interfaces.Histogram {
	h.base.WithTags(tags)
	return h
}
