package listener

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/framework"
)

// Listener handles listening for incomming connections and passing them to the correct Protocols
type Listener struct {
	Sockets     []net.Listener
	Protocols   []framework.ProtocolInstance
	Connections []*Connection
	State       State
	Mutex       sync.Mutex
	WaitGroup   sync.WaitGroup
}

// CreateListener creates a new Listener
func CreateListener(addresses []config.Connection, protocols []framework.ProtocolInstance) (*Listener, error) {
	l := &Listener{
		Sockets:     make([]net.Listener, len(addresses)),
		Protocols:   protocols,
		Connections: []*Connection{},
		State:       LoadedState,
	}
	for i, addr := range addresses {
		var ip net.IP
		if len(addr.Host) > 0 {
			ips, err := net.LookupIP(addr.Host)
			if err != nil {
				return nil, ErrorHostLookup(addr.Host, err)
			}
			if len(ips) == 0 {
				return nil, ErrorNoHostRecords(addr.Host)
			}
			ip = ips[0]
		} else {
			ip = addr.IP
		}
		listener, err := net.ListenTCP("tcp", &net.TCPAddr{
			IP:   ip,
			Port: addr.Port,
		})
		if err != nil {
			return nil, ErrorListenStart(addr, err)
		}
		l.Sockets[i] = listener
	}
	return l, nil
}

// Start the listener
func (l *Listener) Start() {
	if l.State != LoadedState {
		panic("Invalid state")
	}
	l.State = RunningState
	for _, lSock := range l.Sockets {
		sock := lSock
		l.WaitGroup.Add(1)
		go func() {
			defer l.WaitGroup.Done()
			l.Mutex.Lock()
			for l.State == RunningState {
				l.Mutex.Unlock()
				client, err := sock.Accept()
				l.Mutex.Lock()
				if state := l.State; err != nil || state == StoppingState {
					l.Mutex.Unlock()
					if state != StoppingState {
						fmt.Fprintf(os.Stderr, "Error encountered while running listener on %s:\n", sock.Addr())
						fmt.Fprintln(os.Stderr, err)
						sock.Close()
					}
					if client != nil {
						client.Close()
					}
					return
				}
				l.Connections = append(l.Connections, WrapConnection(client, l.Protocols, &l.WaitGroup))
			}
		}()
	}
}

// Stop the listener (and free resources)
func (l *Listener) Stop() {
	if l.State != RunningState {
		panic("Invalid state")
	}
	l.Mutex.Lock()
	l.State = StoppingState
	for _, sock := range l.Sockets {
		sock.Close()
	}
	for _, sock := range l.Connections {
		sock.Stop()
	}
	l.Mutex.Unlock()
	l.WaitGroup.Wait()
}
