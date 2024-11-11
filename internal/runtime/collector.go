package runtime

import (
	"context"
	"github.com/liangweijiang/go-metric/pkg/config"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"github.com/liangweijiang/go-metric/pkg/utils"
	"runtime"
	"runtime/metrics"
)

type collector struct {
	cfg     *config.Config
	meter   interfaces.Meter
	running int32
	closeCh chan struct{}
	// runtime cached info
	msLast *runtime.MemStats
}

func (c *collector) Start() {
	//TODO implement me
	panic("implement me")
}

func (c *collector) Stop() {
	//TODO implement me
	panic("implement me")
}

// https://runebook.dev/cn/docs/go/runtime/metrics/index
// collectRuntimeMetric
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
			c.newSystemUpDownCounter(utils.SanitizeMetricName(name)).Update(context.Background(), float64(sample.Value.Uint64()))
		case metrics.KindFloat64Histogram:

		case metrics.KindBad:

		default:
		}
	}
}

func (c *collector) newSystemGauge(metricName string) interfaces.Gauge {
	return c.meter.NewGauge(metricName, "system metric", "").AddTag("metric_type", "base")
}

func (c *collector) newSystemUpDownCounter(metricName string) interfaces.UpDownCounter {
	return c.meter.NewUpDownCounter(metricName, "system metric", "").AddTag("metric_type", "base")
}

func (c *collector) newSystemCounter(metricName string) interfaces.Counter {
	return c.meter.NewCounter(metricName, "system metric", "").AddTag("metric_type", "base")
}
