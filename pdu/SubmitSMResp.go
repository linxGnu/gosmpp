package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// SubmitSMResp represents deliver_sm resp.
type SubmitSMResp struct {
	base
	Request   SubmitSM
	MessageID string
}

// NewSubmitSMResp returns new SubmitSMResp.
func NewSubmitSMResp() PDU {
	c := &SubmitSMResp{
		base:      newBase(),
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.SUBMIT_SM_RESP
	return c
}

// NewSubmitSMNewSubmitSMRespFromReqResp returns new SubmitSMResp.
func NewSubmitSMRespFromReq(req SubmitSM) (c *SubmitSMResp) {
	c = &SubmitSMResp{
		base:      newBase(),
		Request:   req,
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.SUBMIT_SM_RESP
	return
}

// CanResponse implements PDU interface.
func (c *SubmitSMResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *SubmitSMResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *SubmitSMResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		_ = b.WriteCString(c.MessageID)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitSMResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		c.MessageID, err = b.ReadCString()
		return
	})
}
