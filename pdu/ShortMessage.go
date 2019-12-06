package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"
	"github.com/linxGnu/gosmpp/utils"
)

// ShortMessage message.
type ShortMessage struct {
	SmDefaultMsgID byte
	message        string
	enc            data.Encoding
	messageData    []byte
}

// NewShortMessage returns new ShortMessage.
func NewShortMessage(message string) (s ShortMessage, err error) {
	err = s.SetMessage(message)
	return
}

// NewShortMessageWithEncoding returns new ShortMessage with predefined encoding.
func NewShortMessageWithEncoding(message string, enc data.Encoding) (s ShortMessage, err error) {
	err = s.SetMessageWithEncoding(message, enc)
	return
}

// SetMessage set message and its encoded data.
func (c *ShortMessage) SetMessage(message string) (err error) {
	err = c.SetMessageWithEncoding(message, data.GSM7BIT)
	if err != nil {
		err = c.SetMessageWithEncoding(message, data.ASCII)
	}
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

// GetMessage returns underlying message.
func (c *ShortMessage) GetMessage() (string, error) {
	enc := c.enc
	if enc == nil {
		enc = data.GSM7BIT
	}

	tmp, err := c.GetMessageWithEncoding(enc)
	if err != nil {
		return c.GetMessageWithEncoding(data.ASCII)
	}

	return tmp, err
}

// GetMessageWithEncoding returns (decoded) underlying message.
func (c *ShortMessage) GetMessageWithEncoding(enc data.Encoding) (string, error) {
	if len(c.messageData) == 0 {
		return "", nil
	}

	if enc == nil {
		return string(c.messageData), nil
	}

	return enc.Decode(c.messageData)
}

// Marshal implements PDU interface.
func (c *ShortMessage) Marshal(b *utils.ByteBuffer) {
	n := byte(len(c.messageData))
	b.Grow(int(n) + 3)

	var coding byte
	if c.enc == nil {
		coding = data.GSM7BITCoding
	} else {
		coding = c.enc.DataCoding()
	}
	_ = b.WriteByte(coding)

	_ = b.WriteByte(c.SmDefaultMsgID)

	_ = b.WriteByte(n)
	_, _ = b.Write(c.messageData[:n])
}

// Unmarshal implements PDU interface.
func (c *ShortMessage) Unmarshal(b *utils.ByteBuffer) (err error) {
	var dataCoding, n byte
	if dataCoding, err = b.ReadByte(); err == nil {
		if c.SmDefaultMsgID, err = b.ReadByte(); err == nil {
			if n, err = b.ReadByte(); err == nil {
				if c.messageData, err = b.ReadN(int(n)); err == nil {
					c.enc = data.FromDataCoding(dataCoding)
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
