package tracing

import (
	"github.com/raaaaaaaay86/go-project-structure/pkg/configs"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"os"
)

func newTraceResource(serviceName string) (*resource.Resource, error) {
	attrs := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)
	return resource.Merge(resource.Default(), attrs)
}

func newStdOutExporter() (*stdouttrace.Exporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(os.Stdout),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
}

func NewJaegerExporter(config *configs.YamlConfig) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Jaeger.Endpoint)))
}

func NewStdOutTracerProvider(serviceName string) *traceSdk.TracerProvider {
	traceResource, err := newTraceResource(serviceName)
	if err != nil {
		panic(err)
	}
	traceExporter, err := newStdOutExporter()
	if err != nil {
		panic(err)
	}
	return traceSdk.NewTracerProvider(traceSdk.WithResource(traceResource), traceSdk.WithBatcher(traceExporter))
}

func NewEmptyTracerProvider() *traceSdk.TracerProvider {
	traceResource, err := newTraceResource("")
	if err != nil {
		panic(err)
	}
	return traceSdk.NewTracerProvider(traceSdk.WithResource(traceResource))
}

func NewJaegerTracerProvider(serviceName string, exporter *jaeger.Exporter) (trace.TracerProvider, error) {
	traceResource, err := newTraceResource(serviceName)
	if err != nil {
		return nil, err
	}

	return traceSdk.NewTracerProvider(traceSdk.WithResource(traceResource), traceSdk.WithBatcher(exporter)), nil
}
