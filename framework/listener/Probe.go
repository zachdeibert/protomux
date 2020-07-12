package listener

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/zachdeibert/protomux/framework"
)

// Probe represents a single probing attempt
type Probe struct {
	Local     net.Addr
	Remote    net.Addr
	WaitGroup *sync.WaitGroup
	ReadData  []byte
	WriteData []byte
	Closed    bool
	Result    framework.ProtocolState
}

// AsyncProbe starts a probe asynchronously
func AsyncProbe(data []byte, localAddr, remoteAddr net.Addr, waitGroup *sync.WaitGroup, proto framework.ProtocolInstance) *Probe {
	probe := &Probe{
		Local:     localAddr,
		Remote:    remoteAddr,
		WaitGroup: waitGroup,
		ReadData:  data,
		WriteData: []byte{},
		Closed:    false,
	}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		probe.Result = proto.Handle(probe, framework.ProbingState, waitGroup)
	}()
	return probe
}

func (p *Probe) Read(b []byte) (int, error) {
	var read int
	if len(p.ReadData) > len(b) {
		read = len(b)
	} else {
		read = len(p.ReadData)
	}
	if read == 0 {
		return 0, io.EOF
	}
	copy(b[0:read], p.ReadData[0:read])
	p.ReadData = p.ReadData[read:]
	return read, nil
}

func (p *Probe) Write(b []byte) (int, error) {
	p.WriteData = append(p.WriteData, b...)
	return len(b), nil
}

// Close closes the connection
func (p *Probe) Close() error {
	p.Closed = true
	return nil
}

// LocalAddr returns the local network address
func (p *Probe) LocalAddr() net.Addr {
	return p.Local
}

// RemoteAddr returns the remote network address
func (p *Probe) RemoteAddr() net.Addr {
	return p.Remote
}

// SetDeadline sets the read and write deadlines associated with the connection
func (p *Probe) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for Read calls
func (p *Probe) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for Write calls
func (p *Probe) SetWriteDeadline(t time.Time) error {
	return nil
}
