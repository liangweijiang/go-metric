package prom

import (
	"github.com/liangweijiang/go-metric/internal/tag"
	"go.opentelemetry.io/otel/attribute"
	"sync/atomic"
)

// Base represents a foundational structure within a metrics system, embedding common attributes like a name, tags for context, and a completion status.
type Base struct {
	name      string
	tags      tag.Tags
	completed int32
}

// ready checks if the Base instance is ready for operations by atomically swapping the completed status from 0 to 1.
// It returns true if the swap was successful, indicating the Base is ready; otherwise, false.
// This method ensures thread-safe initialization status checking.
func (b *Base) ready() bool {
	return atomic.CompareAndSwapInt32(&b.completed, 0, 1)
}

// AddTag adds a tag with the specified key and value to the Base's tags collection.
// It appends a new attribute.KeyValue pair to the tags slice.
func (b *Base) AddTag(key, value string) {
	b.tags = append(b.tags, attribute.String(key, value))
}

// WithTags sets the provided tags on the Base instance, appending them to existing tags.
// If the input map is nil or empty, the function does nothing.
// This method is intended to be used to add contextual metadata to metrics.
// Parameters:
//
//	tags: A map of tags to set on the Base instance.
//
// The function modifies the Base instance in place and has no return value.
func (b *Base) WithTags(tags map[string]string) {
	if tags == nil || len(tags) == 0 {
		return
	}
	for k, v := range tags {
		b.AddTag(k, v)
	}
}
