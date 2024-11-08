package interfaces

import (
	"context"
	"time"
)

// Counter is an interface for incrementing a metric by a given delta and managing tags.
// It provides methods to increment the counter, add tags individually or in bulk, adhering to naming constraints.
type Counter interface {
	Incr(ctx context.Context, delta float64)
	IncrOne(ctx context.Context)
	// AddTag 单次增加一组tag
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	AddTag(key string, value string) Counter
	// WithTags 以map全量初始化所有tags
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	WithTags(tags map[string]string) Counter
}

// UpDownCounter represents an instrument that supports incrementing and decrementing a value.
// It is designed to track quantities that can go both up and down, such as the number of active users in a system.
// The interface includes methods to update the counter by a given delta, increment or decrement by one, and manage tags for added context.
type UpDownCounter interface {
	Update(ctx context.Context, delta float64)
	IncrOne(ctx context.Context)
	DecrOne(ctx context.Context)
	// AddTag 单次增加一组tag
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	AddTag(key string, value string) UpDownCounter
	// WithTags 以map全量初始化所有tags
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	WithTags(tags map[string]string) UpDownCounter
}

// Histogram defines an interface for recording the distribution of values, such as timing events or other measured values.
// It supports updating with different time units and offers tagging capabilities for adding metadata to measurements.
type Histogram interface {
	// Update 记录一段时间耗时
	Update(ctx context.Context, d time.Duration)
	// UpdateInSeconds 记录一段单位秒时间耗时
	UpdateInSeconds(ctx context.Context, s float64)
	// UpdateInMilliseconds 记录一段单位毫秒时间耗时
	UpdateInMilliseconds(ctx context.Context, m float64)
	// UpdateSine 记录从某个时间开始的耗时
	UpdateSine(ctx context.Context, start time.Time)
	// Time 记录函数执行的耗时
	Time(f func())
	// AddTag 单次增加一组tag
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	AddTag(key string, value string) Histogram
	// WithTags 以map全量初始化所有tags
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	WithTags(tags map[string]string) Histogram
}

// Gauge is an interface representing a metric gauge which can be updated to track the current value of a measurable attribute.
// It supports adding tags to provide additional context to the gauge readings dynamically.
type Gauge interface {
	Update(ctx context.Context, v float64)
	// AddTag 单次增加一组tag
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	AddTag(key string, value string) Gauge
	// WithTags 以map全量初始化所有tags
	// 不能以 __ 双下划线开头, 否则会自动转义，(^[a-zA-Z_][a-zA-Z0-9_]*$)
	WithTags(tags map[string]string) Gauge
}
