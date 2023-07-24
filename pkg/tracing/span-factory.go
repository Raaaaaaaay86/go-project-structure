package tracing

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func HttpSpanFactory(tracerProvider trace.TracerProvider, ctx *gin.Context, fullPackage string) (context.Context, trace.Span) {
	return tracerProvider.Tracer(fullPackage).Start(ctx, fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.Path))
}

func RecordHttpError(span trace.Span, code int, err error) {
	span.RecordError(err)
	span.SetAttributes(attribute.Int("http.status_code", code))
}

func ApplicationSpanFactory(tracerProvider trace.TracerProvider, ctx context.Context, fullPackage string, method string) (context.Context, trace.Span) {
	return tracerProvider.Tracer(fullPackage).Start(ctx, method)
}

func RepositorySpanFactory(tracerProvider trace.TracerProvider, ctx context.Context, fullPackage string, method string) (context.Context, trace.Span) {
	return tracerProvider.Tracer(fullPackage).Start(ctx, method)
}
