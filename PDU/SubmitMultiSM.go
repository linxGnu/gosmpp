package PDU

import (
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/PDU/TLV"
	"github.com/linxGnu/gosmpp/Utils"
)

type SubmitMultiSM struct {
	Request
	serviceType string
	sourceAddr  *Address
	destAddr    *Address
	esmClass    byte

	protocolId           byte
	priorityFlag         byte
	scheduleDeliveryTime string
	validityPeriod       string
	replaceIfPresentFlag byte
	smDefaultMsgId       byte
	smLength             int16
	shortMessage         *ShortMessage

	registeredDelivery byte
	dataCoding         byte

	// optional params
	userMessageReference *TLV.TLVShort
	sourcePort           *TLV.TLVShort
	sourceAddrSubunit    *TLV.TLVByte
	destinationPort      *TLV.TLVShort
	destAddrSubunit      *TLV.TLVByte
	sarMsgRefNum         *TLV.TLVShort
	sarTotalSegments     *TLV.TLVUByte
	sarSegmentSeqnum     *TLV.TLVUByte
	payloadType          *TLV.TLVByte
	messagePayload       *TLV.TLVOctets
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
	displayTime         *TLV.TLVByte
	smsSignal           *TLV.TLVShort
	msValidity          *TLV.TLVByte
	msMsgWaitFacilities *TLV.TLVByte
	alertOnMsgDelivery  *TLV.TLVEmpty
	languageIndicator   *TLV.TLVByte
}

func NewSubmitMultiSM() *SubmitMultiSM {
	a := &SubmitMultiSM{}
	a.Construct()

	return a
}

func (a *SubmitMultiSM) Construct() {
	defer a.SetRealReference(a)
	a.Request.Construct()

	a.SetCommandId(Data.SUBMIT_MULTI)

	a.scheduleDeliveryTime = Data.DFLT_SCHEDULE
	a.validityPeriod = Data.DFLT_VALIDITY
	a.priorityFlag = Data.DFLT_PRIORITY_FLAG
	a.smDefaultMsgId = Data.DFLT_DFLTMSGID
	a.smLength = int16(Data.DFLT_MSG_LEN)
	a.protocolId = Data.DFLT_PROTOCOLID
	a.replaceIfPresentFlag = Data.DFTL_REPLACE_IFP
	a.shortMessage = NewShortMessageWithMaxLength(int32(Data.SM_MSG_LEN))

	a.serviceType = Data.DFLT_SRVTYPE
	a.sourceAddr = NewAddressWithMaxLength(Data.SM_DATA_ADDR_LEN)
	a.destAddr = NewAddressWithMaxLength(Data.SM_DATA_ADDR_LEN)
	a.esmClass = Data.DFLT_ESM_CLASS
	a.registeredDelivery = Data.DFLT_REG_DELIVERY
	a.dataCoding = Data.DFLT_DATA_CODING
	a.userMessageReference = TLV.NewTLVShortWithTag(Data.OPT_PAR_USER_MSG_REF)
	a.sourcePort = TLV.NewTLVShortWithTag(Data.OPT_PAR_SRC_PORT)
	a.sourceAddrSubunit = TLV.NewTLVByteWithTag(Data.OPT_PAR_SRC_ADDR_SUBUNIT)
	a.destinationPort = TLV.NewTLVShortWithTag(Data.OPT_PAR_DST_PORT)
	a.destAddrSubunit = TLV.NewTLVByteWithTag(Data.OPT_PAR_DST_ADDR_SUBUNIT)
	a.sarMsgRefNum = TLV.NewTLVShortWithTag(Data.OPT_PAR_SAR_MSG_REF_NUM)
	a.sarTotalSegments = TLV.NewTLVUByteWithTag(Data.OPT_PAR_SAR_TOT_SEG)
	a.sarSegmentSeqnum = TLV.NewTLVUByteWithTag(Data.OPT_PAR_SAR_SEG_SNUM)
	a.payloadType = TLV.NewTLVByteWithTag(Data.OPT_PAR_PAYLOAD_TYPE)
	a.messagePayload = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_MSG_PAYLOAD, int(Data.OPT_PAR_MSG_PAYLOAD_MIN), int(Data.OPT_PAR_MSG_PAYLOAD_MAX))
	a.privacyIndicator = TLV.NewTLVByteWithTag(Data.OPT_PAR_PRIV_IND)
	a.callbackNum = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_CALLBACK_NUM, int(Data.OPT_PAR_CALLBACK_NUM_MIN), int(Data.OPT_PAR_CALLBACK_NUM_MAX))
	a.callbackNumPresInd = TLV.NewTLVByteWithTag(Data.OPT_PAR_CALLBACK_NUM_PRES_IND)
	a.callbackNumAtag = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_CALLBACK_NUM_ATAG, int(Data.OPT_PAR_CALLBACK_NUM_ATAG_MIN), int(Data.OPT_PAR_CALLBACK_NUM_ATAG_MAX))
	a.sourceSubaddress = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_SRC_SUBADDR, int(Data.OPT_PAR_SRC_SUBADDR_MIN), int(Data.OPT_PAR_SRC_SUBADDR_MAX))
	a.destSubaddress = TLV.NewTLVOctetsWithTagLength(Data.OPT_PAR_DEST_SUBADDR, int(Data.OPT_PAR_DEST_SUBADDR_MIN), int(Data.OPT_PAR_DEST_SUBADDR_MAX))
	a.displayTime = TLV.NewTLVByteWithTag(Data.OPT_PAR_DISPLAY_TIME)
	a.smsSignal = TLV.NewTLVShortWithTag(Data.OPT_PAR_SMS_SIGNAL)
	a.msValidity = TLV.NewTLVByteWithTag(Data.OPT_PAR_MS_VALIDITY)
	a.msMsgWaitFacilities = TLV.NewTLVByteWithTag(Data.OPT_PAR_MSG_WAIT) // bit mask
	a.alertOnMsgDelivery = TLV.NewTLVEmptyWithTag(Data.OPT_PAR_ALERT_ON_MSG_DELIVERY)
	a.languageIndicator = TLV.NewTLVByteWithTag(Data.OPT_PAR_LANG_IND)

	a.registerOptional(a.userMessageReference)
	a.registerOptional(a.sourcePort)
	a.registerOptional(a.sourceAddrSubunit)
	a.registerOptional(a.destinationPort)
	a.registerOptional(a.destAddrSubunit)
	a.registerOptional(a.sarMsgRefNum)
	a.registerOptional(a.sarTotalSegments)
	a.registerOptional(a.sarSegmentSeqnum)
	a.registerOptional(a.payloadType)
	a.registerOptional(a.messagePayload)
	a.registerOptional(a.privacyIndicator)
	a.registerOptional(a.callbackNum)
	a.registerOptional(a.callbackNumPresInd)
	a.registerOptional(a.callbackNumAtag)
	a.registerOptional(a.sourceSubaddress)
	a.registerOptional(a.destSubaddress)
	a.registerOptional(a.displayTime)
	a.registerOptional(a.smsSignal)
	a.registerOptional(a.msValidity)
	a.registerOptional(a.msMsgWaitFacilities)
	a.registerOptional(a.alertOnMsgDelivery)
	a.registerOptional(a.languageIndicator)
}

func (c *SubmitMultiSM) GetInstance() (IPDU, error) {
	return NewSubmitMultiSM(), nil
}

func (c *SubmitMultiSM) CreateResponse() (IResponse, error) {
	return NewSubmitMultiSMResp(), nil
}

func (c *SubmitMultiSM) SetBody(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("SubmitMultiSM: set body buffer is nil")
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

	st, err := buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetScheduleDeliveryTime(st)
	if err != nil {
		return
	}

	st, err = buf.Read_CString()
	if err != nil {
		return
	}
	err = c.SetValidityPeriod(st)
	if err != nil {
		return
	}

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetRegisteredDelivery(dat)

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetReplaceIfPresentFlag(dat)

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetDataCoding(dat)

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetSmDefaultMsgId(dat)

	dat, err = buf.Read_Byte()
	if err != nil {
		return
	}
	c.SetSmLength(Common.DecodeUnsigned(dat))

	tmpBuf, err := buf.Read_Bytes(int(c.GetSmLength()))
	if err != nil {
		return
	}

	err = c.shortMessage.SetData(tmpBuf)
	return
}

func (c *SubmitMultiSM) GetBody() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			buf = nil
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

	buf = Utils.NewBuffer(make([]byte, 0, len(c.GetServiceType())+1+src.Len()+des.Len()+3*Utils.SZ_BYTE+len(c.GetScheduleDeliveryTime())+1+len(c.GetValidityPeriod())+1+5*Utils.SZ_BYTE+shortMessage.Len()))

	buf.Write_CString(c.GetServiceType())
	buf.Write_Buffer(src)
	buf.Write_Buffer(des)
	buf.Write_UnsafeByte(c.GetEsmClass())
	buf.Write_UnsafeByte(c.GetProtocolId())
	buf.Write_UnsafeByte(c.GetPriorityFlag())
	buf.Write_CString(c.GetScheduleDeliveryTime())
	buf.Write_CString(c.GetValidityPeriod())
	buf.Write_UnsafeByte(c.GetRegisteredDelivery())
	buf.Write_UnsafeByte(c.GetReplaceIfPresentFlag())
	buf.Write_UnsafeByte(c.GetDataCoding())
	buf.Write_UnsafeByte(c.GetSmDefaultMsgId())
	buf.Write_UnsafeByte(Common.EncodeUnsigned(int16(c.GetSmLength())))

	err = buf.Write_Buffer(shortMessage)
	return
}

func (c *SubmitMultiSM) SetEsmClass(dat byte) {
	c.esmClass = dat
}

func (c *SubmitMultiSM) GetEsmClass() byte {
	return c.esmClass
}

func (c *SubmitMultiSM) SetRegisteredDelivery(dat byte) {
	c.registeredDelivery = dat
}

func (c *SubmitMultiSM) GetRegisteredDelivery() byte {
	return c.registeredDelivery
}

func (c *SubmitMultiSM) SetDataCoding(dat byte) {
	c.dataCoding = dat
}

func (c *SubmitMultiSM) GetDataCoding() byte {
	return c.dataCoding
}

func (c *SubmitMultiSM) SetSmLength(value int16) {
	c.smLength = value
}

func (c *SubmitMultiSM) GetSmLength() int16 {
	return c.smLength
}

func (c *SubmitMultiSM) SetSmDefaultMsgId(value byte) {
	c.smDefaultMsgId = value
}

func (c *SubmitMultiSM) GetSmDefaultMsgId() byte {
	return c.smDefaultMsgId
}

func (c *SubmitMultiSM) SetScheduleDeliveryTime(value string) *Exception.Exception {
	err := c.CheckDate(value)
	if err != nil {
		return err
	}

	c.scheduleDeliveryTime = value
	return nil
}

func (c *SubmitMultiSM) GetScheduleDeliveryTime() string {
	return c.scheduleDeliveryTime
}

func (c *SubmitMultiSM) SetValidityPeriod(value string) *Exception.Exception {
	err := c.CheckDate(value)
	if err != nil {
		return err
	}

	c.validityPeriod = value
	return nil
}

func (c *SubmitMultiSM) GetValidityPeriod() string {
	return c.validityPeriod
}

func (c *SubmitMultiSM) SetServiceType(value string) *Exception.Exception {
	err := c.CheckCStringMax(value, int(Data.SM_SRVTYPE_LEN))
	if err != nil {
		return err
	}

	c.serviceType = value
	return nil
}

func (c *SubmitMultiSM) SetSourceAddr(value *Address) {
	c.sourceAddr = value
}

func (c *SubmitMultiSM) SetSourceAddrFromStr(st string) *Exception.Exception {
	tmp, err := NewAddressWithAddr(st)
	if err != nil {
		return err
	}

	c.sourceAddr = tmp
	return nil
}

func (c *SubmitMultiSM) SetSourceAddrFromStrTon(ton, npi byte, st string) *Exception.Exception {
	tmp, err := NewAddressWithTonNpiAddr(ton, npi, st)
	if err != nil {
		return err
	}

	c.sourceAddr = tmp
	return nil
}

func (c *SubmitMultiSM) SetDestAddr(value *Address) {
	c.destAddr = value
}

func (c *SubmitMultiSM) SetDestAddrFromStr(st string) *Exception.Exception {
	tmp, err := NewAddressWithAddr(st)
	if err != nil {
		return err
	}

	c.destAddr = tmp
	return nil
}

func (c *SubmitMultiSM) SetDestAddrFromStrTon(ton, npi byte, st string) *Exception.Exception {
	tmp, err := NewAddressWithTonNpiAddr(ton, npi, st)
	if err != nil {
		return err
	}

	c.destAddr = tmp
	return nil
}

func (c *SubmitMultiSM) GetServiceType() string {
	return c.serviceType
}

func (c *SubmitMultiSM) GetSourceAddr() *Address {
	return c.sourceAddr
}

func (c *SubmitMultiSM) GetDestAddr() *Address {
	return c.destAddr
}

func (c *SubmitMultiSM) HasUserMessageReference() bool {
	return c.userMessageReference.HasValue()
}

func (c *SubmitMultiSM) HasSourcePort() bool {
	return c.sourcePort.HasValue()
}

func (c *SubmitMultiSM) HasSourceAddrSubunit() bool {
	return c.sourceAddrSubunit.HasValue()
}

func (c *SubmitMultiSM) HasDestinationPort() bool {
	return c.destinationPort.HasValue()
}

func (c *SubmitMultiSM) HasDestAddrSubunit() bool {
	return c.destAddrSubunit.HasValue()
}

func (c *SubmitMultiSM) HasSarMsgRefNum() bool {
	return c.sarMsgRefNum.HasValue()
}

func (c *SubmitMultiSM) HasSarTotalSegments() bool {
	return c.sarTotalSegments.HasValue()
}

func (c *SubmitMultiSM) HasSarSegmentSeqnum() bool {
	return c.sarSegmentSeqnum.HasValue()
}

func (c *SubmitMultiSM) HasPayloadType() bool {
	return c.payloadType.HasValue()
}

func (c *SubmitMultiSM) HasMessagePayload() bool {
	return c.messagePayload.HasValue()
}

func (c *SubmitMultiSM) HasPrivacyIndicator() bool {
	return c.privacyIndicator.HasValue()
}

func (c *SubmitMultiSM) HasCallbackNum() bool {
	return c.callbackNum.HasValue()
}

func (c *SubmitMultiSM) HasCallbackNumPresInd() bool {
	return c.callbackNumPresInd.HasValue()
}

func (c *SubmitMultiSM) HasCallbackNumAtag() bool {
	return c.callbackNumAtag.HasValue()
}

func (c *SubmitMultiSM) HasSourceSubaddress() bool {
	return c.sourceSubaddress.HasValue()
}

func (c *SubmitMultiSM) HasDestSubaddress() bool {
	return c.destSubaddress.HasValue()
}

func (c *SubmitMultiSM) HasDisplayTime() bool {
	return c.displayTime.HasValue()
}

func (c *SubmitMultiSM) HasSmsSignal() bool {
	return c.smsSignal.HasValue()
}

func (c *SubmitMultiSM) HasMsValidity() bool {
	return c.msValidity.HasValue()
}

func (c *SubmitMultiSM) HasMsMsgWaitFacilities() bool {
	return c.msMsgWaitFacilities.HasValue()
}

func (c *SubmitMultiSM) HasAlertOnMsgDelivery() bool {
	return c.alertOnMsgDelivery.HasValue()
}

func (c *SubmitMultiSM) HasLanguageIndicator() bool {
	return c.languageIndicator.HasValue()
}

func (c *SubmitMultiSM) SetUserMessageReference(value int16) *Exception.Exception {
	return c.userMessageReference.SetValue(value)
}

func (c *SubmitMultiSM) SetSourcePort(value int16) *Exception.Exception {
	return c.sourcePort.SetValue(value)
}

func (c *SubmitMultiSM) SetSourceAddrSubunit(value byte) *Exception.Exception {
	return c.sourceAddrSubunit.SetValue(value)
}

func (c *SubmitMultiSM) SetDestinationPort(value int16) *Exception.Exception {
	return c.destinationPort.SetValue(value)
}

func (c *SubmitMultiSM) SetDestAddrSubunit(value byte) *Exception.Exception {
	return c.destAddrSubunit.SetValue(value)
}

func (c *SubmitMultiSM) SetSarMsgRefNum(value int16) *Exception.Exception {
	return c.sarMsgRefNum.SetValue(value)
}

func (c *SubmitMultiSM) SetSarTotalSegments(value uint8) *Exception.Exception {
	return c.sarTotalSegments.SetValue(value)
}

func (c *SubmitMultiSM) SetSarSegmentSeqnum(value uint8) *Exception.Exception {
	return c.sarSegmentSeqnum.SetValue(value)
}

func (c *SubmitMultiSM) SetPayloadType(value byte) *Exception.Exception {
	return c.payloadType.SetValue(value)
}

func (c *SubmitMultiSM) SetMessagePayload(value *Utils.ByteBuffer) *Exception.Exception {
	return c.messagePayload.SetValue(value)
}

func (c *SubmitMultiSM) SetPrivacyIndicator(value byte) *Exception.Exception {
	return c.privacyIndicator.SetValue(value)
}

func (c *SubmitMultiSM) SetCallbackNum(value *Utils.ByteBuffer) *Exception.Exception {
	return c.callbackNum.SetValue(value)
}

func (c *SubmitMultiSM) SetCallbackNumPresInd(value byte) *Exception.Exception {
	return c.callbackNumPresInd.SetValue(value)
}

func (c *SubmitMultiSM) SetCallbackNumAtag(value *Utils.ByteBuffer) *Exception.Exception {
	return c.callbackNumAtag.SetValue(value)
}

func (c *SubmitMultiSM) SetSourceSubaddress(value *Utils.ByteBuffer) *Exception.Exception {
	return c.sourceSubaddress.SetValue(value)
}

func (c *SubmitMultiSM) SetDestSubaddress(value *Utils.ByteBuffer) *Exception.Exception {
	return c.destSubaddress.SetValue(value)
}

func (c *SubmitMultiSM) SetDisplayTime(value byte) *Exception.Exception {
	return c.displayTime.SetValue(value)
}

func (c *SubmitMultiSM) SetSmsSignal(value int16) *Exception.Exception {
	return c.smsSignal.SetValue(value)
}

func (c *SubmitMultiSM) SetMsValidity(value byte) *Exception.Exception {
	return c.msValidity.SetValue(value)
}

func (c *SubmitMultiSM) SetMsMsgWaitFacilities(value byte) *Exception.Exception {
	return c.msMsgWaitFacilities.SetValue(value)
}

func (c *SubmitMultiSM) SetAlertOnMsgDelivery(value bool) *Exception.Exception {
	return c.alertOnMsgDelivery.SetValue(value)
}

func (c *SubmitMultiSM) SetLanguageIndicator(value byte) *Exception.Exception {
	return c.languageIndicator.SetValue(value)
}

func (c *SubmitMultiSM) GetUserMessageReference() (int16, *Exception.Exception) {
	return c.userMessageReference.GetValue()
}

func (c *SubmitMultiSM) GetSourcePort() (int16, *Exception.Exception) {
	return c.sourcePort.GetValue()
}

func (c *SubmitMultiSM) GetSourceAddrSubunit() (byte, *Exception.Exception) {
	return c.sourceAddrSubunit.GetValue()
}

func (c *SubmitMultiSM) GetDestinationPort() (int16, *Exception.Exception) {
	return c.destinationPort.GetValue()
}

func (c *SubmitMultiSM) GetDestAddrSubunit() (byte, *Exception.Exception) {
	return c.destAddrSubunit.GetValue()
}

func (c *SubmitMultiSM) GetSarMsgRefNum() (int16, *Exception.Exception) {
	return c.sarMsgRefNum.GetValue()
}

func (c *SubmitMultiSM) GetSarTotalSegments() (byte, *Exception.Exception) {
	return c.sarTotalSegments.GetValue()
}

func (c *SubmitMultiSM) GetSarSegmentSeqnum() (byte, *Exception.Exception) {
	return c.sarSegmentSeqnum.GetValue()
}

func (c *SubmitMultiSM) GetPayloadType() (byte, *Exception.Exception) {
	return c.payloadType.GetValue()
}

func (c *SubmitMultiSM) GetMessagePayload() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.messagePayload.GetValue()
}

func (c *SubmitMultiSM) GetPrivacyIndicator() (byte, *Exception.Exception) {
	return c.privacyIndicator.GetValue()
}

func (c *SubmitMultiSM) GetCallbackNum() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.callbackNum.GetValue()
}

func (c *SubmitMultiSM) GetCallbackNumPresInd() (byte, *Exception.Exception) {
	return c.callbackNumPresInd.GetValue()
}

func (c *SubmitMultiSM) GetCallbackNumAtag() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.callbackNumAtag.GetValue()
}

func (c *SubmitMultiSM) GetSourceSubaddress() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.sourceSubaddress.GetValue()
}

func (c *SubmitMultiSM) GetDestSubaddress() (*Utils.ByteBuffer, *Exception.Exception) {
	return c.destSubaddress.GetValue()
}

func (c *SubmitMultiSM) GetDisplayTime() (byte, *Exception.Exception) {
	return c.displayTime.GetValue()
}

func (c *SubmitMultiSM) GetSmsSignal() (int16, *Exception.Exception) {
	return c.smsSignal.GetValue()
}

func (c *SubmitMultiSM) GetMsValidity() (byte, *Exception.Exception) {
	return c.msValidity.GetValue()
}

func (c *SubmitMultiSM) GetMsMsgWaitFacilities() (byte, *Exception.Exception) {
	return c.msMsgWaitFacilities.GetValue()
}

func (c *SubmitMultiSM) GetAlertOnMsgDelivery() (bool, *Exception.Exception) {
	return c.alertOnMsgDelivery.GetValue()
}

func (c *SubmitMultiSM) GetLanguageIndicator() (byte, *Exception.Exception) {
	return c.languageIndicator.GetValue()
}

func (c *SubmitMultiSM) GetReplaceIfPresentFlag() byte {
	return c.replaceIfPresentFlag
}

func (c *SubmitMultiSM) SetReplaceIfPresentFlag(dat byte) {
	c.replaceIfPresentFlag = dat
}

func (c *SubmitMultiSM) SetProtocolId(dat byte) {
	c.protocolId = dat
}

func (c *SubmitMultiSM) GetProtocolId() byte {
	return c.protocolId
}

func (c *SubmitMultiSM) SetPriorityFlag(dat byte) {
	c.priorityFlag = dat
}

func (c *SubmitMultiSM) GetPriorityFlag() byte {
	return c.priorityFlag
}
