package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// UnbindResp PDU.
type UnbindResp struct {
	base
	Request Unbind
}

// NewUnbindResp returns UnbindResp.
func NewUnbindResp() PDU {
	c := &UnbindResp{
		base: newBase(),
	}
	c.CommandID = data.UNBIND_RESP
	return c
}

// NewUnbindRespFromReq returns UnbindResp.
func NewUnbindRespFromReq(req Unbind) (c *UnbindResp) {
	c = &UnbindResp{
		base:    newBase(),
		Request: req,
	}
	c.CommandID = data.UNBIND_RESP
	return
}

// CanResponse implements PDU interface.
func (c *UnbindResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *UnbindResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *UnbindResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *UnbindResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
