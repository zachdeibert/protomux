package engine

import (
	"net"
	"sync"
	"time"

	"github.com/zachdeibert/protomux/framework"
)

// Connection represents a client that has connected to this Listener
type Connection struct {
	Remote        *RemoteConnection
	Engine        *Engine
	LocalAddress  net.Addr
	RemoteAddress net.Addr
	Closed        bool
	WaitGroup     sync.WaitGroup
	Mutex         sync.Mutex
	ReadBuffer    []byte
	WriteBuffer   []byte
	Priority      int
}

// CreateConnection creates a new Connection
func CreateConnection(remote *RemoteConnection, proto framework.ProtocolInstance, engine *Engine) *Connection {
	c := &Connection{
		Remote:        remote,
		Engine:        engine,
		LocalAddress:  remote.Socket.LocalAddr(),
		RemoteAddress: remote.Socket.RemoteAddr(),
		Closed:        false,
		ReadBuffer:    []byte{},
		WriteBuffer:   []byte{},
		Priority:      -1,
	}
	c.WaitGroup.Add(1)
	go func() {
		defer c.WaitGroup.Done()
		if err := proto.Handle(c); err != nil && err != ErrorClosed {
			c.Engine.NonCriticalError(err)
		}
		c.Close()
	}()
	return c
}

func (c *Connection) Read(b []byte) (int, error) {
	c.Remote.Mutex.Lock()
	for len(c.Remote.Connections) > 1 && len(c.ReadBuffer) == 0 && !c.Closed {
		c.Remote.Cond.Wait()
	}
	if c.Closed {
		c.Remote.Mutex.Unlock()
		return 0, ErrorClosed
	}
	if len(c.ReadBuffer) > 0 {
		l := len(b)
		if l > len(c.ReadBuffer) {
			l = len(c.ReadBuffer)
		}
		copy(b[:l], c.ReadBuffer[:l])
		c.ReadBuffer = c.ReadBuffer[l:]
		if len(c.ReadBuffer) == 0 {
			c.Remote.Cond.Broadcast()
		}
		c.Remote.Mutex.Unlock()
		return l, nil
	}
	if len(c.Remote.Connections) == 1 {
		c.Remote.Mutex.Unlock()
		return c.Remote.Socket.Read(b)
	}
	panic("Control should never reach here")
}

func (c *Connection) Write(b []byte) (int, error) {
	c.Remote.Mutex.Lock()
	c.WriteBuffer = b
	c.Remote.Cond.Broadcast()
	for len(c.Remote.Connections) > 1 && len(c.WriteBuffer) > 0 && !c.Closed {
		c.Remote.Cond.Wait()
	}
	if c.Closed {
		c.Remote.Mutex.Unlock()
		return 0, ErrorClosed
	}
	if len(c.Remote.Connections) == 1 {
		c.Remote.Mutex.Unlock()
		return c.Remote.Socket.Write(b)
	}
	if len(c.WriteBuffer) == 0 {
		c.Remote.Mutex.Unlock()
		return len(b), nil
	}
	panic("Control should never reach here")
}

// RequireExclusive notes that this Connection should be the only Connection on the RemoteConnection at this point
func (c *Connection) RequireExclusive(priority int) error {
	c.Remote.Mutex.Lock()
	c.Priority = priority
	for {
		var winner *Connection = nil
		maxPriority := -1
		for _, c := range c.Remote.Connections {
			if c.Priority < 0 {
				break
			}
			if c.Priority > maxPriority {
				maxPriority = c.Priority
				winner = c
			}
		}
		if winner == nil {
			c.Remote.Cond.Wait()
		} else {
			c.Remote.Mutex.Unlock()
			if winner == c {
				return nil
			}
			c.Close()
			return ErrorClosed
		}
	}
}

// Close the Connection
func (c *Connection) Close() error {
	c.Mutex.Lock()
	if !c.Closed {
		c.Closed = true
		c.Mutex.Unlock()
		c.Remote.ReleaseConnection(c, &c.WaitGroup)
	} else {
		c.Mutex.Unlock()
	}
	return nil
}

// LocalAddr returns the local network address
func (c *Connection) LocalAddr() net.Addr {
	return c.LocalAddress
}

// RemoteAddr returns the remote network address
func (c *Connection) RemoteAddr() net.Addr {
	return c.RemoteAddress
}

// SetDeadline sets the read and write deadlines
func (c *Connection) SetDeadline(t time.Time) error {
	return ErrorDeadlinesNotSupported
}

// SetReadDeadline sets the read deadline
func (c *Connection) SetReadDeadline(t time.Time) error {
	return ErrorDeadlinesNotSupported
}

// SetWriteDeadline sets the write deadline
func (c *Connection) SetWriteDeadline(t time.Time) error {
	return ErrorDeadlinesNotSupported
}
