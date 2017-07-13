package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/TLV"
	"github.com/linxGnu/gosmpp/Utils"
)

type BindResponse struct {
	Response
	systemId           string
	scInterfaceVersion *TLV.TLVByte
}

func NewBindResponseWithCmdId(cmdId int32) *BindResponse {
	a := &BindResponse{}
	a.Construct()
	a.SetCommandId(cmdId)

	return a
}

func (c *BindResponse) Construct() {
	defer c.SetRealReference(c)
	c.Response.Construct()

	c.scInterfaceVersion = TLV.NewTLVByteWithTag(Data.OPT_PAR_SC_IF_VER)
	c.registerOptional(c.scInterfaceVersion)
}

func (c *BindResponse) GetInstance() (IPDU, error) {
	return &BindResponse{}, nil
}

func (c *BindResponse) SetBody(buffer *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)
	err = nil

	if buffer == nil {
		err = Exception.NewExceptionFromStr("BindResponse: set body is nil")
		return
	}

	if c.GetCommandStatus() == 0 {
		sysId, err1 := buffer.Read_CString()
		if err1 != nil {
			err = err1
			return
		}

		err = c.SetSystemId(sysId)
		return
	}

	return
}

func (c *BindResponse) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetSystemId())<<1+1))
	err = buf.Write_CString(c.GetSystemId())

	return
}

func (c *BindResponse) SetSystemId(sysId string) *Exception.Exception {
	err := c.CheckStringMax(sysId, int(Data.SM_SYSID_LEN))
	if err != nil {
		return err
	}

	c.systemId = sysId
	return nil
}

func (c *BindResponse) GetSystemId() string {
	return c.systemId
}

func (c *BindResponse) HasScInterfaceVersion() bool {
	return c.scInterfaceVersion.HasValue()
}

func (c *BindResponse) SetScInterfaceVersion(value byte) {
	c.scInterfaceVersion.SetValue(value)
}

func (c *BindResponse) GetScInterfaceVersion() (byte, *Exception.Exception) {
	return c.scInterfaceVersion.GetValue()
}
