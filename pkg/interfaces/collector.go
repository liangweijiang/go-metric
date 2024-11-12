package interfaces

// MetricCollector defines an interface for collecting and managing metrics, providing methods to start and stop the collection process.
// Implementations of this interface should handle the gathering and recording of metrics data.
type MetricCollector interface {
	Start()
	Stop()
}
