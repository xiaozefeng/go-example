package recovery

import (
	"github.com/xiaozefeng/go-example/advanced/web/bee"
	"log"
)

type Builder struct {
	StatusCode int
	ErrMsg     string
	LogFunc    func(ctx *bee.Context)
}

func (b Builder) Build() bee.Middleware {
	return func(next bee.HandleFunc) bee.HandleFunc {
		return func(ctx *bee.Context) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("进入recover 逻辑")
					ctx.RespStatusCode = b.StatusCode
					ctx.RespData = []byte(b.ErrMsg)
					b.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}
