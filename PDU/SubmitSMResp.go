package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type SubmitSMResp struct {
	Response
	messageId string
}

func NewSubmitSMResp() *SubmitSMResp {
	a := &SubmitSMResp{}
	a.Construct()

	return a
}

func (c *SubmitSMResp) Construct() {
	defer c.SetRealReference(c)
	c.Response.Construct()

	c.SetCommandId(Data.SUBMIT_SM_RESP)
}

func (c *SubmitSMResp) GetInstance() (IPDU, error) {
	return NewSubmitSMResp(), nil
}

func (c *SubmitSMResp) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("SubmitSMResp: set body buffer is nil")
		return
	}

	if c.GetCommandStatus() == 0 {
		st, err1 := buf.Read_CString()
		if err1 != nil {
			err = err1
			return
		}

		err = c.SetMessageId(st)
		return
	}

	return
}

func (c *SubmitSMResp) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			buf = nil
		}
	}()

	source = c.This.(IPDU)

	buf = Utils.NewBuffer(make([]byte, 0, 16))

	if c.GetCommandStatus() == 0 {
		err1 := buf.Write_CString(c.GetMessageId())
		if err1 != nil {
			err = err1
			return
		}
	}

	return
}

func (c *SubmitSMResp) SetMessageId(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_MSGID_LEN))
	if err != nil {
		return err
	}

	c.messageId = value
	return nil
}

func (c *SubmitSMResp) GetMessageId() string {
	return c.messageId
}
