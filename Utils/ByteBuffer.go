package Utils

import (
	"bytes"
	"fmt"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
)

const (
	SZ_BYTE  int = 1
	SZ_SHORT int = 2
	SZ_INT   int = 4
	SZ_LONG  int = 8
)

type ByteBuffer struct {
	*bytes.Buffer
}

// NewBuffer create new buffer from preallocated buffer array
func NewBuffer(inp []byte) *ByteBuffer {
	if inp == nil {
		return &ByteBuffer{bytes.NewBuffer([]byte{})}
	}

	return &ByteBuffer{bytes.NewBuffer(inp)}
}

func (c *ByteBuffer) read(n int) (r []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	if c.Buffer == nil {
		return nil, fmt.Errorf("Buffer not init!")
	}

	if n == 0 {
		return []byte{}, nil
	}

	if n < 0 {
		n = c.Buffer.Len()
	}

	res := make([]byte, n)
	num, e := c.Read(res)
	if e != nil {
		return nil, e
	}

	if num < n {
		return nil, fmt.Errorf("Buffer is not enough byte to read!")
	}

	return res, nil
}

// ReadBytes read n byte with checking error
func (c *ByteBuffer) Read_Bytes(n int) (*ByteBuffer, *Exception.Exception) {
	res, e := c.read(n)
	if e != nil {
		return nil, Exception.NewException(e)
	}

	return NewBuffer(res), nil
}

// Read_UnsafeByte read byte without knowing about error
func (c *ByteBuffer) Read_UnsafeByte() byte {
	if c.Buffer == nil {
		return 0
	}

	b, _ := c.ReadByte()

	return b
}

// Read_Byte read byte with checking error
func (c *ByteBuffer) Read_Byte() (b byte, ex *Exception.Exception) {
	if c.Buffer == nil {
		return 0, Exception.NewExceptionFromStr("Buffer not init!")
	}

	b, err := c.ReadByte()
	if err != nil {
		ex = Exception.NewException(err)
	} else {
		ex = nil
	}

	return
}

// Read_UnsafeShort read short without knowing about error
func (c *ByteBuffer) Read_UnsafeShort() int16 {
	if c.Buffer == nil {
		return 0
	}

	b, _ := c.Read_Short()

	return b
}

// Read_Short read short with checking error
func (c *ByteBuffer) Read_Short() (int16, *Exception.Exception) {
	if c.Buffer == nil {
		return 0, Exception.NewExceptionFromStr("Buffer not init!")
	}

	b, err := c.read(SZ_SHORT)
	if err != nil {
		return 0, Exception.NewException(err)
	}

	var result int16
	result |= int16(b[0] & 0xff)
	result <<= 8
	result |= int16(b[1] & 0xff)

	return result, nil
}

// Read_UnsafeInt read short without knowing about error
func (c *ByteBuffer) Read_UnsafeInt() int32 {
	if c.Buffer == nil {
		return 0
	}

	b, _ := c.Read_Int()

	return b
}

// Read_Int read int with checking error
func (c *ByteBuffer) Read_Int() (int32, *Exception.Exception) {
	if c.Buffer == nil {
		return 0, Exception.NewExceptionFromStr("Buffer not init!")
	}

	b, err := c.read(SZ_INT)
	if err != nil {
		return 0, Exception.NewException(err)
	}

	var result int32
	result |= int32(b[0] & 0xff)
	result <<= 8
	result |= int32(b[1] & 0xff)
	result <<= 8
	result |= int32(b[2] & 0xff)
	result <<= 8
	result |= int32(b[3] & 0xff)

	return result, nil
}

func (c *ByteBuffer) write(data []byte, n int) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	if c.Buffer == nil {
		c.Buffer = &bytes.Buffer{}
	}

	num, err := c.Write(data)
	if err != nil {
		return err
	}

	if num != n {
		return fmt.Errorf("Can not write more byte to buffer")
	}

	return nil
}

// Write_Bytes write array of byte to buffer
func (c *ByteBuffer) Write_Bytes(data []byte) *Exception.Exception {
	if data == nil {
		return Exception.NewExceptionFromStr("Data is nil for writing to buffer!")
	}

	e := c.write(data, len(data))
	if e != nil {
		return Exception.NewException(e)
	}

	return nil
}

// Write_UnsafeByte write a byte to buffer without checking error
func (c *ByteBuffer) Write_UnsafeByte(data byte) {
	if c.Buffer == nil {
		c.Buffer = &bytes.Buffer{}
	}

	c.WriteByte(data)
}

// Write_Byte write a byte to buffer with checking error
func (c *ByteBuffer) Write_Byte(data byte) *Exception.Exception {
	if c.Buffer == nil {
		c.Buffer = &bytes.Buffer{}
	}

	err := c.WriteByte(data)
	if err != nil {
		return Exception.NewException(err)
	}

	return nil
}

// Write_UnsafeShort write short without checking error
func (c *ByteBuffer) Write_UnsafeShort(data int16) {
	if c.Buffer == nil {
		c.Buffer = &bytes.Buffer{}
	}

	// buf[1] = byte(data & 0xff)
	// buf[0] = byte((data >> 8) & 0xff)
	c.write([]byte{byte((data >> 8) & 0xff), byte(data & 0xff)}, SZ_SHORT)
}

// Write_Short write short with checking error
func (c *ByteBuffer) Write_Short(data int16) *Exception.Exception {
	if c.Buffer == nil {
		c.Buffer = &bytes.Buffer{}
	}

	err := c.write([]byte{byte((data >> 8) & 0xff), byte(data & 0xff)}, SZ_SHORT)
	if err != nil {
		return Exception.NewException(err)
	}

	return nil
}

// Write_UnsafeInt write int without checking error
func (c *ByteBuffer) Write_UnsafeInt(data int32) {
	if c.Buffer == nil {
		c.Buffer = &bytes.Buffer{}
	}

	// buf[3] = byte(data & 0xff)
	// buf[2] = byte((data >> 8) & 0xff)
	// buf[1] = byte((data >> 16) & 0xff)
	// buf[0] = byte((data >> 24) & 0xff)
	c.write([]byte{byte((data >> 24) & 0xff), byte((data >> 16) & 0xff), byte((data >> 8) & 0xff), byte(data & 0xff)}, SZ_INT)
}

// Write_Int write int with checking error
func (c *ByteBuffer) Write_Int(data int32) *Exception.Exception {
	if c.Buffer == nil {
		c.Buffer = &bytes.Buffer{}
	}

	err := c.write([]byte{byte((data >> 24) & 0xff), byte((data >> 16) & 0xff), byte((data >> 8) & 0xff), byte(data & 0xff)}, SZ_INT)
	if err != nil {
		return Exception.NewException(err)
	}

	return nil
}

// Write_Buffer append buffer
func (c *ByteBuffer) Write_Buffer(d *ByteBuffer) *Exception.Exception {
	if c.Buffer == nil {
		c.Buffer = &bytes.Buffer{}
	}

	tmp := d.Bytes()
	err := c.write(tmp, len(tmp))
	if err != nil {
		return Exception.NewException(err)
	}

	return nil
}

func (c *ByteBuffer) write_String0(st string, isCString bool, enc Data.Encoding) (err *Exception.Exception) {
	defer func() {
		if e := recover(); e != nil {
			err = Exception.NewException(fmt.Errorf("%v", e))
		}
	}()

	if len(st) > 0 {
		var stringBuf []byte
		if enc == nil {
			stringBuf = []byte(st)
		} else {
			stB, err := enc.Encode(st)
			if err != nil {
				return Exception.UnsupportedEncodingException
			}
			stringBuf = stB
		}

		err := c.Write_Bytes(stringBuf)
		if err != nil {
			return err
		}
	}

	if isCString {
		return c.Write_Byte(0)
	}

	return nil
}

func (c *ByteBuffer) Write_CString(String string) *Exception.Exception {
	return c.write_String0(String, true, Data.ENC_ASCII)
}

func (c *ByteBuffer) Write_CStringWithEnc(String string, enc Data.Encoding) *Exception.Exception {
	return c.write_String0(String, true, enc)
}

func (c *ByteBuffer) Write_String(String string) *Exception.Exception {
	return c.Write_StringWithEnc(String, Data.ENC_ASCII)
}

func (c *ByteBuffer) Write_StringWithEnc(String string, enc Data.Encoding) *Exception.Exception {
	return c.write_String0(String, false, enc)
}

func (c *ByteBuffer) Write_StringWithCount(String string, count int) *Exception.Exception {
	return c.Write_StringWithCountEnc(String, count, Data.ENC_ASCII)
}

func (c *ByteBuffer) Write_StringWithCountEnc(String string, count int, enc Data.Encoding) *Exception.Exception {
	if count <= 0 {
		return nil
	}

	stLength := GetStringLength(String)
	if count > stLength {
		return Exception.WrongLengthOfStringException
	}
	String = Substring(String, count)

	return c.Write_StringWithEnc(String, enc)
}

// Read_CString read cstring with terminating zero
func (c *ByteBuffer) Read_CString() (string, *Exception.Exception) {
	if c.Buffer == nil {
		return "", nil
	}

	buf, err := c.ReadBytes(0)
	if err != nil {
		return "", Exception.TerminatingZeroNotFoundException
	}

	return string(buf[:len(buf)-1]), nil
}

// Read_String read string with limited size
func (c *ByteBuffer) Read_String(size int) (string, *Exception.Exception) {
	buf, err := c.read(size)
	if err != nil {
		return "", Exception.NewException(err)
	}

	return string(buf), nil
}

// GetHexDump get hex dump
func (c *ByteBuffer) GetHexDump() string {
	return fmt.Sprintf("%x", c.Buffer.Bytes())
}
