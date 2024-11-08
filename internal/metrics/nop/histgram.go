package nop

import (
	"context"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"time"
)

var _ interfaces.Histogram = (*nopHistogram)(nil)

type nopHistogram struct{}

var Histogram = &nopHistogram{}

func (n *nopHistogram) Update(_ context.Context, _ time.Duration) {}

func (n *nopHistogram) UpdateInSeconds(_ context.Context, _ float64) {}

func (n *nopHistogram) UpdateInMilliseconds(_ context.Context, _ float64) {}

func (n *nopHistogram) UpdateSine(_ context.Context, _ time.Time) {}

func (n *nopHistogram) Time(_ func()) {}

func (n *nopHistogram) AddTag(_ string, _ string) interfaces.Histogram { return n }

func (n *nopHistogram) WithTags(_ map[string]string) interfaces.Histogram { return n }
