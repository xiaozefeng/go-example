package pool

import (
	"fmt"
	"log"
	"net"
	"testing"
)

const (
	addr = "localhost:9090"
)

func TestPool(t *testing.T) {
	pool, err := NewPool(func() (net.Conn, error) {
		return net.Dial("tcp", addr)
	},
		WithInitialPoolSize(5),
		WithMaxPoolSize(30),
	)
	if err != nil {
		t.Error(err)
	}
	size := 50

	conns := make([]net.Conn, 0, size)
	for i := 0; i < size; i++ {
		conn, err := pool.Get()
		if err != nil {
			t.Error(err)
		}
		conns = append(conns, conn)
	}
	for _, conn := range conns {
		conn.Close()
	}

	fmt.Println("pool size", pool.Size()) // 30
	pool.Close()

}

func TestMockTCPServer(t *testing.T) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("listen tcp server on", addr)
	for {
		conn, err := listen.Accept()
		log.Println("conn has accept")
		if err != nil {
			log.Println("accept conn error:", err)
		}
		go func(conn net.Conn) {
			defer func(conn net.Conn) {
				_ = conn.Close()
			}(conn)
			buf := make([]byte, 0, 20)
			_, err = conn.Read(buf)
			if err != nil {
				log.Println("handle conn error:", err)
			}
			_, _ = conn.Write([]byte("hi,back\n"))
		}(conn)
	}
}
