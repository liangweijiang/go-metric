package constant

type MeterProviderType int

const (
	MeterProviderTypePrometheus MeterProviderType = iota + 1
)

const (
	SdkVersion          = "0.1"
	PrometheusMeterName = "go-metrics/prometheus-meter"
)
