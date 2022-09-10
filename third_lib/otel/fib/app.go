package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"strconv"
)

const name = "fib"

type App struct {
	reader io.Reader
	logger *log.Logger
}

func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{
		reader: r,
		logger: l,
	}
}

func (a *App) Run(ctx context.Context) error {
	for {
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")
		n, err := a.Poll(newCtx)
		if err != nil {
			span.End()
			return err
		}

		a.Write(newCtx, n)
		span.End()
	}
}

func (a *App) Poll(ctx context.Context) (uint, error) {
	// trace
	_, span := otel.Tracer(name).Start(ctx, "Poll")

	// print accesslog
	a.logger.Print("please input fibonacci number: ")

	// read fib number and fill in app.reader
	var n uint
	_, err := fmt.Fscanf(a.reader, "%d\n", &n)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	// convert fib number to string , reject to overflow an int64
	nStr := strconv.FormatUint(uint64(n), 10)

	// store fib number to span attributes
	span.SetAttributes(attribute.String("request.n", nStr))
	return n, nil
}

func (a *App) Write(ctx context.Context, n uint) {
	// trace
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	// define anonymous func
	var f = func(ctx context.Context) (uint64, error) {
		_, span = otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		fibonacci, err := Fibonacci(n)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return 0, err
		}
		return fibonacci, nil
	}
	// call anonymous func
	result, err := f(ctx)
	if err != nil {
		a.logger.Printf("fibonacci(%d)= %v\n", n, err)
	} else {
		a.logger.Printf("fibonacci(%d)= %v\n", n, result)
	}
}
