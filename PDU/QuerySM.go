package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type QuerySM struct {
	Request
	messageId  string
	sourceAddr *Address
}

func NewQuerySM() *QuerySM {
	a := &QuerySM{}
	a.Construct()

	return a
}

func (c *QuerySM) Construct() {
	defer c.SetRealReference(c)
	c.Request.Construct()

	c.SetCommandId(Data.QUERY_SM)
	c.sourceAddr = NewAddress()
}

func (c *QuerySM) GetInstance() (IPDU, error) {
	return NewQuerySM(), nil
}

func (c *QuerySM) CreateResponse() (IResponse, error) {
	return NewQuerySMResp(), nil
}

func (c *QuerySM) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("QuerySM: set body buffer is nil")
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

	err = c.sourceAddr.SetData(buf)
	return
}

func (c *QuerySM) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			buf = nil
		}
	}()

	source = c.This.(IPDU)

	addr, err := c.GetSourceAddr().GetData()
	if err != nil {
		return
	}

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetMessageId())+1+addr.Len()))

	buf.Write_CString(c.GetMessageId())
	err = buf.Write_Buffer(addr)

	return
}

func (c *QuerySM) SetSourceAddr(value *Address) {
	c.sourceAddr = value
}

func (c *QuerySM) SetSourceAddrWithAddr(value string) *Exception.Exception {
	a, err := NewAddressWithAddr(value)
	if err != nil {
		return err
	}

	c.sourceAddr = a
	return nil
}

func (c *QuerySM) SetSourceAddrWithTonNpiAddr(ton, npi byte, value string) *Exception.Exception {
	a, err := NewAddressWithTonNpiAddr(ton, npi, value)
	if err != nil {
		return err
	}

	c.sourceAddr = a
	return nil
}

func (c *QuerySM) SetMessageId(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_MSGID_LEN))
	if err != nil {
		return err
	}

	c.messageId = value
	return nil
}

func (c *QuerySM) GetMessageId() string {
	return c.messageId
}

func (c *QuerySM) GetSourceAddr() *Address {
	return c.sourceAddr
}
