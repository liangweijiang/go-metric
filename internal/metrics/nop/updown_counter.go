package nop

import (
	"context"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
)

// _ is a blank identifier assignment to assert that (*nopUpDownCounter)(nil) implements the interfaces.UpDownCounter interface.
var _ interfaces.UpDownCounter = (*nopUpDownCounter)(nil)

// nopUpDownCounter represents a no-operation implementation of the UpDownCounter interface.
// It provides empty methods for updating, incrementing, decrementing, and managing tags,
// effectively serving as a placeholder or a disabled counter.
type nopUpDownCounter struct{}

// UpDownCounter is an instance of nopUpDownCounter, providing no-op implementations for updating, incrementing, decrementing, and adding tags to up-down counter operations.
var UpDownCounter = &nopUpDownCounter{}

// Update adjusts the counter by the given delta. This method is a no-op and does nothing.
func (n *nopUpDownCounter) Update(_ context.Context, _ float64) {}

// IncrOne increments the up-down counter by one. This is a no-operation implementation.
func (n *nopUpDownCounter) IncrOne(_ context.Context) {}

// DecrOne decrements the up-down counter by one. This method is a no-operation implementation.
func (n *nopUpDownCounter) DecrOne(_ context.Context) {}

// AddTag adds a tag to the up-down counter instance.
// It returns the same nopUpDownCounter instance for method chaining.
// Tags are ignored in this no-operation implementation.
func (n *nopUpDownCounter) AddTag(_ string, _ string) interfaces.UpDownCounter { return n }

// WithTags returns a new UpDownCounter with the provided tags set. This operation is a no-op and the original instance is returned unmodified.
func (n *nopUpDownCounter) WithTags(_ map[string]string) interfaces.UpDownCounter { return n }
