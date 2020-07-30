package engine

import (
	"net"
	"sync"

	"github.com/zachdeibert/protomux/config"
)

// Listener handles listening for incomming connections
type Listener struct {
	Socket    net.Listener
	Service   *Service
	Engine    *Engine
	WaitGroup sync.WaitGroup
	Cleanup   bool
}

// CreateListener creates a new Listener
func CreateListener(address config.Connection, service *Service, engine *Engine) (*Listener, error) {
	l := &Listener{
		Service: service,
		Engine:  engine,
		Cleanup: false,
	}
	var ip net.IP
	if len(address.Host) > 0 {
		ips, err := net.LookupIP(address.Host)
		if err != nil {
			return nil, ErrorHostLookup(address.Host, err)
		}
		if len(ips) == 0 {
			return nil, ErrorNoHostRecords(address.Host)
		}
		ip = ips[0]
	} else {
		ip = address.IP
	}
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   ip,
		Port: address.Port,
	})
	if err != nil {
		return nil, ErrorListenStart(address, err)
	}
	l.Socket = listener
	return l, nil
}

// Start the listener
func (l *Listener) Start() {
	l.WaitGroup.Add(1)
	go func() {
		defer l.WaitGroup.Done()
		for !l.Cleanup {
			conn, err := l.Socket.Accept()
			if err != nil || l.Cleanup {
				if err != nil && !l.Cleanup {
					l.Engine.NormalError(err)
				}
				if l.Cleanup {
					return
				}
			} else {
				l.Service.AddRemote(conn)
			}
		}
	}()
}

// Stop the listener (and free resources)
func (l *Listener) Stop() {
	l.Cleanup = true
	if err := l.Socket.Close(); err != nil {
		l.Engine.NonCriticalError(err)
	}
	l.WaitGroup.Wait()
}
