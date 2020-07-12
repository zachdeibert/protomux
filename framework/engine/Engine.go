package engine

import (
	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/framework"
	"github.com/zachdeibert/protomux/framework/listener"
)

// Engine runs the main part of the program
type Engine struct {
	Listeners []*listener.Listener
}

// CreateEngine creates a new Engine
func CreateEngine(cfg config.Config) (*Engine, error) {
	eng := &Engine{
		Listeners: []*listener.Listener{},
	}
	for _, srv := range cfg.Services {
		protos := []framework.ProtocolInstance{}
		for _, p := range srv.Protocols {
			impl, ok := framework.Protocols[p.Name]
			if !ok {
				return nil, ErrorUnknownProtocol(p.Name)
			}
			for _, remote := range p.Remotes {
				inst, err := impl.Configure(p.Parameters, remote.Name, remote.Parameters)
				if err != nil {
					return nil, err
				}
				protos = append(protos, inst)
			}
		}
		l, err := listener.CreateListener(srv.ListenAddresses, protos)
		if err != nil {
			return nil, err
		}
		eng.Listeners = append(eng.Listeners, l)
	}
	return eng, nil
}

// Start starts the Engine
func (e *Engine) Start() {
	for _, l := range e.Listeners {
		l.Start()
	}
}

// Stop stops the Engine
func (e *Engine) Stop() {
	for _, l := range e.Listeners {
		l.Stop()
	}
}
