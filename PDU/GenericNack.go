package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type GenericNack struct {
	Response
}

func NewGenericNack() *GenericNack {
	a := &GenericNack{}
	a.Construct()

	return a
}

func NewGenericNackWithCmStatusSeqNum(cmdStatus, seqNumber int32) *GenericNack {
	a := NewGenericNack()
	a.SetCommandStatus(cmdStatus)
	a.SetSequenceNumber(seqNumber)

	return a
}

func (c *GenericNack) Construct() {
	defer c.SetRealReference(c)
	c.Response.Construct()

	c.SetCommandId(Data.GENERIC_NACK)
}

func (c *GenericNack) GetInstance() (IPDU, error) {
	return NewGenericNack(), nil
}

func (c *GenericNack) SetBody(buf *Utils.ByteBuffer) (*Exception.Exception, IPDU) {
	return nil, nil
}

func (c *GenericNack) GetBody() (*Utils.ByteBuffer, *Exception.Exception, IPDU) {
	return nil, nil, nil
}
