package test

import (
	"bytes"
	"testing"

	"github.com/linxGnu/gosmpp/Utils"
)

const (
	ABC     string = "ABC"
	ASCII   string = "ASCII"
	INVALID string = "INVALID"
	NULL    byte   = 0x00
	A       byte   = 0x41
	B       byte   = 0x42
	C       byte   = 0x43
)

var t_bite byte = 0x1f
var t_short int16 = 666

func TestAppendByte0(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Byte(t_bite)

	expected := []byte{t_bite}
	if bytes.Compare(buffer.Bytes(), expected) != 0 {
		t.Fail()
	}
}

func TestAppendShort0(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Short(t_short)

	expected := []byte{0x02, 0x9a}
	if bytes.Compare(buffer.Bytes(), expected) != 0 {
		t.Fail()
	}
}

func TestAppendInt0(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Int(666)

	expected := []byte{NULL, NULL, 0x02, 0x9a}
	if bytes.Compare(buffer.Bytes(), expected) != 0 {
		t.Fail()
	}
}

func TestAppendCString0(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_CString(ABC)

	expected := []byte{A, B, C, NULL}
	if bytes.Compare(buffer.Bytes(), expected) != 0 {
		t.Fail()
	}
}

func TestAppendString(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_String(ABC)

	expected := []byte{A, B, C}
	if bytes.Compare(buffer.Bytes(), expected) != 0 {
		t.Fail()
	}
}

func TestAppendCStringWithNULL(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_CString("")

	expected := []byte{NULL}
	if bytes.Compare(buffer.Bytes(), expected) != 0 {
		t.Fail()
	}
}

func TestAppendBytesWithNull(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	err := buffer.Write_Bytes(nil)

	if err == nil {
		t.Fail()
	}
}

func TestAppendBytes(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{A, B, C})

	expected := []byte{A, B, C}
	if bytes.Compare(buffer.Bytes(), expected) != 0 {
		t.Fail()
	}
}

func TestRemoveByte(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{A, B, C})

	bite, err := buffer.Read_Byte()
	if err != nil || bite != A || buffer.Len() != 2 {
		t.Fail()
	}
}

func TestRemoveByteNegative(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{A, B, C})

	_, err := buffer.Read_Bytes(-10)
	if err != nil || buffer.Len() != 0 {
		t.Fail()
	}
}

func TestRemoveByteNotEnough(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{A})

	_, err := buffer.Read_Short()
	if err == nil {
		t.Fail()
	}
}

func TestRemoveShort(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{0x01, 0x02, 0x03, 0x04})

	res, err := buffer.Read_Short()
	if err != nil || res != (1<<8)+2 || buffer.Len() != 2 {
		t.Fail()
	}
}

func TestRemoveInt(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{0x01, 0x02, 0x03, 0x04})

	res, err := buffer.Read_Int()
	if err != nil || res != (1<<24)+(2<<16)+(3<<8)+4 || buffer.Len() != 0 {
		t.Fail()
	}
}

func TestReadInt(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{0x01, 0x02, 0x03, 0x04, 0x05})

	res, err := buffer.Read_Int()
	if err != nil || res != (1<<24)+(2<<16)+(3<<8)+4 || buffer.Len() != 1 {
		t.Fail()
	}
}

func TestRemoveCStringWithoutTerminator(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{0x01, 0x02, 0x03, 0x04, 0x05})

	_, err := buffer.Read_CString()
	if err == nil {
		t.Fail()
	}
}

func TestRemoveCString(t *testing.T) {
	buffer := Utils.NewBufferDefault()
	buffer.Write_Bytes([]byte{A, B, NULL, C, NULL})

	res, err := buffer.Read_CString()
	if err != nil || "AB" != res || buffer.Len() != 2 {
		t.Fail()
	}
}
