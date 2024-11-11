package meter

import (
	"github.com/liangweijiang/go-metric/internal/global"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
)

// GetGlobalMeter retrieves the globally configured instance of interfaces.Meter, enabling creation of metric instruments and management of observability features across the application.
func GetGlobalMeter() interfaces.Meter {
	return global.GetMeter()
}

// SetGlobalMeter sets the provided meter as the global meter if it's not nil, enabling metric instrumentation's.
// This function ensures thread safety during the update process using atomic operations.
// Parameters:
//
//	meter (interfaces.Meter): The meter implementation to set as global, extending BaseMeter capabilities with component retrieval.
func SetGlobalMeter(meter interfaces.Meter) {
	global.SetMeter(meter)
}
