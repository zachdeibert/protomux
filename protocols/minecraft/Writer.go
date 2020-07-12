package minecraft

import (
	"bufio"
	"bytes"
	"io"
)

// Writer writes primitives to the Minecraft protocol
type Writer struct {
	Base   io.Writer
	Parent *Writer
}

// CreateWriter creates a new Writer
func CreateWriter(base io.Writer) *Writer {
	return &Writer{
		Base:   bufio.NewWriter(base),
		Parent: nil,
	}
}

// Flush the currently written data into a packet
func (w *Writer) Flush() error {
	if stream, ok := w.Base.(*bufio.Writer); ok {
		return stream.Flush()
	}
	return nil
}

// WriteVarInt writes a VarInt
func (w *Writer) WriteVarInt(val int) error {
	if val == 0 {
		w.Base.Write([]byte{0})
	}
	for val != 0 {
		b := byte(val & 0x7F)
		val >>= 7
		if val != 0 {
			b |= 0x80
		}
		if _, err := w.Base.Write([]byte{b}); err != nil {
			return err
		}
	}
	return nil
}

// WriteString writes a String
func (w *Writer) WriteString(val string) error {
	if err := w.WriteVarInt(len(val)); err != nil {
		return err
	}
	if _, err := w.Base.Write([]byte(val)); err != nil {
		return err
	}
	return nil
}

// WriteUByte writes an unsigned byte
func (w *Writer) WriteUByte(val uint8) error {
	_, err := w.Base.Write([]byte{val})
	if err != nil {
		return err
	}
	return nil
}

// WriteUShort writes an unsigned short
func (w *Writer) WriteUShort(val uint16) error {
	if err := w.WriteUByte(uint8((val & 0xFF00) >> 8)); err != nil {
		return err
	}
	if err := w.WriteUByte(uint8(val & 0x00FF)); err != nil {
		return err
	}
	return nil
}

// WriteUInt writes an unsigned int
func (w *Writer) WriteUInt(val uint32) error {
	if err := w.WriteUShort(uint16((val & 0xFFFF0000) >> 16)); err != nil {
		return err
	}
	if err := w.WriteUShort(uint16(val & 0x0000FFFF)); err != nil {
		return err
	}
	return nil
}

// WriteULong writes an unsigned long
func (w *Writer) WriteULong(val uint64) error {
	if err := w.WriteUInt(uint32((val & 0xFFFFFFFF00000000) >> 32)); err != nil {
		return err
	}
	if err := w.WriteUInt(uint32(val & 0x00000000FFFFFFFF)); err != nil {
		return err
	}
	return nil
}

// WriteUncompressedPacket writes a packet
func (w *Writer) WriteUncompressedPacket(id int) *Writer {
	n := &Writer{
		Base:   &bytes.Buffer{},
		Parent: w,
	}
	n.WriteVarInt(id)
	return n
}

// Close flushes the packet to the network
func (w *Writer) Close() error {
	if w.Parent != nil {
		data := w.Base.(*bytes.Buffer).Bytes()
		if err := w.Parent.WriteVarInt(len(data)); err != nil {
			return err
		}
		for len(data) > 0 {
			n, err := w.Parent.Base.Write(data)
			if err != nil {
				return err
			}
			data = data[n:]
		}
		if err := w.Parent.Flush(); err != nil {
			return err
		}
	}
	return nil
}
