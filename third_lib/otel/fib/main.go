package main

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"io"
	"log"
	"os"
	"os/signal"
)

func main() {
	logger := log.New(os.Stdout, "", 0)
	f, err := os.Create("traces.txt")
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	exporter, err := newExporter(f)
	if err != nil {
		logger.Fatal(err)
	}
	provider := trace.NewTracerProvider(trace.WithBatcher(exporter), trace.WithResource(newResource()))
	defer func() {
		if err = provider.Shutdown(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()

	otel.SetTracerProvider(provider)

	errCh := make(chan error, 1)
	app := NewApp(os.Stdin, logger)
	go func() {
		errCh <- app.Run(context.Background())
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	select {
	case err = <-errCh:
		if err != nil {
			logger.Fatal(err)
		}
	case <-sig:
		logger.Println("\n1" +
			"goodbye")
		return
	}
}

func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithoutTimestamps(),
	)
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
