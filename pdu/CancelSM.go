package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// CancelSM PDU.
type CancelSM struct {
	base
	ServiceType string
	MessageID   string
	SourceAddr  Address
	DestAddr    Address
}

// NewCancelSM returns CancelSM PDU.
func NewCancelSM() (c *CancelSM) {
	c = &CancelSM{
		base:        newBase(),
		ServiceType: data.DFLT_SRVTYPE,
		MessageID:   data.DFLT_MSGID,
		SourceAddr:  NewAddress(),
		DestAddr:    NewAddress(),
	}
	c.CommandID = data.CANCEL_SM
	return
}

// CanResponse implements PDU interface.
func (c *CancelSM) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *CancelSM) GetResponse() PDU {
	return NewCancelSMResp(c)
}

// Marshal implements PDU interface.
func (c *CancelSM) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.ServiceType) + len(c.MessageID) + 2)

		b.WriteCString(c.ServiceType)
		b.WriteCString(c.MessageID)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *CancelSM) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
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
