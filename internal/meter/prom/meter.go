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

const (
	sdkVersion          = "1.0"
	prometheusMeterName = "go-metrics/prometheus-meter"
)

// PrometheusMeter 基于prometheus 度量器
type PrometheusMeter struct {
	cfg     *config.Config
	running int32
	onCh    chan struct{}
	offCh   chan struct{}
	meter   api.Meter
	servers []interfaces.MeterServer
	handler http.Handler
}

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

	resource, err := ResourceWithAttr(cfg.WithAttributes())
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

func (p *PrometheusMeter) GetHandler() http.Handler {
	return p.handler
}

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

func (p *PrometheusMeter) isRunning() bool {
	return atomic.LoadInt32(&p.running) == 1
}
