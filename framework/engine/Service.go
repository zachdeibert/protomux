package engine

import (
	"net"
	"sync"

	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/framework"
)

// Service represents a set of listeners that all perform the same task
type Service struct {
	Protocols []framework.ProtocolInstance
	Listeners []*Listener
	Engine    *Engine
	Remotes   []*RemoteConnection
	Mutex     sync.Mutex
}

// CreateService creates a new Service
func CreateService(addresses []config.Connection, protocols []framework.ProtocolInstance, engine *Engine) (*Service, error) {
	srv := &Service{
		Protocols: protocols,
		Listeners: make([]*Listener, len(addresses)),
		Engine:    engine,
		Remotes:   []*RemoteConnection{},
	}
	for i, addr := range addresses {
		l, err := CreateListener(addr, srv, engine)
		if err != nil {
			return nil, err
		}
		srv.Listeners[i] = l
	}
	return srv, nil
}

// Start the Service
func (s *Service) Start() {
	for _, l := range s.Listeners {
		l.Start()
	}
}

// Stop the Service
func (s *Service) Stop() {
	for _, l := range s.Listeners {
		l.Stop()
	}
	var wg sync.WaitGroup
	s.Mutex.Lock()
	for _, remote := range s.Remotes {
		wg.Add(1)
		go func(remote *RemoteConnection) {
			defer wg.Done()
			remote.Close()
		}(remote)
	}
	s.Mutex.Unlock()
	wg.Wait()
}

// ReleaseRemote removes a remove connection from this Service
func (s *Service) ReleaseRemote(r *RemoteConnection) {
	s.Mutex.Lock()
	n := 0
	for _, x := range s.Remotes {
		if x != r {
			s.Remotes[n] = x
			n++
		}
	}
	s.Remotes = s.Remotes[:n]
	s.Mutex.Unlock()
}

// AddRemote adds a new remote connection to this Service
func (s *Service) AddRemote(conn net.Conn) {
	s.Mutex.Lock()
	s.Remotes = append(s.Remotes, CreateRemoteConnection(conn, s.Protocols, s, s.Engine))
	s.Mutex.Unlock()
}
