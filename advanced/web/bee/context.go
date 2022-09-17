package bee

import "net/http"

type Context struct {
	Request *http.Request
	Resp    http.ResponseWriter

	// 缓存的响应部分
	RespStatusCode int
	RespData       []byte

	PathParams map[string]string
	// 命中的路由
	MatchedRoute string
}

func (c *Context) WriteString(content string) error {
	//_, err := c.Resp.Write([]byte(content))
	c.RespStatusCode = http.StatusOK
	c.RespData = []byte(content)
	return nil
}

func (c *Context) BadRequestErr(msg string) {
	c.RespStatusCode = http.StatusBadRequest
	c.RespData = []byte(msg)
}

func (c *Context) NotFoundErr(msg string) {
	c.RespStatusCode = http.StatusNotFound
	c.RespData = []byte(msg)
}

func (c *Context) InternalServerErr(msg string) {
	c.RespStatusCode = http.StatusInternalServerError
	c.RespData = []byte(msg)
}

func (c *Context) QueryValue(key string) (string, error) {
	params := c.Request.URL.Query()
	vals, ok := params[key]
	if !ok || len(vals) == 0 {
		return "", nil
	}

	return vals[0], nil
}
