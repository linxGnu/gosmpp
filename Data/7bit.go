package Data

// Below source code belongs to github.com/xlab/at/pdu/7bit.go
import (
	"bytes"
	"errors"
	"fmt"
)

// Esc code.
const Esc byte = 0x1B

// Ctrl+Z code.
const Sub byte = 0x1A

// <CR> code.
const CR byte = 0x0D

var crcr = []byte{CR, CR}
var cr = []byte{CR}

const (
	max     byte = 0x7F
	unknown rune = '?'
)

// ErrUnexpectedByte happens when someone tries to decode non GSM 7-bit encoded string.
var ErrUnexpectedByte = errors.New("7bit decode: met an unexpected byte")

// Encode7Bit encodes the given UTF-8 text into GSM 7-bit (3GPP TS 23.038) encoding with packing.
func Encode7Bit(str string) []byte {
	raw7 := make([]byte, 0, len(str))
	for _, r := range str {
		i := gsmTable.Index(r)
		if i < 0 {
			b := gsmEscapes.to7Bit(r)
			if b != byte(unknown) {
				raw7 = append(raw7, Esc, b)
			} else {
				raw7 = append(raw7, b)
			}
			continue
		}
		raw7 = append(raw7, byte(i))
	}
	return pack7Bit(raw7)
}

// Decode7Bit decodes the given GSM 7-bit packed octet data (3GPP TS 23.038) into a UTF-8 encoded string.
func Decode7Bit(octets []byte) (str string, err error) {
	raw7 := unpack7Bit(octets)
	var escaped bool
	var r rune
	for _, b := range raw7 {
		if b > max {
			err = ErrUnexpectedByte
			return
		} else if escaped {
			r = gsmEscapes.from7Bit(b)
			escaped = false
		} else if b == Esc {
			escaped = true
			continue
		} else {
			r = gsmTable.Rune(int(b))
		}
		str += string(r)
	}
	return
}

func pad(n, block int) int {
	if n%block == 0 {
		return n
	}
	return (n/block + 1) * block
}

func blocks(n, block int) int {
	if n%block == 0 {
		return n / block
	}
	return n/block + 1
}

func pack7Bit(raw7 []byte) []byte {
	pack7 := make([]byte, blocks(len(raw7)*7, 8))
	pack := func(out []byte, b byte, oct int, bit uint8) (int, uint8) {
		for i := uint8(0); i < 7; i++ {
			out[oct] |= b >> i & 1 << bit
			bit++
			if bit == 8 {
				oct++
				bit = 0
			}
		}
		return oct, bit
	}
	var oct int   // current octet in pack7
	var bit uint8 // current bit in octet
	var b byte    // current byte in raw7
	for i := range raw7 {
		b = raw7[i]
		oct, bit = pack(pack7, b, oct, bit)
	}
	// N.B. in order to not confuse 7 zero-bits with @
	// <CR> code is added to the packed bits.
	if 8-bit == 7 {
		oct, bit = pack(pack7, CR, oct, bit)
	} else if bit == 0 && b == CR {
		// and if data ends with <CR> on the octet boundary,
		// then we add an additional octet with <CR>. See (3GPP TS 23.038).
		pack7 = append(pack7, 0x00)
		oct, bit = pack(pack7, CR, oct, bit)
	}
	return pack7
}

func unpack7Bit(pack7 []byte) []byte {
	raw7 := make([]byte, 0, len(pack7))
	var sep byte  // current septet
	var bit uint8 // current bit in septet
	for _, oct := range pack7 {
		for i := uint8(0); i < 8; i++ {
			sep |= oct >> i & 1 << bit
			bit++
			if bit == 7 {
				raw7 = append(raw7, sep)
				sep = 0
				bit = 0
			}
		}
	}
	if bytes.HasSuffix(raw7, crcr) || bytes.HasSuffix(raw7, cr) {
		raw7 = raw7[:len(raw7)-1]
	}
	return raw7
}

func displayPack(buf []byte) (out string) {
	for i := 0; i < len(buf)*8; i++ {
		b := buf[i/8]
		if i%8 == 0 {
			out += fmt.Sprintf("\n%02X:", b)
		}
		off := 7 - uint8(i%8)
		out += fmt.Sprintf("%4d", b>>off&1)
	}
	return out
}

type escape struct {
	from byte
	to   rune
}

type escapeTable [0x0A]escape

func (et *escapeTable) to7Bit(r rune) byte {
	for _, esc := range et {
		if esc.to == r {
			return esc.from
		}
	}
	return byte(unknown)
}

func (et *escapeTable) from7Bit(b byte) rune {
	for _, esc := range et {
		if esc.from == b {
			return esc.to
		}
	}
	return unknown
}

var gsmEscapes = escapeTable{
	{0x0A, 0x000C}, /* FORM FEED */
	{0x14, 0x005E}, /* CIRCUMFLEX ACCENT */
	{0x28, 0x007B}, /* LEFT CURLY BRACKET */
	{0x29, 0x007D}, /* RIGHT CURLY BRACKET */
	{0x2F, 0x005C}, /* REVERSE SOLIDUS */
	{0x3C, 0x005B}, /* LEFT SQUARE BRACKET */
	{0x3D, 0x007E}, /* TILDE */
	{0x3E, 0x005D}, /* RIGHT SQUARE BRACKET */
	{0x40, 0x007C}, /* VERTICAL LINE */
	{0x65, 0x20AC}, /* EURO SIGN */
}

type runeTable [0x80]rune

func (rt *runeTable) Index(r rune) int {
	for i := range rt {
		if rt[i] == r {
			return i
		}
	}
	return -1
}

func (rt *runeTable) Rune(idx int) rune {
	if idx >= 0 && idx < len(rt) {
		return rt[idx]
	}
	return unknown
}

// Thanks, Jeroen @ Mobile Tidings
var gsmTable = runeTable{
	/* 0x00 */ 0x0040, /* COMMERCIAL AT */
	/* 0x01 */ 0x00A3, /* POUND SIGN */
	/* 0x02 */ 0x0024, /* DOLLAR SIGN */
	/* 0x03 */ 0x00A5, /* YEN SIGN */
	/* 0x04 */ 0x00E8, /* LATIN SMALL LETTER E WITH GRAVE */
	/* 0x05 */ 0x00E9, /* LATIN SMALL LETTER E WITH ACUTE */
	/* 0x06 */ 0x00F9, /* LATIN SMALL LETTER U WITH GRAVE */
	/* 0x07 */ 0x00EC, /* LATIN SMALL LETTER I WITH GRAVE */
	/* 0x08 */ 0x00F2, /* LATIN SMALL LETTER O WITH GRAVE */
	/* 0x09 */ 0x00E7, /* LATIN SMALL LETTER C WITH CEDILLA */
	/* 0x0A */ 0x000A, /* LINE FEED */
	/* 0x0B */ 0x00D8, /* LATIN CAPITAL LETTER O WITH STROKE */
	/* 0x0C */ 0x00F8, /* LATIN SMALL LETTER O WITH STROKE */
	/* 0x0D */ 0x000D, /* CARRIAGE RETURN */
	/* 0x0E */ 0x00C5, /* LATIN CAPITAL LETTER A WITH RING ABOVE */
	/* 0x0F */ 0x00E5, /* LATIN SMALL LETTER A WITH RING ABOVE */
	/* 0x10 */ 0x0394, /* GREEK CAPITAL LETTER DELTA */
	/* 0x11 */ 0x005F, /* LOW LINE */
	/* 0x12 */ 0x03A6, /* GREEK CAPITAL LETTER PHI */
	/* 0x13 */ 0x0393, /* GREEK CAPITAL LETTER GAMMA */
	/* 0x14 */ 0x039B, /* GREEK CAPITAL LETTER LAMDA */
	/* 0x15 */ 0x03A9, /* GREEK CAPITAL LETTER OMEGA */
	/* 0x16 */ 0x03A0, /* GREEK CAPITAL LETTER PI */
	/* 0x17 */ 0x03A8, /* GREEK CAPITAL LETTER PSI */
	/* 0x18 */ 0x03A3, /* GREEK CAPITAL LETTER SIGMA */
	/* 0x19 */ 0x0398, /* GREEK CAPITAL LETTER THETA */
	/* 0x1A */ 0x039E, /* GREEK CAPITAL LETTER XI */
	/* 0x1B */ 0x00A0, /* ESCAPE TO EXTENSION TABLE */
	/* 0x1C */ 0x00C6, /* LATIN CAPITAL LETTER AE */
	/* 0x1D */ 0x00E6, /* LATIN SMALL LETTER AE */
	/* 0x1E */ 0x00DF, /* LATIN SMALL LETTER SHARP S (German) */
	/* 0x1F */ 0x00C9, /* LATIN CAPITAL LETTER E WITH ACUTE */
	/* 0x20 */ 0x0020, /* SPACE */
	/* 0x21 */ 0x0021, /* EXCLAMATION MARK */
	/* 0x22 */ 0x0022, /* QUOTATION MARK */
	/* 0x23 */ 0x0023, /* NUMBER SIGN */
	/* 0x24 */ 0x00A4, /* CURRENCY SIGN */
	/* 0x25 */ 0x0025, /* PERCENT SIGN */
	/* 0x26 */ 0x0026, /* AMPERSAND */
	/* 0x27 */ 0x0027, /* APOSTROPHE */
	/* 0x28 */ 0x0028, /* LEFT PARENTHESIS */
	/* 0x29 */ 0x0029, /* RIGHT PARENTHESIS */
	/* 0x2A */ 0x002A, /* ASTERISK */
	/* 0x2B */ 0x002B, /* PLUS SIGN */
	/* 0x2C */ 0x002C, /* COMMA */
	/* 0x2D */ 0x002D, /* HYPHEN-MINUS */
	/* 0x2E */ 0x002E, /* FULL STOP */
	/* 0x2F */ 0x002F, /* SOLIDUS */
	/* 0x30 */ 0x0030, /* DIGIT ZERO */
	/* 0x31 */ 0x0031, /* DIGIT ONE */
	/* 0x32 */ 0x0032, /* DIGIT TWO */
	/* 0x33 */ 0x0033, /* DIGIT THREE */
	/* 0x34 */ 0x0034, /* DIGIT FOUR */
	/* 0x35 */ 0x0035, /* DIGIT FIVE */
	/* 0x36 */ 0x0036, /* DIGIT SIX */
	/* 0x37 */ 0x0037, /* DIGIT SEVEN */
	/* 0x38 */ 0x0038, /* DIGIT EIGHT */
	/* 0x39 */ 0x0039, /* DIGIT NINE */
	/* 0x3A */ 0x003A, /* COLON */
	/* 0x3B */ 0x003B, /* SEMICOLON */
	/* 0x3C */ 0x003C, /* LESS-THAN SIGN */
	/* 0x3D */ 0x003D, /* EQUALS SIGN */
	/* 0x3E */ 0x003E, /* GREATER-THAN SIGN */
	/* 0x3F */ 0x003F, /* QUESTION MARK */
	/* 0x40 */ 0x00A1, /* INVERTED EXCLAMATION MARK */
	/* 0x41 */ 0x0041, /* LATIN CAPITAL LETTER A */
	/* 0x42 */ 0x0042, /* LATIN CAPITAL LETTER B */
	/* 0x43 */ 0x0043, /* LATIN CAPITAL LETTER C */
	/* 0x44 */ 0x0044, /* LATIN CAPITAL LETTER D */
	/* 0x45 */ 0x0045, /* LATIN CAPITAL LETTER E */
	/* 0x46 */ 0x0046, /* LATIN CAPITAL LETTER F */
	/* 0x47 */ 0x0047, /* LATIN CAPITAL LETTER G */
	/* 0x48 */ 0x0048, /* LATIN CAPITAL LETTER H */
	/* 0x49 */ 0x0049, /* LATIN CAPITAL LETTER I */
	/* 0x4A */ 0x004A, /* LATIN CAPITAL LETTER J */
	/* 0x4B */ 0x004B, /* LATIN CAPITAL LETTER K */
	/* 0x4C */ 0x004C, /* LATIN CAPITAL LETTER L */
	/* 0x4D */ 0x004D, /* LATIN CAPITAL LETTER M */
	/* 0x4E */ 0x004E, /* LATIN CAPITAL LETTER N */
	/* 0x4F */ 0x004F, /* LATIN CAPITAL LETTER O */
	/* 0x50 */ 0x0050, /* LATIN CAPITAL LETTER P */
	/* 0x51 */ 0x0051, /* LATIN CAPITAL LETTER Q */
	/* 0x52 */ 0x0052, /* LATIN CAPITAL LETTER R */
	/* 0x53 */ 0x0053, /* LATIN CAPITAL LETTER S */
	/* 0x54 */ 0x0054, /* LATIN CAPITAL LETTER T */
	/* 0x55 */ 0x0055, /* LATIN CAPITAL LETTER U */
	/* 0x56 */ 0x0056, /* LATIN CAPITAL LETTER V */
	/* 0x57 */ 0x0057, /* LATIN CAPITAL LETTER W */
	/* 0x58 */ 0x0058, /* LATIN CAPITAL LETTER X */
	/* 0x59 */ 0x0059, /* LATIN CAPITAL LETTER Y */
	/* 0x5A */ 0x005A, /* LATIN CAPITAL LETTER Z */
	/* 0x5B */ 0x00C4, /* LATIN CAPITAL LETTER A WITH DIAERESIS */
	/* 0x5C */ 0x00D6, /* LATIN CAPITAL LETTER O WITH DIAERESIS */
	/* 0x5D */ 0x00D1, /* LATIN CAPITAL LETTER N WITH TILDE */
	/* 0x5E */ 0x00DC, /* LATIN CAPITAL LETTER U WITH DIAERESIS */
	/* 0x5F */ 0x00A7, /* SECTION SIGN */
	/* 0x60 */ 0x00BF, /* INVERTED QUESTION MARK */
	/* 0x61 */ 0x0061, /* LATIN SMALL LETTER A */
	/* 0x62 */ 0x0062, /* LATIN SMALL LETTER B */
	/* 0x63 */ 0x0063, /* LATIN SMALL LETTER C */
	/* 0x64 */ 0x0064, /* LATIN SMALL LETTER D */
	/* 0x65 */ 0x0065, /* LATIN SMALL LETTER E */
	/* 0x66 */ 0x0066, /* LATIN SMALL LETTER F */
	/* 0x67 */ 0x0067, /* LATIN SMALL LETTER G */
	/* 0x68 */ 0x0068, /* LATIN SMALL LETTER H */
	/* 0x69 */ 0x0069, /* LATIN SMALL LETTER I */
	/* 0x6A */ 0x006A, /* LATIN SMALL LETTER J */
	/* 0x6B */ 0x006B, /* LATIN SMALL LETTER K */
	/* 0x6C */ 0x006C, /* LATIN SMALL LETTER L */
	/* 0x6D */ 0x006D, /* LATIN SMALL LETTER M */
	/* 0x6E */ 0x006E, /* LATIN SMALL LETTER N */
	/* 0x6F */ 0x006F, /* LATIN SMALL LETTER O */
	/* 0x70 */ 0x0070, /* LATIN SMALL LETTER P */
	/* 0x71 */ 0x0071, /* LATIN SMALL LETTER Q */
	/* 0x72 */ 0x0072, /* LATIN SMALL LETTER R */
	/* 0x73 */ 0x0073, /* LATIN SMALL LETTER S */
	/* 0x74 */ 0x0074, /* LATIN SMALL LETTER T */
	/* 0x75 */ 0x0075, /* LATIN SMALL LETTER U */
	/* 0x76 */ 0x0076, /* LATIN SMALL LETTER V */
	/* 0x77 */ 0x0077, /* LATIN SMALL LETTER W */
	/* 0x78 */ 0x0078, /* LATIN SMALL LETTER X */
	/* 0x79 */ 0x0079, /* LATIN SMALL LETTER Y */
	/* 0x7A */ 0x007A, /* LATIN SMALL LETTER Z */
	/* 0x7B */ 0x00E4, /* LATIN SMALL LETTER A WITH DIAERESIS */
	/* 0x7C */ 0x00F6, /* LATIN SMALL LETTER O WITH DIAERESIS */
	/* 0x7D */ 0x00F1, /* LATIN SMALL LETTER N WITH TILDE */
	/* 0x7E */ 0x00FC, /* LATIN SMALL LETTER U WITH DIAERESIS */
	/* 0x7F */ 0x00E0, /* LATIN SMALL LETTER A WITH GRAVE */
}
