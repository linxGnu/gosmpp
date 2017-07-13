package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

type IPDUHeader interface {
	GetSequenceNumber() int32
	SetSequenceNumber(seq int32)
	GetCommandId() int32
	SetCommandId(cmdId int32)
	GetCommandLength() int32
	SetCommandLength(length int32)
	GetCommandStatus() int32
	SetCommandStatus(status int32)
}

type PDUHeader struct {
	Common.ByteData
	CommandLength  int32
	CommandId      int32
	CommandStatus  int32
	SequenceNumber int32
}

func NewPDUHeader() *PDUHeader {
	a := &PDUHeader{}
	a.Construct()

	return a
}

func (c *PDUHeader) Construct() {
	defer c.SetRealReference(c)
	c.ByteData.Construct()

	c.SequenceNumber = 1
}

func (c *PDUHeader) GetData() (res *Utils.ByteBuffer, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			res = nil
		}
	}()

	buf := Utils.NewBuffer(make([]byte, 0, Utils.SZ_INT*4))

	buf.Write_UnsafeInt(c.CommandLength)
	buf.Write_UnsafeInt(c.CommandId)
	buf.Write_UnsafeInt(c.CommandStatus)
	buf.Write_UnsafeInt(c.SequenceNumber)

	return buf, nil
}

func (c *PDUHeader) GetCommandLength() int32 {
	return c.CommandLength
}

func (c *PDUHeader) SetCommandLength(length int32) {
	c.CommandLength = length
}

func (c *PDUHeader) GetCommandId() int32 {
	return c.CommandId
}

func (c *PDUHeader) SetCommandId(cmdId int32) {
	c.CommandId = cmdId
}

func (c *PDUHeader) GetCommandStatus() int32 {
	return c.CommandStatus
}

func (c *PDUHeader) SetCommandStatus(status int32) {
	c.CommandStatus = status
}

func (c *PDUHeader) GetSequenceNumber() int32 {
	return c.SequenceNumber
}

func (c *PDUHeader) SetSequenceNumber(seq int32) {
	c.SequenceNumber = seq
}

func (c *PDUHeader) SetData(buf *Utils.ByteBuffer) *Exception.Exception {
	if buf == nil || buf.Buffer == nil {
		return Exception.NewExceptionFromStr("PDUHeader: buffer passing is nil")
	}

	val, err := buf.Read_Int()
	if err != nil {
		return err
	}
	c.CommandLength = val

	val, err = buf.Read_Int()
	if err != nil {
		return err
	}
	c.CommandId = val

	val, err = buf.Read_Int()
	if err != nil {
		return err
	}
	c.CommandStatus = val

	val, err = buf.Read_Int()
	if err != nil {
		return err
	}
	c.SequenceNumber = val

	return nil
}
