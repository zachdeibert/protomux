package engine

import (
	"net"
	"sync"

	"github.com/zachdeibert/protomux/framework"
)

const bufferSize = 4096

// RemoteConnection represents one remote connection that can be connected to multiple Connection objects if multiple protocols are being evaluated
type RemoteConnection struct {
	Socket      net.Conn
	Connections []*Connection
	Service     *Service
	Engine      *Engine
	Mutex       sync.Mutex
	Cond        *sync.Cond
	WaitGroup   sync.WaitGroup
	Closed      bool
}

// CreateRemoteConnection creates a new RemoteConnection
func CreateRemoteConnection(conn net.Conn, protocols []framework.ProtocolInstance, service *Service, engine *Engine) *RemoteConnection {
	rc := &RemoteConnection{
		Socket:      conn,
		Connections: make([]*Connection, len(protocols)),
		Service:     service,
		Engine:      engine,
		Closed:      false,
	}
	rc.Cond = sync.NewCond(&rc.Mutex)
	for i, proto := range protocols {
		rc.Connections[i] = CreateConnection(rc, proto, engine)
	}
	rc.WaitGroup.Add(1)
	go func() {
		defer rc.WaitGroup.Done()
		rc.Mutex.Lock()
		defer rc.Mutex.Unlock()
		buffer := make([]byte, bufferSize)
		for !rc.Closed && len(rc.Connections) > 1 {
			ready := true
			for _, c := range rc.Connections {
				if len(c.ReadBuffer) > 0 {
					ready = false
					break
				}
			}
			if ready {
				n, err := rc.Socket.Read(buffer)
				if err != nil {
					engine.NormalError(err)
					go rc.Close()
					return
				}
				for _, c := range rc.Connections {
					c.ReadBuffer = buffer[:n]
				}
				rc.Cond.Broadcast()
			} else {
				rc.Cond.Wait()
			}
		}
	}()
	rc.WaitGroup.Add(1)
	go func() {
		defer rc.WaitGroup.Done()
		rc.Mutex.Lock()
		defer rc.Mutex.Unlock()
		for !rc.Closed && len(rc.Connections) > 1 {
			n := 0
			for i, c := range rc.Connections {
				if i == 0 {
					n = len(c.WriteBuffer)
				} else {
					if len(c.WriteBuffer) < n {
						n = len(c.WriteBuffer)
					}
					for i, b := range c.WriteBuffer[:n] {
						if rc.Connections[0].WriteBuffer[i] != b {
							n = i
							break
						}
					}
				}
				if n == 0 {
					break
				}
			}
			if n == 0 {
				rc.Cond.Wait()
			} else {
				n, err := rc.Socket.Write(rc.Connections[0].WriteBuffer[:n])
				if err != nil {
					engine.NormalError(err)
					go rc.Close()
					return
				}
				for _, c := range rc.Connections {
					c.WriteBuffer = c.WriteBuffer[n:]
				}
				rc.Cond.Broadcast()
			}
		}
	}()
	return rc
}

// ReleaseConnection marks a Connection as no longer being a valid part of this RemoteConnection
func (rc *RemoteConnection) ReleaseConnection(c *Connection, wg *sync.WaitGroup) {
	rc.Mutex.Lock()
	rc.WaitGroup.Add(1)
	go func() {
		defer rc.WaitGroup.Done()
		wg.Wait()
	}()
	n := 0
	for _, x := range rc.Connections {
		if x != c {
			rc.Connections[n] = x
			n++
		}
	}
	rc.Connections = rc.Connections[:n]
	rc.Cond.Broadcast()
	if len(rc.Connections) == 0 {
		rc.Mutex.Unlock()
		rc.Close()
	} else {
		rc.Mutex.Unlock()
	}
}

// Close the RemoteConnection
func (rc *RemoteConnection) Close() {
	if err := rc.Socket.Close(); err != nil {
		rc.Engine.NonCriticalError(err)
	}
	rc.Mutex.Lock()
	rc.Closed = true
	rc.Cond.Broadcast()
	for len(rc.Connections) > 0 {
		rc.Cond.Wait()
	}
	rc.Mutex.Unlock()
	rc.WaitGroup.Wait()
	rc.Service.ReleaseRemote(rc)
}
