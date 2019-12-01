package pdu

import (
	"github.com/linxGnu/gosmpp/data"
)

// ShortMessage message.
type ShortMessage struct {
	Message     string
	Enc         data.Encoding
	MessageData []byte
}

// NewShortMessage returns new ShortMessage.
func NewShortMessage(message string) (s ShortMessage, err error) {
	err = s.SetMessage(message)
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
	if c.MessageData, err = enc.Encode(message); err == nil {
		c.Message = message
		c.Enc = enc
	}
	return
}

// GetMessage returns underlying message.
func (c *ShortMessage) GetMessage() (string, error) {
	enc := c.Enc
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
	if len(c.MessageData) == 0 {
		return "", nil
	}

	if enc == nil {
		return string(c.MessageData), nil
	}

	return enc.Decode(c.MessageData)
}

// func (c *ShortMessage) SetData(buf *Utils.ByteBuffer) *Exception.Exception {
// 	if buf == nil || buf.Buffer == nil {
// 		c.messageData = nil
// 		c.length = 0
// 		c.message = ""
// 		return nil
// 	}

// 	c.messageData = buf.Bytes()
// 	if c.messageData == nil {
// 		c.length = 0
// 	} else {
// 		c.length = int32(len(c.messageData))
// 	}

// 	if c.length < c.minLength || c.length > c.maxLength {
// 		return Exception.WrongLengthException
// 	}

// 	return nil
// }

// func (c *ShortMessage) GetData() (*Utils.ByteBuffer, *Exception.Exception) {
// 	return Utils.NewBuffer(c.messageData), nil
// }

// func (c *ShortMessage) SetEncoding(enc Data.Encoding) *Exception.Exception {
// 	if enc == nil {
// 		return Exception.NewExceptionFromStr("ShortMessage: encoding object is nil")
// 	}

// 	tmp, err1 := enc.Decode(c.messageData)
// 	if err1 != nil {
// 		return Exception.UnsupportedEncodingException
// 	}

// 	c.message = tmp
// 	c.enc = enc
// 	return nil
// }
