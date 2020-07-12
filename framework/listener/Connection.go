package listener

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/zachdeibert/protomux/framework"
)

const bufferSize = 4096

// Connection represents a client that has connected to this Listener
type Connection struct {
	Socket      net.Conn
	WaitGroup   *sync.WaitGroup
	State       framework.HandlerState
	ReadPrefix  []byte
	WritePrefix []byte
	Local       net.Addr
	Remote      net.Addr
}

// WrapConnection creeates a new Connection and starts the protocol multiplexing process
func WrapConnection(sock net.Conn, protos []framework.ProtocolInstance, waitGroup *sync.WaitGroup) *Connection {
	conn := &Connection{
		Socket:      sock,
		WaitGroup:   waitGroup,
		State:       framework.ProbingState,
		ReadPrefix:  []byte{},
		WritePrefix: []byte{},
		Local:       sock.LocalAddr(),
		Remote:      sock.RemoteAddr(),
	}
	waitGroup.Add(1)
	go func() {
		fmt.Println("Starting connection.")
		defer waitGroup.Done()
		var subGroup sync.WaitGroup
		defer subGroup.Wait()
		defer sock.Close()
		buf := make([]byte, bufferSize)
		probes := make([]*Probe, len(protos))
		possible := protos
		for conn.State != framework.StoppingState {
			n, err := sock.Read(buf)
			if err != nil {
				return
			}
			conn.ReadPrefix = append(conn.ReadPrefix, buf[:n]...)
			probes = probes[:len(possible)]
			for i, proto := range possible {
				probes[i] = AsyncProbe(conn.ReadPrefix, conn.Local, conn.Remote, &subGroup, proto)
			}
			subGroup.Wait()
			matched := []framework.ProtocolInstance{}
			n = 0
			for i, probe := range probes {
				valid := probe.Result != framework.ProtocolNotMatched && len(probe.WriteData) >= len(conn.WritePrefix)
				if valid {
					for i, b := range conn.WritePrefix {
						if probe.WriteData[i] != b {
							valid = false
							break
						}
					}
				}
				if valid {
					probes[n] = probe
					possible[n] = possible[i]
					if probe.Result == framework.ProtocolMatched {
						matched = append(matched, possible[i])
					}
					n++
				}
			}
			probes = probes[:n]
			possible = possible[:n]
			if len(matched) == 1 {
				possible = matched
			}
			switch len(possible) {
			case 0:
				fmt.Println("No valid protocols...")
				return
			case 1:
				fmt.Println("Protocol matched.")
				if probes[0].Closed {
					sock.Write(probes[0].WriteData)
				} else {
					possible[0].Handle(conn, framework.HandoffState, waitGroup)
				}
				return
			}
			fmt.Println("Multiple protocols matched.")
			commonWrite := probes[0].WriteData[len(conn.WritePrefix):]
			for _, probe := range probes {
				data := probe.WriteData[len(conn.WritePrefix):]
				if len(data) < len(commonWrite) {
					commonWrite = commonWrite[:len(data)]
				}
				for i, b := range commonWrite {
					if data[i] != b {
						commonWrite = commonWrite[:i]
						break
					}
				}
			}
			if len(commonWrite) > 0 {
				sock.Write(commonWrite)
				conn.WritePrefix = append(conn.WritePrefix, commonWrite...)
			}
		}
	}()
	return conn
}

// Stop is called when a listener is stopping
func (c *Connection) Stop() {
	c.State = framework.StoppingState
	switch sock := c.Socket.(type) {
	case *net.TCPConn:
		sock.CloseRead()
	}
}

func (c *Connection) Read(b []byte) (int, error) {
	if len(c.ReadPrefix) > 0 {
		var read int
		if len(c.ReadPrefix) > len(b) {
			read = len(b)
		} else {
			read = len(c.ReadPrefix)
		}
		copy(b[0:read], c.ReadPrefix[0:read])
		c.ReadPrefix = c.ReadPrefix[read:]
		return read, nil
	}
	if c.State == framework.StoppingState {
		return 0, io.EOF
	}
	return c.Socket.Read(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	write := 0
	if len(c.WritePrefix) > 0 {
		if len(c.WritePrefix) > len(b) {
			write = len(b)
		} else {
			write = len(c.WritePrefix)
		}
		for i := 0; i < write; i++ {
			if c.WritePrefix[i] != b[i] {
				c.State = framework.StoppingState
				c.Socket.Close()
				return c.Socket.Write(b)
			}
		}
		c.WritePrefix = c.WritePrefix[:write]
		b = b[write:]
		if len(b) == 0 {
			return write, nil
		}
	}
	w, err := c.Socket.Write(b)
	return w + write, err
}

// Close closes the connection
func (c *Connection) Close() error {
	return c.Socket.Close()
}

// LocalAddr returns the local network address
func (c *Connection) LocalAddr() net.Addr {
	return c.Local
}

// RemoteAddr returns the remote network address
func (c *Connection) RemoteAddr() net.Addr {
	return c.Remote
}

// SetDeadline sets the read and write deadlines associated with the connection
func (c *Connection) SetDeadline(t time.Time) error {
	return c.Socket.SetDeadline(t)
}

// SetReadDeadline sets the deadline for Read calls
func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.Socket.SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for Write calls
func (c *Connection) SetWriteDeadline(t time.Time) error {
	return c.Socket.SetWriteDeadline(t)
}
