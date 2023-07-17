package tracing

import (
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
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

func newJaegerExporter() (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
}

func NewStdOutTracerProvider(serviceName string) *trace.TracerProvider {
	traceResource, err := newTraceResource(serviceName)
	if err != nil {
		panic(err)
	}
	traceExporter, err := newStdOutExporter()
	if err != nil {
		panic(err)
	}
	return trace.NewTracerProvider(trace.WithResource(traceResource), trace.WithBatcher(traceExporter))
}

func NewEmptyTracerProvider(serviceName string) *trace.TracerProvider {
	traceResource, err := newTraceResource(serviceName)
	if err != nil {
		panic(err)
	}
	return trace.NewTracerProvider(trace.WithResource(traceResource))
}

func NewJaegerTracerProvider(serviceName string) *trace.TracerProvider {
	traceResource, err := newTraceResource(serviceName)
	if err != nil {
		panic(err)
	}
	traceExporter, err := newJaegerExporter()
	if err != nil {
		panic(err)
	}
	return trace.NewTracerProvider(trace.WithResource(traceResource), trace.WithBatcher(traceExporter))
}
