package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// DeliverSMResp PDU.
type DeliverSMResp struct {
	base
	MessageID string
}

// NewDeliverSMResp returns new DeliverSMResp.
func NewDeliverSMResp() PDU {
	c := &DeliverSMResp{
		base:      newBase(),
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.DELIVER_SM_RESP
	return c
}

// CanResponse implements PDU interface.
func (c *DeliverSMResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *DeliverSMResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *DeliverSMResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		_ = b.WriteCString(c.MessageID)
	})
}

// Unmarshal implements PDU interface.
func (c *DeliverSMResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		c.MessageID, err = b.ReadCString()
		return
	})
}
