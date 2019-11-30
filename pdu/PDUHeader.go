package pdu

import (
	"encoding/binary"

	"github.com/linxGnu/gosmpp/utils"
)

// Header represents PDU header.
type Header struct {
	CommandLength  int32
	CommandID      int32
	CommandStatus  int32
	SequenceNumber int32
}

// NewHeader creates new PDU Header.
func NewHeader() *Header {
	return &Header{SequenceNumber: 1}
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
func (c *Header) Unmarshal(b *utils.ByteBuffer) (err error) {
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

// AssignSequenceNumber assigns sequence number.
func (c *Header) AssignSequenceNumber() {
	c.SequenceNumber = nextSequenceNumber()
}

// ResetSequenceNumber resets sequence number.
func (c *Header) ResetSequenceNumber() {
	c.SequenceNumber = 1
}

// Marshal to buffer.
func (c *Header) Marshal(b *utils.ByteBuffer) {
	b.Grow(16)
	b.WriteInt(c.CommandLength)
	b.WriteInt(c.CommandID)
	b.WriteInt(c.CommandStatus)
	b.WriteInt(c.SequenceNumber)
}
