package global

import (
	"github.com/liangweijiang/go-metric/internal/meter/nop"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"sync/atomic"
)

// globalMeter is an atomic value storing a meterStore, providing thread-safe access to a globally accessible meter instance.
// It can be used to set and get the current global meter implementation for observability purposes like monitoring and distributed tracing.
var globalMeter = atomic.Value{}

// meterStore holds a reference to an interfaces.Meter instance, facilitating storage and retrieval operations, typically within a concurrency-safe context.
type meterStore struct {
	meter interfaces.Meter
}

// init sets the global meter to a no-operation (NOP) meter, effectively disabling metric collection.
func init() {
	SetNopMeter()
}

// SetNopMeter replaces the current global meter with a no-operation (NOP) meter instance.
// This function is useful for disabling metric collection without altering code that relies on the global meter.
func SetNopMeter() {
	globalMeter.Store(meterStore{
		meter: nop.NewNopMeter(),
	})
}

// SetMeter sets the global meter if the provided meter is not nil, enabling metric instrumentations.
// It utilizes atomic operations to ensure thread safety when updating the global meter store.
func SetMeter(meter interfaces.Meter) {
	if meter == nil {
		return
	}
	globalMeter.Store(meterStore{
		meter: meter,
	})
}

// GetMeter returns the globally stored instance of interfaces.Meter.
// It utilizes atomic loading to safely retrieve the meter from the globalMeter atomic value.
// This function is designed for accessing the shared meter for creating metric instruments
// and managing observability aspects within an application.
func GetMeter() interfaces.Meter {
	return globalMeter.Load().(meterStore).meter
}
