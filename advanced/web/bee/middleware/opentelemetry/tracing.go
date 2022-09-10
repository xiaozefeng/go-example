package opentelemetry

import (
	"github.com/xiaozefeng/go-example/advanced/web/bee"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
)

const defaultInstrumentationName = "go-example/advanced/web/bee/middleware/opentelemetry"

type Builder struct {
	Tracer trace.Tracer
}

func (b Builder) Build() bee.Middleware {
	if b.Tracer == nil {
		b.Tracer = otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	}
	return func(next bee.HandleFunc) bee.HandleFunc {
		return func(ctx *bee.Context) {
			reqCtx := ctx.Request.Context()
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Request.Header))
			reqCtx, span := b.Tracer.Start(reqCtx, "unknown", trace.WithAttributes())
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.Request.Method))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Request.Host))
			span.SetAttributes(attribute.String("http.url", ctx.Request.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Request.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Request.RemoteAddr))
			span.SetAttributes(attribute.String("http.proto", ctx.Request.Proto))

			ctx.Request = ctx.Request.WithContext(reqCtx)
			defer func() {
				if ctx.MatchedRoute != "" {
					span.SetName(ctx.MatchedRoute)
				}
				log.Println("tracing: http.status", ctx.RespStatusCode)
				span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
			}()
			next(ctx)
		}
	}
}
