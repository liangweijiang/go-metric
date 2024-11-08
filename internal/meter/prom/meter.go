package prom

import (
	"github.com/liangweijiang/go-metric/internal/meter/prom/server"
	"github.com/liangweijiang/go-metric/internal/metrics/nop"
	"github.com/liangweijiang/go-metric/internal/metrics/prom"
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

// sdkVersion denotes the current version of the SDK.
// prometheusMeterName is the metric name used for Prometheus instrumentation.
const (
	sdkVersion          = "1.0"
	prometheusMeterName = "go-metrics/prometheus-meter"
)

// PrometheusMeter is a struct that encapsulates the configuration and runtime state for a Prometheus-based metrics system.
// It includes channels to manage the running state, references to meter servers, and an HTTP handler for exposing metrics.
// The struct utilizes a provided configuration to initialize and manage Prometheus metrics collection.
type PrometheusMeter struct {
	cfg     *config.Config
	running int32
	onCh    chan struct{}
	offCh   chan struct{}
	meter   api.Meter
	servers []interfaces.MeterServer
	handler http.Handler
}

// NewPrometheusMeter initializes and configures a Prometheus-based meter provider
// using the provided configuration. It sets up a meter with explicit bucket
// histograms for instrument kind Histogram as defined in the config, along with
// a Prometheus exporter and HTTP handler. Optionally, it can also configure a
// PushGateway server based on the config settings.
//
// Parameters:
//   - cfg (*config.Config): Configuration object for meter setup.
//
// Returns:
//   - interfaces.Meter: An instance of PrometheusMeter ready to record metrics.
//   - error: If there's an error during initialization, such as failed exporter creation or resource setup.
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
	return promMeter, nil
}

// GetHandler returns the HTTP handler for serving Prometheus metrics exposed by the PrometheusMeter.
func (p *PrometheusMeter) GetHandler() http.Handler {
	return p.handler
}

// WithRunning sets the running state of the PrometheusMeter.
// If 'on' is true, it sends a signal on the 'onCh' channel.
// If 'on' is false, it sends a signal on the 'offCh' channel.
// Channels are used with a non-blocking send to avoid potential blocking if the channels are full.
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

// NewCounter creates a new Counter metric if the PrometheusMeter is running.
// It requires the metric name, description, and unit as parameters.
// Returns a no-op Counter if the meter is not active.
// Parameters:
// metricName (string): The name of the counter metric.
// desc (string): Description for the counter metric.
// unit (string): The unit of measurement for the counter.
// Returns:
// interfaces.Counter: An interface representing the counter metric.
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

// NewUpDownCounter creates a new UpDownCounter metric if the PrometheusMeter is running.
// It requires metricName, description, and unit as parameters to define the metric.
// Returns an interfaces.UpDownCounter or a no-op counter if the meter is not running.
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

// NewGauge creates a new Gauge metric if the PrometheusMeter is running.
// It takes the metric name, description, and unit as input parameters and returns an interfaces.Gauge.
// If the meter is not running, it returns a no-op Gauge.
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

// NewHistogram creates a new Histogram metric if the PrometheusMeter is running.
// It takes the metric name, description, and unit as input parameters and returns an interfaces.Histogram.
// If the meter is not running, it returns a no-op Histogram.
func (p *PrometheusMeter) NewHistogram(metricName, desc, unit string) interfaces.Histogram {
	if !p.isRunning() {
		return nop.Histogram
	}
	histogram, err := p.meter.Float64Histogram(metricName,
		api.WithDescription(desc),
		api.WithUnit(unit))
	if err != nil {
		p.cfg.WriteInfoOrNot("failed to create prometheus histogram: " + err.Error())
		return nop.Histogram
	}
	return prom.NewHistogram(metricName, histogram)
}

// isRunning checks if the PrometheusMeter is currently running.
// It reads the atomic 'running' flag to determine the state.
// Returns true if the meter is running, false otherwise.
func (p *PrometheusMeter) isRunning() bool {
	return atomic.LoadInt32(&p.running) == 1
}
