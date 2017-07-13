package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/TLV"
	"github.com/linxGnu/gosmpp/Utils"
)

type DataSMResp struct {
	Response
	messageId                string
	deliveryFailureReason    *TLV.TLVByte
	networkErrorCode         *TLV.TLVOctets
	additionalStatusInfoText *TLV.TLVString
	dpfResult                *TLV.TLVByte
}

func NewDataSMResp() *DataSMResp {
	a := &DataSMResp{}
	a.Construct()

	return a
}

func (a *DataSMResp) Construct() {
	defer a.SetRealReference(a)
	a.Response.Construct()

	a.SetCommandId(Data.DATA_SM_RESP)
	a.messageId = Data.DFLT_MSGID
	a.deliveryFailureReason = TLV.NewTLVByteWithTag(Data.OPT_PAR_DEL_FAIL_RSN)
	a.networkErrorCode = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_NW_ERR_CODE, int(Data.OPT_PAR_NW_ERR_CODE_MIN), int(Data.OPT_PAR_NW_ERR_CODE_MAX))
	a.additionalStatusInfoText = TLV.NewTLVStringWithTagLength(Data.OPT_PAR_ADD_STAT_INFO, int(Data.OPT_PAR_ADD_STAT_INFO_MIN), int(Data.OPT_PAR_ADD_STAT_INFO_MAX))
	a.dpfResult = TLV.NewTLVByteWithTag(Data.OPT_PAR_DPF_RES)
}

func (c *DataSMResp) GetInstance() (IPDU, error) {
	return NewDataSMResp(), nil
}

func (c *DataSMResp) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("DataSMResp: set body buffer is nil")
		return
	}

	val, err := buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetMessageId(val)

	return
}

func (c *DataSMResp) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	buf = Utils.NewBuffer(make([]byte, 0, 16))
	err = buf.Write_CString(c.GetMessageId())

	return
}

func (c *DataSMResp) SetMessageId(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_MSGID_LEN))
	if err != nil {
		return err
	}

	c.messageId = value
	return nil
}

func (c *DataSMResp) GetMessageId() string {
	return c.messageId
}

func (c *DataSMResp) GetDeliveryFailureReason() (byte, *Exception.Exception) {
	return c.deliveryFailureReason.GetValue()
}

func (c *DataSMResp) GetNetworkErrorCode() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.networkErrorCode.GetValue()
}

func (c *DataSMResp) GetAdditionalStatusInfoText() (string, *Exception.Exception) {
	return c.additionalStatusInfoText.GetValue()
}

func (c *DataSMResp) GetDpfResult() (byte, *Exception.Exception) {
	return c.dpfResult.GetValue()
}

func (c *DataSMResp) HasDeliveryFailureReason() bool {
	return c.deliveryFailureReason.HasValue()
}

func (c *DataSMResp) HasNetworkErrorCode() bool {
	return c.networkErrorCode.HasValue()
}

func (c *DataSMResp) HasAdditionalStatusInfoText() bool {
	return c.additionalStatusInfoText.HasValue()
}

func (c *DataSMResp) HasDpfResult() bool {
	return c.dpfResult.HasValue()
}

func (c *DataSMResp) SetDeliveryFailureReason(val byte) *Exception.Exception {
	return c.deliveryFailureReason.SetValue(val)
}

func (c *DataSMResp) SetNetworkErrorCode(val *Utils.ByteBuffer) *Exception.Exception {
	return c.networkErrorCode.SetValue(val)
}

func (c *DataSMResp) SetAdditionalStatusInfoText(val string) *Exception.Exception {
	return c.additionalStatusInfoText.SetValue(val)
}

func (c *DataSMResp) SetDpfResult(val byte) *Exception.Exception {
	return c.dpfResult.SetValue(val)
}
