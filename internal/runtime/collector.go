package runtime

import (
	"context"
	"github.com/liangweijiang/go-metric/pkg/config"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"github.com/liangweijiang/go-metric/pkg/utils"
	"runtime"
	"runtime/metrics"
	"sync/atomic"
	"time"
)

// defaultRuntimeCollectInterval defines the default interval at which runtime metrics are collected.
// It is set to 10 seconds.
const defaultRuntimeCollectInterval = time.Second * 10

// collector encapsulates the logic for collecting and managing runtime metrics based on a provided configuration.
// It holds onto configuration settings, a metrics Meter instance, an atomic flag indicating its running state,
// and a channel to signal closure for clean shutdown. Additionally, it caches the last collected runtime memory statistics.
type collector struct {
	cfg     *config.Config
	meter   interfaces.Meter
	running int32
	closeCh chan struct{}
	// runtime cached info
	msLast *runtime.MemStats
}

// NewRuntimeCollector initializes and returns a new runtime metric collector.
// It takes a configuration pointer and a meter interface to set up the collector.
// The collector is designed to gather runtime metrics based on the provided configuration settings.
func NewRuntimeCollector(cfg *config.Config, meter interfaces.Meter) interfaces.MetricCollector {
	return &collector{
		cfg:     cfg,
		meter:   meter,
		running: 0,
		closeCh: make(chan struct{}),
	}
}

// Start initiates the collection of runtime metrics if they are enabled in the configuration.
// It sets the running state to prevent multiple starts and spawns a goroutine to execute the Collect method.
// If the metrics collection is already running or disabled, it logs the appropriate message and exits.
func (c *collector) Start() {
	if !c.cfg.RuntimeMetricsCollect {
		c.cfg.WriteErrorOrNot("runtime metrics collect is disabled")
		return
	}
	c.cfg.WriteInfoOrNot("runtime metrics collect is enabled")
	if !atomic.CompareAndSwapInt32(&c.running, 0, 1) {
		c.cfg.WriteErrorOrNot("runtime metrics collect is already running")
		return
	}
	go c.Collect()
}

// Collect continuously fetches runtime metrics at a predefined interval until a stop signal is received.
// It initiates a ticker that triggers the collection process, which involves calling `collectRuntimeMetric`.
// The method stops when a signal is sent through `closeCh`.
func (c *collector) Collect() {
	c.cfg.WriteInfoOrNot("start runtime metrics collect")
	ticker := time.NewTicker(defaultRuntimeCollectInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.closeCh:
			c.cfg.WriteInfoOrNot("stop runtime metrics collect")
			return
		case <-ticker.C:
			c.collectRuntimeMetric()
		}
	}
}

// Stop halts the runtime metrics collection process.
// It atomically sets the running state to stopped and signals the collection goroutine to terminate.
// Returns without action if the collector is not currently running.
func (c *collector) Stop() {
	if !atomic.CompareAndSwapInt32(&c.running, 1, 0) {
		c.cfg.WriteErrorOrNot("runtime metrics collect is not running")
		return
	}
	c.closeCh <- struct{}{}
	c.cfg.WriteErrorOrNot("stop runtime metrics collect")
}

// collectRuntimeMetric fetches current readings for all available runtime metrics,
// converts them into the appropriate OpenTelemetry metric types (Gauge, Counter, UpDownCounter),
// and updates them within the collector's meter, ensuring metric names are sanitized for compatibility.
func (c *collector) collectRuntimeMetric() {
	// Get descriptions for all supported metrics.
	descs := metrics.All()
	samples := make([]metrics.Sample, len(descs))
	for i := range samples {
		samples[i].Name = descs[i].Name
	}

	// Sample the metrics. Re-use the samples slice if you can!
	metrics.Read(samples)

	for i, sample := range samples {
		name, value := sample.Name, sample.Value
		if !descs[i].Cumulative {
			switch value.Kind() {
			case metrics.KindUint64:
				c.newSystemGauge(utils.SanitizeMetricName(name)).Update(context.Background(), float64(sample.Value.Uint64()))
			default:
			}
			continue
		}

		switch value.Kind() {
		case metrics.KindUint64:
			c.newSystemCounter(utils.SanitizeMetricName(name)).Incr(context.Background(), float64(sample.Value.Uint64()))
		case metrics.KindFloat64:
			c.newSystemUpDownCounter(utils.SanitizeMetricName(name)).Update(context.Background(), float64(sample.Value.Float64()))
		case metrics.KindFloat64Histogram:

		case metrics.KindBad:

		default:
		}
	}
}

// newSystemGauge creates a new system Gauge metric with the specified name and tags it as a base metric type.
// It utilizes the collector's meter to instantiate the Gauge.
// param metricName: The name of the gauge metric.
// return: An interfaces.Gauge instance configured as a system metric.
func (c *collector) newSystemGauge(metricName string) interfaces.Gauge {
	return c.meter.NewGauge(metricName, "system metric", "").AddTag("metric_type", "base")
}

// newSystemUpDownCounter creates a new UpDownCounter instrument for system metrics with a specified name.
// It adds a default tag "metric_type" with the value "base" to provide context about the counter's nature.
// This counter is capable of both incrementing and decrementing to track values that can rise and fall.
// Parameters:
// - metricName: The name of the metric for which the UpDownCounter is being created.
// Returns:
// - An instance of UpDownCounter configured for system metrics use.
func (c *collector) newSystemUpDownCounter(metricName string) interfaces.UpDownCounter {
	return c.meter.NewUpDownCounter(metricName, "system metric", "").AddTag("metric_type", "base")
}

func (c *collector) newSystemCounter(metricName string) interfaces.Counter {
	return c.meter.NewCounter(metricName, "system metric", "").AddTag("metric_type", "base")
}
