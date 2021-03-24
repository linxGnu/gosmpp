package data

// Source code in this file is copied from: https://github.com/fiorix/go-smpp/master/smpp/encoding/gsm7.go
import (
	"bytes"
	"math"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

/*
GSM 7-bit default alphabet and extension table
Source: https://en.wikipedia.org/wiki/GSM_03.38#GSM_7-bit_default_alphabet_and_extension_table_of_3GPP_TS_23.038_/_GSM_03.38
*/

// GSM7Custom returns a GSM 7-bit Bit Encoding with custom alphabets
func GSM7Custom(alphabet map[rune]byte, alphabetEx map[rune]byte, escape byte, packed bool) encoding.Encoding {
	e := gsm7Custom{
		alphabet:       make(map[rune]byte),
		alphabetEx:     make(map[rune]byte),
		escapeSequence: escape,
		packed:         packed,
	}
	for r, b := range alphabet {
		e.forward[r] = b
		e.reverse[b] = r
	}
	for r, b := range alphabetEx {
		e.forwardEsc[r] = b
		e.reverseEsc[b] = r
	}
	return nil
}

type gsm7Custom struct {
	alphabet map[rune]byte
	// alphabetEx the ex alphabet requires escape char being put in front,
	// so each character in this alphabet takes 2 bytes to encoded
	alphabetEx     map[rune]byte
	escapeSequence byte
	packed         bool

	// internal uses, this assume symmetric mapping
	forward    map[rune]byte
	reverse    map[byte]rune
	forwardEsc map[rune]byte
	reverseEsc map[byte]rune
}

func (g gsm7Custom) String() string {
	if g.packed {
		return "GSM 7-bit Customized Alphabet (Packed)"
	}
	return "GSM 7-bit Customized Alphabet (Unpacked)"
}

func (g *gsm7Custom) DecodeBytes(dst, src []byte) (nDst, nSrc int, err error) {
	if len(src) == 0 {
		return 0, 0, nil
	}

	septets := unpack(src, g.packed)

	nSeptet := 0
	builder := bytes.NewBufferString("")
	for nSeptet < len(septets) {
		b := septets[nSeptet]
		if b == g.escapeSequence {
			nSeptet++
			if nSeptet >= len(septets) {
				return 0, 0, ErrInvalidByte
			}
			e := septets[nSeptet]
			if r, ok := g.reverseEsc[e]; ok {
				builder.WriteRune(r)
			} else {
				return 0, 0, ErrInvalidByte
			}
		} else if r, ok := g.reverse[b]; ok {
			builder.WriteRune(r)
		} else {
			return 0, 0, ErrInvalidByte
		}
		nSeptet++
	}
	text := builder.Bytes()
	nDst = len(text)

	if len(dst) < nDst {
		return 0, 0, transform.ErrShortDst
	}

	copy(dst, text)
	return
}

func (g *gsm7Custom) EncodeBytes(dst, src []byte) (nDst, nSrc int, err error) {
	if len(src) == 0 {
		return 0, 0, nil
	}

	text := string(src) // work with []rune (a.k.a string) instead of []byte
	septets := make([]byte, 0, len(text))
	for _, r := range text {
		if v, ok := g.forward[r]; ok {
			septets = append(septets, v)
		} else if v, ok := g.forwardEsc[r]; ok {
			septets = append(septets, escapeSequence, v)
		} else {
			return 0, 0, ErrInvalidCharacter
		}
		nSrc++
	}

	nDst = len(septets)
	if g.packed {
		nDst = int(math.Ceil(float64(len(septets)) * 7 / 8))
	}
	if len(dst) < nDst {
		return 0, 0, transform.ErrShortDst
	}

	if !g.packed {
		copy(dst, septets)
		return nDst, nSrc, nil
	}

	nDst = pack(dst, septets)
	return
}

type GSM7BCustom struct {
	gsm7Custom
}

func (g *GSM7BCustom) Encode(w string) (dst []byte, err error) {
	_, _, err = g.EncodeBytes(dst, []byte(w))
	return
}

func (g *GSM7BCustom) Decode(src []byte) (w string, err error) {
	z := []byte{}
	_, _, err = g.DecodeBytes(z, src)
	return string(z), err
}
