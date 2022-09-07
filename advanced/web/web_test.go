package web

import (
	"fmt"
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	s := NewServer()
	s.Get("/user/profile", func(c *Context) {
		_ = c.WriteString("match /userprofile\n")
	})
	s.Get("/order/detail", func(c *Context) {
		_ = c.WriteString("match /order/detail\n")
	})
	s.Get("/user/*", func(ctx *Context) {
		_ = ctx.WriteString("match /user/*\n")
	})
	s.Get("/order/detail/:id", func(c *Context) {
		_ = c.WriteString(fmt.Sprintf("math /order/detail/%s\n", c.PathParams["id"]))
	})

	g := s.Group("/v1/product")
	g.AddRoute(http.MethodGet, "/list", func(ctx *Context) {
		_ = ctx.WriteString("match /v1/product/list\n")
	})
	err := s.Start(":8080")
	if err != nil {
		return
	}
}
