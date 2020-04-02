package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// CancelSM PDU is issued by the ESME to cancel one or more previously submitted short messages
// that are still pending delivery. The command may specify a particular message to cancel, or
// all messages for a particular source, destination and service_type are to be cancelled.
type CancelSM struct {
	base
	ServiceType string
	MessageID   string
	SourceAddr  Address
	DestAddr    Address
}

// NewCancelSM returns CancelSM PDU.
func NewCancelSM() PDU {
	c := &CancelSM{
		base:        newBase(),
		ServiceType: data.DFLT_SRVTYPE,
		MessageID:   data.DFLT_MSGID,
		SourceAddr:  NewAddress(),
		DestAddr:    NewAddress(),
	}
	c.CommandID = data.CANCEL_SM
	return c
}

// CanResponse implements PDU interface.
func (c *CancelSM) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *CancelSM) GetResponse() PDU {
	return NewCancelSMRespFromReq(c)
}

// Marshal implements PDU interface.
func (c *CancelSM) Marshal(b *ByteBuffer) {
	c.base.marshal(b, func(b *ByteBuffer) {
		b.Grow(len(c.ServiceType) + len(c.MessageID) + 2)

		_ = b.WriteCString(c.ServiceType)
		_ = b.WriteCString(c.MessageID)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *CancelSM) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, func(b *ByteBuffer) (err error) {
		if c.ServiceType, err = b.ReadCString(); err == nil {
			if c.MessageID, err = b.ReadCString(); err == nil {
				if err = c.SourceAddr.Unmarshal(b); err == nil {
					err = c.DestAddr.Unmarshal(b)
				}
			}
		}
		return
	})
}
