package interfaces

import "net/http"

// BaseMeter defines an interface for creating and managing metric instruments like counters, up-down counters, gauges, and histograms.
// It also allows controlling the SDKS's running state and provides an HTTP handler for metric exposition.
type BaseMeter interface {
	//GetHandler 返回http handler
	GetHandler() http.Handler
	// WithRunning 设置为false，SDK切换为空实现，关闭指标的收集功能
	WithRunning(on bool)
	NewCounter(metricName, desc, unit string) Counter
	NewUpDownCounter(metricName, desc, unit string) UpDownCounter
	NewGauge(metricName, desc, unit string) Gauge
	NewHistogram(metricName, desc, unit string) Histogram
}

// Meter extends the BaseMeter interface, adding the capability to retrieve the components
// associated with the middleware for observability purposes, such as monitoring and distributed tracing.
type Meter interface {
	BaseMeter
	// Components() Components // 返回中间件埋点方法
}

// MeterServer defines an interface for a metric server that can start and stop its service.
// Implementations of this interface should handle the lifecycle of a metrics collection and reporting endpoint.
type MeterServer interface {
	Start()
	Stop()
}
