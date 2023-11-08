package pdu

import (
	"errors"
	"github.com/linxGnu/gosmpp/data"
	"io"
)

// SubmitSMResp PDU.
type SubmitSMResp struct {
	base
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

// NewSubmitSMRespFromReq returns new SubmitSMResp.
func NewSubmitSMRespFromReq(req *SubmitSM) PDU {
	c := NewSubmitSMResp().(*SubmitSMResp)
	if req != nil {
		c.SequenceNumber = req.SequenceNumber
	}
	return c
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
func (c *SubmitSMResp) Marshal(b *ByteBuffer) {
	c.base.marshal(b, func(b *ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		_ = b.WriteCString(c.MessageID)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitSMResp) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, func(b *ByteBuffer) (err error) {
		c.MessageID, err = b.ReadCString()
		if errors.Is(err, io.EOF) {
			return nil
		}
		return
	})
}
