package pdu

import (
	"encoding/binary"
	"sync/atomic"
)

func nextSequenceNumber(s *int32) (v int32) {
	// & 0x7FFFFFFF: cater for integer overflow
	// Allowed range is 0x01 to 0x7FFFFFFF. This
	// will still result in a single invalid value
	// of 0x00 every ~2 billion PDUs (not too bad):
	if v = atomic.AddInt32(s, 1) & 0x7FFFFFFF; v <= 0 {
		v = 1
	}
	return
}

// Header represents PDU header.
type Header struct {
	CommandLength  int32
	CommandID      int32
	CommandStatus  int32
	SequenceNumber int32
}

// ParseHeader parses PDU header.
func ParseHeader(v [16]byte) (h Header) {
	h.CommandLength = int32(binary.BigEndian.Uint32(v[:]))
	h.CommandID = int32(binary.BigEndian.Uint32(v[4:]))
	h.CommandStatus = int32(binary.BigEndian.Uint32(v[8:]))
	h.SequenceNumber = int32(binary.BigEndian.Uint32(v[12:]))
	return
}

// Unmarshal from buffer.
func (c *Header) Unmarshal(b *ByteBuffer) (err error) {
	c.CommandLength, err = b.ReadInt()
	if err == nil {
		c.CommandID, err = b.ReadInt()
		if err == nil {
			if c.CommandStatus, err = b.ReadInt(); err == nil {
				c.SequenceNumber, err = b.ReadInt()
			}
		}
	}
	return
}

var sequenceNumber int32

// AssignSequenceNumber assigns sequence number auto-incrementally.
func (c *Header) AssignSequenceNumber() {
	c.SetSequenceNumber(nextSequenceNumber(&sequenceNumber))
}

// ResetSequenceNumber resets sequence number.
func (c *Header) ResetSequenceNumber() {
	c.SequenceNumber = 1
}

// GetSequenceNumber returns assigned sequence number.
func (c *Header) GetSequenceNumber() int32 {
	return c.SequenceNumber
}

// SetSequenceNumber manually sets sequence number.
func (c *Header) SetSequenceNumber(v int32) {
	c.SequenceNumber = v
}

// Marshal to buffer.
func (c *Header) Marshal(b *ByteBuffer) {
	b.Grow(16)
	b.WriteInt(c.CommandLength)
	b.WriteInt(c.CommandID)
	b.WriteInt(c.CommandStatus)
	b.WriteInt(c.SequenceNumber)
}
