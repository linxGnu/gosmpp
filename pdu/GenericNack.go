package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// GenericNack PDU is a generic negative acknowledgement to an SMPP PDU submitted
// with an invalid message header. A generic_nack response is returned in the following cases:
//
// - Invalid command_length
//   If the receiving SMPP entity, on decoding an SMPP PDU, detects an invalid command_length
//   (either too short or too long), it should assume that the data is corrupt. In such cases
//   a generic_nack PDU must be returned to the message originator.
//
// - Unknown command_id
//   If an unknown or invalid command_id is received, a generic_nack PDU must also be returned to the originator.
type GenericNack struct {
	base
}

// NewGenericNack returns new GenericNack PDU.
func NewGenericNack() PDU {
	c := &GenericNack{
		base: newBase(),
	}
	c.CommandID = data.GENERIC_NACK
	return c
}

// CanResponse implements PDU interface.
func (c *GenericNack) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *GenericNack) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *GenericNack) Marshal(b *ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *GenericNack) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
