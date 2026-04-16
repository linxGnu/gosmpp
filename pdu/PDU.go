package pdu

import (
	"io"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"
)

// PDU represents PDU interface.
type PDU interface {
	// Marshal PDU to buffer.
	Marshal(*ByteBuffer)

	// Unmarshal PDU from buffer.
	Unmarshal(*ByteBuffer) error

	// CanResponse indicates that PDU could response to SMSC.
	CanResponse() bool

	// GetResponse PDU.
	GetResponse() PDU

	// RegisterOptionalParam assigns an optional param.
	RegisterOptionalParam(Field)

	// GetHeader returns PDU header.
	GetHeader() Header

	// IsOk returns true if command status is OK.
	IsOk() bool

	// IsGNack returns true if PDU is GNack.
	IsGNack() bool

	// AssignSequenceNumber assigns sequence number auto-incrementally.
	AssignSequenceNumber()

	// ResetSequenceNumber resets sequence number.
	ResetSequenceNumber()

	// GetSequenceNumber returns assigned sequence number.
	GetSequenceNumber() int32

	// SetSequenceNumber manually sets sequence number.
	SetSequenceNumber(int32)
}

type base struct {
	Header
	OptionalParameters map[Tag]Field
}

func newBase() (v base) {
	v.OptionalParameters = make(map[Tag]Field)
	v.AssignSequenceNumber()
	return
}

// GetHeader returns pdu header.
func (c *base) GetHeader() Header {
	return c.Header
}

func (c *base) unmarshal(b *ByteBuffer, bodyReader func(*ByteBuffer) error) (err error) {
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

			// body < command_length, still have optional parameters ?
			if got < cmdLength {
				var optParam []byte
				if optParam, err = b.ReadN(cmdLength - got); err == nil {
					err = c.unmarshalOptionalParam(optParam)
				}
				if err != nil {
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

func (c *base) unmarshalOptionalParam(optParam []byte) (err error) {
	buf := NewBuffer(optParam)
	for buf.Len() > 0 {
		var field Field
		if err = field.Unmarshal(buf); err == nil {
			c.OptionalParameters[field.Tag] = field
		} else {
			return
		}
	}
	return
}

// Marshal to buffer.
func (c *base) marshal(b *ByteBuffer, bodyWriter func(*ByteBuffer)) {
	bodyBuf := NewBuffer(nil)

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
func (c *base) RegisterOptionalParam(tlv Field) {
	c.OptionalParameters[tlv.Tag] = tlv
}

// IsOk is status ok.
func (c *base) IsOk() bool {
	return c.CommandStatus == data.ESME_ROK
}

// IsGNack is generic n-ack.
func (c *base) IsGNack() bool {
	return c.CommandID == data.GENERIC_NACK
}

// GetSARConcatInfo checks optional TLV params for SAR fields and returns concatenation info.
func (c *base) GetSARConcatInfo() (totalParts, partNum byte, mref uint16, found bool) {
	var foundRef, foundTot, foundSeq bool
	// Iterate over all optional TLV fields attached to this PDU
	for _, tlv := range c.OptionalParameters { // (Assume c.OptionalParams holds the list of TLV Field structs)
		switch tlv.Tag {
		case TagSarMsgRefNum: // SAR reference number (should be 2 bytes)
			if len(tlv.Data) == 2 {
				// Combine two bytes into a 16-bit reference (big-endian as per SMPP spec)
				mref = uint16(tlv.Data[0])<<8 | uint16(tlv.Data[1])
				foundRef = true
			}
		case TagSarTotalSegments: // Total number of segments (1 byte)
			if len(tlv.Data) == 1 {
				totalParts = tlv.Data[0]
				foundTot = true
			}
		case TagSarSegmentSeqnum: // Segment sequence number (1 byte)
			if len(tlv.Data) == 1 {
				partNum = tlv.Data[0]
				foundSeq = true
			}
		}
	}
	// All three must be found to consider the data complete
	found = foundRef && foundTot && foundSeq
	if !found {
		// If any part is missing or lengths were incorrect, return with 'found' = false and zeros
		totalParts, partNum, mref = 0, 0, 0
	}
	return
}

// Parse PDU from reader.
func Parse(r io.Reader) (pdu PDU, err error) {
	var headerBytes [16]byte

	if _, err = io.ReadFull(r, headerBytes[:]); err != nil {
		return
	}

	header := ParseHeader(headerBytes)
	if header.CommandLength < 16 || header.CommandLength > data.MAX_PDU_LEN {
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
		buf := NewBuffer(make([]byte, 0, header.CommandLength))
		_, _ = buf.Write(headerBytes[:])
		if len(bodyBytes) > 0 {
			_, _ = buf.Write(bodyBytes)
		}
		err = pdu.Unmarshal(buf)
	}

	return
}
