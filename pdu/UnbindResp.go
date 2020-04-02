package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// UnbindResp PDU.
type UnbindResp struct {
	base
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
func NewUnbindRespFromReq(req *Unbind) PDU {
	c := NewUnbindResp().(*UnbindResp)
	if req != nil {
		c.SequenceNumber = req.SequenceNumber
	}
	return c
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
func (c *UnbindResp) Marshal(b *ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *UnbindResp) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
