package tracing

import (
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/trace"
)

type ApplicationTracer trace.TracerProvider

func NewApplicationTracer(exporter *jaeger.Exporter) (ApplicationTracer, error) {
	return NewJaegerTracerProvider("application", exporter)
}

type RepositoryTracer trace.TracerProvider

func NewRepositoryTracer(exporter *jaeger.Exporter) (RepositoryTracer, error) {
	return NewJaegerTracerProvider("repository", exporter)
}

type HttpTracer trace.TracerProvider

func NewHttpTracer(exporter *jaeger.Exporter) (HttpTracer, error) {
	return NewJaegerTracerProvider("http", exporter)
}

type GormTracer trace.TracerProvider

func NewGormTracer(exporter *jaeger.Exporter) (GormTracer, error) {
	provider, err := NewJaegerTracerProvider("gorm", exporter)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
