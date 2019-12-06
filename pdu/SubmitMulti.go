package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// SubmitMulti PDU.
type SubmitMulti struct {
	base
	ServiceType          string
	SourceAddr           Address
	DestAddrs            DestinationAddresses
	EsmClass             byte
	ProtocolID           byte
	PriorityFlag         byte
	ScheduleDeliveryTime string
	ValidityPeriod       string // not used
	RegisteredDelivery   byte
	ReplaceIfPresentFlag byte // not used
	Message              ShortMessage
}

// NewSubmitMulti returns NewSubmitMulti PDU.
func NewSubmitMulti() PDU {
	message, _ := NewShortMessage("")
	c := &SubmitMulti{
		base:                 newBase(),
		ServiceType:          data.DFLT_SRVTYPE,
		SourceAddr:           NewAddress(),
		DestAddrs:            NewDestinationAddresses(),
		EsmClass:             data.DFLT_ESM_CLASS,
		ProtocolID:           data.DFLT_PROTOCOLID,
		PriorityFlag:         data.DFLT_PRIORITY_FLAG,
		ScheduleDeliveryTime: data.DFLT_SCHEDULE,
		ValidityPeriod:       data.DFLT_VALIDITY,
		RegisteredDelivery:   data.DFLT_REG_DELIVERY,
		ReplaceIfPresentFlag: data.DFTL_REPLACE_IFP,
		Message:              message,
	}
	c.CommandID = data.SUBMIT_MULTI
	return c
}

// CanResponse implements PDU interface.
func (c *SubmitMulti) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *SubmitMulti) GetResponse() PDU {
	return NewSubmitMultiRespFromReq(*c)
}

// Marshal implements PDU interface.
func (c *SubmitMulti) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.ServiceType) + len(c.ScheduleDeliveryTime) + len(c.ValidityPeriod) + 10)

		_ = b.WriteCString(c.ServiceType)
		c.SourceAddr.Marshal(b)
		c.DestAddrs.Marshal(b)
		_ = b.WriteByte(c.EsmClass)
		_ = b.WriteByte(c.ProtocolID)
		_ = b.WriteByte(c.PriorityFlag)
		_ = b.WriteCString(c.ScheduleDeliveryTime)
		_ = b.WriteCString(c.ValidityPeriod)
		_ = b.WriteByte(c.RegisteredDelivery)
		_ = b.WriteByte(c.ReplaceIfPresentFlag)
		c.Message.Marshal(b)
	})
}

// Unmarshal implements PDU interface.
func (c *SubmitMulti) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.ServiceType, err = b.ReadCString(); err == nil {
			if err = c.SourceAddr.Unmarshal(b); err == nil {
				if err = c.DestAddrs.Unmarshal(b); err == nil {
					if c.EsmClass, err = b.ReadByte(); err == nil {
						if c.ProtocolID, err = b.ReadByte(); err == nil {
							if c.PriorityFlag, err = b.ReadByte(); err == nil {
								if c.ScheduleDeliveryTime, err = b.ReadCString(); err == nil {
									if c.ValidityPeriod, err = b.ReadCString(); err == nil {
										if c.RegisteredDelivery, err = b.ReadByte(); err == nil {
											if c.ReplaceIfPresentFlag, err = b.ReadByte(); err == nil {
												err = c.Message.Unmarshal(b)
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
