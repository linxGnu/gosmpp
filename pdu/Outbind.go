package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// Outbind PDU is used by the SMSC to signal an ESME to originate a bind_receiver request to the SMSC.
type Outbind struct {
	base
	SystemID string
	Password string
}

// NewOutbind returns Outbind PDU.
func NewOutbind() PDU {
	c := &Outbind{
		base: newBase(),
	}
	c.CommandID = data.OUTBIND
	return c
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
func (c *Outbind) Marshal(b *ByteBuffer) {
	c.base.marshal(b, func(b *ByteBuffer) {
		b.Grow(len(c.SystemID) + len(c.Password) + 2)

		_ = b.WriteCString(c.SystemID)
		_ = b.WriteCString(c.Password)
	})
}

// Unmarshal implements PDU interface.
func (c *Outbind) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, func(b *ByteBuffer) (err error) {
		if c.SystemID, err = b.ReadCString(); err == nil {
			c.Password, err = b.ReadCString()
		}
		return
	})
}
