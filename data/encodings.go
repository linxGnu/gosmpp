package data

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

// Encoding interface.
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

type gsm7bit struct{}

func (c gsm7bit) Encode(str string) ([]byte, error) {
	return Encode7Bit(str), nil
}

func (c gsm7bit) Decode(data []byte) (string, error) {
	return Decode7Bit(data)
}

type ascii struct{}

func (c ascii) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (c ascii) Decode(data []byte) (string, error) {
	return string(data), nil
}

type utf8 struct{}

func (c utf8) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (c utf8) Decode(data []byte) (string, error) {
	return string(data), nil
}

type cp1252 struct{}

func (c cp1252) Encode(str string) ([]byte, error) {
	return encode(str, charmap.Windows1252.NewEncoder())
}

func (c cp1252) Decode(data []byte) (string, error) {
	return decode(data, charmap.Windows1252.NewDecoder())
}

type iso8859 struct{}

func (c iso8859) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_1.NewEncoder())
}

func (c iso8859) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_1.NewDecoder())
}

type utf16BEM struct{}

func (c utf16BEM) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	return encode(str, tmp.NewEncoder())
}

func (c utf16BEM) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	return decode(data, tmp.NewDecoder())
}

type utf16LEM struct{}

func (c utf16LEM) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	return encode(str, tmp.NewEncoder())
}

func (c utf16LEM) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	return decode(data, tmp.NewDecoder())
}

type utf16BE struct{}

func (c utf16BE) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return encode(str, tmp.NewEncoder())
}

func (c utf16BE) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return decode(data, tmp.NewDecoder())
}

type utf16LE struct{}

func (c utf16LE) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return encode(str, tmp.NewEncoder())
}

func (c utf16LE) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return decode(data, tmp.NewDecoder())
}

type utf16 struct{}

func (c utf16) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
	return encode(str, tmp.NewEncoder())
}

func (c utf16) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
	return decode(data, tmp.NewDecoder())
}

var (
	// ASCII is ascii encoding.
	ASCII Encoding = &ascii{}

	// GSM7BIT is gsm-7bit encoding.
	GSM7BIT Encoding = &gsm7bit{}

	// CHAR encoding.
	CHAR Encoding = ASCII

	// CP1252 is Windows Latin-1
	CP1252 Encoding = &cp1252{}

	// UTF8 encoding.
	UTF8 Encoding = &utf8{}

	// ISO8859 encoding.
	ISO8859 Encoding = &iso8859{}

	// UTF16BEM is UTF-16 Big Endian with BOM (byte order mark).
	UTF16BEM Encoding = &utf16BEM{}

	// UTF16LEM is UTF-16 Little Endian with BOM.
	UTF16LEM Encoding = &utf16LEM{}

	// UTF16BE is UTF-16 Big Endian without BOM.
	UTF16BE Encoding = &utf16BE{}

	// UTF16LE is UTF-16 Little Endian without BOM.
	UTF16LE Encoding = &utf16LE{}

	// UTF16 encoding.
	UTF16 = &utf16{}
)
