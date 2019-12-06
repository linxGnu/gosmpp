package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// ReplaceSMResp PDU.
type ReplaceSMResp struct {
	base
	Request ReplaceSM
}

// NewReplaceSMResp returns ReplaceSMResp.
func NewReplaceSMResp() PDU {
	c := &ReplaceSMResp{
		base: newBase(),
	}
	c.CommandID = data.REPLACE_SM_RESP
	return c
}

// NewReplaceSMRespFromReq returns ReplaceSMResp.
func NewReplaceSMRespFromReq(req ReplaceSM) (c *ReplaceSMResp) {
	c = &ReplaceSMResp{
		base:    newBase(),
		Request: req,
	}
	c.CommandID = data.REPLACE_SM_RESP
	return
}

// CanResponse implements PDU interface.
func (c *ReplaceSMResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *ReplaceSMResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *ReplaceSMResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, nil)
}

// Unmarshal implements PDU interface.
func (c *ReplaceSMResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, nil)
}
