package nop

import (
	"context"
	"go-mertric/pkg/interfaces"
)

// _ is a blank identifier used for type assertion to ensure that nopCounter satisfies the interfaces.Counter interface requirements.
var _ interfaces.Counter = (*nopCounter)(nil)

// nopCounter represents a no-operation (NOP) counter that implements the Counter interface.
// It is designed to not perform any operations, effectively acting as a placeholder or disabled counter.
type nopCounter struct{}

// Counter is a no-operation counter instance, useful as a default or placeholder.
// It implements the interfaces.Counter interface, providing empty methods for incrementing
// and adding tags, which have no effect.
var Counter = &nopCounter{}

// Incr increments the counter by the given value. This method does nothing as it's part of a no-operation (NOP) counter.
func (n *nopCounter) Incr(_ context.Context, _ float64) {}

// IncrOne increments the counter by one. This method is a part of the `nopCounter` struct and does not perform any operation, serving as a no-op.
func (n *nopCounter) IncrOne(_ context.Context) {}

// AddTag adds a tag to the counter instance, returning the counter itself.
// It adheres to the tag key-value format validation rules defined by the Counter interface.
func (n *nopCounter) AddTag(_ string, _ string) interfaces.Counter { return n }

// WithTags initializes all tags for the counter using the provided map. It adheres to the same tag key-value format validation rules. This method is part of the no-operation logic and returns the receiver as is.
func (n *nopCounter) WithTags(_ map[string]string) interfaces.Counter { return n }
