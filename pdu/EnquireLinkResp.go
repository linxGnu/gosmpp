package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// EnquireLinkResp PDU.
type EnquireLinkResp struct {
	base
}

// NewEnquireLinkResp returns EnquireLinkResp.
func NewEnquireLinkResp() PDU {
	c := &EnquireLinkResp{
		base: newBase(),
	}
	c.CommandID = data.ENQUIRE_LINK_RESP
	return c
}

// CanResponse implements PDU interface.
func (c *EnquireLinkResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *EnquireLinkResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *EnquireLinkResp) Marshal(b *ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *EnquireLinkResp) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
