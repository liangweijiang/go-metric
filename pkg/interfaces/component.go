package interfaces

import (
	"context"
	"time"
)

// ComponentHistogram 用来记录中间件的直方图
type ComponentHistogram interface {
	// Update 记录一段时间耗时
	Update(ctx context.Context, d time.Duration)
	// UpdateInSeconds 记录一段单位秒时间耗时
	UpdateInSeconds(ctx context.Context)
	// UpdateInMilliseconds 记录一段单位毫秒时间耗时
	UpdateInMilliseconds(ctx context.Context)
	// UpdateSine 记录从某个时间开始的耗时
	UpdateSine(ctx context.Context, start time.Time)
}

// ComponentCounter 计数器
type ComponentCounter interface {
	Incr(ctx context.Context, delta float64)
	IncrOne(ctx context.Context)
}

// ComponentUpDownCounter 增减计数器可以增加和减少
type ComponentUpDownCounter interface {
	Update(ctx context.Context, delta float64)
	IncrOne(ctx context.Context)
	DecrOne(ctx context.Context)
}
