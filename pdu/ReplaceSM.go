package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// ReplaceSM PDU.
type ReplaceSM struct {
	base
	MessageID            string
	SourceAddr           Address
	ScheduleDeliveryTime string
	ValidityPeriod       string
	RegisteredDelivery   byte
	SmDefaultMsgID       byte
	ShortMessage         ShortMessage
}

// NewReplaceSM returns ReplaceSM PDU.
func NewReplaceSM() (c *ReplaceSM) {
	message, _ := NewShortMessage("")
	c = &ReplaceSM{
		base:                 newBase(),
		SourceAddr:           NewAddress(),
		ScheduleDeliveryTime: data.DFLT_SCHEDULE,
		ValidityPeriod:       data.DFLT_VALIDITY,
		RegisteredDelivery:   data.DFLT_REG_DELIVERY,
		SmDefaultMsgID:       data.DFLT_DFLTMSGID,
		ShortMessage:         message,
	}
	c.CommandID = data.REPLACE_SM
	return
}

// CanResponse implements PDU interface.
func (c *ReplaceSM) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *ReplaceSM) GetResponse() PDU {
	return NewReplaceSMResp(c)
}

// Marshal implements PDU interface.
func (c *ReplaceSM) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.MessageID) + len(c.ScheduleDeliveryTime) + len(c.ValidityPeriod) + 5)

		b.WriteCString(c.MessageID)
		c.SourceAddr.Marshal(b)
		b.WriteCString(c.ScheduleDeliveryTime)
		b.WriteCString(c.ValidityPeriod)
		_ = b.WriteByte(c.RegisteredDelivery)
		_ = b.WriteByte(c.SmDefaultMsgID)
		c.ShortMessage.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *ReplaceSM) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.MessageID, err = b.ReadCString(); err == nil {
			if err = c.SourceAddr.Unmarshal(b); err == nil {
				if c.ScheduleDeliveryTime, err = b.ReadCString(); err == nil {
					if c.ValidityPeriod, err = b.ReadCString(); err == nil {
						if c.RegisteredDelivery, err = b.ReadByte(); err == nil {
							if c.SmDefaultMsgID, err = b.ReadByte(); err == nil {
								c.ShortMessage.Unmarshal(b)
							}
						}
					}
				}
			}
		}
		return
	})
}
