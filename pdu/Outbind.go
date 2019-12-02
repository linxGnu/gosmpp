package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// Outbind PDU.
type Outbind struct {
	base
	SystemID string
	Password string
}

// NewOutbind returns Outbind PDU.
func NewOutbind() (c *Outbind) {
	c = &Outbind{
		base: newBase(),
	}
	c.CommandID = data.CANCEL_SM
	return
}

// CanResponse implements PDU interface.
func (c *Outbind) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *Outbind) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *Outbind) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.SystemID) + len(c.Password) + 2)
		b.WriteCString(c.SystemID)
		b.WriteCString(c.Password)
	})
}

// Unmarshal implements PDU interface.
func (c *Outbind) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.SystemID, err = b.ReadCString(); err == nil {
			c.Password, err = b.ReadCString()
		}
		return
	})
}
