package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// SubmitMultiResp PDU.
type SubmitMultiResp struct {
	base
	Request       SubmitMulti
	MessageID     string
	UnsuccessSMEs UnsuccessSMEs
}

// NewSubmitMultiResp returns new SubmitMultiResp.
func NewSubmitMultiResp() PDU {
	c := &SubmitMultiResp{
		base:      newBase(),
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.SUBMIT_MULTI_RESP
	return c
}

// NewSubmitMultiRespFromReq returns new SubmitMultiResp.
func NewSubmitMultiRespFromReq(req SubmitMulti) (c *SubmitMultiResp) {
	c = &SubmitMultiResp{
		base:      newBase(),
		Request:   req,
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.SUBMIT_MULTI_RESP
	return
}

// CanResponse implements PDU interface.
func (c *SubmitMultiResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *SubmitMultiResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *SubmitMultiResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		_ = b.WriteCString(c.MessageID)
		c.UnsuccessSMEs.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitMultiResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.MessageID, err = b.ReadCString(); err == nil {
			err = c.UnsuccessSMEs.Unmarshal(b)
		}
		return
	})
}
