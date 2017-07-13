package PDU

import (
	"fmt"
	"sync"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/PDU/TLV"
	"github.com/linxGnu/gosmpp/Utils"
)

const (
	VALID_NONE   byte = 0
	VALID_HEADER byte = 1
	VALID_BODY   byte = 2
	VALID_ALL    byte = 3
)

var SequenceNumber int32
var l sync.Mutex

func nextSequenceNumber() int32 {
	l.Lock()
	defer l.Unlock()

	// & 0x7FFFFFFF: cater for integer overflow
	// Allowed range is 0x01 to 0x7FFFFFFF. This
	// will still result in a single invalid value
	// of 0x00 every ~2 billion PDUs (not too bad):
	SequenceNumber++
	return SequenceNumber & 0x7FFFFFFF
}

type IPDU interface {
	IPDUHeader
	GetBody() (*Utils.ByteBuffer, *Exception.Exception, IPDU)
	SetBody(buffer *Utils.ByteBuffer) (*Exception.Exception, IPDU)
	GetInstance() (IPDU, error)
	SetData(buffer *Utils.ByteBuffer) (*Exception.Exception, IPDU)
	GetData() (*Utils.ByteBuffer, *Exception.Exception, IPDU)
	AssignSequenceNumber()
	IsEquals(pdu IPDU) bool
}

type PDU struct {
	Common.ByteData
	Header                  *PDUHeader
	OptionalParameters      []TLV.ITLV
	ExtraOptionalParameters []TLV.ITLV
	SequenceNumberChanged   bool
	Valid                   byte
	ApplicationSpecificInfo map[interface{}]interface{}
	RealRef                 IPDU
}

func NewPDU() *PDU {
	a := &PDU{}
	a.Construct()

	return a
}

func NewPDUWithCommand(commandId int32) *PDU {
	a := NewPDU()
	a.CheckHeader()
	a.SetCommandId(commandId)

	return a
}

func (c *PDU) Construct() {
	defer c.SetRealReference(c)
	c.ByteData.Construct()

	c.OptionalParameters = make([]TLV.ITLV, 0)
	c.ExtraOptionalParameters = make([]TLV.ITLV, 0)
	c.Valid = VALID_ALL
}

func (c *PDU) CheckHeader() {
	if c.Header == nil {
		c.Header = NewPDUHeader()
	}
}

func (c *PDU) CanResponse() bool {
	return false
}

func (c *PDU) IsRequest() bool {
	return false
}

func (c *PDU) IsResponse() bool {
	return false
}

func (c *PDU) AssignSequenceNumber() {
	c.AssignSequenceNumber0(false)
}

func (c *PDU) AssignSequenceNumber0(always bool) {
	if !c.SequenceNumberChanged || always {
		c.SetSequenceNumber(nextSequenceNumber())
	}
}

func (c *PDU) ResetSequenceNumber() {
	c.SetSequenceNumber(0)
	c.SequenceNumberChanged = false
}

func (c *PDU) SetData(buf *Utils.ByteBuffer) (err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	source = c.This.(IPDU)

	if buf == nil || buf.Buffer == nil {
		err = Exception.NewExceptionFromStr("PDU: set data buffer is nil")
		return
	}

	initialBufLen := buf.Len()

	c.SetValid(VALID_NONE)
	if initialBufLen < int(Data.PDU_HEADER_SIZE) {
		err = Exception.NotEnoughDataInByteBufferException
		return
	}

	headerBuf, err := buf.Read_Bytes(int(Data.PDU_HEADER_SIZE))
	if err != nil {
		return
	}

	err = c.SetHeader(headerBuf)
	if err != nil {
		return
	}
	c.SetValid(VALID_HEADER)

	err, _ = source.SetBody(buf)
	if err != nil {
		return
	}
	c.SetValid(VALID_BODY)

	got := initialBufLen - buf.Len()
	if got < int(c.GetCommandLength()) {
		b1, e1 := buf.Read_Bytes(int(c.GetCommandLength()) - got)
		if e1 != nil {
			err = e1
			return
		}

		err = c.SetOptionalBody(b1)
		if err != nil {
			return
		}
	}
	c.SetValid(VALID_ALL)

	if buf.Len() != initialBufLen-int(c.GetCommandLength()) {
		err = Exception.InvalidPDUException
		return
	}

	return
}

func (c *PDU) GetData() (buf *Utils.ByteBuffer, err *Exception.Exception, source IPDU) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			buf = nil
		}
	}()

	source = c.This.(IPDU)

	body, err, _ := source.GetBody()
	if err != nil {
		return
	} else if body == nil {
		body = Utils.NewBuffer([]byte{})
	}

	opbody, err := c.GetOptionalBody()
	if err != nil {
		return
	} else if opbody == nil {
		opbody = Utils.NewBuffer([]byte{})
	}

	buf = Utils.NewBuffer(make([]byte, 0, body.Len()+opbody.Len()))

	buf.Write_Buffer(body)
	buf.Write_Buffer(opbody)

	c.SetCommandLength(int32(buf.Len()) + Data.PDU_HEADER_SIZE)

	pduBuf, err := c.GetHeader()
	if err != nil {
		return
	}

	err = pduBuf.Write_Buffer(buf)
	return pduBuf, err, source
}

func (c *PDU) SetValid(valid byte) {
	c.Valid = valid
}

func (c *PDU) GetValid() byte {
	return c.Valid
}

func (c *PDU) IsValid() bool {
	return c.GetValid() == VALID_ALL
}

func (c *PDU) IsInValid() bool {
	return c.GetValid() == VALID_NONE
}

func (c *PDU) IsHeaderValid() bool {
	return c.GetValid() == VALID_HEADER
}

func (c *PDU) SetHeader(header *Utils.ByteBuffer) *Exception.Exception {
	c.CheckHeader()

	err := c.Header.SetData(header)
	if err != nil {
		return err
	}

	c.SequenceNumberChanged = true
	return nil
}

func (c *PDU) GetHeader() (*Utils.ByteBuffer, *Exception.Exception) {
	c.CheckHeader()
	return c.Header.GetData()
}

func (c *PDU) SetOptionalBody(buf *Utils.ByteBuffer) *Exception.Exception {
	if buf == nil || buf.Buffer == nil {
		return Exception.NewExceptionFromStr("PDU: optional body buffer is nil")
	}

	for buf.Len() > 0 {
		tlvHeader, err := buf.Read_Bytes(int(Data.TLV_HEADER_SIZE))
		if err != nil {
			return err
		}

		tag, err := tlvHeader.Read_Short()
		if err != nil {
			return err
		}

		tlv := c.findOptional(c.OptionalParameters, tag)
		if tlv == nil {
			tlv = TLV.NewTLVOctetsWithTag(tag)
			c.registerExtraOptional(tlv)
		}

		length, err := tlvHeader.Read_Short()
		if err != nil {
			return err
		}

		tlvBuf, err := buf.Read_Bytes(int(length))
		if err != nil {
			return err
		}

		err = tlv.SetValueData(tlvBuf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *PDU) GetOptionalBody() (res *Utils.ByteBuffer, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			res = nil
		}
	}()

	dat1, err := c.GetOptionalBodyBuffer(c.OptionalParameters)
	if err != nil {
		return nil, err
	}

	dat2, err := c.GetOptionalBodyBuffer(c.ExtraOptionalParameters)
	if err != nil {
		return nil, err
	}

	optBody := Utils.NewBuffer(make([]byte, 0, dat1.Len()+dat2.Len()))

	optBody.Write_Buffer(dat1)
	optBody.Write_Buffer(dat2)

	return optBody, nil
}

func (c *PDU) GetOptionalBodyBuffer(optionalParams []TLV.ITLV) (res *Utils.ByteBuffer, err *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			err = Exception.NewException(fmt.Errorf("%v", errs))
			res = nil
		}
	}()

	optBody := Utils.NewBuffer(make([]byte, 0, 64))

	if optionalParams == nil {
		return optBody, nil
	}

	for _, tlv := range optionalParams {
		if tlv != nil && tlv.HasValue() {
			dat, err := tlv.GetData()
			if err == nil {
				err = optBody.Write_Buffer(dat)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return optBody, nil
}

func (c *PDU) registerOptional(tlv TLV.ITLV) {
	if tlv != nil && c.OptionalParameters != nil {
		c.OptionalParameters = append(c.OptionalParameters, tlv)
	}
}

func (c *PDU) registerExtraOptional(tlv TLV.ITLV) {
	if tlv != nil && c.ExtraOptionalParameters != nil {
		c.ExtraOptionalParameters = append(c.ExtraOptionalParameters, tlv)
	}
}

func (c *PDU) findOptional(optionalParams []TLV.ITLV, tag int16) TLV.ITLV {
	if optionalParams == nil {
		return nil
	}

	for _, tlv := range optionalParams {
		if tlv != nil {
			if tlv.GetTag() == tag {
				return tlv
			}
		}
	}

	return nil
}

func (c *PDU) SetExtraOptional(tlv TLV.ITLV) {
	c.replaceExtraOptional(tlv)
}

func (c *PDU) GetExtraOptional(tag int16) TLV.ITLV {
	return c.findOptional(c.ExtraOptionalParameters, tag)
}

func (c *PDU) replaceExtraOptional(tlv TLV.ITLV) {
	if tlv == nil {
		return
	}

	tlvTag := tlv.GetTag()

	for ind, existing := range c.ExtraOptionalParameters {
		if existing != nil && existing.GetTag() == tlvTag {
			c.ExtraOptionalParameters[ind] = tlv
			return
		}
	}

	c.registerExtraOptional(tlv)
}

func (c *PDU) GetCommandLength() int32 {
	c.CheckHeader()
	return c.Header.GetCommandLength()
}

func (c *PDU) SetCommandLength(length int32) {
	c.CheckHeader()
	c.Header.SetCommandLength(length)
}

func (c *PDU) GetCommandId() int32 {
	c.CheckHeader()
	return c.Header.GetCommandId()
}

func (c *PDU) SetCommandId(cmdid int32) {
	c.CheckHeader()
	c.Header.SetCommandId(cmdid)
}

func (c *PDU) GetCommandStatus() int32 {
	c.CheckHeader()
	return c.Header.GetCommandStatus()
}

func (c *PDU) SetCommandStatus(status int32) {
	c.CheckHeader()
	c.Header.SetCommandStatus(status)
}

func (c *PDU) GetSequenceNumber() int32 {
	c.CheckHeader()
	return c.Header.GetSequenceNumber()
}

func (c *PDU) SetSequenceNumber(seq int32) {
	c.CheckHeader()
	c.Header.SetSequenceNumber(seq)
	c.SequenceNumberChanged = true
}

func (c *PDU) IsOk() bool {
	return c.GetCommandStatus() == int32(Data.ESME_ROK)
}

func (c *PDU) IsGNack() bool {
	return c.GetCommandId() == int32(Data.GENERIC_NACK)
}

func CreatePDU(buf *Utils.ByteBuffer) (IPDU, *Exception.Exception, IPDUHeader) {
	if buf == nil || buf.Buffer == nil {
		return nil, Exception.NewExceptionFromStr("Can not create PDU with nil buffer!"), nil
	}

	headerBuf, err := buf.Read_Bytes(int(Data.PDU_HEADER_SIZE))
	if err != nil {
		return nil, Exception.HeaderIncompleteException, nil
	}

	header := NewPDUHeader()
	err = header.SetData(headerBuf)
	if err != nil {
		return nil, err, nil
	}

	if buf.Len()+int(Data.PDU_HEADER_SIZE) < int(header.GetCommandLength()) {
		return nil, Exception.MessageIncompleteException, header
	}

	pdu, err1 := CreatePDUWithCmdId(header.GetCommandId())
	if err1 != nil {
		return nil, Exception.UnknownCommandIdException, header
	}

	restBuf, err := buf.Read_Bytes(int(header.GetCommandLength()) - int(Data.PDU_HEADER_SIZE))
	if err != nil {
		return nil, err, nil
	}

	headerBuf, err = header.GetData()
	if err != nil {
		return nil, err, nil
	}

	err = headerBuf.Write_Buffer(restBuf)
	if err != nil {
		return nil, err, nil
	}

	err, _ = pdu.SetData(headerBuf)
	if err != nil {
		return pdu, err, header
	}

	return pdu, nil, header
}

func (c *PDU) IsEquals(a IPDU) bool {
	if a == nil {
		return false
	}

	return c.GetSequenceNumber() == a.GetSequenceNumber() && c.GetCommandId() == a.GetCommandId()
}
