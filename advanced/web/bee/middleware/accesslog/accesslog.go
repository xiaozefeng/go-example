package accesslog

import (
	"encoding/json"
	"github.com/xiaozefeng/go-example/advanced/web/bee"
	"io"
)

type Builder struct {
	logFunc func(content string)
}

func (b Builder) LogFunc(logFunc func(content string)) Builder {
	b.logFunc = logFunc
	return b
}

func (b Builder) Build() bee.Middleware {
	return func(next bee.HandleFunc) bee.HandleFunc {
		return func(ctx *bee.Context) {
			body := io.NopCloser(ctx.Request.Body)
			all, err := io.ReadAll(body)
			if err == nil {
				log := accessLog{
					Method:     ctx.Request.Method,
					Body:       string(all),
					HTTPMethod: ctx.Request.Method,
					Path:       ctx.Request.URL.Path,
				}
				bs, err := json.Marshal(&log)
				if err == nil {
					b.logFunc(string(bs))
				}
			}
			next(ctx)
		}
	}
}

type accessLog struct {
	Method     string
	Body       string
	Path       string
	HTTPMethod string `json:"httpMethod"`
}
