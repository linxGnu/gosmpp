package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

type ReplaceSM struct {
	Request
	messageId            string
	sourceAddr           *Address
	scheduleDeliveryTime string
	validityPeriod       string
	registeredDelivery   byte
	smDefaultMsgId       byte
	smLength             int16
	shortMessage         *ShortMessage
}

func NewReplaceSM() *ReplaceSM {
	a := &ReplaceSM{}
	a.Construct()

	return a
}

func (a *ReplaceSM) Construct() {
	defer a.SetRealReference(a)
	a.Request.Construct()

	a.messageId = Data.DFLT_MSGID
	a.sourceAddr = NewAddress()
	a.scheduleDeliveryTime = Data.DFLT_SCHEDULE
	a.validityPeriod = Data.DFLT_VALIDITY
	a.registeredDelivery = Data.DFLT_REG_DELIVERY
	a.smDefaultMsgId = Data.DFLT_DFLTMSGID
	a.smLength = int16(Data.DFLT_MSG_LEN)
	a.shortMessage = NewShortMessageWithMaxLength(int32(Data.SM_MSG_LEN))
}

func (c *ReplaceSM) GetInstance() (IPDU, error) {
	return NewReplaceSM(), nil
}

func (c *ReplaceSM) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("ReplaceSM: set body buffer is nil")
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
	if err != nil {
		return
	}

	dat, err = buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetScheduleDeliveryTime(dat)
	if err != nil {
		return
	}

	dat, err = buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetValidityPeriod(dat)
	if err != nil {
		return
	}

	byt, err := buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetRegisteredDelivery(byt)

	byt, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetSmDefaultMsgId(byt)

	byt, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetSmLength(Common.DecodeUnsigned(byt))

	tmp, err := buf.Read_Bytes(int(c.GetSmLength()))
	if err != nil {
		return
	}

	err = c.shortMessage.SetData(tmp)
	return
}

func (c *ReplaceSM) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			buf = nil
		}
	}()

	source = c.This.(IPDU)

	dat, err := c.GetSourceAddr().GetData()
	if err != nil {
		return
	}

	shortMessage, err := c.shortMessage.GetData()
	if err != nil {
		return
	}

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetMessageId())+1+dat.Len()+len(c.GetScheduleDeliveryTime())+1+len(c.GetValidityPeriod())+1+Utils.SZ_BYTE*3+shortMessage.Len()))

	buf.Write_CString(c.GetMessageId())
	buf.Write_Buffer(dat)
	buf.Write_CString(c.GetScheduleDeliveryTime())
	buf.Write_CString(c.GetValidityPeriod())
	buf.Write_UnsafeByte(c.GetRegisteredDelivery())
	buf.Write_UnsafeByte(c.GetSmDefaultMsgId())
	buf.Write_UnsafeByte(Common.EncodeUnsigned(c.GetSmLength()))

	err = buf.Write_Buffer(shortMessage)
	return
}

func (c *ReplaceSM) CreateResponse() (IResponse, error) {
	return NewReplaceSMResp(), nil
}

func (c *ReplaceSM) SetSourceAddr(value *Address) {
	c.sourceAddr = value
}

func (c *ReplaceSM) SetSourceAddrWithAddr(value string) *Exception.Exception {
	a, err := NewAddressWithAddr(value)
	if err != nil {
		return err
	}

	c.sourceAddr = a
	return nil
}

func (c *ReplaceSM) SetSourceAddrWithTonNpiAddr(ton, npi byte, value string) *Exception.Exception {
	a, err := NewAddressWithTonNpiAddr(ton, npi, value)
	if err != nil {
		return err
	}

	c.sourceAddr = a
	return nil
}

func (c *ReplaceSM) SetMessageId(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_MSGID_LEN))
	if err != nil {
		return err
	}

	c.messageId = value
	return nil
}

func (c *ReplaceSM) GetMessageId() string {
	return c.messageId
}

func (c *ReplaceSM) GetSourceAddr() *Address {
	return c.sourceAddr
}

func (c *ReplaceSM) SetScheduleDeliveryTime(value string) *Exception.Exception {
	err := c.CheckDate(value)
	if err != nil {
		return err
	}

	c.scheduleDeliveryTime = value
	return nil
}

func (c *ReplaceSM) GetScheduleDeliveryTime() string {
	return c.scheduleDeliveryTime
}

func (c *ReplaceSM) SetValidityPeriod(value string) *Exception.Exception {
	err := c.CheckDate(value)
	if err != nil {
		return err
	}

	c.validityPeriod = value
	return nil
}

func (c *ReplaceSM) GetValidityPeriod() string {
	return c.validityPeriod
}

func (c *ReplaceSM) SetShortMessage(value string) *Exception.Exception {
	err := c.shortMessage.SetMessage(value)
	if err != nil {
		return err
	}

	c.SetSmLength(int16(c.shortMessage.GetLength()))
	return nil
}

func (c *ReplaceSM) SetShortMessageWithEncoding(value string, enc Data.Encoding) *Exception.Exception {
	err := c.shortMessage.SetMessageWithEncoding(value, enc)
	if err != nil {
		return err
	}

	c.SetSmLength(int16(c.shortMessage.GetLength()))
	return nil
}

func (c *ReplaceSM) SetSmLength(value int16) {
	c.smLength = value
}

func (c *ReplaceSM) GetSmLength() int16 {
	return c.smLength
}

func (c *ReplaceSM) SetRegisteredDelivery(value byte) {
	c.registeredDelivery = value
}

func (c *ReplaceSM) SetSmDefaultMsgId(value byte) {
	c.smDefaultMsgId = value
}

func (c *ReplaceSM) GetRegisteredDelivery() byte {
	return c.registeredDelivery
}

func (c *ReplaceSM) GetSmDefaultMsgId() byte {
	return c.smDefaultMsgId
}
