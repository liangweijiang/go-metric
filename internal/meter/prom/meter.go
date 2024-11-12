package prom

import (
	"github.com/liangweijiang/go-metric/internal/meter/prom/server"
	"github.com/liangweijiang/go-metric/internal/metrics/nop"
	"github.com/liangweijiang/go-metric/internal/metrics/prom"
	"github.com/liangweijiang/go-metric/internal/runtime"
	"github.com/liangweijiang/go-metric/pkg/config"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	cliprom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"net/http"
	"sync/atomic"
)

// sdkVersion represents the current version of the SDK.
// prometheusMeterName is the name used for Prometheus metrics meter.
const (
	sdkVersion          = "1.0"
	prometheusMeterName = "go-metrics/prometheus-meter"
)

// PrometheusMeter encapsulates the configuration and components necessary for managing Prometheus metrics.
// It includes channels for controlling the meter's lifecycle, the primary meter instance,
// a collection of meter servers, an HTTP handler for metrics exposure, and a runtime metric collector.
// This structure facilitates starting and stopping metric collection and export functionalities dynamically.
type PrometheusMeter struct {
	cfg              *config.Config
	running          int32
	onCh             chan struct{}
	offCh            chan struct{}
	meter            api.Meter
	servers          []interfaces.MeterServer
	handler          http.Handler
	runtimeCollector interfaces.MetricCollector
}

// NewPrometheusMeter initializes and configures a Prometheus-based meter for metric collection.
// It sets up a metric registry, exporter, resource, and meter provider based on the provided configuration.
// Additionally, it configures a histogram view and starts a runtime collector.
// If configured, it also sets up servers for pushing metrics to a gateway and serving HTTP requests for metrics.
// Returns a PrometheusMeter instance and an error if any occur during setup.
func NewPrometheusMeter(cfg *config.Config) (interfaces.Meter, error) {
	registry := cliprom.NewRegistry()
	exporter, err := prometheus.New(
		prometheus.WithRegisterer(registry),
		prometheus.WithoutScopeInfo(),
	)
	if err != nil {
		cfg.WriteErrorOrNot("failed to create prometheus exporter: " + err.Error())
		return nil, err
	}

	resource, err := ResourceWithAttr(cfg.WithBaseTags())
	if err != nil {
		cfg.WriteErrorOrNot("failed to create resource: " + err.Error())
		return nil, err
	}
	provider := metric.NewMeterProvider(
		metric.WithResource(resource),
		metric.WithReader(exporter),
		metric.WithView(
			metric.NewView(
				metric.Instrument{
					Kind: metric.InstrumentKindHistogram,
				},
				metric.Stream{
					Aggregation: metric.AggregationExplicitBucketHistogram{
						Boundaries: cfg.HistogramBoundaries,
					},
				},
			),
		),
	)

	meter := provider.Meter(prometheusMeterName, api.WithInstrumentationVersion(sdkVersion), api.WithInstrumentationAttributes())
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	promMeter := &PrometheusMeter{
		cfg:     cfg,
		running: 1,
		onCh:    make(chan struct{}),
		offCh:   make(chan struct{}),
		meter:   meter,
		handler: handler,
	}
	if cfg.PushGateway != nil {
		promMeter.servers = append(promMeter.servers, server.NewPromPushGatewayServer(cfg, registry))
	}
	if cfg.PrometheusPort > 0 {
		promMeter.servers = append(promMeter.servers, server.NewPromHttpServer(cfg, promMeter.GetHandler()))
	}

	promMeter.runtimeCollector = runtime.NewRuntimeCollector(cfg, promMeter)
	promMeter.runtimeCollector.Start()
	for _, meterServer := range promMeter.servers {
		meterServer.Start()
	}

	go promMeter.signalListener()
	return promMeter, nil
}

// signalListener monitors channels to start or stop the PrometheusMeter and its components.
// It listens for signals on `onCh` to start and `offCh` to stop the meter, managing the runtime collector
// and all meter servers accordingly. The method ensures the meter can only be started once and stopped once.
func (p *PrometheusMeter) signalListener() {
	for {
		select {
		case <-p.onCh:
			if !atomic.CompareAndSwapInt32(&p.running, 0, 1) {
				p.cfg.WriteInfoOrNot("prometheus meter is already running")
				return
			}
			p.cfg.WriteInfoOrNot("prometheus meter is started")
			p.runtimeCollector.Start()
			for _, meterServer := range p.servers {
				meterServer.Start()
			}
		case <-p.offCh:
			if !atomic.CompareAndSwapInt32(&p.running, 1, 0) {
				p.cfg.WriteInfoOrNot("prometheus meter is already stopped")
				return
			}
			p.cfg.WriteInfoOrNot("prometheus meter is stopped")
			p.runtimeCollector.Stop()
			for _, meterServer := range p.servers {
				meterServer.Stop()
			}
		}
	}
}

// GetHandler returns the HTTP handler for exposing Prometheus metrics.
// This handler can be used to integrate with HTTP servers to serve metrics data.
// It retrieves the pre-configured http.Handler instance associated with the PrometheusMeter.
func (p *PrometheusMeter) GetHandler() http.Handler {
	return p.handler
}

// WithRunning sets the running state of the PrometheusMeter to the specified boolean value.
// When `on` is true, it attempts to send a signal on the `onCh` channel to start the meter.
// When `on` is false, it tries to send a signal on the `offCh` channel to stop the meter.
// Channels are used with a non-blocking send to avoid blocking the caller if the signals are not immediately processed.
func (p *PrometheusMeter) WithRunning(on bool) {
	if on {
		select {
		case p.onCh <- struct{}{}:
		default:

		}
	} else {
		select {
		case p.offCh <- struct{}{}:
		default:

		}
	}
}

// NewCounter creates a new Counter metric with the specified name, description, and unit.
// It returns a no-op counter if the PrometheusMeter is not running.
// This method uses the underlying meter to create a Float64Counter and wraps it with a custom Counter implementation.
// In case of failure creating the counter, a log message is emitted and a no-op counter is returned.
func (p *PrometheusMeter) NewCounter(metricName, desc, unit string) interfaces.Counter {
	if !p.isRunning() {
		return nop.Counter
	}
	counter, err := p.meter.Float64Counter(
		metricName,
		api.WithDescription(desc),
		api.WithUnit(unit),
	)
	if err != nil {
		p.cfg.WriteInfoOrNot("failed to create prometheus counter: " + err.Error())
		return nop.Counter
	}
	return prom.NewCounter(metricName, counter)
}

// NewUpDownCounter creates a new UpDownCounter metric within the PrometheusMeter.
// It requires a metric name, description, and unit of measure.
// If the PrometheusMeter is not running, it returns a no-op UpDownCounter.
// Otherwise, it initializes a new UpDownCounter with the provided parameters and adds it to the meter.
// Returns an error if the UpDownCounter creation fails within the underlying meter.
func (p *PrometheusMeter) NewUpDownCounter(metricName, desc, unit string) interfaces.UpDownCounter {
	if !p.isRunning() {
		return nop.UpDownCounter
	}
	udCounter, err := p.meter.Float64UpDownCounter(metricName,
		api.WithDescription(desc),
		api.WithUnit(unit),
	)
	if err != nil {
		p.cfg.WriteInfoOrNot("failed to create prometheus upDownCounter: " + err.Error())
		return nop.UpDownCounter
	}
	return prom.NewUpDownCounter(metricName, udCounter)
}

// NewGauge creates a new Gauge metric with the specified name, description, and unit within the PrometheusMeter.
// Returns a no-op Gauge if the PrometheusMeter is not currently running.
// It uses the provided metricName, description, and unit to configure the gauge via the underlying meter.
// In case of an error during gauge creation, a log is emitted and a no-op Gauge is returned.
func (p *PrometheusMeter) NewGauge(metricName, desc, unit string) interfaces.Gauge {
	if !p.isRunning() {
		return nop.Gauge
	}
	gauge, err := p.meter.Float64Gauge(metricName,
		api.WithDescription(desc),
		api.WithUnit(unit))
	if err != nil {
		p.cfg.WriteInfoOrNot("failed to create prometheus gauge: " + err.Error())
		return nop.Gauge
	}
	return prom.NewGauge(metricName, gauge)
}

// NewHistogram creates a new Histogram metric with the specified name, description, and unit within the PrometheusMeter.
// If the PrometheusMeter is not running, it returns a no-op Histogram.
// The method configures the histogram using the underlying meter with explicit bucket boundaries.
// In case of an error during histogram creation, a log message is emitted, and a no-op Histogram is returned.
func (p *PrometheusMeter) NewHistogram(metricName, desc, unit string) interfaces.Histogram {
	if !p.isRunning() {
		return nop.Histogram
	}
	histogram, err := p.meter.Float64Histogram(metricName,
		api.WithDescription(desc),
		api.WithUnit(unit),
		api.WithExplicitBucketBoundaries())
	if err != nil {
		p.cfg.WriteInfoOrNot("failed to create prometheus histogram: " + err.Error())
		return nop.Histogram
	}
	return prom.NewHistogram(metricName, histogram)
}

// isRunning checks if the PrometheusMeter is currently running.
// It returns true if the meter is running, false otherwise.
func (p *PrometheusMeter) isRunning() bool {
	return atomic.LoadInt32(&p.running) == 1
}
