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

// UDHL returns length (octets) of encoded UDH, including the UDHL byte.
//
// If there is no InfoElement (IE), returns 0. If total length exceed 255, return -1.
func (u UDH) UDHL() (l int) {
	if len(u) == 0 {
		return
	}

	for i := range u {
		if len(u[i].Data) > 255 {
			return -1
		}

		// to account for id and type bytes
		if l += 2 + len(u[i].Data); l > 255 {
			return -1
		}
	}

	// include the udhlength byte itself
	l++

	return
}

// MarshalBinary marshal UDH into bytes array
// The first byte is UDHL
// MarshalBinary preserve InformationElement order as they appears in the UDH
//
// If the total UDHL is larger than what length(byte) can specified,
// this will truncate IE until total length fit within 256, if you want to
// check if any IE has been truncated, see if UDHL() > 2^8
func (u UDH) MarshalBinary() (b []byte, err error) {
	reservedLength := u.UDHL()
	if reservedLength == -1 {
		err = fmt.Errorf("header limit (255 in marshal size) exceeds")
		return
	}
	if reservedLength == 0 {
		return
	}

	// reserve the first byte for UDHL
	buf := bytes.NewBuffer(make([]byte, 1, reservedLength))

	// marshalling elements
	length := 0
L:
	for i := 0; i < len(u); i++ {
		// Begin marshaling UDH data, each IE is composed of 3 parts:
		//		[ ID_1, LENGTH_1, DATA_N ]
		// When adding a new IE, if total length ID + LEN + DATA
		// exceed 256, we skip that IE altogether
		addition := 2 + len(u[i].Data)

		// limit exceeded, break loop?
		if length += addition; length > 255 {
			length -= addition
			break L
		}

		buf.WriteByte(u[i].ID)
		buf.WriteByte(byte(len(u[i].Data)))
		buf.Write(u[i].Data)
	}

	// only set buffer when UDHL length is not zero
	if length > 0 {
		// final assignment and encode length
		b = buf.Bytes()
		b[0] = byte(length)
	}

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
func (u *UDH) UnmarshalBinary(src []byte) (int, error) {
	if len(src) < 1 {
		return 0, fmt.Errorf("decode error UDHL %d underflow", 0)
	}

	udhl := int(src[0])
	if udhl == 0 {
		return 0, fmt.Errorf("error: UDHL length is 0, probably sender mistake forgot to include UDH but still set UDH flag in ESME_CLASS")
	}

	// check length, excluding first UDHL byte
	if len(src)-1 < udhl {
		return 0, fmt.Errorf("decode error UDH underflow, expect len %d got %d", udhl, len(src))
	}

	// count number of bytes which are read
	var (
		read = 1 // UDHL byte
	)

	ies := []InfoElement{}
	for read < udhl { // loop until we still have data to read
		ie := InfoElement{}

		r, err := ie.UnmarshalBinary(src[read:])
		if err != nil {
			return 0, err
		}

		// moving forward
		read += r

		// add info elements
		ies = append(ies, ie)
	}

	*u = UDH(ies)
	return read, nil
}

// FindInfoElement find the first occurrence of the Information Element with id
func (u UDH) FindInfoElement(id byte) (ie *InfoElement, found bool) {
	for i := range u {
		if u[i].ID == id {
			return &u[i], true
		}
	}
	return nil, false
}

// GetConcatInfo return the FIRST concatenated message IE,
func (u UDH) GetConcatInfo() (totalParts, partNum, mref byte, found bool) {
	if len(u) == 0 {
		found = false
		return
	}

	if ie, ok := u.FindInfoElement(data.UDH_CONCAT_MSG_8_BIT_REF); ok && len(ie.Data) == 3 {
		mref = ie.Data[0]
		totalParts = ie.Data[1]
		partNum = ie.Data[2]
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
func NewIEConcatMessage(totalParts, partNum, mref byte) InfoElement {
	return InfoElement{
		ID:   data.UDH_CONCAT_MSG_8_BIT_REF,
		Data: []byte{byte(mref), byte(totalParts), byte(partNum)},
	}
}

// UnmarshalBinary unmarshal IE from binary in src, only read a single IE,
// expect src at least of length 2 with correct IE format:
//		[ ID_1, LENGTH_1, DATA_N ]
func (ie *InfoElement) UnmarshalBinary(src []byte) (int, error) {
	if len(src) < 2 {
		return 0, fmt.Errorf("decode error InfoElement underflow, len = %d", len(src))
	}

	// second byte is len
	ieLen := int(src[1])

	// check length, excluding first 2 bytes
	if len(src)-2 < ieLen {
		return 0, fmt.Errorf("decode error InfoElement underflow, expect length %d, got %d", ieLen, len(src))
	}

	ie.ID = src[0]               // first byte is ID
	ie.Data = src[2:(ieLen + 2)] // 3rd byte onward is data

	return 2 + ieLen, nil
}
