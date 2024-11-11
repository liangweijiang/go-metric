package meter

import (
	"github.com/liangweijiang/go-metric/pkg/config"
	"github.com/liangweijiang/go-metric/pkg/interfaces"
	"time"
)

// envOption encapsulates a MeterEnv to be used as an option for configuring environments in a Config structure.
type envOption struct {
	evn config.MeterEnv
}

// ApplyConfig sets the environment in the provided config to the one stored in the envOption instance.
func (e *envOption) ApplyConfig(cfg *config.Config) {
	cfg.Env = e.evn
}

// WithEnv returns an Option that sets the environment in the Config to the specified MeterEnv.
func WithEnv(env config.MeterEnv) interfaces.Option {
	return &envOption{
		evn: env,
	}
}

// reportMetricOption defines a configuration option for reporting metrics with a specific local IP address and port.
type reportMetricOption struct {
	localIp string
	port    int
}

// ApplyConfig applies the report metric option settings to the provided configuration.
// It sets the ReportMetricPort and LocalIP fields of the given config.Config instance.
// Parameters:
//   - cfg: Pointer to the config.Config to be updated with report metric options.
//
// Returns:
//
//	None
func (p *reportMetricOption) ApplyConfig(cfg *config.Config) {
	cfg.ReportMetricPort = p.port
	cfg.LocalIP = p.localIp
}

// WithReportMetric creates an Option to configure the reporting of metrics with a specified local IP and port.
// This option sets the 'LocalIP' and 'ReportMetricPort' fields in the provided Config.
// Parameters:
//
//	localIp (string): The local IP address to use for metric reporting.
//	port (int): The network port for metric reporting.
//
// Returns:
//
//	interfaces.Option: An Option instance to apply the report metric configuration.
func WithReportMetric(localIp string, port int) interfaces.Option {
	return &reportMetricOption{
		localIp: localIp,
		port:    port,
	}
}

// meterProviderOption encapsulates configuration for selecting a specific meter provider type.
type meterProviderOption struct {
	providerType config.MeterProviderType
}

// ApplyConfig applies the meter provider type from the option to the given configuration's MeterProvider field.
// It sets the MeterProvider field of the provided config to the value stored in the meterProviderOption instance.
// This method is used to configure the desired meter provider type within the config structure.
// Parameters:
// cfg (*config.Config): The configuration object to update with the meter provider type.
// Returns:
// None
func (m *meterProviderOption) ApplyConfig(cfg *config.Config) {
	cfg.MeterProvider = m.providerType
}

// WithProviderType returns an Option that sets the meter provider type in a Config.
// It wraps the specified providerType within a meterProviderOption which implements the Option interface.
// The ApplyConfig method of this option updates the MeterProvider field of a config.Config instance.
// Parameters:
// providerType (config.MeterProviderType): The type of meter provider to configure.
// Returns:
// interfaces.Option: An Option to apply the meter provider type to a config.Config during setup.
func WithProviderType(providerType config.MeterProviderType) interfaces.Option {
	return &meterProviderOption{
		providerType: providerType,
	}
}

// baseTagsOption holds a set of base tags to be applied to configurations.
type baseTagsOption struct {
	baseTags map[string]string
}

// ApplyConfig sets the base tags from the baseTagsOption instance into the provided config.Config's BaseTags field.
func (b *baseTagsOption) ApplyConfig(cfg *config.Config) {
	cfg.BaseTags = b.baseTags
}

// WithBaseTags creates an Option that sets the base tags for metric configuration.
// It takes a map of string keys to string values which represent the base tags.
// These tags will be applied to all metrics by the config consumer.
// Returns an interfaces.Option instance that can be used to configure a config.Config instance.
func WithBaseTags(baseTags map[string]string) interfaces.Option {
	return &baseTagsOption{
		baseTags: baseTags,
	}
}

// pushGatewayOption holds configuration parameters for a Push Gateway integration, including the gateway address and the push period.
type pushGatewayOption struct {
	address string
	period  time.Duration
}

// ApplyConfig applies the push gateway configuration options to the provided config instance.
// It sets the GatewayAddress and PushPeriod within the config's PushGateway field.
// Parameters:
// cfg (*config.Config): The configuration to be updated with push gateway settings.
// Returns:
// None
func (p *pushGatewayOption) ApplyConfig(cfg *config.Config) {
	cfg.PushGateway = &config.PushGatewayCfg{
		GatewayAddress: p.address,
		PushPeriod:     p.period,
	}
}

// WithPushGateway creates an Option that configures the address and push period for a Push Gateway integration.
// It returns an Option interface that can be applied to a config instance to set the GatewayAddress and PushPeriod fields within the PushGateway configuration.
// Parameters:
// address (string): The address of the Push Gateway.
// period (time.Duration): The interval at which metrics should be pushed to the gateway.
// Returns:
// interfaces.Option: An Option to configure Push Gateway settings in a Config.
func WithPushGateway(address string, period time.Duration) interfaces.Option {
	return &pushGatewayOption{
		address: address,
		period:  period,
	}
}

// histogramBoundariesOption is a configuration option for setting histogram boundary values used to define data buckets in a metrics setup.
type histogramBoundariesOption struct {

	// boundaries contains the boundary values for histogram data buckets. It is used to configure the histogram boundaries in a metrics configuration.
	boundaries []float64
}

// ApplyConfig applies the histogram boundary values stored in the histogramBoundariesOption to the provided config.Config instance.
func (h *histogramBoundariesOption) ApplyConfig(cfg *config.Config) {
	cfg.HistogramBoundaries = h.boundaries
}

// WithHistogramBoundaries creates an Option to set custom histogram bucket boundaries for metric configurations.
// It returns an interfaces.Option that applies the provided float64 slice boundaries to the HistogramBoundaries field of a config.Config when applied.
func WithHistogramBoundaries(boundaries []float64) interfaces.Option {
	return &histogramBoundariesOption{
		boundaries: boundaries,
	}
}

// infoLogOption allows customization of the info log write function within a configuration.
// It holds a function that accepts a string message intended for informational logging.
type infoLogOption struct {

	// infoLogFunc is a function type that accepts a string argument and is used for writing informational log messages.
	infoLogFunc func(s string)
}

func (i *infoLogOption) ApplyConfig(cfg *config.Config) {
	cfg.InfoLogWrite = i.infoLogFunc
}

func WithInfoLogWrite(logFunc func(s string)) interfaces.Option {
	return &infoLogOption{
		infoLogFunc: logFunc,
	}
}

// errorLogOption holds a function to handle error logging.
// It is used as an option to configure the error logging behavior within a Config instance.
type errorLogOption struct {

	// errorLogFunc is a function type that accepts a string argument for logging error messages. It is used to customize error logging behavior within the application.
	errorLogFunc func(s string)
}

// ApplyConfig applies the error logging function from errorLogOption to the provided config's ErrorLogWrite field.
// This sets the custom error logging behavior for the application.
// Parameters:
// cfg (*config.Config): The configuration object to which the error logging function will be applied.
// Returns:
// None
func (e *errorLogOption) ApplyConfig(cfg *config.Config) {
	cfg.ErrorLogWrite = e.errorLogFunc
}

// WithErrorLogWrite returns an Option that sets the error logging function for a Config instance.
// The provided logFunc will be used to handle error messages within the application.
// Parameters:
// logFunc (func(s string)): A function that takes a string message and logs the error.
// Returns:
// interfaces.Option: An Option interface to apply the error logging function to a Config.
func WithErrorLogWrite(logFunc func(s string)) interfaces.Option {
	return &errorLogOption{
		errorLogFunc: logFunc,
	}
}
