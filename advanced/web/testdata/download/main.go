package main

import (
	"fmt"
	"github.com/xiaozefeng/go-example/advanced/web/bee"
	"github.com/xiaozefeng/go-example/advanced/web/bee/middleware/accesslog"
	"github.com/xiaozefeng/go-example/advanced/web/bee/middleware/opentelemetry"
	"github.com/xiaozefeng/go-example/advanced/web/bee/middleware/prometheus"
	"github.com/xiaozefeng/go-example/advanced/web/bee/middleware/recovery"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

func RepeatBody() bee.Middleware {
	return func(next bee.HandleFunc) bee.HandleFunc {
		return func(ctx *bee.Context) {
			ctx.Request.Body = io.NopCloser(ctx.Request.Body)
			next(ctx)
		}
	}
}

func initZipKinProvider() error {
	exporter, err := zipkin.New("http://uat20:31732/api/v2/spans", zipkin.WithLogger(log.New(os.Stderr, "bee-server", log.Ldate|log.Ltime|log.Llongfile)))
	if err != nil {
		return err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("bee-server"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func initJaegerProvider() error {
	url := "http://uat20:31305/api/traces"
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(

			semconv.SchemaURL,
			semconv.ServiceNameKey.String("bee-server"),
			attribute.String("environment", "dev"),
			attribute.Int64("id", 1),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func main() {
	err := initZipKinProvider()
	if err != nil {
		log.Fatal(err)
	}
	tracingMiddleware := opentelemetry.Builder{}.Build()
	prometheusMiddleware := prometheus.Builder{
		Name:        "bee-server",
		Subsystem:   "",
		ConstLabels: nil,
		Help:        "",
	}.Build()

	s := bee.NewServer()
	recoverMiddleware := recovery.Builder{
		StatusCode: 500,
		ErrMsg:     "Internal Server Error\n",
		LogFunc: func(ctx *bee.Context) {
			log.Println(ctx.Request.URL.Path, "??????panic")
		},
	}.Build()
	s.Use(recoverMiddleware)
	s.Use(tracingMiddleware)
	s.Use(prometheusMiddleware)
	s.Use(RepeatBody())
	s.Use(accesslog.Builder{}.LogFunc(func(content string) {
		log.Println("access log: " + content)
	}).Build())
	s.Get("/user/profile", func(c *bee.Context) {
		_ = c.WriteString("match /userprofile\n")
	})
	s.Get("/order/detail", func(c *bee.Context) {
		_ = c.WriteString("match /order/detail\n")
	})
	s.Get("/user/*", func(ctx *bee.Context) {
		_ = ctx.WriteString("match /user/*\n")
	})
	s.Get("/order/detail/:id", func(c *bee.Context) {
		_ = c.WriteString(fmt.Sprintf("math /order/detail/%s\n", c.PathParams["id"]))
	})

	s.Get("/reg/:name([a-z]+)", func(c *bee.Context) {
		_ = c.WriteString(fmt.Sprintf("match /reg: %s\n", c.PathParams["name"]))
	})
	s.Get("/reg1/:name([a-z]+)/detail", func(c *bee.Context) {
		_ = c.WriteString(fmt.Sprintf("match /reg/%s/detail\n", c.PathParams["name"]))
	})

	s.Get("/panic", func(ctx *bee.Context) {
		panic("???????????? panic")
	})

	fileUploader := &bee.FileUploader{
		FileField: "file",
		DestPathFunc: func(fh *multipart.FileHeader) string {
			uploadPath, err := filepath.Abs("./advanced/web/testdata/upload")
			if err != nil {
				return ""
			}
			return path.Join(uploadPath, fh.Filename)
		},
	}
	s.Post("/upload", fileUploader.Handle())

	downloader := &bee.FileDownloader{
		Dir: "./testdata/download",
	}
	s.Post("/download", downloader.Handle())

	s.Get("/md", func(ctx *bee.Context) {
		_ = ctx.WriteString("??????md??????\n")
	}, func(next bee.HandleFunc) bee.HandleFunc {
		return func(ctx *bee.Context) {
			fmt.Println("??????md???Middleware")
			next(ctx)
		}
	})

	g := s.Group("/v1/product")
	g.Post("/list", func(ctx *bee.Context) {
		_ = ctx.WriteString("match /v1/product/list\n")
	})
	log.Println("started http server at :8080")
	err = s.Start(":8080")
	if err != nil {
		return
	}
}
