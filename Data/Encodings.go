package Data

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

type Encoding interface {
	Encode(str string) ([]byte, error)
	Decode([]byte) (string, error)
}

func encode(str string, encoder *encoding.Encoder) ([]byte, error) {
	return encoder.Bytes([]byte(str))
}

func decode(data []byte, decoder *encoding.Decoder) (string, error) {
	tmp, err := decoder.Bytes(data)
	if err != nil {
		return "", err
	}

	return string(tmp), nil
}

// ENC_GSM7BIT_s ...
type ENC_GSM7BIT_s struct {
}

func (c ENC_GSM7BIT_s) Encode(str string) ([]byte, error) {
	return Encode7Bit(str), nil
}

func (c ENC_GSM7BIT_s) Decode(data []byte) (string, error) {
	return Decode7Bit(data)
}

// ENC_ASCII_s ..
type ENC_ASCII_s struct {
}

func (c ENC_ASCII_s) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (c ENC_ASCII_s) Decode(data []byte) (string, error) {
	return string(data), nil
}

// ENC_UTF8_s ..
type ENC_UTF8_s struct {
}

func (c ENC_UTF8_s) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (c ENC_UTF8_s) Decode(data []byte) (string, error) {
	return string(data), nil
}

// ENC_CP1252_s ...
type ENC_CP1252_s struct{}

func (c ENC_CP1252_s) Encode(str string) ([]byte, error) {
	return encode(str, charmap.Windows1252.NewEncoder())
}

func (c ENC_CP1252_s) Decode(data []byte) (string, error) {
	return decode(data, charmap.Windows1252.NewDecoder())
}

// ENC_ISO8859_1_s ...
type ENC_ISO8859_1_s struct{}

func (c ENC_ISO8859_1_s) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_1.NewEncoder())
}

func (c ENC_ISO8859_1_s) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_1.NewDecoder())
}

// ENC_UTF16_BEM_s ...
type ENC_UTF16_BEM_s struct{}

func (c ENC_UTF16_BEM_s) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	return encode(str, tmp.NewEncoder())
}

func (c ENC_UTF16_BEM_s) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	return decode(data, tmp.NewDecoder())
}

// ENC_UTF16_LEM_s ...
type ENC_UTF16_LEM_s struct{}

func (c ENC_UTF16_LEM_s) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	return encode(str, tmp.NewEncoder())
}

func (c ENC_UTF16_LEM_s) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	return decode(data, tmp.NewDecoder())
}

// ENC_UTF16_BE_s ...
type ENC_UTF16_BE_s struct{}

func (c ENC_UTF16_BE_s) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return encode(str, tmp.NewEncoder())
}

func (c ENC_UTF16_BE_s) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return decode(data, tmp.NewDecoder())
}

// ENC_UTF16_LE_s ...
type ENC_UTF16_LE_s struct{}

func (c ENC_UTF16_LE_s) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return encode(str, tmp.NewEncoder())
}

func (c ENC_UTF16_LE_s) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return decode(data, tmp.NewDecoder())
}

// ENC_UTF16_s ...
type ENC_UTF16_s struct{}

func (c ENC_UTF16_s) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
	return encode(str, tmp.NewEncoder())
}

func (c ENC_UTF16_s) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
	return decode(data, tmp.NewDecoder())
}

// Ascii
var ENC_ASCII Encoding = ENC_ASCII_s{}

// Ascii
var CHAR_ENC Encoding = ENC_ASCII

// Windows Latin-1
var ENC_CP1252 Encoding = ENC_CP1252_s{}

// GSM 7-bit unpacked
var ENC_GSM7BIT Encoding = ENC_GSM7BIT_s{}

// Eight-bit Unicode Transformation Format
var ENC_UTF8 Encoding = ENC_UTF8_s{}

// ISO 8859-1, Latin alphabet No. 1
var ENC_ISO8859_1 Encoding = ENC_ISO8859_1_s{}

var ENC_UTF16_BEM Encoding = ENC_UTF16_BEM_s{}

var ENC_UTF16_BE Encoding = ENC_UTF16_BE_s{}

var ENC_UTF16_LEM Encoding = ENC_UTF16_LEM_s{}

var ENC_UTF16_LE Encoding = ENC_UTF16_LE_s{}

var ENC_UTF16 Encoding = ENC_UTF16_s{}
