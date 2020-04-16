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
	// LATIN1Coding is iso-8859-1 coding
	LATIN1Coding byte = 0x03
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

type gsm7bit struct {
	packed bool
}

func (c gsm7bit) Encode(str string) ([]byte, error) {
	return encode(str, GSM7(c.packed).NewEncoder())
}

func (c gsm7bit) Decode(data []byte) (string, error) {
	return decode(data, GSM7(c.packed).NewDecoder())
}

func (c gsm7bit) DataCoding() byte { return GSM7BITCoding }

func (c gsm7bit) ShouldSplit(text string, octetLimit uint) (shouldSplit bool) {
	return uint(len(text)) > octetLimit
}

func (c gsm7bit) EncodeSplit(text string, octetLimit uint) (allSeg [][]byte, err error) {
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

func (c ascii) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (c ascii) Decode(data []byte) (string, error) {
	return string(data), nil
}

func (c ascii) DataCoding() byte { return ASCIICoding }

type iso88591 struct{}

func (c iso88591) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_1.NewEncoder())
}

func (c iso88591) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_1.NewDecoder())
}

func (c iso88591) DataCoding() byte { return LATIN1Coding }

type iso88595 struct{}

func (c iso88595) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_5.NewEncoder())
}

func (c iso88595) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_5.NewDecoder())
}

func (c iso88595) DataCoding() byte { return CYRILLICCoding }

type iso88598 struct{}

func (c iso88598) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_8.NewEncoder())
}

func (c iso88598) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_8.NewDecoder())
}

func (c iso88598) DataCoding() byte { return HEBREWCoding }

type ucs2 struct{}

func (c ucs2) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return encode(str, tmp.NewEncoder())
}

func (c ucs2) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return decode(data, tmp.NewDecoder())
}

func (c ucs2) ShouldSplit(text string, octetLimit uint) (shouldSplit bool) {
	return uint(len(text)*2) > octetLimit
}

func (c ucs2) EncodeSplit(text string, octetLimit uint) (allSeg [][]byte, err error) {
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

func (c ucs2) DataCoding() byte { return UCS2Coding }

var (
	// GSM7BIT is gsm-7bit encoding.
	GSM7BIT Encoding = &gsm7bit{packed: false}

	// GSM7BITPACKED is packed gsm-7bit encoding.
	// Most of SMSC(s) use unpack version.
	// Should be tested before using.
	GSM7BITPACKED Encoding = &gsm7bit{packed: true}

	// ASCII is ascii encoding.
	ASCII Encoding = &ascii{}

	// LATIN1 encoding.
	LATIN1 Encoding = &iso88591{}

	// CYRILLIC encoding.
	CYRILLIC Encoding = &iso88595{}

	// HEBREW encoding.
	HEBREW Encoding = &iso88598{}

	// UCS2 encoding.
	UCS2 Encoding = &ucs2{}
)

var codingMap = map[byte]Encoding{
	GSM7BITCoding:  GSM7BIT,
	ASCIICoding:    ASCII,
	LATIN1Coding:   LATIN1,
	CYRILLICCoding: CYRILLIC,
	HEBREWCoding:   HEBREW,
	UCS2Coding:     UCS2,
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
