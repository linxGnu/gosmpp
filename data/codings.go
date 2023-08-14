package data

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

const (
	// GSM7BITCoding is gsm-7bit coding
	GSM7BITCoding byte = 0x00
	// ASCIICoding is ascii coding
	ASCIICoding byte = 0x01
	// BINARY8BIT1Coding is 8-bit binary coding
	BINARY8BIT1Coding byte = 0x02
	// LATIN1Coding is iso-8859-1 coding
	LATIN1Coding byte = 0x03
	// BINARY8BIT2Coding is 8-bit binary coding
	BINARY8BIT2Coding byte = 0x04
	// CYRILLICCoding is iso-8859-5 coding
	CYRILLICCoding byte = 0x06
	// HEBREWCoding is iso-8859-8 coding
	HEBREWCoding byte = 0x07
	// UCS2Coding is UCS2 coding
	UCS2Coding byte = 0x08
)

// EncDec wraps encoder and decoder interface.
type EncDec interface {
	Encode(str string) ([]byte, error)
	Decode([]byte) (string, error)
}

// Encoding interface.
type Encoding interface {
	EncDec
	DataCoding() byte
}

func encode(str string, encoder *encoding.Encoder) ([]byte, error) {
	return encoder.Bytes([]byte(str))
}

func decode(data []byte, decoder *encoding.Decoder) (st string, err error) {
	tmp, err := decoder.Bytes(data)
	if err == nil {
		st = string(tmp)
	}
	return
}

// CustomEncoding is wrapper for user-defined data encoding.
type CustomEncoding struct {
	encDec EncDec
	coding byte
}

// NewCustomEncoding creates new custom encoding.
func NewCustomEncoding(coding byte, encDec EncDec) Encoding {
	return &CustomEncoding{
		coding: coding,
		encDec: encDec,
	}
}

// Encode string.
func (c *CustomEncoding) Encode(str string) ([]byte, error) {
	return c.encDec.Encode(str)
}

// Decode data to string.
func (c *CustomEncoding) Decode(data []byte) (string, error) {
	return c.encDec.Decode(data)
}

// DataCoding flag.
func (c *CustomEncoding) DataCoding() byte {
	return c.coding
}

type gsm7bit struct {
	packed bool
}

func (c *gsm7bit) Encode(str string) ([]byte, error) {
	return encode(str, GSM7(c.packed).NewEncoder())
}

func (c *gsm7bit) Decode(data []byte) (string, error) {
	return decode(data, GSM7(c.packed).NewDecoder())
}

func (c *gsm7bit) DataCoding() byte { return GSM7BITCoding }

func (c *gsm7bit) ShouldSplit(text string, octetLimit uint) (shouldSplit bool) {
	if c.packed {
		return uint((len(text)*7+7)/8) > octetLimit
	} else {
		return uint(len(text)) > octetLimit
	}
}

func (c *gsm7bit) EncodeSplit(text string, octetLimit uint) (allSeg [][]byte, err error) {
	if octetLimit < 64 {
		octetLimit = 134
	}

	allSeg = [][]byte{}
	runeSlice := []rune(text)

	fr, to := 0, int(octetLimit)
	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		seg, err := c.Encode(string(runeSlice[fr:to]))
		if err != nil {
			return nil, err
		}
		allSeg = append(allSeg, seg)
		fr, to = to, to+int(octetLimit)
	}

	return
}

type ascii struct{}

func (*ascii) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (*ascii) Decode(data []byte) (string, error) {
	return string(data), nil
}

func (*ascii) DataCoding() byte { return ASCIICoding }

type iso88591 struct{}

func (*iso88591) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_1.NewEncoder())
}

func (*iso88591) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_1.NewDecoder())
}

func (*iso88591) DataCoding() byte { return LATIN1Coding }

type binary8bit1 struct{}

func (*binary8bit1) Encode(_ string) ([]byte, error) {
	return []byte{}, ErrNotImplEncode
}

func (*binary8bit1) Decode(_ []byte) (string, error) {
	return "", ErrNotImplDecode
}

func (*binary8bit1) DataCoding() byte { return BINARY8BIT1Coding }

type binary8bit2 struct{}

func (*binary8bit2) Encode(_ string) ([]byte, error) {
	return []byte{}, ErrNotImplEncode
}

func (*binary8bit2) Decode(_ []byte) (string, error) {
	return "", ErrNotImplDecode
}

func (*binary8bit2) DataCoding() byte { return BINARY8BIT2Coding }

type iso88595 struct{}

func (*iso88595) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_5.NewEncoder())
}

func (*iso88595) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_5.NewDecoder())
}

func (*iso88595) DataCoding() byte { return CYRILLICCoding }

type iso88598 struct{}

func (*iso88598) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_8.NewEncoder())
}

func (*iso88598) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_8.NewDecoder())
}

func (*iso88598) DataCoding() byte { return HEBREWCoding }

type ucs2 struct{}

func (*ucs2) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return encode(str, tmp.NewEncoder())
}

func (*ucs2) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return decode(data, tmp.NewDecoder())
}

func (*ucs2) ShouldSplit(text string, octetLimit uint) (shouldSplit bool) {
	return uint(len(text)*2) > octetLimit
}

func (c *ucs2) EncodeSplit(text string, octetLimit uint) (allSeg [][]byte, err error) {
	if octetLimit < 64 {
		octetLimit = 134
	}

	allSeg = [][]byte{}
	runeSlice := []rune(text)
	hextetLim := int(octetLimit / 2) // round down

	// hextet = 16 bits, the correct terms should be hexadectet
	fr, to := 0, hextetLim
	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}

		seg, err := c.Encode(string(runeSlice[fr:to]))
		if err != nil {
			return nil, err
		}
		allSeg = append(allSeg, seg)

		fr, to = to, to+hextetLim
	}

	return
}

func (*ucs2) DataCoding() byte { return UCS2Coding }

var (
	// GSM7BIT is gsm-7bit encoding.
	GSM7BIT Encoding = &gsm7bit{packed: false}

	// GSM7BITPACKED is packed gsm-7bit encoding.
	// Most of SMSC(s) use unpack version.
	// Should be tested before using.
	GSM7BITPACKED Encoding = &gsm7bit{packed: true}

	// ASCII is ascii encoding.
	ASCII Encoding = &ascii{}

	// BINARY8BIT1 is binary 8-bit encoding.
	BINARY8BIT1 Encoding = &binary8bit1{}

	// LATIN1 encoding.
	LATIN1 Encoding = &iso88591{}

	// BINARY8BIT2 is binary 8-bit encoding.
	BINARY8BIT2 Encoding = &binary8bit2{}

	// CYRILLIC encoding.
	CYRILLIC Encoding = &iso88595{}

	// HEBREW encoding.
	HEBREW Encoding = &iso88598{}

	// UCS2 encoding.
	UCS2 Encoding = &ucs2{}
)

var codingMap = map[byte]Encoding{
	GSM7BITCoding:     GSM7BIT,
	ASCIICoding:       ASCII,
	BINARY8BIT1Coding: BINARY8BIT1,
	LATIN1Coding:      LATIN1,
	BINARY8BIT2Coding: BINARY8BIT2,
	CYRILLICCoding:    CYRILLIC,
	HEBREWCoding:      HEBREW,
	UCS2Coding:        UCS2,
}

// FromDataCoding returns encoding from DataCoding value.
func FromDataCoding(code byte) (enc Encoding) {
	enc = codingMap[code]
	return
}

// Splitter extend encoding object by defining a split function
// that split a string into multiple segments
// Each segment string, when encoded, must be within a certain octet limit
type Splitter interface {
	// ShouldSplit check if the encoded data of given text should be splitted under octetLimit
	ShouldSplit(text string, octetLimit uint) (should bool)
	EncodeSplit(text string, octetLimit uint) ([][]byte, error)
}
