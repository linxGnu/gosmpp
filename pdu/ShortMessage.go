package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"
)

// ShortMessage message.
type ShortMessage struct {
	SmDefaultMsgID    byte
	message           string
	enc               data.Encoding
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
	if len(c.messageData) > 0 {
		st, err = enc.Decode(c.messageData)
	}
	return
}

// Marshal implements PDU interface.
func (c *ShortMessage) Marshal(b *ByteBuffer) {
	n := byte(len(c.messageData))
	b.Grow(int(n) + 3)

	var coding byte
	if c.enc == nil {
		coding = data.GSM7BITCoding
	} else {
		coding = c.enc.DataCoding()
	}

	if !c.withoutDataCoding {
		_ = b.WriteByte(coding)
	}

	_ = b.WriteByte(c.SmDefaultMsgID)

	_ = b.WriteByte(n)
	_, _ = b.Write(c.messageData[:n])
}

// Unmarshal implements PDU interface.
func (c *ShortMessage) Unmarshal(b *ByteBuffer) (err error) {
	var dataCoding, n byte
	if !c.withoutDataCoding {
		if dataCoding, err = b.ReadByte(); err == nil {
			if c.SmDefaultMsgID, err = b.ReadByte(); err == nil {
				if n, err = b.ReadByte(); err == nil {
					if c.messageData, err = b.ReadN(int(n)); err == nil {
						c.enc = data.FromDataCoding(dataCoding)
					}
				}
			}
		}
	} else {
		if c.SmDefaultMsgID, err = b.ReadByte(); err == nil {
			if n, err = b.ReadByte(); err == nil {
				if c.messageData, err = b.ReadN(int(n)); err == nil {
					c.enc = data.FromDataCoding(0)
				}
			}
		}
	}
	return
}

// Encoding returns message encoding.
func (c *ShortMessage) Encoding() data.Encoding {
	return c.enc
}
