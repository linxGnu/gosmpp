package pdu

import (
	"io"
	"sync/atomic"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"
	"github.com/linxGnu/gosmpp/pdu/tlv"
	"github.com/linxGnu/gosmpp/utils"
)

var sequenceNumber int32

func nextSequenceNumber() (v int32) {
	// & 0x7FFFFFFF: cater for integer overflow
	// Allowed range is 0x01 to 0x7FFFFFFF. This
	// will still result in a single invalid value
	// of 0x00 every ~2 billion PDUs (not too bad):
	if v = atomic.AddInt32(&sequenceNumber, 1) & 0x7FFFFFFF; v <= 0 {
		v = 1
	}
	return
}

// PDU represents PDU interface.
type PDU interface {
	Marshal(*utils.ByteBuffer)
	Unmarshal(*utils.ByteBuffer) error
	CanResponse() bool
	GetResponse() PDU
}

type base struct {
	Header
	OptionalParameters map[tlv.Tag]tlv.Field
}

func newBase() (v base) {
	v.OptionalParameters = make(map[tlv.Tag]tlv.Field)
	v.AssignSequenceNumber()
	return
}

func (c *base) unmarshal(b *utils.ByteBuffer, bodyReader func(*utils.ByteBuffer) error) (err error) {
	fullLen := b.Len()

	if err = c.Header.Unmarshal(b); err == nil {

		// try to unmarshal body
		if bodyReader != nil {
			err = bodyReader(b)
		}

		if err == nil {
			// command length
			cmdLength := int(c.CommandLength)

			// got - total read byte(s)
			got := fullLen - b.Len()
			if got > cmdLength {
				err = errors.ErrInvalidPDU
				return
			}

			// have optional body?
			if got < cmdLength {

				// the rest is optional body
				var optionalBody []byte
				optionalBody, err = b.ReadN(cmdLength - got)
				if err != nil {
					return
				}

				if err = c.unmarshalOptionalBody(optionalBody); err != nil {
					return
				}
			}

			// validate again
			if b.Len() != fullLen-cmdLength {
				err = errors.ErrInvalidPDU
			}
		}
	}

	return
}

func (c *base) unmarshalOptionalBody(body []byte) (err error) {
	buf := utils.NewBuffer(body)
	for buf.Len() > 0 {
		var field tlv.Field
		if err = field.Unmarshal(buf); err != nil {
			return
		}
		c.OptionalParameters[field.Tag] = field
	}
	return
}

// Marshal to buffer.
func (c *base) marshal(b *utils.ByteBuffer, bodyWriter func(*utils.ByteBuffer)) {
	bodyBuf := utils.NewBuffer(nil)

	// body
	if bodyWriter != nil {
		bodyWriter(bodyBuf)
	}

	// optional body
	for _, v := range c.OptionalParameters {
		v.Marshal(bodyBuf)
	}

	// write header
	c.CommandLength = int32(data.PDU_HEADER_SIZE + bodyBuf.Len())
	c.Header.Marshal(b)

	// write body and its optional params
	b.WriteBuffer(bodyBuf)
}

// RegisterOptionalParam register optional param.
func (c *base) RegisterOptionalParam(tlv tlv.Field) {
	c.OptionalParameters[tlv.Tag] = tlv
}

// IsOk is ok message.
func (c *base) IsOk() bool {
	return c.CommandStatus == int32(data.ESME_ROK)
}

// IsGNack is generic n-ack.
func (c *base) IsGNack() bool {
	return c.CommandID == int32(data.GENERIC_NACK)
}

// Parse PDU from reader.
func Parse(r io.Reader) (pdu PDU, err error) {
	var headerBytes [16]byte

	if _, err = io.ReadFull(r, headerBytes[:]); err != nil {
		return
	}

	header := ParseHeader(headerBytes)
	if header.CommandLength < 16 {
		err = errors.ErrInvalidPDU
		return
	}

	// too large?
	if header.CommandLength > data.MAX_PDU_LEN {
		err = errors.ErrInvalidPDU
		return
	}

	// read pdu body
	bodyBytes := make([]byte, header.CommandLength-16)
	if len(bodyBytes) > 0 {
		if _, err = io.ReadFull(r, bodyBytes); err != nil {
			return
		}
	}

	// try to create pdu
	if pdu, err = CreatePDUFromCmdID(header.CommandID); err == nil {
		buf := utils.NewBuffer(make([]byte, 0, header.CommandLength))
		_, _ = buf.Write(headerBytes[:])
		_, _ = buf.Write(bodyBytes)

		err = pdu.Unmarshal(buf)
	}

	return
}
