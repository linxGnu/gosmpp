package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// Unbind PDU.
type Unbind struct {
	base
}

// NewUnbind returns Unbind PDU.
func NewUnbind() (c *Unbind) {
	c = &Unbind{
		base: newBase(),
	}
	c.CommandID = data.UNBIND
	return
}

// CanResponse implements PDU interface.
func (c *Unbind) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *Unbind) GetResponse() PDU {
	return NewUnbindResp(c)
}

// Marshal implements PDU interface.
func (c *Unbind) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *Unbind) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
