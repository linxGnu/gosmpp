package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type Outbind struct {
	Request
	systemId string
	password string
}

func NewOutbind() *Outbind {
	a := &Outbind{}
	a.Construct()

	return a
}

func (c *Outbind) Construct() {
	defer c.SetRealReference(c)
	c.Request.Construct()

	c.SetCommandId(Data.OUTBIND)
}

func (c *Outbind) GetInstance() (IPDU, error) {
	return NewOutbind(), nil
}

func (c *Outbind) CreateResponse() (IResponse, error) {
	return nil, nil
}

func (c *Outbind) CanResponse() bool {
	return false
}

func (c *Outbind) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("Outbind: set body buffer is nil")
		return
	}

	tmp, err := buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetSystemId(tmp)
	if err != nil {
		return
	}

	tmp, err = buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetPassword(tmp)
	if err != nil {
		return
	}

	return
}

func (c *Outbind) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetSystemId())+1+len(c.GetPassword())+1))

	buf.Write_CString(c.GetSystemId())
	err = buf.Write_CString(c.GetPassword())

	return
}

func (c *Outbind) SetSystemId(str string) *Exception.Exception {
	err := c.CheckStringMax(str, int(Data.SM_SYSID_LEN))
	if err != nil {
		return err
	}

	c.systemId = str
	return nil
}

func (c *Outbind) GetSystemId() string {
	return c.systemId
}

func (c *Outbind) SetPassword(str string) *Exception.Exception {
	err := c.CheckStringMax(str, int(Data.SM_PASS_LEN))
	if err != nil {
		return err
	}

	c.password = str
	return nil
}

func (c *Outbind) GetPassword() string {
	return c.password
}
