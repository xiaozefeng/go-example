package pool

import (
	"github.com/fatih/pool"
	"net"
	"testing"
)

func TestPool(t *testing.T) {
	factory := func() (net.Conn, error) { return net.Dial("tcp", "127.0.0.1:9090") }

	p, err := pool.NewChannelPool(5, 40, factory)
	if err != nil {
		t.Error(err)
	}
	conn, err := p.Get()
	if err != nil {
		t.Error(err)
	}
	conn.Write([]byte("Hello\n"))
	conn.Close()

}
