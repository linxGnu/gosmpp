package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// QuerySMResp PDU.
type QuerySMResp struct {
	base
	Request      *QuerySM
	MessageID    string
	FinalDate    string
	MessageState byte
	ErrorCode    byte
}

// NewQuerySMResp returns new QuerySM PDU.
func NewQuerySMResp(req *QuerySM) (c *QuerySMResp) {
	c = &QuerySMResp{
		base:         newBase(),
		Request:      req,
		MessageID:    req.MessageID,
		FinalDate:    data.DFLT_DATE,
		MessageState: data.DFLT_MSG_STATE,
		ErrorCode:    data.DFLT_ERR,
	}
	c.CommandID = data.QUERY_SM_RESP
	return
}

// CanResponse implements PDU interface.
func (c *QuerySMResp) CanResponse() bool {
	return false
}

// GetResponse implements PDU interface.
func (c *QuerySMResp) GetResponse() PDU {
	return nil
}

// Marshal implements PDU interface.
func (c *QuerySMResp) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + len(c.FinalDate) + 4)

		b.WriteCString(c.MessageID)
		b.WriteCString(c.FinalDate)
		_ = b.WriteByte(c.MessageState)
		_ = b.WriteByte(c.ErrorCode)
	})
}

// Unmarshal implements PDU interface.
func (c *QuerySMResp) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.MessageID, err = b.ReadCString(); err == nil {
			if c.FinalDate, err = b.ReadCString(); err == nil {
				if c.MessageState, err = b.ReadByte(); err == nil {
					c.ErrorCode, err = b.ReadByte()
				}
			}
		}
		return
	})
}
