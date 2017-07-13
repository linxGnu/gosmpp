package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type IBindRequest interface {
	IRequest
	IsTransmitter() bool
	IsReceiver() bool
}

type BindRequest struct {
	Request
	systemId         string
	password         string
	systemType       string
	addressRange     *AddressRange
	interfaceVersion byte
}

func NewBindRequest() *BindRequest {
	a := &BindRequest{}
	a.Construct()

	return a
}

func NewBindRequestWithCmdId(cmdId int32) *BindRequest {
	a := NewBindRequest()
	a.SetCommandId(cmdId)

	return a
}

func (c *BindRequest) Construct() {
	defer c.SetRealReference(c)
	c.Request.Construct()

	c.systemId = Data.DFLT_SYSID
	c.password = Data.DFLT_PASS
	c.systemType = Data.DFLT_SYSTYPE
	c.addressRange = NewAddressRange()
	c.interfaceVersion = Data.SMPP_V34
}

func (c *BindRequest) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("BindRequest: buffer input is nil")
		return
	}

	dat, err := buf.Read_CString()
	if err != nil {
		return
	}

	err = c.SetSystemId(dat)
	if err != nil {
		return
	}

	dat, err = buf.Read_CString()
	if err != nil {
		return
	}

	err = c.SetPassword(dat)
	if err != nil {
		return
	}

	dat, err = buf.Read_CString()
	if err != nil {
		return
	}

	err = c.SetSystemType(dat)
	if err != nil {
		return
	}

	bt, err := buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetInterfaceVersion(bt)

	err = c.addressRange.SetData(buf)
	return
}

func (c *BindRequest) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	dat, err := c.addressRange.GetData()
	if err != nil {
		return
	}

	buf = Utils.NewBuffer(make([]byte, 0, len(c.systemId)+1+len(c.password)+1+len(c.systemType)+1+Utils.SZ_BYTE+dat.Len()))

	buf.Write_CString(c.systemId)
	buf.Write_CString(c.password)
	buf.Write_CString(c.systemType)
	buf.Write_UnsafeByte(c.interfaceVersion)
	buf.Write_Buffer(dat)

	return
}

func (c *BindRequest) SetSystemId(sysId string) *Exception.Exception {
	err := c.CheckStringMax(sysId, int(Data.SM_SYSID_LEN))
	if err != nil {
		return err
	}

	c.systemId = sysId
	return nil
}

func (c *BindRequest) SetPassword(pass string) *Exception.Exception {
	err := c.CheckStringMax(pass, int(Data.SM_PASS_LEN))
	if err != nil {
		return err
	}

	c.password = pass
	return nil
}

func (c *BindRequest) SetSystemType(t string) *Exception.Exception {
	err := c.CheckStringMax(t, int(Data.SM_SYSTYPE_LEN))
	if err != nil {
		return err
	}

	c.systemType = t
	return nil
}

func (c *BindRequest) SetInterfaceVersion(version byte) {
	c.interfaceVersion = version
}

func (c *BindRequest) SetAddressRange(addr *AddressRange) {
	c.addressRange = addr
}

func (c *BindRequest) SetAddressRangeFromTonNpiString(ton, npi byte, rangeString string) error {
	addrRange, err := NewAddressRangeWithTonNpiAddr(ton, npi, rangeString)
	if err != nil {
		return err
	}

	c.SetAddressRange(addrRange)
	return nil
}

func (c *BindRequest) SetAddressRangeFromString(rangeString string) error {
	addrRange, err := NewAddressRangeWithAddr(rangeString)
	if err != nil {
		return err
	}

	c.SetAddressRange(addrRange)
	return nil
}

func (c *BindRequest) GetSystemId() string {
	return c.systemId
}

func (c *BindRequest) GetPassword() string {
	return c.password
}

func (c *BindRequest) GetSystemType() string {
	return c.systemType
}

func (c *BindRequest) GetInterfaceVersion() byte {
	return c.interfaceVersion
}

func (c *BindRequest) GetAddressRange() *AddressRange {
	return c.addressRange
}

func (c *BindRequest) IsTransmitter() bool {
	return false
}

func (c *BindRequest) IsReceiver() bool {
	return false
}
