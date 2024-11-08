package nop

import (
	"context"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
)

// _ is a blank identifier used for type assertion to ensure that nopGauge implements the interfaces.Gauge interface.
var _ interfaces.Gauge = (*nopGauge)(nil)

// nopGauge represents a no-operation gauge that implements the Gauge interface.
// It is designed to be a passive placeholder, ignoring all update calls and tag manipulations.
// This can be useful in scenarios where gauge functionality is optional or disabled.
type nopGauge struct{}

// Gauge is a no-operation gauge metric implementation.
// It provides empty methods for updating and tagging, useful as a default or placeholder.
// It implements the interfaces.Gauge interface.
var Gauge = &nopGauge{}

// Update is a no-operation method for updating the gauge value.
// It takes a context and a float64 value but does nothing with them.
// This method is part of the Gauge interface implementation.
func (n *nopGauge) Update(_ context.Context, _ float64) {}

// AddTag adds a single tag to the gauge instance and returns the modified gauge.
// The key and value are used to associate metadata with the gauge.
// It follows the same naming convention as WithTags for keys.
func (n *nopGauge) AddTag(_ string, _ string) interfaces.Gauge { return n }

// WithTags initializes all tags from a map for the gauge instance, returning the gauge itself.
// It follows the same tag naming constraints as AddTag.
// Tags starting with __ will be automatically escaped.
// Parameters:
//
//	tags - A map of string key-value pairs representing tags to be added to the gauge.
//
// Returns:
//
//	The gauge instance with updated tags.
func (n *nopGauge) WithTags(_ map[string]string) interfaces.Gauge { return n }
