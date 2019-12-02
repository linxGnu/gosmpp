package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// SubmitSM PDU.
type SubmitSM struct {
	base
	ServiceType          string
	SourceAddr           Address
	DestAddr             Address
	EsmClass             byte
	ProtocolID           byte
	PriorityFlag         byte
	ScheduleDeliveryTime string // not used
	ValidityPeriod       string // not used
	RegisteredDelivery   byte
	ReplaceIfPresentFlag byte // not used
	DataCoding           byte
	SmDefaultMsgID       byte
	ShortMessage         ShortMessage
}

// NewSubmitSM returns SubmitSM PDU.
func NewSubmitSM() (c *SubmitSM) {
	message, _ := NewShortMessage("")
	c = &SubmitSM{
		base:                 newBase(),
		ServiceType:          data.DFLT_SRVTYPE,
		SourceAddr:           NewAddress(),
		DestAddr:             NewAddress(),
		EsmClass:             data.DFLT_ESM_CLASS,
		ProtocolID:           data.DFLT_PROTOCOLID,
		PriorityFlag:         data.DFLT_PRIORITY_FLAG,
		ScheduleDeliveryTime: data.DFLT_SCHEDULE,
		ValidityPeriod:       data.DFLT_VALIDITY,
		RegisteredDelivery:   data.DFLT_REG_DELIVERY,
		ReplaceIfPresentFlag: data.DFTL_REPLACE_IFP,
		DataCoding:           data.DFLT_DATA_CODING,
		SmDefaultMsgID:       data.DFLT_DFLTMSGID,
		ShortMessage:         message,
	}
	c.CommandID = data.SUBMIT_SM
	return
}

// CanResponse implements PDU interface.
func (c *SubmitSM) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *SubmitSM) GetResponse() PDU {
	return NewSubmitSMResp(c)
}

// Marshal implements PDU interface.
func (c *SubmitSM) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.ServiceType) + len(c.ScheduleDeliveryTime) + len(c.ValidityPeriod) + 10)

		b.WriteCString(c.ServiceType)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
		_ = b.WriteByte(c.EsmClass)
		_ = b.WriteByte(c.ProtocolID)
		_ = b.WriteByte(c.PriorityFlag)
		b.WriteCString(c.ScheduleDeliveryTime)
		b.WriteCString(c.ValidityPeriod)
		_ = b.WriteByte(c.RegisteredDelivery)
		_ = b.WriteByte(c.ReplaceIfPresentFlag)
		_ = b.WriteByte(c.DataCoding)
		_ = b.WriteByte(c.SmDefaultMsgID)
		c.ShortMessage.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitSM) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.ServiceType, err = b.ReadCString(); err == nil {
			if err = c.SourceAddr.Unmarshal(b); err == nil {
				if err = c.DestAddr.Unmarshal(b); err == nil {
					if c.EsmClass, err = b.ReadByte(); err == nil {
						if c.ProtocolID, err = b.ReadByte(); err == nil {
							if c.PriorityFlag, err = b.ReadByte(); err == nil {
								if c.ScheduleDeliveryTime, err = b.ReadCString(); err == nil {
									if c.ValidityPeriod, err = b.ReadCString(); err == nil {
										if c.RegisteredDelivery, err = b.ReadByte(); err == nil {
											if c.ReplaceIfPresentFlag, err = b.ReadByte(); err == nil {
												if c.DataCoding, err = b.ReadByte(); err == nil {
													if c.SmDefaultMsgID, err = b.ReadByte(); err == nil {
														c.ShortMessage.Unmarshal(b)
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return
	})
}
