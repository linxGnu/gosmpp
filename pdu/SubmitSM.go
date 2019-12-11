package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// SubmitSM PDU is used by an ESME to submit a short message to the SMSC for onward
// transmission to a specified short message entity (SME). The submit_sm PDU does
// not support the transaction message mode.
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
	Message              ShortMessage
}

// NewSubmitSM returns SubmitSM PDU.
func NewSubmitSM() PDU {
	message, _ := NewShortMessage("")
	c := &SubmitSM{
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
		Message:              message,
	}
	c.CommandID = data.SUBMIT_SM
	return c
}

// CanResponse implements PDU interface.
func (c *SubmitSM) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *SubmitSM) GetResponse() PDU {
	return NewSubmitSMResp()
}

// Marshal implements PDU interface.
func (c *SubmitSM) Marshal(b *ByteBuffer) {
	c.base.marshal(b, func(b *ByteBuffer) {
		b.Grow(len(c.ServiceType) + len(c.ScheduleDeliveryTime) + len(c.ValidityPeriod) + 10)

		_ = b.WriteCString(c.ServiceType)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
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
func (c *SubmitSM) Unmarshal(b *ByteBuffer) error {
	return c.base.unmarshal(b, func(b *ByteBuffer) (err error) {
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
