package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// Unbind PDU is to deregister an instance of an ESME from the SMSC and inform the SMSC
// that the ESME no longer wishes to use this network connection for the submission or
// delivery of messages.
type Unbind struct {
	base
}

// NewUnbind returns Unbind PDU.
func NewUnbind() PDU {
	c := &Unbind{
		base: newBase(),
	}
	c.CommandID = data.UNBIND
	return c
}

// CanResponse implements PDU interface.
func (c *Unbind) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *Unbind) GetResponse() PDU {
	return NewUnbindRespFromReq(c)
}

// Marshal implements PDU interface.
func (c *Unbind) Marshal(b *ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *Unbind) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
