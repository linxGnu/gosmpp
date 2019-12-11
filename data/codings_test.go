package data

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func fromHex(h string) (v []byte) {
	var err error
	v, err = hex.DecodeString(h)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func testEncoding(t *testing.T, enc EncDec, original, expected string) {
	encoded, err := enc.Encode(original)
	require.Nil(t, err)
	require.Equal(t, fromHex(expected), encoded)

	decoded, err := enc.Decode(encoded)
	require.Nil(t, err)
	require.Equal(t, original, decoded)
}

func TestCoding(t *testing.T) {
	require.Nil(t, FromDataCoding(12))
	require.Equal(t, GSM7BIT, FromDataCoding(0))
	require.Equal(t, ASCII, FromDataCoding(1))
	require.Equal(t, UCS2, FromDataCoding(8))
	require.Equal(t, LATIN1, FromDataCoding(3))
	require.Equal(t, CYRILLIC, FromDataCoding(6))
	require.Equal(t, HEBREW, FromDataCoding(7))
}

func TestGSM7Bit(t *testing.T) {
	require.EqualValues(t, 0, GSM7BIT.DataCoding())
	testEncoding(t, GSM7BIT, "gjwklgjkwP", "676a776b6c676a6b7750")
}

func TestAscii(t *testing.T) {
	require.EqualValues(t, 1, ASCII.DataCoding())
	testEncoding(t, ASCII, "agjwklgjkwP", "61676a776b6c676a6b7750")
}

func TestUCS2(t *testing.T) {
	require.EqualValues(t, 8, UCS2.DataCoding())
	testEncoding(t, UCS2, "agjwklgjkwP", "00610067006a0077006b006c0067006a006b00770050")
}

func TestLatin1(t *testing.T) {
	require.EqualValues(t, 3, LATIN1.DataCoding())
	testEncoding(t, LATIN1, "agjwklgjkwPÓ", "61676a776b6c676a6b7750d3")
}

func TestCYRILLIC(t *testing.T) {
	require.EqualValues(t, 6, CYRILLIC.DataCoding())
	testEncoding(t, CYRILLIC, "agjwklgjkwPф", "61676A776B6C676A6B7750E4")
}

func TestHebrew(t *testing.T) {
	require.EqualValues(t, 7, HEBREW.DataCoding())
	testEncoding(t, HEBREW, "agjwklgjkwPץ", "61676A776B6C676A6B7750F5")
}

func TestOtherCodings(t *testing.T) {
	testEncoding(t, UTF16BEM, "ngưỡng cứa cuỗc đợi", "feff006e006701b01ee1006e0067002000631ee900610020006300751ed70063002001111ee30069")
	testEncoding(t, UTF16LEM, "ngưỡng cứa cuỗc đợi", "fffe6e006700b001e11e6e00670020006300e91e6100200063007500d71e630020001101e31e6900")
	testEncoding(t, UTF16BE, "ngưỡng cứa cuỗc đợi", "006e006701b01ee1006e0067002000631ee900610020006300751ed70063002001111ee30069")
	testEncoding(t, UTF16LE, "ngưỡng cứa cuỗc đợi", "6e006700b001e11e6e00670020006300e91e6100200063007500d71e630020001101e31e6900")
}
