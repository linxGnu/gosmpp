package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// DataSM PDU.
type DataSM struct {
	base
	ServiceType        string
	SourceAddr         Address
	DestAddr           Address
	EsmClass           byte
	RegisteredDelivery byte
	DataCoding         byte
}

// NewDataSM returns new data sm pdu.
func NewDataSM() PDU {
	c := &DataSM{
		base:               newBase(),
		ServiceType:        data.DFLT_SRVTYPE,
		SourceAddr:         NewAddressWithMaxLength(data.SM_DATA_ADDR_LEN),
		DestAddr:           NewAddressWithMaxLength(data.SM_DATA_ADDR_LEN),
		EsmClass:           data.DFLT_ESM_CLASS,
		RegisteredDelivery: data.DFLT_REG_DELIVERY,
		DataCoding:         data.DFLT_DATA_CODING,
	}
	c.CommandID = data.DATA_SM
	return c
}

// CanResponse implements PDU interface.
func (c *DataSM) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (c *DataSM) GetResponse() PDU {
	return NewDataSMRespFromReq(*c)
}

// Marshal implements PDU interface.
func (c *DataSM) Marshal(b *utils.ByteBuffer) {
	c.base.marshal(b, func(b *utils.ByteBuffer) {
		b.Grow(len(c.ServiceType) + 4)

		b.WriteCString(c.ServiceType)
		c.SourceAddr.Marshal(b)
		c.DestAddr.Marshal(b)
		_ = b.WriteByte(c.EsmClass)
		_ = b.WriteByte(c.RegisteredDelivery)
		_ = b.WriteByte(c.DataCoding)
	})
}

// Unmarshal implements PDU interface.
func (c *DataSM) Unmarshal(b *utils.ByteBuffer) error {
	return c.base.unmarshal(b, func(b *utils.ByteBuffer) (err error) {
		if c.ServiceType, err = b.ReadCString(); err == nil {
			if err = c.SourceAddr.Unmarshal(b); err == nil {
				if err = c.DestAddr.Unmarshal(b); err == nil {
					if c.EsmClass, err = b.ReadByte(); err == nil {
						if c.RegisteredDelivery, err = b.ReadByte(); err == nil {
							c.DataCoding, err = b.ReadByte()
						}
					}
				}
			}
		}
		return
	})
}
