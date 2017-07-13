package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/TLV"
	"github.com/linxGnu/gosmpp/Utils"
)

type DataSM struct {
	Request
	serviceType        string
	sourceAddr         *Address
	destAddr           *Address
	esmClass           byte
	registeredDelivery byte
	dataCoding         byte

	// optional params
	userMessageReference *TLV.TLVShort
	sourcePort           *TLV.TLVShort
	sourceAddrSubunit    *TLV.TLVByte
	sourceNetworkType    *TLV.TLVByte
	sourceBearerType     *TLV.TLVByte
	sourceTelematicsId   *TLV.TLVByte
	destinationPort      *TLV.TLVShort
	destAddrSubunit      *TLV.TLVByte
	destNetworkType      *TLV.TLVByte
	destBearerType       *TLV.TLVByte
	destTelematicsId     *TLV.TLVShort
	sarMsgRefNum         *TLV.TLVShort
	sarTotalSegments     *TLV.TLVUByte
	sarSegmentSeqnum     *TLV.TLVUByte
	moreMsgsToSend       *TLV.TLVByte
	qosTimeToLive        *TLV.TLVInt
	payloadType          *TLV.TLVByte
	messagePayload       *TLV.TLVOctets
	setDpf               *TLV.TLVByte
	receiptedMessageId   *TLV.TLVString
	messageState         *TLV.TLVByte
	networkErrorCode     *TLV.TLVOctets
	// exactly 3
	privacyIndicator *TLV.TLVByte
	callbackNum      *TLV.TLVOctets
	// 4-19
	callbackNumPresInd *TLV.TLVByte
	callbackNumAtag    *TLV.TLVOctets
	// 1-65
	sourceSubaddress *TLV.TLVOctets
	// 2-23
	destSubaddress      *TLV.TLVOctets
	userResponseCode    *TLV.TLVByte
	displayTime         *TLV.TLVByte
	smsSignal           *TLV.TLVShort
	msValidity          *TLV.TLVByte
	msMsgWaitFacilities *TLV.TLVByte
	numberOfMessages    *TLV.TLVByte
	alertOnMsgDelivery  *TLV.TLVEmpty
	languageIndicator   *TLV.TLVByte
	itsReplyType        *TLV.TLVByte
	itsSessionInfo      *TLV.TLVShort
}

func NewDataSM() *DataSM {
	a := &DataSM{}
	a.Construct()

	return a
}

func (a *DataSM) Construct() {
	defer a.SetRealReference(a)
	a.Request.Construct()

	a.SetCommandId(Data.DATA_SM)

	a.serviceType = Data.DFLT_SRVTYPE
	a.sourceAddr = NewAddressWithMaxLength(Data.SM_DATA_ADDR_LEN)
	a.destAddr = NewAddressWithMaxLength(Data.SM_DATA_ADDR_LEN)
	a.esmClass = Data.DFLT_ESM_CLASS
	a.registeredDelivery = Data.DFLT_REG_DELIVERY
	a.dataCoding = Data.DFLT_DATA_CODING
	a.userMessageReference = TLV.NewTLVShortWithTag(Data.OPT_PAR_USER_MSG_REF)
	a.sourcePort = TLV.NewTLVShortWithTag(Data.OPT_PAR_SRC_PORT)
	a.sourceAddrSubunit = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_ADDR_SUBUNIT)
	a.sourceNetworkType = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_NW_TYPE)
	a.sourceBearerType = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_BEAR_TYPE)
	a.sourceTelematicsId = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_TELE_ID)
	a.destinationPort = TLV.NewTLVShortWithTag(Data.OPT_PAR_DST_PORT)
	a.destAddrSubunit = TLV.NewTLVByteWithTag(Data.OPT_PAR_DST_ADDR_SUBUNIT)
	a.destNetworkType = TLV.NewTLVByteWithTag(Data.OPT_PAR_DST_NW_TYPE)
	a.destBearerType = TLV.NewTLVByteWithTag(Data.OPT_PAR_DST_BEAR_TYPE)
	a.destTelematicsId = TLV.NewTLVShortWithTag(Data.OPT_PAR_DST_TELE_ID)
	a.sarMsgRefNum = TLV.NewTLVShortWithTag(Data.OPT_PAR_SAR_MSG_REF_NUM)
	a.sarTotalSegments = TLV.NewTLVUByteWithTag(Data.OPT_PAR_SAR_TOT_SEG)
	a.sarSegmentSeqnum = TLV.NewTLVUByteWithTag(Data.OPT_PAR_SAR_SEG_SNUM)
	a.moreMsgsToSend = TLV.NewTLVByteWithTag(Data.OPT_PAR_MORE_MSGS)
	a.qosTimeToLive = TLV.NewTLVIntWithTag(Data.OPT_PAR_QOS_TIME_TO_LIVE)
	a.payloadType = TLV.NewTLVByteWithTag(Data.OPT_PAR_PAYLOAD_TYPE)
	a.messagePayload = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_MSG_PAYLOAD, int(Data.OPT_PAR_MSG_PAYLOAD_MIN), int(Data.OPT_PAR_MSG_PAYLOAD_MAX))
	a.setDpf = TLV.NewTLVByteWithTag(Data.OPT_PAR_SET_DPF)
	a.receiptedMessageId = TLV.NewTLVStringWithTagLength(Data.OPT_PAR_RECP_MSG_ID, int(Data.OPT_PAR_RECP_MSG_ID_MIN), int(Data.OPT_PAR_RECP_MSG_ID_MAX))
	a.messageState = TLV.NewTLVByteWithTag(Data.OPT_PAR_MSG_STATE)
	a.networkErrorCode = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_NW_ERR_CODE, int(Data.OPT_PAR_NW_ERR_CODE_MIN), int(Data.OPT_PAR_NW_ERR_CODE_MAX))
	a.privacyIndicator = TLV.NewTLVByteWithTag(Data.OPT_PAR_PRIV_IND)
	a.callbackNum = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_CALLBACK_NUM, int(Data.OPT_PAR_CALLBACK_NUM_MIN), int(Data.OPT_PAR_CALLBACK_NUM_MAX))
	a.callbackNumPresInd = TLV.NewTLVByteWithTag(Data.OPT_PAR_CALLBACK_NUM_PRES_IND)
	a.callbackNumAtag = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_CALLBACK_NUM_ATAG, int(Data.OPT_PAR_CALLBACK_NUM_ATAG_MIN), int(Data.OPT_PAR_CALLBACK_NUM_ATAG_MAX))
	a.sourceSubaddress = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_SRC_SUBADDR, int(Data.OPT_PAR_SRC_SUBADDR_MIN), int(Data.OPT_PAR_SRC_SUBADDR_MAX))
	a.destSubaddress = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_DEST_SUBADDR, int(Data.OPT_PAR_DEST_SUBADDR_MIN), int(Data.OPT_PAR_DEST_SUBADDR_MAX))
	a.userResponseCode = TLV.NewTLVByteWithTag(Data.OPT_PAR_USER_RESP_CODE)
	a.displayTime = TLV.NewTLVByteWithTag(Data.OPT_PAR_DISPLAY_TIME)
	a.smsSignal = TLV.NewTLVShortWithTag(Data.OPT_PAR_SMS_SIGNAL)
	a.msValidity = TLV.NewTLVByteWithTag(Data.OPT_PAR_MS_VALIDITY)
	a.msMsgWaitFacilities = TLV.NewTLVByteWithTag(Data.OPT_PAR_MSG_WAIT) // bit mask
	a.numberOfMessages = TLV.NewTLVByteWithTag(Data.OPT_PAR_NUM_MSGS)
	a.alertOnMsgDelivery = TLV.NewTLVEmptyWithTag(Data.OPT_PAR_ALERT_ON_MSG_DELIVERY)
	a.languageIndicator = TLV.NewTLVByteWithTag(Data.OPT_PAR_LANG_IND)
	a.itsReplyType = TLV.NewTLVByteWithTag(Data.OPT_PAR_ITS_REPLY_TYPE)
	a.itsSessionInfo = TLV.NewTLVShortWithTag(Data.OPT_PAR_ITS_SESSION_INFO)

	a.registerOptional(a.userMessageReference)
	a.registerOptional(a.sourcePort)
	a.registerOptional(a.sourceAddrSubunit)
	a.registerOptional(a.sourceNetworkType)
	a.registerOptional(a.sourceBearerType)
	a.registerOptional(a.sourceTelematicsId)
	a.registerOptional(a.destinationPort)
	a.registerOptional(a.destAddrSubunit)
	a.registerOptional(a.destNetworkType)
	a.registerOptional(a.destBearerType)
	a.registerOptional(a.destTelematicsId)
	a.registerOptional(a.sarMsgRefNum)
	a.registerOptional(a.sarTotalSegments)
	a.registerOptional(a.sarSegmentSeqnum)
	a.registerOptional(a.moreMsgsToSend)
	a.registerOptional(a.qosTimeToLive)
	a.registerOptional(a.payloadType)
	a.registerOptional(a.messagePayload)
	a.registerOptional(a.setDpf)
	a.registerOptional(a.receiptedMessageId)
	a.registerOptional(a.messageState)
	a.registerOptional(a.networkErrorCode)
	a.registerOptional(a.privacyIndicator)
	a.registerOptional(a.callbackNum)
	a.registerOptional(a.callbackNumPresInd)
	a.registerOptional(a.callbackNumAtag)
	a.registerOptional(a.sourceSubaddress)
	a.registerOptional(a.destSubaddress)
	a.registerOptional(a.userResponseCode)
	a.registerOptional(a.displayTime)
	a.registerOptional(a.smsSignal)
	a.registerOptional(a.msValidity)
	a.registerOptional(a.msMsgWaitFacilities)
	a.registerOptional(a.numberOfMessages)
	a.registerOptional(a.alertOnMsgDelivery)
	a.registerOptional(a.languageIndicator)
	a.registerOptional(a.itsReplyType)
	a.registerOptional(a.itsSessionInfo)
}

func (c *DataSM) GetInstance() (IPDU, error) {
	return NewDataSM(), nil
}

func (c *DataSM) CreateResponse() (IResponse, error) {
	return NewDataSMResp(), nil
}

func (c *DataSM) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("DataSM: set body buffer is nil")
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

	err = c.sourceAddr.SetData(buf)
	if err != nil {
		return
	}

	err = c.destAddr.SetData(buf)
	if err != nil {
		return
	}

	dat, err := buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetEsmClass(dat)

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetRegisteredDelivery(dat)

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetDataCoding(dat)

	return
}

func (c *DataSM) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
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

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetServiceType())+1+src.Len()+des.Len()+Utils.SZ_BYTE*3))

	buf.Write_CString(c.GetServiceType())
	buf.Write_Buffer(src)
	buf.Write_Buffer(des)
	buf.Write_UnsafeByte(c.GetEsmClass())
	buf.Write_UnsafeByte(c.GetRegisteredDelivery())
	buf.Write_UnsafeByte(c.GetDataCoding())

	return
}

func (c *DataSM) SetEsmClass(dat byte) {
	c.esmClass = dat
}

func (c *DataSM) GetEsmClass() byte {
	return c.esmClass
}

func (c *DataSM) SetRegisteredDelivery(dat byte) {
	c.registeredDelivery = dat
}

func (c *DataSM) GetRegisteredDelivery() byte {
	return c.registeredDelivery
}

func (c *DataSM) SetDataCoding(dat byte) {
	c.dataCoding = dat
}

func (c *DataSM) GetDataCoding() byte {
	return c.dataCoding
}

func (c *DataSM) SetServiceType(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_SRVTYPE_LEN))
	if err != nil {
		return err
	}

	c.serviceType = value
	return nil
}

func (c *DataSM) SetSourceAddr(value *Address) {
	c.sourceAddr = value
}

func (c *DataSM) SetSourceAddrFromStr(st string) *Exception.Exception {
	tmp, err := NewAddressWithAddr(st)
	if err != nil {
		return err
	}

	c.sourceAddr = tmp
	return nil
}

func (c *DataSM) SetSourceAddrFromStrTon(ton, npi byte, st string) *Exception.Exception {
	tmp, err := NewAddressWithTonNpiAddr(ton, npi, st)
	if err != nil {
		return err
	}

	c.sourceAddr = tmp
	return nil
}

func (c *DataSM) SetDestAddr(value *Address) {
	c.destAddr = value
}

func (c *DataSM) SetDestAddrFromStr(st string) *Exception.Exception {
	tmp, err := NewAddressWithAddr(st)
	if err != nil {
		return err
	}

	c.destAddr = tmp
	return nil
}

func (c *DataSM) SetDestAddrFromStrTon(ton, npi byte, st string) *Exception.Exception {
	tmp, err := NewAddressWithTonNpiAddr(ton, npi, st)
	if err != nil {
		return err
	}

	c.destAddr = tmp
	return nil
}

func (c *DataSM) GetServiceType() string {
	return c.serviceType
}

func (c *DataSM) GetSourceAddr() *Address {
	return c.sourceAddr
}

func (c *DataSM) GetDestAddr() *Address {
	return c.destAddr
}

func (c *DataSM) HasUserMessageReference() bool {
	return c.userMessageReference.HasValue()
}

func (c *DataSM) HasSourcePort() bool {
	return c.sourcePort.HasValue()
}

func (c *DataSM) HasSourceAddrSubunit() bool {
	return c.sourceAddrSubunit.HasValue()
}

func (c *DataSM) HasSourceNetworkType() bool {
	return c.sourceNetworkType.HasValue()
}

func (c *DataSM) HasSourceBearerType() bool {
	return c.sourceBearerType.HasValue()
}

func (c *DataSM) HasSourceTelematicsId() bool {
	return c.sourceTelematicsId.HasValue()
}

func (c *DataSM) HasDestinationPort() bool {
	return c.destinationPort.HasValue()
}

func (c *DataSM) HasDestAddrSubunit() bool {
	return c.destAddrSubunit.HasValue()
}

func (c *DataSM) HasDestNetworkType() bool {
	return c.destNetworkType.HasValue()
}

func (c *DataSM) HasDestBearerType() bool {
	return c.destBearerType.HasValue()
}

func (c *DataSM) HasDestTelematicsId() bool {
	return c.destTelematicsId.HasValue()
}

func (c *DataSM) HasSarMsgRefNum() bool {
	return c.sarMsgRefNum.HasValue()
}

func (c *DataSM) HasSarTotalSegments() bool {
	return c.sarTotalSegments.HasValue()
}

func (c *DataSM) HasSarSegmentSeqnum() bool {
	return c.sarSegmentSeqnum.HasValue()
}

func (c *DataSM) HasMoreMsgsToSend() bool {
	return c.moreMsgsToSend.HasValue()
}

func (c *DataSM) HasQosTimeToLive() bool {
	return c.qosTimeToLive.HasValue()
}

func (c *DataSM) HasPayloadType() bool {
	return c.payloadType.HasValue()
}

func (c *DataSM) HasMessagePayload() bool {
	return c.messagePayload.HasValue()
}

func (c *DataSM) HasSetDpf() bool {
	return c.setDpf.HasValue()
}

func (c *DataSM) HasReceiptedMessageId() bool {
	return c.receiptedMessageId.HasValue()
}

func (c *DataSM) HasMessageState() bool {
	return c.messageState.HasValue()
}

func (c *DataSM) HasNetworkErrorCode() bool {
	return c.networkErrorCode.HasValue()
}

func (c *DataSM) HasPrivacyIndicator() bool {
	return c.privacyIndicator.HasValue()
}

func (c *DataSM) HasCallbackNum() bool {
	return c.callbackNum.HasValue()
}

func (c *DataSM) HasCallbackNumPresInd() bool {
	return c.callbackNumPresInd.HasValue()
}

func (c *DataSM) HasCallbackNumAtag() bool {
	return c.callbackNumAtag.HasValue()
}

func (c *DataSM) HasSourceSubaddress() bool {
	return c.sourceSubaddress.HasValue()
}

func (c *DataSM) HasDestSubaddress() bool {
	return c.destSubaddress.HasValue()
}

func (c *DataSM) HasUserResponseCode() bool {
	return c.userResponseCode.HasValue()
}

func (c *DataSM) HasDisplayTime() bool {
	return c.displayTime.HasValue()
}

func (c *DataSM) HasSmsSignal() bool {
	return c.smsSignal.HasValue()
}

func (c *DataSM) HasMsValidity() bool {
	return c.msValidity.HasValue()
}

func (c *DataSM) HasMsMsgWaitFacilities() bool {
	return c.msMsgWaitFacilities.HasValue()
}

func (c *DataSM) HasNumberOfMessages() bool {
	return c.numberOfMessages.HasValue()
}

func (c *DataSM) HasAlertOnMsgDelivery() bool {
	return c.alertOnMsgDelivery.HasValue()
}

func (c *DataSM) HasLanguageIndicator() bool {
	return c.languageIndicator.HasValue()
}

func (c *DataSM) HasItsReplyType() bool {
	return c.itsReplyType.HasValue()
}

func (c *DataSM) HasItsSessionInfo() bool {
	return c.itsSessionInfo.HasValue()
}

func (c *DataSM) SetUserMessageReference(value int16) *Exception.Exception {
	return c.userMessageReference.SetValue(value)
}

func (c *DataSM) SetSourcePort(value int16) *Exception.Exception {
	return c.sourcePort.SetValue(value)
}

func (c *DataSM) SetSourceAddrSubunit(value byte) *Exception.Exception {
	return c.sourceAddrSubunit.SetValue(value)
}

func (c *DataSM) SetSourceNetworkType(value byte) *Exception.Exception {
	return c.sourceNetworkType.SetValue(value)
}

func (c *DataSM) SetSourceBearerType(value byte) *Exception.Exception {
	return c.sourceBearerType.SetValue(value)
}

func (c *DataSM) SetSourceTelematicsId(value byte) *Exception.Exception {
	return c.sourceTelematicsId.SetValue(value)
}

func (c *DataSM) SetDestinationPort(value int16) *Exception.Exception {
	return c.destinationPort.SetValue(value)
}

func (c *DataSM) SetDestAddrSubunit(value byte) *Exception.Exception {
	return c.destAddrSubunit.SetValue(value)
}

func (c *DataSM) SetDestNetworkType(value byte) *Exception.Exception {
	return c.destNetworkType.SetValue(value)
}

func (c *DataSM) SetDestBearerType(value byte) *Exception.Exception {
	return c.destBearerType.SetValue(value)
}

func (c *DataSM) SetDestTelematicsId(value int16) *Exception.Exception {
	return c.destTelematicsId.SetValue(value)
}

func (c *DataSM) SetSarMsgRefNum(value int16) *Exception.Exception {
	return c.sarMsgRefNum.SetValue(value)
}

func (c *DataSM) SetSarTotalSegments(value uint8) *Exception.Exception {
	return c.sarTotalSegments.SetValue(value)
}

func (c *DataSM) SetSarSegmentSeqnum(value uint8) *Exception.Exception {
	return c.sarSegmentSeqnum.SetValue(value)
}

func (c *DataSM) SetMoreMsgsToSend(value byte) *Exception.Exception {
	return c.moreMsgsToSend.SetValue(value)
}

func (c *DataSM) SetQosTimeToLive(value int32) *Exception.Exception {
	return c.qosTimeToLive.SetValue(value)
}

func (c *DataSM) SetPayloadType(value byte) *Exception.Exception {
	return c.payloadType.SetValue(value)
}

func (c *DataSM) SetMessagePayload(value *Utils.ByteBuffer) *Exception.Exception {
	return c.messagePayload.SetValue(value)
}

func (c *DataSM) SetSetDpf(value byte) *Exception.Exception {
	return c.setDpf.SetValue(value)
}

func (c *DataSM) SetReceiptedMessageId(value string) *Exception.Exception {
	return c.receiptedMessageId.SetValue(value)
}

func (c *DataSM) SetMessageState(value byte) *Exception.Exception {
	return c.messageState.SetValue(value)
}

func (c *DataSM) SetNetworkErrorCode(value *Utils.ByteBuffer) *Exception.Exception {
	return c.networkErrorCode.SetValue(value)
}

func (c *DataSM) SetPrivacyIndicator(value byte) *Exception.Exception {
	return c.privacyIndicator.SetValue(value)
}

func (c *DataSM) SetCallbackNum(value *Utils.ByteBuffer) *Exception.Exception {
	return c.callbackNum.SetValue(value)
}

func (c *DataSM) SetCallbackNumPresInd(value byte) *Exception.Exception {
	return c.callbackNumPresInd.SetValue(value)
}

func (c *DataSM) SetCallbackNumAtag(value *Utils.ByteBuffer) *Exception.Exception {
	return c.callbackNumAtag.SetValue(value)
}

func (c *DataSM) SetSourceSubaddress(value *Utils.ByteBuffer) *Exception.Exception {
	return c.sourceSubaddress.SetValue(value)
}

func (c *DataSM) SetDestSubaddress(value *Utils.ByteBuffer) *Exception.Exception {
	return c.destSubaddress.SetValue(value)
}

func (c *DataSM) SetUserResponseCode(value byte) *Exception.Exception {
	return c.userResponseCode.SetValue(value)
}

func (c *DataSM) SetDisplayTime(value byte) *Exception.Exception {
	return c.displayTime.SetValue(value)
}

func (c *DataSM) SetSmsSignal(value int16) *Exception.Exception {
	return c.smsSignal.SetValue(value)
}

func (c *DataSM) SetMsValidity(value byte) *Exception.Exception {
	return c.msValidity.SetValue(value)
}

func (c *DataSM) SetMsMsgWaitFacilities(value byte) *Exception.Exception {
	return c.msMsgWaitFacilities.SetValue(value)
}

func (c *DataSM) SetNumberOfMessages(value byte) *Exception.Exception {
	return c.numberOfMessages.SetValue(value)
}

func (c *DataSM) SetAlertOnMsgDelivery(value bool) *Exception.Exception {
	return c.alertOnMsgDelivery.SetValue(value)
}

func (c *DataSM) SetLanguageIndicator(value byte) *Exception.Exception {
	return c.languageIndicator.SetValue(value)
}

func (c *DataSM) SetItsReplyType(value byte) *Exception.Exception {
	return c.itsReplyType.SetValue(value)
}

func (c *DataSM) SetItsSessionInfo(value int16) *Exception.Exception {
	return c.itsSessionInfo.SetValue(value)
}

func (c *DataSM) GetUserMessageReference() (int16, *Exception.Exception) {
	return c.userMessageReference.GetValue()
}

func (c *DataSM) GetSourcePort() (int16, *Exception.Exception) {
	return c.sourcePort.GetValue()
}

func (c *DataSM) GetSourceAddrSubunit() (byte, *Exception.Exception) {
	return c.sourceAddrSubunit.GetValue()
}

func (c *DataSM) GetSourceNetworkType() (byte, *Exception.Exception) {
	return c.sourceNetworkType.GetValue()
}

func (c *DataSM) GetSourceBearerType() (byte, *Exception.Exception) {
	return c.sourceBearerType.GetValue()
}

func (c *DataSM) GetSourceTelematicsId() (byte, *Exception.Exception) {
	return c.sourceTelematicsId.GetValue()
}

func (c *DataSM) GetDestinationPort() (int16, *Exception.Exception) {
	return c.destinationPort.GetValue()
}

func (c *DataSM) GetDestAddrSubunit() (byte, *Exception.Exception) {
	return c.destAddrSubunit.GetValue()
}

func (c *DataSM) GetDestNetworkType() (byte, *Exception.Exception) {
	return c.destNetworkType.GetValue()
}

func (c *DataSM) GetDestBearerType() (byte, *Exception.Exception) {
	return c.destBearerType.GetValue()
}

func (c *DataSM) GetDestTelematicsId() (int16, *Exception.Exception) {
	return c.destTelematicsId.GetValue()
}

func (c *DataSM) GetSarMsgRefNum() (int16, *Exception.Exception) {
	return c.sarMsgRefNum.GetValue()
}

func (c *DataSM) GetSarTotalSegments() (byte, *Exception.Exception) {
	return c.sarTotalSegments.GetValue()
}

func (c *DataSM) GetSarSegmentSeqnum() (byte, *Exception.Exception) {
	return c.sarSegmentSeqnum.GetValue()
}

func (c *DataSM) GetMoreMsgsToSend() (byte, *Exception.Exception) {
	return c.moreMsgsToSend.GetValue()
}

func (c *DataSM) GetQosTimeToLive() (int32, *Exception.Exception) {
	return c.qosTimeToLive.GetValue()
}

func (c *DataSM) GetPayloadType() (byte, *Exception.Exception) {
	return c.payloadType.GetValue()
}

func (c *DataSM) GetMessagePayload() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.messagePayload.GetValue()
}

func (c *DataSM) GetSetDpf() (byte, *Exception.Exception) {
	return c.setDpf.GetValue()
}

func (c *DataSM) GetReceiptedMessageId() (string, *Exception.Exception) {
	return c.receiptedMessageId.GetValue()
}

func (c *DataSM) GetMessageState() (byte, *Exception.Exception) {
	return c.messageState.GetValue()
}

func (c *DataSM) GetNetworkErrorCode() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.networkErrorCode.GetValue()
}

func (c *DataSM) GetPrivacyIndicator() (byte, *Exception.Exception) {
	return c.privacyIndicator.GetValue()
}

func (c *DataSM) GetCallbackNum() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.callbackNum.GetValue()
}

func (c *DataSM) GetCallbackNumPresInd() (byte, *Exception.Exception) {
	return c.callbackNumPresInd.GetValue()
}

func (c *DataSM) GetCallbackNumAtag() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.callbackNumAtag.GetValue()
}

func (c *DataSM) GetSourceSubaddress() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.sourceSubaddress.GetValue()
}

func (c *DataSM) GetDestSubaddress() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.destSubaddress.GetValue()
}

func (c *DataSM) GetUserResponseCode() (byte, *Exception.Exception) {
	return c.userResponseCode.GetValue()
}

func (c *DataSM) GetDisplayTime() (byte, *Exception.Exception) {
	return c.displayTime.GetValue()
}

func (c *DataSM) GetSmsSignal() (int16, *Exception.Exception) {
	return c.smsSignal.GetValue()
}

func (c *DataSM) GetMsValidity() (byte, *Exception.Exception) {
	return c.msValidity.GetValue()
}

func (c *DataSM) GetMsMsgWaitFacilities() (byte, *Exception.Exception) {
	return c.msMsgWaitFacilities.GetValue()
}

func (c *DataSM) GetNumberOfMessages() (byte, *Exception.Exception) {
	return c.numberOfMessages.GetValue()
}

func (c *DataSM) GetAlertOnMsgDelivery() (bool, *Exception.Exception) {
	return c.alertOnMsgDelivery.GetValue()
}

func (c *DataSM) GetLanguageIndicator() (byte, *Exception.Exception) {
	return c.languageIndicator.GetValue()
}

func (c *DataSM) GetItsReplyType() (byte, *Exception.Exception) {
	return c.itsReplyType.GetValue()
}

func (c *DataSM) GetItsSessionInfo() (int16, *Exception.Exception) {
	return c.itsSessionInfo.GetValue()
}
