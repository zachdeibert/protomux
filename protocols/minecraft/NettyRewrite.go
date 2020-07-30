package minecraft

import (
	"bufio"
	"fmt"
	"net"

	"github.com/zachdeibert/protomux/config"
	"github.com/zachdeibert/protomux/framework"
)

// HandleNettyRewrite handles the protocol for clients that have the Netty rewrite
func (p ProtocolInstance) HandleNettyRewrite(conn framework.Connection, stream *bufio.Reader) error {
	reader := CreateReader(stream)
	writer := CreateWriter(conn)
	pkt, id, err := reader.ReadUncompressedPacket()
	if err != nil {
		return err
	}
	if id != 0 {
		return ErrorProtocol("Expected handshake packet")
	}
	version, err := pkt.ReadVarInt()
	if err != nil {
		return err
	}
	addr, err := pkt.ReadString()
	if err != nil {
		return err
	}
	port, err := pkt.ReadUShort()
	if err != nil {
		return err
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
		return ErrorProtocol("Filter mismatch")
	}
	priority := 1
	if filterData.IsEmpty() {
		priority = 0
	}
	if err = conn.RequireExclusive(priority); err != nil {
		return err
	}
	nextState, err := pkt.ReadVarInt()
	if err != nil {
		return err
	}
	switch nextState {
	case 1: // status
		if p.Action.MOTD == nil {
			// TODO proxy
			return ErrorProtocol("Status proxying")
		}
		_, id, err = reader.ReadUncompressedPacket()
		if err != nil {
			return err
		}
		if id != 0 {
			return ErrorProtocol("Expected status packet")
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
			return err
		}
		if id != 1 {
			return ErrorProtocol("Expected ping packet")
		}
		payload, err := pkt.ReadULong()
		if err != nil {
			return err
		}
		wpkt = writer.WriteUncompressedPacket(1)
		wpkt.WriteULong(payload)
		wpkt.Close()
		return nil
	case 2: // login
		if p.Action.Remote != nil {
			// TODO proxy
			return ErrorProtocol("Login not supported")
		}
		pkt, id, err := reader.ReadUncompressedPacket()
		if err != nil {
			return err
		}
		if id != 0 {
			return ErrorProtocol("Expected login packet")
		}
		username, err := pkt.ReadString()
		if err != nil {
			return err
		}
		wpkt := writer.WriteUncompressedPacket(0x02)
		wpkt.WriteString("00000000-0000-0000-0000-000000000000")
		wpkt.WriteString(username)
		wpkt.Close()
		reason := fmt.Sprintf(`{"text":"%s"}`, *p.Action.Kick)
		wpkt = writer.WriteUncompressedPacket(0x1B)
		wpkt.WriteString(reason)
		wpkt.Close()
		return nil
	default:
		return ErrorProtocol("Unknown next state")
	}
}
