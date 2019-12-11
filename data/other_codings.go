package data

import (
	"golang.org/x/text/encoding/unicode"
)

var (
	// UTF16BEM is UTF-16 Big Endian with BOM (byte order mark).
	UTF16BEM = &utf16BEM{}

	// UTF16LEM is UTF-16 Little Endian with BOM.
	UTF16LEM = &utf16LEM{}

	// UTF16BE is UTF-16 Big Endian without BOM.
	UTF16BE = &utf16BE{}

	// UTF16LE is UTF-16 Little Endian without BOM.
	UTF16LE = &utf16LE{}
)

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
