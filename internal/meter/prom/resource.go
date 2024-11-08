package prom

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

// ResourceWithAttr creates a new OpenTelemetry resource with the provided custom attributes.
// This function allows you to add additional resource attributes to the OpenTelemetry resource.
//
// The function takes a slice of attribute.KeyValue as input, where each KeyValue represents a custom attribute.
// The function returns a pointer to the created resource.Resource and an error if any.
//
// The created resource includes attributes discovered from environment variables (OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME),
// information about the OpenTelemetry SDK used, process information, OS information, container information, host information,
// and the custom attributes provided as input.
//
// Note: You can optionally add your own external Detector implementation by uncommenting the corresponding line in the function.
func ResourceWithAttr(attributes []attribute.KeyValue) (*resource.Resource, error) {
	res, err := resource.New(
		context.Background(),
		resource.WithFromEnv(),                 // Discover and provide attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables.
		resource.WithTelemetrySDK(),            // Discover and provide information about the OpenTelemetry SDK used.
		resource.WithProcess(),                 // Discover and provide process information.
		resource.WithOS(),                      // Discover and provide OS information.
		resource.WithContainer(),               // Discover and provide container information.
		resource.WithHost(),                    // Discover and provide host information.
		resource.WithAttributes(attributes...), // Add custom resource attributes.
		// resource.WithDetectors(third_party.Detector{}),           // Bring your own external Detector implementation.
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}
