package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// EnquireLink PDU.
type EnquireLink struct {
	base
}

// NewEnquireLink returns new EnquireLink PDU.
func NewEnquireLink() PDU {
	c := &EnquireLink{
		base: newBase(),
	}
	c.CommandID = data.ENQUIRE_LINK
	return c
}

// CanResponse implements PDU interface.
func (c *EnquireLink) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *EnquireLink) GetResponse() PDU {
	return NewEnquireLinkResp()
}

// Marshal implements PDU interface.
func (c *EnquireLink) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *EnquireLink) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
