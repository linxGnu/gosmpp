package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

type SubmitMultiSMResp struct {
	Response
	messageId     string
	unsuccessSMEs *UnsuccessSMEsList
}

func NewSubmitMultiSMResp() *SubmitMultiSMResp {
	a := &SubmitMultiSMResp{}
	a.Construct()

	return a
}

func (c *SubmitMultiSMResp) Construct() {
	defer c.SetRealReference(c)
	c.Response.Construct()

	c.SetCommandId(Data.SUBMIT_MULTI_RESP)
	c.messageId = Data.DFLT_MSGID
	c.unsuccessSMEs = NewUnsuccessSMEsList()
}

func (c *SubmitMultiSMResp) GetInstance() (IPDU, error) {
	return NewSubmitMultiSMResp(), nil
}

func (c *SubmitMultiSMResp) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("SubmitMultiResp: set body buffer is nil")
		return
	}

	dat, err := buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetMessageId(dat)
	if err != nil {
		return
	}

	err = c.unsuccessSMEs.SetData(buf)
	return
}

func (c *SubmitMultiSMResp) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			buf = nil
		}
	}()

	source = c.This.(IPDU)

	dat, err := c.unsuccessSMEs.GetData()
	if err != nil {
		return
	}

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetMessageId())+1+dat.Len()))

	buf.Write_CString(c.GetMessageId())
	err = buf.Write_Buffer(dat)

	return
}

func (c *SubmitMultiSMResp) SetMessageId(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_MSGID_LEN))
	if err != nil {
		return err
	}

	c.messageId = value
	return nil
}

func (c *SubmitMultiSMResp) GetMessageId() string {
	return c.messageId
}

func (c *SubmitMultiSMResp) GetNoUnsuccess() int16 {
	return int16(c.unsuccessSMEs.GetCount())
}

func (c *SubmitMultiSMResp) AddUnsuccessSME(sme *UnsuccessSME) *Exception.Exception {
	return c.unsuccessSMEs.AddValue(sme)
}

func (c *SubmitMultiSMResp) GetUnsuccessSME(index int) Common.IByteData {
	return c.unsuccessSMEs.GetValue(index)
}

type UnsuccessSMEsList struct {
	Common.ByteDataList
}

func NewUnsuccessSMEsList() *UnsuccessSMEsList {
	a := &UnsuccessSMEsList{}
	a.Construct()

	return a
}

func (c *UnsuccessSMEsList) Construct() {
	defer c.SetRealReference(c)
	c.ByteDataList.Construct()

	c.MaxSize = int(Data.SM_MAX_CNT_DEST_ADDR)
	c.LengthOfSize = byte(1)
}

func (c *UnsuccessSMEsList) CreateValue() Common.IByteData {
	return NewUnsuccessSME()
}
