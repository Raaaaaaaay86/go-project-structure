package tracing

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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

func newJaegerExporter() (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
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

func NewEmptyTracerProvider(serviceName string) *traceSdk.TracerProvider {
	traceResource, err := newTraceResource(serviceName)
	if err != nil {
		panic(err)
	}
	return traceSdk.NewTracerProvider(traceSdk.WithResource(traceResource))
}

func NewJaegerTracerProvider(serviceName string) *traceSdk.TracerProvider {
	traceResource, err := newTraceResource(serviceName)
	if err != nil {
		panic(err)
	}
	traceExporter, err := newJaegerExporter()
	if err != nil {
		panic(err)
	}
	return traceSdk.NewTracerProvider(traceSdk.WithResource(traceResource), traceSdk.WithBatcher(traceExporter))
}

func HttpSpanFactory(tracerProvider *traceSdk.TracerProvider, ctx *gin.Context, fullPackage string) (context.Context, trace.Span) {
	return tracerProvider.Tracer(fullPackage).Start(ctx, fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.Path))
}
