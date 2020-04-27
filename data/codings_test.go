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

func testEncodingSplit(t *testing.T, enc EncDec, octetLim uint, original string, expected []string, expectDecode []string) {
	splitter, ok := enc.(Splitter)
	require.Truef(t, ok, "Encoding must implement Splitter interface")

	segEncoded, err := splitter.EncodeSplit(original, octetLim)
	require.Nil(t, err)

	for i, seg := range segEncoded {
		require.Equal(t, fromHex(expected[i]), seg)
		require.LessOrEqualf(t, uint(len(seg)), octetLim,
			"Segment len must be less than or equal to %d, got %d", octetLim, len(seg))
		decoded, err := enc.Decode(seg)
		require.Nil(t, err)
		require.Equal(t, expectDecode[i], decoded)
	}
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
	require.EqualValues(t, 0, GSM7BITPACKED.DataCoding())
	testEncoding(t, GSM7BITPACKED, "gjwklgjkwP123+?", "67f57dcd3eabd777684c365bfd00")
}

func TestSplit(t *testing.T) {
	require.EqualValues(t, 00, GSM7BITPACKED.DataCoding())

	t.Run("testShouldSplitGSM7", func(t *testing.T) {
		octetLim := uint(140)
		expect := map[string]bool{
			"":  false,
			"1": false,
			"12312312311231231231123123123112312312311231231231123123123112312312311231231231123123123112312312311231231231123123123112312312311234121212":  false,
			"123123123112312312311231231231123123123112312312311231231231123123123112312312311231231231123123123112312312311231231231123123123112342212121": true,
		}

		splitter, _ := GSM7BIT.(Splitter)
		for k, v := range expect {
			ok := splitter.ShouldSplit(k, octetLim)
			require.Equalf(t, ok, v, "Test case %s", k)
		}
	})

	t.Run("testShouldSplitUCS2", func(t *testing.T) {
	})

	t.Run("testSplitEscapeGSM7", func(t *testing.T) {
		testEncodingSplit(t, GSM7BITPACKED,
			134,
			"gjwklgjkwP123+?sasdasdaqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdqwdqwDQWdqwdqwdqwdqwwqwdqwdqwddqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdwqdqwqwdqwdqwqwdqw{",
			[]string{
				"67f57dcd3eabd777684c365bfde6e139393c2787e37772fc4e8edfc9f13b397e27c7efe4f89d1cbf93e37772fc4e8edfc9f13b397e27c7efe4f89d1cbf93e377729c1cbf93e37762f44a8edfc9f13b397e27c7eff7f89d1cbf93e37732397e27c7efe4f89d1cbf93e37772fc4e8edfc9f13b397e2703",
				"f13b397e27c7c9f738397e8fdfc9f13b397e8fdfc9f1fb0605",
			},
			[]string{
				"gjwklgjkwP123+?sasdasdaqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdqwdqwDQWdqwdqwdqwdqwwqwdqwdqwddqwdqwdqwdqwdqwdqwdqwdqwdqwd",
				"qwdqwdqdwqdqwqwdqwdqwqwdqw{",
			})
	})

	t.Run("testSplitGSM7Empty", func(t *testing.T) {
		testEncodingSplit(t, GSM7BIT,
			134,
			"",
			[]string{
				"",
			},
			[]string{
				"",
			})
	})

	t.Run("testSplitUCS2", func(t *testing.T) {
		testEncodingSplit(t, UCS2,
			134,
			"biggest gift của Christmas là có nhiều big/challenging/meaningful problems để sấp mặt làm",
			[]string{
				"006200690067006700650073007400200067006900660074002000631ee700610020004300680072006900730074006d006100730020006c00e00020006300f30020006e006800691ec100750020006200690067002f006300680061006c006c0065006e00670069006e0067002f006d00650061006e0069006e006700660075006c00200070",
				"0072006f0062006c0065006d0073002001111ec3002000731ea500700020006d1eb700740020006c00e0006d",
			},
			[]string{
				"biggest gift của Christmas là có nhiều big/challenging/meaningful p",
				"roblems để sấp mặt làm",
			})
	})

	t.Run("testSplitUCS2Empty", func(t *testing.T) {
		testEncodingSplit(t, UCS2,
			134,
			"",
			[]string{
				"",
			},
			[]string{
				"",
			})
	})

	// UCS2 character should not be splitted in the middle
	// here 54 character is encoded to 108 octet, but since there are 107 octet limit,
	// a whole 2 octet has to be carried over to the next segment
	t.Run("testSplit_Middle_UCS2", func(t *testing.T) {
		testEncodingSplit(t, UCS2,
			107,
			"biggest gift của Christmas là có nhiều big/challenging",
			[]string{
				"006200690067006700650073007400200067006900660074002000631ee700610020004300680072006900730074006d006100730020006c00e00020006300f30020006e006800691ec100750020006200690067002f006300680061006c006c0065006e00670069006e",
				"0067", // 0x00 0x67 is "g"
			},
			[]string{
				"biggest gift của Christmas là có nhiều big/challengin",
				"g",
			})
	})
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
