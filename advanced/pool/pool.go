package pool

import "net"

type Pool interface {
	Get() (net.Conn, error)
	Close()
	Size() int
}
