package nop

import (
	"github.com/liangweijiang/go-metric/internal/metrics/nop"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"net/http"
)

var _ interfaces.Meter = (*Meter)(nil)

type Meter struct{}

func NewNopMeter() interfaces.Meter {
	return &Meter{}
}

func (n *Meter) GetHandler() http.Handler {
	return nil
}

func (n *Meter) WithRunning(_ bool) {

}

func (n *Meter) NewCounter(_, _, _ string) interfaces.Counter {
	return nop.Counter
}

func (n *Meter) NewUpDownCounter(_, _, _ string) interfaces.UpDownCounter {
	return nop.UpDownCounter
}

func (n *Meter) NewGauge(_, _, _ string) interfaces.Gauge {
	return nop.Gauge
}

func (n *Meter) NewHistogram(_, _, _ string) interfaces.Histogram {
	return nop.Histogram
}
