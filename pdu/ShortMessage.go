package pdu

import (
	"sync/atomic"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"
)

var ref = uint32(0)

// ShortMessage message.
type ShortMessage struct {
	SmDefaultMsgID    byte
	message           string
	enc               data.Encoding
	udHeader          UDH
	messageData       []byte
	withoutDataCoding bool // purpose of ReplaceSM usage
}

// NewShortMessage returns new ShortMessage.
func NewShortMessage(message string) (s ShortMessage, err error) {
	err = s.SetMessageWithEncoding(message, data.GSM7BIT)
	return
}

// NewShortMessageWithEncoding returns new ShortMessage with predefined encoding.
func NewShortMessageWithEncoding(message string, enc data.Encoding) (s ShortMessage, err error) {
	err = s.SetMessageWithEncoding(message, enc)
	return
}

// NewBinaryShortMessage returns new ShortMessage.
func NewBinaryShortMessage(messageData []byte) (s ShortMessage, err error) {
	err = s.SetMessageDataWithEncoding(messageData, data.BINARY8BIT2)
	return
}

// NewBinaryShortMessageWithEncoding returns new ShortMessage with predefined encoding.
func NewBinaryShortMessageWithEncoding(messageData []byte, enc data.Encoding) (s ShortMessage, err error) {
	err = s.SetMessageDataWithEncoding(messageData, enc)
	return
}

// NewLongMessage returns long message splitted into multiple short message
func NewLongMessage(message string) (s []*ShortMessage, err error) {
	return NewLongMessageWithEncoding(message, data.GSM7BIT)
}

// NewLongMessageWithEncoding returns long message splitted into multiple short message with encoding of choice
func NewLongMessageWithEncoding(message string, enc data.Encoding) (s []*ShortMessage, err error) {
	sm := &ShortMessage{
		message: message,
		enc:     enc,
	}
	return sm.split()
}

// SetMessageWithEncoding sets message with encoding.
func (c *ShortMessage) SetMessageWithEncoding(message string, enc data.Encoding) (err error) {
	if c.messageData, err = enc.Encode(message); err == nil {
		if len(c.messageData) > data.SM_MSG_LEN {
			err = errors.ErrShortMessageLengthTooLarge
		} else {
			c.message = message
			c.enc = enc
		}

		if c.enc == data.GSM7BITPACKED { // to prevent unwanted "@"
			runeSlice := []rune(c.message)
			tLen := len(runeSlice)
			escCharsLen := len(data.GetEscapeChars(runeSlice))
			regCharsLen := tLen - escCharsLen
			nSeptet := escCharsLen*2 + regCharsLen
			if (nSeptet+1)%8 == 0 {
				c.messageData[len(c.messageData)-1] = (c.messageData[len(c.messageData)-1] & 0x01) | (0x0D << 1) /* https://en.wikipedia.org/wiki/GSM_03.38 Ref tekst: "..When there are 7 spare bits in the last octet of a message..."*/
			}
		}
	}
	return
}

// SetLongMessageWithEnc sets ShortMessage with message longer than  256 bytes
// callers are expected to call Split() after this
func (c *ShortMessage) SetLongMessageWithEnc(message string, enc data.Encoding) (err error) {
	c.message = message
	c.enc = enc
	return
}

// UDH gets user data header for short message
func (c *ShortMessage) UDH() UDH {
	return c.udHeader
}

// SetUDH sets user data header for short message
// also appends udh to the beginning of messageData
func (c *ShortMessage) SetUDH(udh UDH) {
	c.udHeader = udh
}

// SetMessageDataWithEncoding sets underlying raw data which is used for pdu marshalling.
func (c *ShortMessage) SetMessageDataWithEncoding(d []byte, enc data.Encoding) (err error) {
	if len(d) > data.SM_MSG_LEN {
		err = errors.ErrShortMessageLengthTooLarge
	} else {
		c.messageData = d
		c.enc = enc
	}
	return
}

// GetMessageData returns underlying binary message.
func (c *ShortMessage) GetMessageData() (d []byte, err error) {
	return c.messageData, nil
}

// GetMessage returns underlying message.
func (c *ShortMessage) GetMessage() (st string, err error) {
	enc := c.enc
	if enc == nil {
		enc = data.GSM7BIT
	}
	st, err = c.GetMessageWithEncoding(enc)
	return
}

// GetMessageWithEncoding returns (decoded) underlying message.
func (c *ShortMessage) GetMessageWithEncoding(enc data.Encoding) (st string, err error) {
	if len(c.messageData) > 0 {
		st, err = enc.Decode(c.messageData)
	}
	return
}

// split one short message and split into multiple short message, with UDH
// according to 33GP TS 23.040 section 9.2.3.24.1
//
// NOTE: split() will return array of length 1 if data length is still within the limit
// The encoding interface can implement the data.Splitter interface for ad-hoc splitting rule
func (c *ShortMessage) split() (multiSM []*ShortMessage, err error) {
	var encoding data.Encoding
	if c.enc == nil {
		encoding = data.GSM7BIT
	} else {
		encoding = c.enc
	}

	// check if encoding implements data.Splitter
	splitter, ok := encoding.(data.Splitter)
	// check if encoding implements data.Splitter or split is necessary
	if !ok || !splitter.ShouldSplit(c.message, data.SM_GSM_MSG_LEN) {
		err = c.SetMessageWithEncoding(c.message, c.enc)
		multiSM = []*ShortMessage{c}
		return
	}

	// Reserve 6 bytes for concat message UDH
	//
	// Good references:
	// - https://help.goacoustic.com/hc/en-us/articles/360043843154--How-character-encoding-affects-SMS-message-length
	// - https://www.twilio.com/docs/glossary/what-is-gsm-7-character-encoding
	//
	// Limitation is 160 GSM-7 characters and we also need 6 bytes for UDH
	// -> 134 octets per segment
	// -> this leaves 153 GSM-7 characters per segment.
	segments, err := splitter.EncodeSplit(c.message, data.SM_GSM_MSG_LEN-6)
	if err != nil {
		return nil, err
	}

	// prealloc result
	multiSM = make([]*ShortMessage, 0, len(segments))

	// all segments will have the same ref id
	ref := getRefNum()

	// construct SM(s)
	for i, seg := range segments {
		// create new SM, encode data
		multiSM = append(multiSM, &ShortMessage{
			enc: c.enc,
			// message: we don't really care
			messageData:       seg,
			withoutDataCoding: c.withoutDataCoding,
			udHeader:          UDH{NewIEConcatMessage(uint8(len(segments)), uint8(i+1), uint8(ref))},
		})
	}

	return
}

// Marshal implements PDU interface.
func (c *ShortMessage) Marshal(b *ByteBuffer) {
	var (
		udhBin []byte
		n      = byte(len(c.messageData))
	)

	// Prepend UDH to message data if there are any
	if c.udHeader != nil && c.udHeader.UDHL() > 0 {
		udhBin, _ = c.udHeader.MarshalBinary()
	}

	b.Grow(int(n) + 3)

	var coding byte
	if c.enc == nil {
		coding = data.GSM7BITCoding
	} else {
		coding = c.enc.DataCoding()
	}

	// data_coding
	if !c.withoutDataCoding {
		_ = b.WriteByte(coding)
	}

	// sm_default_msg_id
	_ = b.WriteByte(c.SmDefaultMsgID)

	// sm_length
	if udhBin != nil {
		_ = b.WriteByte(byte(int(n) + len(udhBin)))
		b.Write(udhBin)
	} else {
		_ = b.WriteByte(n)
	}

	// short_message
	_, _ = b.Write(c.messageData[:n])
}

// Unmarshal implements PDU interface.
func (c *ShortMessage) Unmarshal(b *ByteBuffer, udhi bool) (err error) {
	var dataCoding, n byte

	if !c.withoutDataCoding {
		if dataCoding, err = b.ReadByte(); err != nil {
			return
		}
	}

	if c.SmDefaultMsgID, err = b.ReadByte(); err != nil {
		return
	}

	if n, err = b.ReadByte(); err != nil {
		return
	}

	if c.messageData, err = b.ReadN(int(n)); err != nil {
		return
	}
	c.enc = data.FromDataCoding(dataCoding)

	// If short message length is non zero, short message contains User-Data Header
	// Else UDH should be in TLV field MessagePayload
	if udhi && n > 0 {
		udh := UDH{}
		_, err = udh.UnmarshalBinary(c.messageData)
		if err != nil {
			return
		}

		c.udHeader = udh

		f := c.udHeader.UDHL()
		if f > len(c.messageData) {
			err = errors.ErrUDHTooLong
			return
		}

		c.messageData = c.messageData[f:]
	}

	return
}

// Encoding returns message encoding.
func (c *ShortMessage) Encoding() data.Encoding {
	return c.enc
}

// returns an atomically incrementing number each time it's called
func getRefNum() uint32 {
	return atomic.AddUint32(&ref, 1)
}

// NOTE:
// When coding splitting function, I have 4 choices of abstraction
// 1. Split the message before encode
// 2. Split the message after encoded
// 3. Split the message DURING encoding (before bit packing)
// 4. Encode, unpack, split
//
// Disadvantages:
// 1. The only way to really know if each segment will fit into 134 octet limit is
//		to do some kind of simulated encoding, where you calculate the total octet
//		by iterating through each character one by one.
//		Too cumbersome
//
// 2. When breaking string at octet position 134, I have to detemeine which
//		character is it ( by doing some kind of decoding)

//		a. If the character code point does not fit in the octet
//		boundary, it has to be carried-over to the next segment.
//		The remaining bits after extracting the carry-over
//		has to be filled with zero.

//		b. If it is an escape character, then I have to backtrack
//		even further since escape chars are not allowed to be splitted
//		in the middle.
//		Since the second bytes of escape chars can be confused with
//		normal chars, I must always lookback 2 character ( repeat step a for at least 2 septet )

//		c. After extracting the carry-on
//		-> Option 2 is very hard when bit packing is already applied
//
// 3. Options 3 require extending Encoding interface,
//	The not good point is not being able to utilize the encoder's Transform() method
//	The good point is you don't have to do bit packing twice

// 4. Terrible option

// All this headaches really only apply to variable length encoding.
// When using fixed length encoding, you can really split the source message BEFORE encodes.
