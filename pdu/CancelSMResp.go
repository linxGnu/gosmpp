package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// CancelSMResp PDU.
type CancelSMResp struct {
	base
}

// NewCancelSMResp returns CancelSMResp.
func NewCancelSMResp() (c *CancelSMResp) {
	c = &CancelSMResp{}
	c.CommandID = data.CANCEL_SM_RESP
	return
}

// CanResponse implements PDU interface.
func (c *CancelSMResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *CancelSMResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *CancelSMResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *CancelSMResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
