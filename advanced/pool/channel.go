package pool

import (
	"errors"
	"net"
	"sync"
)

var (
	ErrPoolClosed = errors.New("pool closed")
)

// Factory 生成连接的工厂方法
type Factory func() (net.Conn, error)

type options struct {
	initialPoolSize int
	maxPoolSize     int
}

type Option func(options *options)

// pool的channel实现
type channelPool struct {
	connChan chan net.Conn

	factory Factory
	mu      sync.RWMutex
}

func NewPool(factory Factory, opts ...Option) (Pool, error) {
	var o = options{
		initialPoolSize: 5,  //默认 5
		maxPoolSize:     20, //默认20
	}
	for _, opt := range opts {
		opt(&o)
	}

	p := &channelPool{
		connChan: make(chan net.Conn, o.maxPoolSize),
		factory:  factory,
	}

	for i := 0; i < o.initialPoolSize; i++ {
		conn, err := factory()
		if err != nil {
			p.Close()
			return nil, err
		}
		p.connChan <- conn
	}
	return p, nil
}

func (c *channelPool) Get() (net.Conn, error) {
	c.mu.RLock()
	connChan := c.connChan
	factory := c.factory
	c.mu.RUnlock()

	if connChan == nil {
		return nil, ErrPoolClosed
	}

	select {
	case conn := <-connChan:
		if conn == nil {
			return nil, ErrPoolClosed
		}
		return wrapConn(c, conn), nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		// 在close的时候加入到pool中
		return wrapConn(c, conn), nil
	}
}

func (c *channelPool) Close() {
	c.mu.Lock()
	connChan := c.connChan
	c.factory = nil
	c.connChan = nil
	c.mu.Unlock()

	close(connChan)
	for conn := range connChan {
		_ = conn.Close()
	}
}

func (c *channelPool) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.connChan)
}

func (c *channelPool) put(conn net.Conn) error {
	if conn == nil {
		return errors.New("conn is closed")
	}
	c.mu.RLock()
	// pool可能已经关闭了
	if c.connChan == nil {
		return errors.New("pool is closed")
	}
	c.mu.RUnlock()

	select {
	case c.connChan <- conn:
		return nil
	default:
		// pool is full
		return conn.Close()
	}
}

func WithInitialPoolSize(initialSize int) Option {
	if initialSize < 0 {
		panic("invalid initial pool size")
	}
	return func(options *options) {
		options.initialPoolSize = initialSize
	}
}

func WithMaxPoolSize(maxPoolSize int) Option {
	if maxPoolSize < 0 {
		panic("invalid max pool size")
	}
	return func(options *options) {
		options.maxPoolSize = maxPoolSize
	}
}

func wrapConn(p *channelPool, conn net.Conn) net.Conn {
	return &wrapperConn{p: p, Conn: conn}
}
