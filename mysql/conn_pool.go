package mysql

import (
	"errors"
	"fmt"
	"sync"
)

var ErrClosed = errors.New("pool is closed")

type ConnPool struct {
	addr     string
	user     string
	password string
	db       string
	mu       sync.Mutex
	conns    chan BackendConn
}

func NewConnPool(initialCap, maxCap int, addr, user, password, db string) (*ConnPool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &ConnPool{
		addr:     addr,
		user:     user,
		password: password,
		db:       db,
		conns:    make(chan BackendConn, maxCap),
	}

	// create initial connections, if something goes wrong,
	// just close the pool error out.
	for i := 0; i < initialCap; i++ {
		conn, err := factory(addr, user, password, db)
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- conn
	}
	return c, nil
}

func (c *ConnPool) Close() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func factory(addr, user, password, db string) (*Conn, error) {
	conn := &Conn{}
	err := conn.Connect(addr, user, password, db)
	return conn, err

}

func (c *ConnPool) getConns() chan BackendConn {
	c.mu.Lock()
	conns := c.conns
	c.mu.Unlock()
	return conns
}

func (c *ConnPool) Get() (BackendConn, error) {
	conns := c.getConns()
	if conns == nil {
		return nil, ErrClosed
	}

	// wrap our connections with out custom net.Conn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}

		return c.wrapConn(conn), nil
	default:
		conn, err := factory(c.addr, c.user, c.password, c.db)
		if err != nil {
			return nil, err
		}

		return c.wrapConn(conn), nil
	}
}

func (c *ConnPool) put(conn BackendConn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case c.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

type PoolConn struct {
	BackendConn
	c        *ConnPool
	unusable bool
}

func (p *PoolConn) Close() error {
	if p.unusable {
		if p.BackendConn != nil {
			return p.BackendConn.Close()
		}
		return nil
	}
	return p.c.put(p.BackendConn)
}

func (p *PoolConn) MarkUnusable() {
	p.unusable = true
}

func (c *ConnPool) wrapConn(conn BackendConn) BackendConn {
	p := &PoolConn{c: c}
	p.BackendConn = conn
	return p
}
