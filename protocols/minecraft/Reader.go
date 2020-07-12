package minecraft

import (
	"bytes"
	"io"
)

// Reader reads primitives from the Minecraft protocol
type Reader struct {
	Base io.Reader
}

// CreateReader creates a new Reader
func CreateReader(base io.Reader) *Reader {
	return &Reader{
		Base: base,
	}
}

// ReadVarInt reads a VarInt
func (r *Reader) ReadVarInt() (int, error) {
	buf := []byte{0x80}
	res := 0
	bitPos := 0
	for (buf[0] & 0x80) != 0 {
		if _, err := r.Base.Read(buf); err != nil {
			return 0, err
		}
		res |= (int(buf[0]&0x7F) << bitPos)
		bitPos += 7
	}
	return res, nil
}

// ReadString reads a string
func (r *Reader) ReadString() (string, error) {
	l, err := r.ReadVarInt()
	if err != nil {
		return "", err
	}
	buf := make([]byte, l)
	it := buf
	for len(it) > 0 {
		n, err := r.Base.Read(it)
		if err != nil {
			return "", err
		}
		it = it[n:]
	}
	return string(buf), nil
}

// ReadUByte reads an unsigned byte
func (r *Reader) ReadUByte() (uint8, error) {
	buf := []byte{0}
	if _, err := r.Base.Read(buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

// ReadUShort reads an unsigned short
func (r *Reader) ReadUShort() (uint16, error) {
	ms, err := r.ReadUByte()
	if err != nil {
		return 0, err
	}
	ls, err := r.ReadUByte()
	if err != nil {
		return 0, err
	}
	return (uint16(ms) << 8) | uint16(ls), nil
}

// ReadUInt reads an unsigned int
func (r *Reader) ReadUInt() (uint32, error) {
	ms, err := r.ReadUShort()
	if err != nil {
		return 0, err
	}
	ls, err := r.ReadUShort()
	if err != nil {
		return 0, err
	}
	return (uint32(ms) << 16) | uint32(ls), nil
}

// ReadULong reads an unsigned long
func (r *Reader) ReadULong() (uint64, error) {
	ms, err := r.ReadUInt()
	if err != nil {
		return 0, err
	}
	ls, err := r.ReadUInt()
	if err != nil {
		return 0, err
	}
	return (uint64(ms) << 32) | uint64(ls), nil
}

// ReadUncompressedPacket reads a packet
func (r *Reader) ReadUncompressedPacket() (*Reader, int, error) {
	l, err := r.ReadVarInt()
	if err != nil {
		return nil, 0, err
	}
	id, err := r.ReadVarInt()
	if err != nil {
		return nil, 0, err
	}
	buf := make([]byte, l-((id>>7)+1))
	it := buf
	for len(it) > 0 {
		n, err := r.Base.Read(it)
		if err != nil {
			return nil, 0, err
		}
		it = it[n:]
	}
	return CreateReader(bytes.NewReader(buf)), id, nil
}
