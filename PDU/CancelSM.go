package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type CancelSM struct {
	Request
	serviceType string
	messageId   string
	sourceAddr  *Address
	destAddr    *Address
}

func NewCancelSM() *CancelSM {
	a := &CancelSM{}
	a.Construct()

	return a
}

func (c *CancelSM) Construct() {
	defer c.SetRealReference(c)
	c.Request.Construct()

	c.SetCommandId(Data.CANCEL_SM)
	c.serviceType = Data.DFLT_SRVTYPE
	c.messageId = Data.DFLT_MSGID
	c.sourceAddr = NewAddress()
	c.destAddr = NewAddress()
}

func (c *CancelSM) GetInstance() (IPDU, error) {
	return NewCancelSM(), nil
}

func (c *CancelSM) CreateResponse() (IResponse, error) {
	return NewCancelSMResp(), nil
}

func (c *CancelSM) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("CancelSM: set body buffer is nil")
		return
	}

	val, err := buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetServiceType(val)
	if err != nil {
		return
	}

	val, err = buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetMessageId(val)
	if err != nil {
		return
	}

	err = c.sourceAddr.SetData(buf)
	if err != nil {
		return
	}

	err = c.destAddr.SetData(buf)
	return
}

func (c *CancelSM) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	src, err := c.GetSourceAddr().GetData()
	if err != nil {
		return
	}

	des, err := c.GetDestAddr().GetData()
	if err != nil {
		return
	}

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetServiceType())+1+len(c.GetMessageId())+1+src.Len()+des.Len()))

	buf.Write_CString(c.GetServiceType())
	buf.Write_CString(c.GetMessageId())

	err = buf.Write_Buffer(src)
	if err != nil {
		return
	}

	err = buf.Write_Buffer(des)

	return
}

func (c *CancelSM) SetServiceType(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_SRVTYPE_LEN))
	if err != nil {
		return err
	}

	c.serviceType = value
	return nil
}

func (c *CancelSM) SetMessageId(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_MSGID_LEN))
	if err != nil {
		return err
	}

	c.messageId = value
	return nil
}

func (c *CancelSM) SetSourceAddr(value *Address) {
	c.sourceAddr = value
}

func (c *CancelSM) SetSourceAddrFromStr(st string) *Exception.Exception {
	tmp, err := NewAddressWithAddr(st)
	if err != nil {
		return err
	}

	c.sourceAddr = tmp
	return nil
}

func (c *CancelSM) SetSourceAddrFromStrTon(ton, npi byte, st string) *Exception.Exception {
	tmp, err := NewAddressWithTonNpiAddr(ton, npi, st)
	if err != nil {
		return err
	}

	c.sourceAddr = tmp
	return nil
}

func (c *CancelSM) SetDestAddr(value *Address) {
	c.destAddr = value
}

func (c *CancelSM) SetDestAddrFromStr(st string) *Exception.Exception {
	tmp, err := NewAddressWithAddr(st)
	if err != nil {
		return err
	}

	c.destAddr = tmp
	return nil
}

func (c *CancelSM) SetDestAddrFromStrTon(ton, npi byte, st string) *Exception.Exception {
	tmp, err := NewAddressWithTonNpiAddr(ton, npi, st)
	if err != nil {
		return err
	}

	c.destAddr = tmp
	return nil
}

func (c *CancelSM) GetServiceType() string {
	return c.serviceType
}

func (c *CancelSM) GetMessageId() string {
	return c.messageId
}

func (c *CancelSM) GetSourceAddr() *Address {
	return c.sourceAddr
}

func (c *CancelSM) GetDestAddr() *Address {
	return c.destAddr
}
