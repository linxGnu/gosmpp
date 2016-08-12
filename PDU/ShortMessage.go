package PDU

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU/Common"
	"github.com/linxGnu/gosmpp/Utils"
)

type ShortMessage struct {
	Common.ByteData
	minLength   int32
	maxLength   int32
	message     string
	enc         Data.Encoding
	length      int32
	messageData []byte
}

func NewShortMessageWithMaxLength(maxLength int32) *ShortMessage {
	a := &ShortMessage{}
	a.Construct()

	a.maxLength = maxLength

	return a
}

func NewShortMessageWithMinMaxLength(minLength, maxLength int32) *ShortMessage {
	a := &ShortMessage{}
	a.Construct()

	a.minLength = minLength
	a.maxLength = maxLength

	return a
}

func (c *ShortMessage) Construct() {
	defer c.SetRealReference(c)
	c.ByteData.Construct()
}

func (c *ShortMessage) SetData(buf *Utils.ByteBuffer) *Exception.Exception {
	if buf == nil || buf.Buffer == nil {
		c.messageData = nil
		c.length = 0
		c.message = ""
		return nil
	}

	c.messageData = buf.Bytes()
	if c.messageData == nil {
		c.length = 0
	} else {
		c.length = int32(len(c.messageData))
	}

	if c.length < c.minLength || c.length > c.maxLength {
		return Exception.WrongLengthException
	}

	return nil
}

func (c *ShortMessage) GetData() (*Utils.ByteBuffer, *Exception.Exception) {
	return Utils.NewBuffer(c.messageData), nil
}

func (c *ShortMessage) SetMessage(message string) *Exception.Exception {
	err := c.SetMessageWithEncoding(message, Data.ENC_GSM7BIT)
	if err != nil {
		err = c.SetMessageWithEncoding(message, Data.ENC_ASCII)
	}

	return err
}

func (c *ShortMessage) SetMessageWithEncoding(message string, enc Data.Encoding) *Exception.Exception {
	if enc == nil {
		return Exception.NewExceptionFromStr("ShortMessage: encoding object is nil")
	}

	err := c.CheckStringMinMaxEncoding(message, int(c.minLength), int(c.maxLength), enc)
	if err != nil {
		return err
	}

	tmp, err1 := enc.Encode(message)
	if err1 != nil {
		return Exception.UnsupportedEncodingException
	}

	c.messageData = tmp
	c.message = message
	c.enc = enc
	c.length = int32(len(c.messageData))

	return nil
}

func (c *ShortMessage) SetEncoding(enc Data.Encoding) *Exception.Exception {
	if enc == nil {
		return Exception.NewExceptionFromStr("ShortMessage: encoding object is nil")
	}

	tmp, err1 := enc.Decode(c.messageData)
	if err1 != nil {
		return Exception.UnsupportedEncodingException
	}

	c.message = tmp
	c.enc = enc
	return nil
}

func (c *ShortMessage) GetMessage() (string, *Exception.Exception) {
	useEncoding := c.enc
	if c.enc == nil {
		useEncoding = Data.ENC_GSM7BIT
	}

	tmp, err := c.GetMessageWithEncoding(useEncoding)
	if err != nil {
		return c.GetMessageWithEncoding(Data.ENC_ASCII)
	}

	return tmp, err
}

func (c *ShortMessage) GetMessageWithEncoding(enc Data.Encoding) (string, *Exception.Exception) {
	if c.messageData == nil {
		return "", nil
	}

	if enc != nil && c.enc != nil && c.enc == enc {
		tmp, err1 := enc.Decode(c.messageData)
		if err1 != nil {
			return "", Exception.UnsupportedEncodingException
		}

		c.message = tmp
		return tmp, nil
	} else if enc != nil {
		t1, err1 := enc.Decode(c.messageData)
		if err1 != nil {
			return "", Exception.UnsupportedEncodingException
		}

		return t1, nil
	} else {
		return string(c.messageData), nil
	}
}

func (c *ShortMessage) GetLength() int32 {
	return c.length
}

func (c *ShortMessage) GetEncoding() Data.Encoding {
	return c.enc
}
