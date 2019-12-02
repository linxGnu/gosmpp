package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// DataSMResp PDU.
type DataSMResp struct {
	base
	Request   DataSM
	MessageID string
}

// NewDataSMResp returns DataSMResp.
func NewDataSMResp() PDU {
	c := &DataSMResp{
		base:      newBase(),
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.DATA_SM_RESP
	return c
}

// NewDataSMRespFromReq returns DataSMResp.
func NewDataSMRespFromReq(req DataSM) (c *DataSMResp) {
	c = &DataSMResp{
		base:      newBase(),
		Request:   req,
		MessageID: data.DFLT_MSGID,
	}
	c.CommandID = data.DATA_SM_RESP
	return
}

// CanResponse implements PDU interface.
func (c *DataSMResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *DataSMResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *DataSMResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + 1)

		_ = b.WriteCString(c.MessageID)
	})
}

// Unmarshal implements PDU interface.
func (c *DataSMResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		c.MessageID, err = b.ReadCString()
		return
	})
}
