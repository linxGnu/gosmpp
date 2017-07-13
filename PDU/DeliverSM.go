package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/PDU/TLV"
	"github.com/linxGnu/gosmpp/Utils"
)

type DeliverSM struct {
	Request
	serviceType          string
	sourceAddr           *Address
	destAddr             *Address
	esmClass             byte
	protocolId           byte
	priorityFlag         byte
	scheduleDeliveryTime string // not used
	validityPeriod       string // not used
	registeredDelivery   byte
	replaceIfPresentFlag byte // not used
	dataCoding           byte
	smDefaultMsgId       byte
	smLength             int16
	shortMessage         *ShortMessage

	// optional params
	userMessageReference *TLV.TLVShort
	sourcePort           *TLV.TLVShort
	// sourceAddrSubunit    *TLV.TLVByte
	// sourceNetworkType    *TLV.TLVByte
	// sourceBearerType     *TLV.TLVByte
	// sourceTelematicsId   *TLV.TLVByte
	destinationPort *TLV.TLVShort
	// destAddrSubunit      *TLV.TLVByte
	// destNetworkType      *TLV.TLVByte
	// destBearerType       *TLV.TLVByte
	// destTelematicsId     *TLV.TLVShort
	sarMsgRefNum     *TLV.TLVShort
	sarTotalSegments *TLV.TLVUByte
	sarSegmentSeqnum *TLV.TLVUByte
	// moreMsgsToSend     *TLV.TLVByte
	// qosTimeToLive      *TLV.TLVInt
	payloadType    *TLV.TLVByte
	messagePayload *TLV.TLVOctets
	// setDpf             *TLV.TLVByte
	receiptedMessageId *TLV.TLVString
	messageState       *TLV.TLVByte
	networkErrorCode   *TLV.TLVOctets
	// exactly 3
	privacyIndicator *TLV.TLVByte
	callbackNum      *TLV.TLVOctets
	// 4-19
	// callbackNumPresInd *TLV.TLVByte
	// callbackNumAtag    *TLV.TLVOctets
	// 1-65
	sourceSubaddress *TLV.TLVOctets
	// 2-23
	destSubaddress   *TLV.TLVOctets
	userResponseCode *TLV.TLVByte
	// displayTime         *TLV.TLVByte
	// smsSignal           *TLV.TLVShort
	// msValidity          *TLV.TLVByte
	// msMsgWaitFacilities *TLV.TLVByte
	// numberOfMessages   *TLV.TLVByte
	// alertOnMsgDelivery *TLV.TLVEmpty
	languageIndicator *TLV.TLVByte
	// itsReplyType       *TLV.TLVByte
	itsSessionInfo *TLV.TLVShort
}

func NewDeliverSM() *DeliverSM {
	a := &DeliverSM{}
	a.Construct()

	return a
}

func (a *DeliverSM) Construct() {
	defer a.SetRealReference(a)
	a.Request.Construct()

	a.SetCommandId(Data.DELIVER_SM)

	a.serviceType = Data.DFLT_SRVTYPE
	a.sourceAddr = NewAddressWithMaxLength(Data.SM_DATA_ADDR_LEN)
	a.destAddr = NewAddressWithMaxLength(Data.SM_DATA_ADDR_LEN)
	a.esmClass = Data.DFLT_ESM_CLASS
	a.protocolId = Data.DFLT_PROTOCOLID
	a.priorityFlag = Data.DFLT_PRIORITY_FLAG
	a.scheduleDeliveryTime = Data.DFLT_SCHEDULE
	a.validityPeriod = Data.DFLT_VALIDITY
	a.registeredDelivery = Data.DFLT_REG_DELIVERY
	a.replaceIfPresentFlag = Data.DFTL_REPLACE_IFP
	a.dataCoding = Data.DFLT_DATA_CODING
	a.smDefaultMsgId = Data.DFLT_DFLTMSGID
	a.smLength = int16(Data.DFLT_MSG_LEN)
	a.shortMessage = NewShortMessageWithMaxLength(Data.SM_MSG_LEN)

	a.userMessageReference = TLV.NewTLVShortWithTag(Data.OPT_PAR_USER_MSG_REF)
	a.sourcePort = TLV.NewTLVShortWithTag(Data.OPT_PAR_SRC_PORT)
	// a.sourceAddrSubunit = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_ADDR_SUBUNIT)
	// a.sourceNetworkType = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_NW_TYPE)
	// a.sourceBearerType = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_BEAR_TYPE)
	// a.sourceTelematicsId = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_TELE_ID)
	a.destinationPort = TLV.NewTLVShortWithTag(Data.OPT_PAR_DST_PORT)
	// a.destAddrSubunit = TLV.NewTLVByteWithTag(Data.OPT_PAR_DST_ADDR_SUBUNIT)
	// a.destNetworkType = TLV.NewTLVByteWithTag(Data.OPT_PAR_DST_NW_TYPE)
	// a.destBearerType = TLV.NewTLVByteWithTag(Data.OPT_PAR_DST_BEAR_TYPE)
	// a.destTelematicsId = TLV.NewTLVShortWithTag(Data.OPT_PAR_DST_TELE_ID)
	a.sarMsgRefNum = TLV.NewTLVShortWithTag(Data.OPT_PAR_SAR_MSG_REF_NUM)
	a.sarTotalSegments = TLV.NewTLVUByteWithTag(Data.OPT_PAR_SAR_TOT_SEG)
	a.sarSegmentSeqnum = TLV.NewTLVUByteWithTag(Data.OPT_PAR_SAR_SEG_SNUM)
	// a.moreMsgsToSend = TLV.NewTLVByteWithTag(Data.OPT_PAR_MORE_MSGS)
	// a.qosTimeToLive = TLV.NewTLVIntWithTag(Data.OPT_PAR_QOS_TIME_TO_LIVE)
	a.payloadType = TLV.NewTLVByteWithTag(Data.OPT_PAR_PAYLOAD_TYPE)
	a.messagePayload = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_MSG_PAYLOAD, int(Data.OPT_PAR_MSG_PAYLOAD_MIN), int(Data.OPT_PAR_MSG_PAYLOAD_MAX))
	// a.setDpf = TLV.NewTLVByteWithTag(Data.OPT_PAR_SET_DPF)
	a.receiptedMessageId = TLV.NewTLVStringWithTagLength(Data.OPT_PAR_RECP_MSG_ID, int(Data.OPT_PAR_RECP_MSG_ID_MIN), int(Data.OPT_PAR_RECP_MSG_ID_MAX))
	a.messageState = TLV.NewTLVByteWithTag(Data.OPT_PAR_MSG_STATE)
	a.networkErrorCode = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_NW_ERR_CODE, int(Data.OPT_PAR_NW_ERR_CODE_MIN), int(Data.OPT_PAR_NW_ERR_CODE_MAX))
	a.privacyIndicator = TLV.NewTLVByteWithTag(Data.OPT_PAR_PRIV_IND)
	a.callbackNum = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_CALLBACK_NUM, int(Data.OPT_PAR_CALLBACK_NUM_MIN), int(Data.OPT_PAR_CALLBACK_NUM_MAX))
	// a.callbackNumPresInd = TLV.NewTLVByteWithTag(Data.OPT_PAR_CALLBACK_NUM_PRES_IND)
	// a.callbackNumAtag = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_CALLBACK_NUM_ATAG, int(Data.OPT_PAR_CALLBACK_NUM_ATAG_MIN), int(Data.OPT_PAR_CALLBACK_NUM_ATAG_MAX))
	a.sourceSubaddress = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_SRC_SUBADDR, int(Data.OPT_PAR_SRC_SUBADDR_MIN), int(Data.OPT_PAR_SRC_SUBADDR_MAX))
	a.destSubaddress = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_DEST_SUBADDR, int(Data.OPT_PAR_DEST_SUBADDR_MIN), int(Data.OPT_PAR_DEST_SUBADDR_MAX))
	a.userResponseCode = TLV.NewTLVByteWithTag(Data.OPT_PAR_USER_RESP_CODE)
	// a.displayTime = TLV.NewTLVByteWithTag(Data.OPT_PAR_DISPLAY_TIME)
	// a.smsSignal = TLV.NewTLVShortWithTag(Data.OPT_PAR_SMS_SIGNAL)
	// a.msValidity = TLV.NewTLVByteWithTag(Data.OPT_PAR_MS_VALIDITY)
	// a.msMsgWaitFacilities = TLV.NewTLVByteWithTag(Data.OPT_PAR_MSG_WAIT) // bit mask
	// a.numberOfMessages = TLV.NewTLVByteWithTag(Data.OPT_PAR_NUM_MSGS)
	// a.alertOnMsgDelivery = TLV.NewTLVEmptyWithTag(Data.OPT_PAR_ALERT_ON_MSG_DELIVERY)
	a.languageIndicator = TLV.NewTLVByteWithTag(Data.OPT_PAR_LANG_IND)
	// a.itsReplyType = TLV.NewTLVByteWithTag(Data.OPT_PAR_ITS_REPLY_TYPE)
	a.itsSessionInfo = TLV.NewTLVShortWithTag(Data.OPT_PAR_ITS_SESSION_INFO)

	a.registerOptional(a.userMessageReference)
	a.registerOptional(a.sourcePort)
	// a.registerOptional(a.sourceAddrSubunit)
	// a.registerOptional(a.sourceNetworkType)
	// a.registerOptional(a.sourceBearerType)
	// a.registerOptional(a.sourceTelematicsId)
	a.registerOptional(a.destinationPort)
	// a.registerOptional(a.destAddrSubunit)
	// a.registerOptional(a.destNetworkType)
	// a.registerOptional(a.destBearerType)
	// a.registerOptional(a.destTelematicsId)
	a.registerOptional(a.sarMsgRefNum)
	a.registerOptional(a.sarTotalSegments)
	a.registerOptional(a.sarSegmentSeqnum)
	// a.registerOptional(a.moreMsgsToSend)
	// a.registerOptional(a.qosTimeToLive)
	a.registerOptional(a.payloadType)
	a.registerOptional(a.messagePayload)
	// a.registerOptional(a.setDpf)
	a.registerOptional(a.receiptedMessageId)
	a.registerOptional(a.messageState)
	a.registerOptional(a.networkErrorCode)
	a.registerOptional(a.privacyIndicator)
	a.registerOptional(a.callbackNum)
	// a.registerOptional(a.callbackNumPresInd)
	// a.registerOptional(a.callbackNumAtag)
	a.registerOptional(a.sourceSubaddress)
	a.registerOptional(a.destSubaddress)
	a.registerOptional(a.userResponseCode)
	// a.registerOptional(a.displayTime)
	// a.registerOptional(a.smsSignal)
	// a.registerOptional(a.msValidity)
	// a.registerOptional(a.msMsgWaitFacilities)
	// a.registerOptional(a.numberOfMessages)
	// a.registerOptional(a.alertOnMsgDelivery)
	a.registerOptional(a.languageIndicator)
	// a.registerOptional(a.itsReplyType)
	a.registerOptional(a.itsSessionInfo)
}

func (c *DeliverSM) GetInstance() (IPDU, error) {
	return NewDeliverSM(), nil
}

func (c *DeliverSM) CreateResponse() (IResponse, error) {
	return NewDeliverSMResp(), nil
}

func (c *DeliverSM) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("DeliverSM: set body buffer is nil")
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
	c.SetProtocolId(dat)

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetPriorityFlag(dat)

	_, err = buf.Read_CString() // default scheduleDeliveryTime
	if err != nil {
		return
	}

	_, err = buf.Read_CString() // default validityPeriod
	if err != nil {
		return
	}

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetRegisteredDelivery(dat)

	_, err = buf.Read_Byte() // dummy byte: replaceIfPresentFlag
	if err != nil {
		return
	}

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetDataCoding(dat)

	_, err = buf.Read_Byte() // dummy byte: smDefaultMsgId
	if err != nil {
		return
	}

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetSmLength(Common.DecodeUnsigned(dat))

	buf2, err := buf.Read_Bytes(int(c.GetSmLength()))
	if err != nil {
		return
	}
	err = c.shortMessage.SetData(buf2)

	return
}

func (c *DeliverSM) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
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

	shortMessage, err := c.shortMessage.GetData()
	if err != nil {
		return
	}

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetServiceType())+1+src.Len()+des.Len()+Utils.SZ_BYTE*3+len(c.GetScheduleDeliveryTime())+1+len(c.GetValidityPeriod())+1+Utils.SZ_BYTE*5+shortMessage.Len()))

	buf.Write_CString(c.GetServiceType())
	buf.Write_Buffer(src)
	buf.Write_Buffer(des)
	buf.Write_UnsafeByte(c.GetEsmClass())
	buf.Write_UnsafeByte(c.GetProtocolId())
	buf.Write_UnsafeByte(c.GetPriorityFlag())
	buf.Write_CString(c.GetScheduleDeliveryTime())
	buf.Write_CString(c.GetValidityPeriod())
	buf.Write_Byte(c.GetRegisteredDelivery())
	buf.Write_Byte(c.GetReplaceIfPresentFlag())
	buf.Write_Byte(c.GetDataCoding())
	buf.Write_Byte(c.GetSmDefaultMsgId())
	buf.Write_Byte(Common.EncodeUnsigned(c.GetSmLength()))
	buf.Write_Buffer(shortMessage)

	return
}

func (c *DeliverSM) SetEsmClass(dat byte) {
	c.esmClass = dat
}

func (c *DeliverSM) GetEsmClass() byte {
	return c.esmClass
}

func (c *DeliverSM) SetRegisteredDelivery(dat byte) {
	c.registeredDelivery = dat
}

func (c *DeliverSM) GetRegisteredDelivery() byte {
	return c.registeredDelivery
}

func (c *DeliverSM) GetReplaceIfPresentFlag() byte {
	return c.replaceIfPresentFlag
}

func (c *DeliverSM) SetReplaceIfPresentFlag(dat byte) {
	c.replaceIfPresentFlag = dat
}

func (c *DeliverSM) GetValidityPeriod() string {
	return c.validityPeriod
}

func (c *DeliverSM) GetSmDefaultMsgId() byte {
	return c.smDefaultMsgId
}

func (c *DeliverSM) SetDataCoding(dat byte) {
	c.dataCoding = dat
}

func (c *DeliverSM) GetDataCoding() byte {
	return c.dataCoding
}

func (c *DeliverSM) SetProtocolId(dat byte) {
	c.protocolId = dat
}

func (c *DeliverSM) GetProtocolId() byte {
	return c.protocolId
}

func (c *DeliverSM) SetPriorityFlag(dat byte) {
	c.priorityFlag = dat
}

func (c *DeliverSM) GetPriorityFlag() byte {
	return c.priorityFlag
}

func (c *DeliverSM) GetScheduleDeliveryTime() string {
	return c.scheduleDeliveryTime
}

func (c *DeliverSM) SetSmLength(value int16) {
	c.smLength = value
}

func (c *DeliverSM) GetSmLength() int16 {
	return c.smLength
}

func (c *DeliverSM) SetShortMessage(value string) *Exception.Exception {
	err := c.shortMessage.SetMessage(value)
	if err != nil {
		return err
	}

	c.SetSmLength(int16(c.shortMessage.GetLength()))
	return nil
}

// GetShortMessage get short message
func (c *DeliverSM) GetShortMessage() (string, *Exception.Exception) {
	return c.shortMessage.GetMessage()
}

// GetShortMessageWithEncoding get short message with encoding
func (c *DeliverSM) GetShortMessageWithEncoding(enc Data.Encoding) (string, *Exception.Exception) {
	return c.shortMessage.GetMessageWithEncoding(enc)
}

func (c *DeliverSM) SetServiceType(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_SRVTYPE_LEN))
	if err != nil {
		return err
	}

	c.serviceType = value
	return nil
}

func (c *DeliverSM) SetSourceAddr(value *Address) {
	c.sourceAddr = value
}

func (c *DeliverSM) SetSourceAddrFromStr(st string) *Exception.Exception {
	tmp, err := NewAddressWithAddr(st)
	if err != nil {
		return err
	}

	c.sourceAddr = tmp
	return nil
}

func (c *DeliverSM) SetSourceAddrFromStrTon(ton, npi byte, st string) *Exception.Exception {
	tmp, err := NewAddressWithTonNpiAddr(ton, npi, st)
	if err != nil {
		return err
	}

	c.sourceAddr = tmp
	return nil
}

func (c *DeliverSM) SetDestAddr(value *Address) {
	c.destAddr = value
}

func (c *DeliverSM) SetDestAddrFromStr(st string) *Exception.Exception {
	tmp, err := NewAddressWithAddr(st)
	if err != nil {
		return err
	}

	c.destAddr = tmp
	return nil
}

func (c *DeliverSM) SetDestAddrFromStrTon(ton, npi byte, st string) *Exception.Exception {
	tmp, err := NewAddressWithTonNpiAddr(ton, npi, st)
	if err != nil {
		return err
	}

	c.destAddr = tmp
	return nil
}

func (c *DeliverSM) GetServiceType() string {
	return c.serviceType
}

func (c *DeliverSM) GetSourceAddr() *Address {
	return c.sourceAddr
}

func (c *DeliverSM) GetDestAddr() *Address {
	return c.destAddr
}

func (c *DeliverSM) HasUserMessageReference() bool {
	return c.userMessageReference.HasValue()
}

func (c *DeliverSM) HasSourcePort() bool {
	return c.sourcePort.HasValue()
}

func (c *DeliverSM) HasDestinationPort() bool {
	return c.destinationPort.HasValue()
}

func (c *DeliverSM) HasSarMsgRefNum() bool {
	return c.sarMsgRefNum.HasValue()
}

func (c *DeliverSM) HasSarTotalSegments() bool {
	return c.sarTotalSegments.HasValue()
}

func (c *DeliverSM) HasSarSegmentSeqnum() bool {
	return c.sarSegmentSeqnum.HasValue()
}

func (c *DeliverSM) HasPayloadType() bool {
	return c.payloadType.HasValue()
}

func (c *DeliverSM) HasMessagePayload() bool {
	return c.messagePayload.HasValue()
}

func (c *DeliverSM) HasReceiptedMessageId() bool {
	return c.receiptedMessageId.HasValue()
}

func (c *DeliverSM) HasMessageState() bool {
	return c.messageState.HasValue()
}

func (c *DeliverSM) HasNetworkErrorCode() bool {
	return c.networkErrorCode.HasValue()
}

func (c *DeliverSM) HasPrivacyIndicator() bool {
	return c.privacyIndicator.HasValue()
}

func (c *DeliverSM) HasCallbackNum() bool {
	return c.callbackNum.HasValue()
}

func (c *DeliverSM) HasSourceSubaddress() bool {
	return c.sourceSubaddress.HasValue()
}

func (c *DeliverSM) HasDestSubaddress() bool {
	return c.destSubaddress.HasValue()
}

func (c *DeliverSM) HasUserResponseCode() bool {
	return c.userResponseCode.HasValue()
}

func (c *DeliverSM) HasLanguageIndicator() bool {
	return c.languageIndicator.HasValue()
}

func (c *DeliverSM) HasItsSessionInfo() bool {
	return c.itsSessionInfo.HasValue()
}

func (c *DeliverSM) SetUserMessageReference(value int16) *Exception.Exception {
	return c.userMessageReference.SetValue(value)
}

func (c *DeliverSM) SetSourcePort(value int16) *Exception.Exception {
	return c.sourcePort.SetValue(value)
}

func (c *DeliverSM) SetDestinationPort(value int16) *Exception.Exception {
	return c.destinationPort.SetValue(value)
}

func (c *DeliverSM) SetSarMsgRefNum(value int16) *Exception.Exception {
	return c.sarMsgRefNum.SetValue(value)
}

func (c *DeliverSM) SetSarTotalSegments(value uint8) *Exception.Exception {
	return c.sarTotalSegments.SetValue(value)
}

func (c *DeliverSM) SetSarSegmentSeqnum(value uint8) *Exception.Exception {
	return c.sarSegmentSeqnum.SetValue(value)
}

func (c *DeliverSM) SetPayloadType(value byte) *Exception.Exception {
	return c.payloadType.SetValue(value)
}

func (c *DeliverSM) SetMessagePayload(value *Utils.ByteBuffer) *Exception.Exception {
	return c.messagePayload.SetValue(value)
}

func (c *DeliverSM) SetReceiptedMessageId(value string) *Exception.Exception {
	return c.receiptedMessageId.SetValue(value)
}

func (c *DeliverSM) SetMessageState(value byte) *Exception.Exception {
	return c.messageState.SetValue(value)
}

func (c *DeliverSM) SetNetworkErrorCode(value *Utils.ByteBuffer) *Exception.Exception {
	return c.networkErrorCode.SetValue(value)
}

func (c *DeliverSM) SetPrivacyIndicator(value byte) *Exception.Exception {
	return c.privacyIndicator.SetValue(value)
}

func (c *DeliverSM) SetCallbackNum(value *Utils.ByteBuffer) *Exception.Exception {
	return c.callbackNum.SetValue(value)
}

func (c *DeliverSM) SetSourceSubaddress(value *Utils.ByteBuffer) *Exception.Exception {
	return c.sourceSubaddress.SetValue(value)
}

func (c *DeliverSM) SetDestSubaddress(value *Utils.ByteBuffer) *Exception.Exception {
	return c.destSubaddress.SetValue(value)
}

func (c *DeliverSM) SetUserResponseCode(value byte) *Exception.Exception {
	return c.userResponseCode.SetValue(value)
}

func (c *DeliverSM) SetLanguageIndicator(value byte) *Exception.Exception {
	return c.languageIndicator.SetValue(value)
}

func (c *DeliverSM) SetItsSessionInfo(value int16) *Exception.Exception {
	return c.itsSessionInfo.SetValue(value)
}

func (c *DeliverSM) GetUserMessageReference() (int16, *Exception.Exception) {
	return c.userMessageReference.GetValue()
}

func (c *DeliverSM) GetSourcePort() (int16, *Exception.Exception) {
	return c.sourcePort.GetValue()
}

func (c *DeliverSM) GetDestinationPort() (int16, *Exception.Exception) {
	return c.destinationPort.GetValue()
}

func (c *DeliverSM) GetSarMsgRefNum() (int16, *Exception.Exception) {
	return c.sarMsgRefNum.GetValue()
}

func (c *DeliverSM) GetSarTotalSegments() (byte, *Exception.Exception) {
	return c.sarTotalSegments.GetValue()
}

func (c *DeliverSM) GetSarSegmentSeqnum() (byte, *Exception.Exception) {
	return c.sarSegmentSeqnum.GetValue()
}

func (c *DeliverSM) GetPayloadType() (byte, *Exception.Exception) {
	return c.payloadType.GetValue()
}

func (c *DeliverSM) GetMessagePayload() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.messagePayload.GetValue()
}

func (c *DeliverSM) GetReceiptedMessageId() (string, *Exception.Exception) {
	return c.receiptedMessageId.GetValue()
}

func (c *DeliverSM) GetMessageState() (byte, *Exception.Exception) {
	return c.messageState.GetValue()
}

func (c *DeliverSM) GetNetworkErrorCode() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.networkErrorCode.GetValue()
}

func (c *DeliverSM) GetPrivacyIndicator() (byte, *Exception.Exception) {
	return c.privacyIndicator.GetValue()
}

func (c *DeliverSM) GetCallbackNum() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.callbackNum.GetValue()
}

func (c *DeliverSM) GetSourceSubaddress() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.sourceSubaddress.GetValue()
}

func (c *DeliverSM) GetDestSubaddress() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.destSubaddress.GetValue()
}

func (c *DeliverSM) GetUserResponseCode() (byte, *Exception.Exception) {
	return c.userResponseCode.GetValue()
}

func (c *DeliverSM) GetLanguageIndicator() (byte, *Exception.Exception) {
	return c.languageIndicator.GetValue()
}

func (c *DeliverSM) GetItsSessionInfo() (int16, *Exception.Exception) {
	return c.itsSessionInfo.GetValue()
}
