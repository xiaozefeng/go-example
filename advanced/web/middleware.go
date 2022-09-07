package web

import (
	"encoding/json"
	"io"
)

type Middleware func(next HandleFunc) HandleFunc

type AccessLogBuilder struct {
	logFunc func(content string)
}

func (a *AccessLogBuilder) LogFunc(logFunc func(content string)) *AccessLogBuilder {
	a.logFunc = logFunc
	return a
}

func (a *AccessLogBuilder) Build() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			body := io.NopCloser(ctx.Request.Body)
			all, err := io.ReadAll(body)
			if err == nil {
				log := accessLog{
					Method: ctx.Request.Method,
					Body:   string(all),
				}
				bs, err := json.Marshal(&log)
				if err == nil {
					a.logFunc(string(bs))
				}
			}

			next(ctx)
		}
	}
}

type accessLog struct {
	Method string
	Body   string
}

func RepeatBody() Middleware {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			ctx.Request.Body = io.NopCloser(ctx.Request.Body)
			next(ctx)
		}
	}
}
