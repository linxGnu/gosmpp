package test

import (
	"math/rand"
	"testing"

	"github.com/tsocial/gosmpp/Data"
	"github.com/tsocial/gosmpp/PDU"
)

var TONS []byte = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
var NPIS []byte = []byte{0x00, 0x01, 0x03, 0x04, 0x06, 0x08, 0x09, 0x0a, 0x0e, 0x0f, 0x12}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestDefaultConstructorInitialValues(t *testing.T) {
	address := PDU.NewAddress()
	if address.GetTon() != 0x00 || address.GetNpi() != 0x00 || address.GetAddress() != "" {
		t.Fail()
	}
}

func TestDefaultMaxAddressLengthAllpws20Digits(t *testing.T) {
	addr := PDU.NewAddressWithMaxLength(20)

	tmp := RandStringRunes(19)
	addr.SetAddress(tmp)

	if addr.GetAddress() != tmp {
		t.Fail()
	}
}

func TestDefaultMaxAddressLengthDissalows21Digits(t *testing.T) {
	addr := PDU.NewAddress()
	err := addr.SetAddress(RandStringRunes(21))
	if err == nil {
		t.Fail()
	}
}

func TestGetAddressWithEncoding(t *testing.T) {
	address := PDU.NewAddress()
	address.SetAddress("ABCD")

	enc := Data.ENC_UTF16_BE
	st, err := address.GetAddressWithEncoding(enc)

	if err != nil || st != "\x41\x42\x43\x44" {
		t.Fail()
	}
}

func TestGetAddressRangeWithEncoding(t *testing.T) {
	address := PDU.NewAddressRange()
	address.SetAddressRange("ABCDE")

	enc := Data.ENC_UTF16_BE
	st, err := address.GetAddressRangeWithEncoding(enc)

	if err != nil || st != "\x41\x42\x43\x44\x45" {
		t.Fail()
	}
}

func TestGetData(t *testing.T) {
	for _, ton := range TONS {
		for _, npi := range NPIS {
			for len := 1; len <= 0x21; len++ {
				tmp := RandStringRunes(len)
				addres, err := PDU.NewAddressWithTonNpiAddrMaxLen(ton, npi, tmp, int32(len+1))
				if err != nil {
					t.Fail()
				}

				buffer, err := addres.GetData()
				if err != nil {
					t.Fail()
				}

				dat, err := buffer.Read_Byte()
				if err != nil {
					t.Fail()
				}
				if dat != ton {
					t.Fail()
				}

				dat, err = buffer.Read_Byte()
				if err != nil {
					t.Fail()
				}
				if dat != npi {
					t.Fail()
				}

				dd, err := buffer.Read_CString()
				if err != nil {
					t.Fail()
				}
				if dd != tmp {
					t.Fail()
				}
			}
		}
	}
}
