package web

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	s := NewServer()
	b := &AccessLogBuilder{}
	s.Use(RepeatBody(), b.LogFunc(func(content string) {
		fmt.Println(content)
	}).Build())
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

	s.Get("/reg/:name([a-z]+)", func(c *Context) {
		_ = c.WriteString(fmt.Sprintf("match /reg: %s\n", c.PathParams["name"]))
	})
	s.Get("/reg1/:name([a-z]+)/detail", func(c *Context) {
		_ = c.WriteString(fmt.Sprintf("match /reg/%s/detail\n", c.PathParams["name"]))
	})

	g := s.Group("/v1/product")
	g.Post("/list", func(ctx *Context) {
		_ = ctx.WriteString("match /v1/product/list\n")
	})
	err := s.Start(":8080")
	if err != nil {
		return
	}
}
