package pdu

import (
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
func (u UDH) UDHL() (l int) {
	for i := range u {
		l += len(u[i].Data)
	}
	return l
}

// MarshalBinary marshal UDH into bytes array
// The first byte is UDHL
// MarshalBinary preserve InformationElement order as they appears in the UDH
func (u *UDH) MarshalBinary() (b []byte, err error) {
	if len(*u) == 0 {
		return nil, nil
	}

	b = []byte{byte(u.UDHL())}

	for _, ie := range *u {
		b = append(b, ie.Data...)
	}
	return b, nil
}

// UnmarshalBinary reads the InformationElements from the binary User Data
// Header.
// Unmarshal preserve InfoElement order as they appears in the raw data
// The src contains the complete UDH, including the UDHL and all IEs.
// The function returns the number of bytes read from src, and any error
// detected while unmarshalling.
func (u *UDH) UnmarshalBinary(src []byte) (int, error) {
	if len(src) < 1 {
		return 0, fmt.Errorf("Decode error UDHL %d underflow", 0)
	}

	udhl := int(src[0])
	udhl++ // so it includes itself
	ri := 1
	if len(src) < udhl {
		return ri, fmt.Errorf("Decode error InfoElement %d underflow", ri)
	}

	ies := []InfoElement(nil)
	for ri < udhl {
		if udhl < ri+2 {
			return ri, fmt.Errorf("Decode error InfoElement %d underflow", ri)
		}
		ie := InfoElement{}
		ie.ID = src[ri]
		ri++
		iedl := int(src[ri])
		ri++
		if len(src) < ri+iedl {
			return ri, fmt.Errorf("Decode error InfoElement %d underflow", ri)
		}
		ie.Data = append([]byte(nil), src[ri:ri+iedl]...)
		ri += iedl
		ies = append(ies, ie)
	}

	*u = UDH(ies)
	return udhl, nil
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
func (u UDH) GetConcatInfo() (totalParts, partNum, mref int, found bool) {
	if len(u) == 0 {
		found = false
		return
	}
	if ie, found := u.FindInfoElement(data.UDH_CONCAT_MSG_8_BIT_REF); found && len(ie.Data) == 3 {
		mref = int(ie.Data[0])
		totalParts = int(ie.Data[1])
		partNum = int(ie.Data[2])
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
// NOTE: in my opinion, using 16 bit reference is very unnecessary, usi
func NewIEConcatMessage(totalParts, partNum, mref int) InfoElement {
	return InfoElement{
		ID:   data.UDH_CONCAT_MSG_8_BIT_REF,
		Data: []byte{data.UDH_CONCAT_MSG_8_BIT_REF, 0x03, byte(mref), byte(totalParts), byte(partNum)},
	}
}
