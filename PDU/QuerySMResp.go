package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type QuerySMResp struct {
	Response
	messageId    string
	finalDate    string
	messageState byte
	errorCode    byte
}

func NewQuerySMResp() *QuerySMResp {
	a := &QuerySMResp{}
	a.Construct()

	return a
}

func (c *QuerySMResp) Construct() {
	defer c.SetRealReference(c)
	c.Response.Construct()

	c.SetCommandId(Data.QUERY_SM_RESP)
	c.messageId = Data.DFLT_MSGID
	c.finalDate = Data.DFLT_DATE
	c.messageState = Data.DFLT_MSG_STATE
	c.errorCode = Data.DFLT_ERR
}

func (c *QuerySMResp) GetInstance() (IPDU, error) {
	return NewQuerySMResp(), nil
}

func (c *QuerySMResp) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("QuerySMResp: set body buffer is nil")
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

	dat, err = buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetFinalDate(dat)
	if err != nil {
		return
	}

	byt, err := buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetMessageState(byt)

	byt, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetErrorCode(byt)

	return
}

func (c *QuerySMResp) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			buf = nil
		}
	}()

	source = c.This.(IPDU)
	err = nil

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetMessageId())+1+len(c.GetFinalDate())+1+(Utils.SZ_BYTE<<1)))

	buf.Write_CString(c.GetMessageId())
	buf.Write_CString(c.GetFinalDate())
	buf.Write_Byte(c.GetMessageState())
	buf.Write_Byte(c.GetErrorCode())

	return
}

func (c *QuerySMResp) SetMessageState(value byte) {
	c.messageState = value
}

func (c *QuerySMResp) GetMessageState() byte {
	return c.messageState
}

func (c *QuerySMResp) SetErrorCode(value byte) {
	c.errorCode = value
}

func (c *QuerySMResp) GetErrorCode() byte {
	return c.errorCode
}

func (c *QuerySMResp) SetMessageId(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_MSGID_LEN))
	if err != nil {
		return err
	}

	c.messageId = value
	return nil
}

func (c *QuerySMResp) GetMessageId() string {
	return c.messageId
}

func (c *QuerySMResp) SetFinalDate(value string) *Exception.Exception {
	err := c.CheckDate(value)
	if err != nil {
		return err
	}

	c.finalDate = value
	return nil
}

func (c *QuerySMResp) GetFinalDate() string {
	return c.finalDate
}
