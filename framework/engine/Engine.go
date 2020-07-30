package engine

import (
	"fmt"

	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/framework"
)

// Engine runs the main part of the program
type Engine struct {
	Services []*Service
}

// CreateEngine creates a new Engine
func CreateEngine(cfg config.Config) (*Engine, error) {
	eng := &Engine{
		Services: []*Service{},
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
		s, err := CreateService(srv.ListenAddresses, protos, eng)
		if err != nil {
			return nil, err
		}
		eng.Services = append(eng.Services, s)
	}
	return eng, nil
}

// Start starts the Engine
func (e *Engine) Start() {
	for _, s := range e.Services {
		s.Start()
	}
}

// Stop stops the Engine
func (e *Engine) Stop() {
	for _, s := range e.Services {
		s.Stop()
	}
}

// NonCriticalError handles a non-critical error
func (e *Engine) NonCriticalError(err error) {
	fmt.Println(err)
}

// NormalError handles an error that isn't critical, but also is more severe than other non-critical errors
func (e *Engine) NormalError(err error) {
	fmt.Println(err)
}
