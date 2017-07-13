package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

type DeliverSMResp struct {
	Response
	messageID string
}

func NewDeliverSMResp() IResponse {
	a := &DeliverSMResp{}
	a.Construct()

	return a
}

func (c *DeliverSMResp) Construct() {
	defer c.SetRealReference(c)
	c.Response.Construct()

	c.SetCommandId(Data.DELIVER_SM_RESP)
	c.messageID = Data.DFLT_MSGID
}

func (c *DeliverSMResp) GetInstance() (IPDU, error) {
	return NewDeliverSMResp(), nil
}

func (c *DeliverSMResp) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("DeliverSMResp: set body buffer is nil")
		return
	}

	messageId, err := buf.Read_CString()
	if err != nil {
		return
	}

	err = c.SetMessageID(messageId)
	return
}

func (c *DeliverSMResp) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	buf = Utils.NewBuffer(make([]byte, 0, 16))
	err = buf.Write_CString(c.messageID)

	return
}

func (c *DeliverSMResp) SetMessageID(value string) *Exception.Exception {
	err := c.CheckStringMax(value, int(Data.SM_MSGID_LEN))
	if err != nil {
		return err
	}

	c.messageID = value
	return nil
}

func (c *DeliverSMResp) GetMessageID() string {
	return c.messageID
}
