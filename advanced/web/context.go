package web

import "net/http"

type Context struct {
	Request    *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
}

func (c *Context) WriteString(content string) error {
	_, err := c.Resp.Write([]byte(content))
	return err
}
