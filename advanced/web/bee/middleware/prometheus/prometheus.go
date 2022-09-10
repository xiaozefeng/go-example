package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/xiaozefeng/go-example/advanced/web/bee"
	"strconv"
	"time"
)

type Builder struct {
	Name        string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func (b Builder) Build() bee.Middleware {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:        b.Name,
		Help:        b.Help,
		Subsystem:   b.Subsystem,
		ConstLabels: b.ConstLabels,
	}, []string{"pattern", "method", "status"})
	return func(next bee.HandleFunc) bee.HandleFunc {
		return func(ctx *bee.Context) {
			startTime := time.Now()
			defer func() {
				endTime := time.Now()
				go func() {
					statusCode := ctx.RespStatusCode
					route := "unknown"
					if len(ctx.MatchedRoute) > 0 {
						route = ctx.MatchedRoute
					}
					ms := endTime.Sub(startTime).Milliseconds()
					summaryVec.WithLabelValues(route, ctx.Request.Method, strconv.Itoa(statusCode)).Observe(float64(ms))
				}()
			}()
			next(ctx)
		}
	}
}
