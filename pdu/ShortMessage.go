package pdu

import (
	"sync/atomic"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"
)

var (
	ref = uint32(0)
)

// ShortMessage message.
type ShortMessage struct {
	SmDefaultMsgID    byte
	message           string
	enc               data.Encoding
	udHeader          UDH
	messageData       []byte
	withoutDataCoding bool
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

// NewLongMessage return long message splitted into multiple short message
func NewLongMessage(message string) (s []*ShortMessage, err error) {
	return NewLongMessageWithEncoding(message, data.GSM7BIT)
}

// NewLongMessage return long message splitted into multiple short message with encoding of choice
func NewLongMessageWithEncoding(message string, enc data.Encoding) (s []*ShortMessage, err error) {
	sm := &ShortMessage{
		message: message,
		enc:     enc,
	}
	return sm.Split()
}

// SetMessageWithEncoding set message with encoding.
func (c *ShortMessage) SetMessageWithEncoding(message string, enc data.Encoding) (err error) {
	if c.messageData, err = enc.Encode(message); err == nil {
		if len(c.messageData) > data.SM_MSG_LEN {
			err = errors.ErrShortMessageLengthTooLarge
		} else {
			c.message = message
			c.enc = enc
		}
	}
	return
}

// SetLongMessageWithEnc set ShortMessage with message longer than  256 bytes
// callers are expected to call Split() after this
func (c *ShortMessage) SetLongMessageWithEnc(message string, enc data.Encoding) (err error) {
	c.message = message
	c.enc = enc
	return
}

// UDH get user data header for short message
func (c *ShortMessage) UDH() UDH {
	return c.udHeader
}

// SetUDH set user data header for short message
// also appends udh to the beginning of messageData
func (c *ShortMessage) SetUDH(udh UDH) {
	c.udHeader = udh
}

// SetMessageData sets underlying raw data which is used for pdu marshalling.
func (c *ShortMessage) SetMessageData(data []byte) {
	c.messageData = data
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
	if len(c.messageData) == 0 {
		return
	}

	f, t := 0, len(c.messageData)

	if c.udHeader.UDHL() > 0 {
		f = c.udHeader.UDHL() + 1
		if f >= t {
			err = errors.ErrUDHTooLong
			return
		}
	}

	st, err = enc.Decode(c.messageData[f:t])
	return
}

// Split split one short message and split into multiple short message, with UDH
// according to 33GP TS 23.040 section 9.2.3.24.1
// NOTE: Split() will return array of length 1 if data length is still within the limit
// The encoding interface can implement the data.Splitter interface for ad-hoc splitting rule
func (c *ShortMessage) Split() (multiSM []*ShortMessage, err error) {
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

	// reserve 6 bytes for concat message UDH
	segments, err := splitter.EncodeSplit(c.message, data.SM_GSM_MSG_LEN-6)
	if err != nil {
		return nil, err
	}

	ref := getRefNum() // all segments will have the same ref id
	multiSM = []*ShortMessage{}
	for i, seg := range segments {
		// create new SM, encode data
		multiSM = append(multiSM, &ShortMessage{
			enc: c.enc,
			// message: we don't really care
			messageData:       seg,
			withoutDataCoding: c.withoutDataCoding,
			udHeader:          UDH{NewIEConcatMessage(len(segments), i+1, int(ref))},
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

	// Prepend UDH to messgae data if there are any
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

	// short message contain User-Data Header
	if udhi {
		udh := UDH{}
		_, err = udh.UnmarshalBinary(c.messageData)
		if err != nil {
			return
		}

		c.udHeader = udh
	}

	return
}

// Encoding returns message encoding.
func (c *ShortMessage) Encoding() data.Encoding {
	return c.enc
}

// String returns message content or content parse error
func (c *ShortMessage) String() string {
	message, err := c.GetMessage()
	if err != nil {
		return err.Error()
	}
	return message
}

// getRefNum return a atomically incrementing number each time it's called
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
