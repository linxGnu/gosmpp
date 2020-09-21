package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/linxGnu/gosmpp/data"
)

// For now, this package only support message uses of UDH for message concatenation
// No plan for supporting other Enhanced Messaging Service
// Credit to https://github.com/warthog618/sms

// UDH represent User Data Header
// as defined in 3GPP TS 23.040 Section 9.2.3.24.
type UDH []InfoElement

// UDHL return the length (number of octet) of the encoded UDH itself
func (u UDH) UDHL() int {
	length := 1
	for index := range u {
		length += 2
		length += len(u[index].Data)
	}
	return length
}

// MarshalBinary marshal UDH into bytes array
// The first byte is UDHL
// MarshalBinary preserve InformationElement order as they appears in the UDH
func (u UDH) MarshalBinary() (b []byte, err error) {
	if len(u) == 0 {
		return
	}
	var buf bytes.Buffer
	buf.Grow(u.UDHL())
	buf.WriteByte(0)
	for i := range u {
		buf.WriteByte(u[i].ID)
		buf.WriteByte(byte(len(u[i].Data)))
		buf.Write(u[i].Data)
	}
	b = buf.Bytes()
	b[0] = byte(len(b)) - 1
	return
}

// UnmarshalBinary reads the InformationElements from the binary User Data
// Header.
// Unmarshal preserve InfoElement order as they appears in the raw data
// The src contains the complete UDH, including the UDHL and all IEs.
// The function returns the number of bytes read from src, and any error
// detected while unmarshalling.
func (u *UDH) UnmarshalBinary(src []byte) (int, error) {
	if length := len(src); length > 1 && length < int(src[0]+1) {
		return 0, fmt.Errorf("Decode error UDHL underflow")
	}
	length := src[0]
	payload := src[1 : length+1]
	header := UDH{}
	var start, end byte
	for i := byte(0); i < length; {
		blockSize := payload[i+1]
		start = i + 2
		end = start + blockSize
		if end <= byte(len(payload)) {
			return 0, fmt.Errorf("Decode error InfoElement underflow")
		}
		header = append(header, InfoElement{
			ID:   payload[i],
			Data: payload[start:end],
		})
		i += blockSize + 2
	}
	if len(header) > 0 {
		*u = header
	}
	return int(length + 1), nil
}

// FindInfoElement find the last occurrence of the Information Element with id
func (u UDH) FindInfoElement(id byte) (ie *InfoElement, found bool) {
	for i := len(u) - 1; i >= 0; i-- {
		if u[i].ID == id {
			return &u[i], true
		}
	}
	return nil, false
}

// GetConcatInfo return concatenated message info, return 0 if
// Concat Message InfoElement is not found in the UDH
func (u UDH) GetConcatInfo() (totalParts, sequence byte, reference uint16, found bool) {
	for i := range u {
		options := u[i].Data
		switch u[i].ID {
		case data.UDH_CONCAT_MSG_8_BIT_REF:
			if len(options) != 3 {
				return
			}
			reference = uint16(options[0])
			totalParts = options[1]
			sequence = options[2]
			found = true
		case data.UDH_CONCAT_MSG_16_BIT_REF:
			if len(options) != 4 {
				return
			}
			reference = binary.BigEndian.Uint16(options[0:2])
			totalParts = options[2]
			sequence = options[3]
			found = true
		}
	}
	return
}

// InfoElement represent a 3 parts Information-Element
// as defined in 3GPP TS 23.040 Section 9.2.3.24
// Each InfoElement is comprised of it's identifier and data
type InfoElement struct {
	ID   byte
	Data []byte
}

// NewIEConcatMessage  turn a new IE element for concat message info
// IE.Data is populated at time of object creation
func NewIEConcatMessage(totalParts, sequence byte, reference uint16) InfoElement {
	var element InfoElement
	if reference > 0xFF {
		element.ID = data.UDH_CONCAT_MSG_16_BIT_REF
		element.Data = []byte{0, 0, totalParts, sequence}
		binary.BigEndian.PutUint16(element.Data[0:2], reference)
	} else {
		element.ID = data.UDH_CONCAT_MSG_8_BIT_REF
		element.Data = []byte{byte(reference), totalParts, sequence}
	}
	return element
}
