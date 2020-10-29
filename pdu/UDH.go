package pdu

import (
	"bytes"
	"fmt"

	"github.com/linxGnu/gosmpp/data"
)

// For now, this package only support message uses of UDH for message concatenation
// No plan for supporting other Enhanced Messaging Service
// Credit to https://github.com/warthog618/sms

// UDH represent User Data Header
// as defined in 3GPP TS 23.040 Section 9.2.3.24.
type UDH []InfoElement

// UDHL returns length (octets) of encoded UDH, including the UDHL byte
// If the total UDHL is larger than what length(byte) can specified,
// this will truncate IE until total length fit within 256, if you want to
// check if any IE has been truncated, see if UDHL() > 2^8
func (u UDH) UDHL() (l int) {
	if len(u) == 0 {
		return 0
	}
	for i := range u {
		l += 2 + len(u[i].Data) // to account for id and type bytes
		if l > 255 {
			l -= 2 + len(u[i].Data)
			break
		}
	}
	return l + 1 // include the udhlength byte itself
}

// MarshalBinary marshal UDH into bytes array
// The first byte is UDHL
// MarshalBinary preserve InformationElement order as they appears in the UDH
//
// If the total UDHL is larger than what length(byte) can specified,
// this will truncate IE until total length fit within 256, if you want to
// check if any IE has been truncated, see if UDHL() > 2^8
func (u UDH) MarshalBinary() (b []byte, err error) {
	if len(u) == 0 {
		return
	}
	buf := new(bytes.Buffer)
	buf.WriteByte(0) // reserve the firt byte for UDHL

	truncLimit := uint32(255) // -1 for the UDHL byte itself
	length := uint32(0)
L:
	for i := 0; i < len(u); i++ {
		// Begin marshaling UDH data, each IE is composed of 3 parts:
		//		[ ID_1, LENGTH_1, DATA_N ]
		// When adding a new IE, if total length ID + LEN + DATA
		// exceed 256, we skip that IE altogether
		length += uint32(2 + len(u[i].Data))
		if length > truncLimit { // limit exceeded, break loop
			length -= uint32(2 + len(u[i].Data))
			break L
		}
		buf.WriteByte(u[i].ID)
		buf.WriteByte(byte(len(u[i].Data)))
		buf.Write(u[i].Data)
	}
	b = buf.Bytes() // final assignment and encode length
	b[0] = byte(length)
	return
}

// UnmarshalBinary reads the InformationElements from raw binary UDH.
// Unmarshal preserve InfoElement order as they appears in the raw data
// The src contains the complete UDH, including the UDHL and all IEs.
// Returns the number of bytes read from src, and the first error
// detected while unmarshalling.
//
// Since UDHL can only represented in 1 byte, UnmarshalBinary
// will only read up to a maximum of 256 byte regardless of src length
func (u *UDH) UnmarshalBinary(src []byte) (read int, err error) {
	if len(src) < 1 {
		err = fmt.Errorf("Decode error UDHL %d underflow", 0)
		return
	}

	udhl, read := int(src[0]), 1
	if len(src) < udhl {
		err = fmt.Errorf("Decode error UDH underflow, expect len %d got %d", udhl, len(src))
		return
	}
	if udhl == 0 {
		err = fmt.Errorf("error: UDHL length is 0, probably sender mistake forgot to include UDH but still set UDH flag in ESME_CLASS")
		return
	}

	ies := []InfoElement{}
	for read < udhl { // loop until we still have data to read
		ie := InfoElement{}
		var r int
		r, err = ie.UnmarshalBinary(src[read:])
		if err != nil {
			return
		}
		read += r
		ies = append(ies, ie)
	}
	*u = UDH(ies)
	return read, nil
}

// FindInfoElement find the first occurrence of the Information Element with id
func (u UDH) FindInfoElement(id byte) (ie *InfoElement, found bool) {
	for i := len(u) - 1; i >= 0; i-- {
		if u[i].ID == id {
			return &u[i], true
		}
	}
	return nil, false
}

// GetConcatInfo return the FIRST concatenated message IE,
func (u UDH) GetConcatInfo() (totalParts, partNum, mref uint8, found bool) {
	if len(u) == 0 {
		found = false
		return
	}
	if ie, ok := u.FindInfoElement(data.UDH_CONCAT_MSG_8_BIT_REF); ok && len(ie.Data) == 3 {
		mref = uint8(ie.Data[0])
		totalParts = uint8(ie.Data[1])
		partNum = uint8(ie.Data[2])
		found = ok
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
func NewIEConcatMessage(totalParts, partNum, mref uint8) InfoElement {
	return InfoElement{
		ID:   data.UDH_CONCAT_MSG_8_BIT_REF,
		Data: []byte{byte(mref), byte(totalParts), byte(partNum)},
	}
}

// UnmarshalBinary unmarshal IE from binary in src, only read a single IE,
// expect src at least of length 2 with correct IE format:
//		[ ID_1, LENGTH_1, DATA_N ]
func (ie *InfoElement) UnmarshalBinary(src []byte) (read int, err error) {
	if len(src) < 2 {
		err = fmt.Errorf("Decode error InfoElement underflow, len = %d", len(src))
		return
	}
	ieLen := int(src[1]) // second byte is len
	if len(src) < ieLen+2 {
		err = fmt.Errorf("Decode error InfoElement underflow, expect length %d, got %d", ieLen, len(src))
		return
	}

	ie.ID = src[0]               // first byte is ID
	ie.Data = src[2:(ieLen + 2)] // 3rd byte onward is data
	read = int(2 + ieLen)
	return
}
