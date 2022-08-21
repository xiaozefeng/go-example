package pool

import (
	"net"
	"sync"
)

type wrapperConn struct {
	net.Conn

	mu      sync.RWMutex
	p       *channelPool
	useless bool
}

func (w *wrapperConn) Close() error {
	w.mu.RLock()
	useless := w.useless
	w.mu.RUnlock()
	if useless {
		// real close conn
		if w.Conn != nil {
			return w.Conn.Close()
		}
	} else {
		// put conn to the pool
		return w.p.put(w.Conn)
	}
	return nil
}

func (w *wrapperConn) MarkUseless() {
	w.mu.Lock()
	w.useless = true
	w.mu.Unlock()
}
