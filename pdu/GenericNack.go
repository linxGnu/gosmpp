package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// GenerickNack PDU.
type GenerickNack struct {
	base
}

// NewGenerickNack returns new GenerickNack PDU.
func NewGenerickNack() (c *GenerickNack) {
	c = &GenerickNack{
		base: newBase(),
	}
	c.CommandID = data.GENERIC_NACK
	return
}

// CanResponse implements PDU interface.
func (c *GenerickNack) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *GenerickNack) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *GenerickNack) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *GenerickNack) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
