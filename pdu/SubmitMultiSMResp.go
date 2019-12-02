package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// SubmitMultiSMResp PDU.
type SubmitMultiSMResp struct {
	base
	Request       PDU
	MessageID     string
	UnsuccessSMEs UnsuccessSMEs
}

// NewSubmitMultiSMResp returns new SubmitMultiSMResp.
func NewSubmitMultiSMResp(req PDU) (c *SubmitMultiSMResp) {
	c = &SubmitMultiSMResp{
		base:      newBase(),
		Request:   req,
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.SUBMIT_MULTI_RESP
	return
}

// CanResponse implements PDU interface.
func (c *SubmitMultiSMResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *SubmitMultiSMResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *SubmitMultiSMResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		b.WriteCString(c.MessageID)
		c.UnsuccessSMEs.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitMultiSMResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.MessageID, err = b.ReadCString(); err == nil {
			err = c.UnsuccessSMEs.Unmarshal(b)
		}
		return
	})
}
