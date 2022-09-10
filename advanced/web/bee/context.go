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
