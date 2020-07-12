package minecraft

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/framework"
)

// HandleNettyRewrite handles the protocol for clients that have the Netty rewrite
func (p ProtocolInstance) HandleNettyRewrite(conn net.Conn, stream *bufio.Reader, state framework.HandlerState, waitGroup *sync.WaitGroup) framework.ProtocolState {
	reader := CreateReader(stream)
	writer := CreateWriter(conn)
	pkt, id, err := reader.ReadUncompressedPacket()
	if err != nil {
		return e(err)
	}
	if id != 0 {
		return framework.ProtocolNotMatched
	}
	version, err := pkt.ReadVarInt()
	if err != nil {
		return e(err)
	}
	addr, err := pkt.ReadString()
	if err != nil {
		if err == io.EOF {
			return framework.ProtocolNeedsMoreData
		}
		return framework.ProtocolNotMatched
	}
	port, err := pkt.ReadUShort()
	if err != nil {
		return e(err)
	}
	c := config.Connection{
		Port: int(port),
	}
	if ip := net.ParseIP(addr); ip == nil {
		c.Host = addr
	} else {
		c.IP = ip
	}
	filterData := FilteringProps{
		Version:       []Version{Version(version)},
		ServerAddress: []config.Connection{c},
	}
	if !filterData.Check(p.Filter) {
		return framework.ProtocolNotMatched
	}
	if state != framework.HandoffState {
		if p.Filter.IsEmpty() {
			return framework.ProtocolNeedsMoreData
		}
		return framework.ProtocolMatched
	}
	nextState, err := pkt.ReadVarInt()
	if err != nil {
		return e(err)
	}
	switch nextState {
	case 1: // status
		if p.Action.MOTD == nil {
			// TODO proxy
			return framework.ProtocolNotMatched
		}
		_, id, err = reader.ReadUncompressedPacket()
		if err != nil {
			return e(err)
		}
		if id != 0 {
			return framework.ProtocolNotMatched
		}
		versionName := ""
		for i, v := range p.Filter.Version {
			if v == Version(version) {
				versionName = p.Filter.VersionName[i]
			}
		}
		if versionName == "" {
			var ok bool
			if versionName, ok = NettyVersionNames[Version(version)]; !ok {
				versionName = "unknown"
			}
		}
		status := fmt.Sprintf(`{"version":{"name":"%s","protocol":%d},"players":{"max":0,"online":0,"sample":[]},"description":{"text":"%s"}}`, versionName, version, *p.Action.MOTD)
		wpkt := writer.WriteUncompressedPacket(0)
		wpkt.WriteString(status)
		wpkt.Close()
		pkt, id, err = reader.ReadUncompressedPacket()
		if err != nil {
			return e(err)
		}
		if id != 1 {
			return framework.ProtocolNotMatched
		}
		payload, err := pkt.ReadULong()
		if err != nil {
			return e(err)
		}
		wpkt = writer.WriteUncompressedPacket(1)
		wpkt.WriteULong(payload)
		wpkt.Close()
		return framework.ProtocolMatched
	case 2: // login
		if p.Action.Remote != nil {
			// TODO proxy
			return framework.ProtocolNotMatched
		}
		pkt, id, err := reader.ReadUncompressedPacket()
		if err != nil {
			return e(err)
		}
		if id != 0 {
			return framework.ProtocolNotMatched
		}
		username, err := pkt.ReadString()
		if err != nil {
			return e(err)
		}
		wpkt := writer.WriteUncompressedPacket(0x02)
		wpkt.WriteString("00000000-0000-0000-0000-000000000000")
		wpkt.WriteString(username)
		wpkt.Close()
		reason := fmt.Sprintf(`{"text":"%s"}`, *p.Action.Kick)
		wpkt = writer.WriteUncompressedPacket(0x1B)
		wpkt.WriteString(reason)
		wpkt.Close()
		return framework.ProtocolMatched
	default:
		return framework.ProtocolNotMatched
	}
}
