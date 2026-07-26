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

		if enc == GSM7BITPACKED {
			seg = shiftBitsOneRight(seg)
		}
		decoded, err := enc.Decode(seg)
		require.Nil(t, err)
		require.Equal(t, expectDecode[i], decoded)
	}
}

func shiftBitsOneRight(input []byte) []byte {
	carry := byte(0)
	for i := len(input) - 1; i >= 0; i-- {
		// Save the carry bit from the previous byte
		nextCarry := input[i] & 0b00000001
		// Shift the current byte to the right
		input[i] >>= 1
		// Apply the carry from the previous byte to the current byte
		input[i] |= carry << 7
		// Update the carry for the next byte
		carry = nextCarry
	}
	return input
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

func TestShouldSplit(t *testing.T) {
	t.Run("testShouldSplit_ASCII", func(t *testing.T) {
		octetLim := uint(140)
		expect := map[string]bool{
			"":  false,
			"1": false,
			"12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890":  false, // exactly 140 chars
			"123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901": true,  // 141 chars
		}

		splitter, ok := ASCII.(Splitter)
		require.True(t, ok, "ASCII must implement Splitter interface")
		for k, v := range expect {
			ok := splitter.ShouldSplit(k, octetLim)
			require.Equalf(t, v, ok, "Test case len=%d", len(k))
		}
	})

	t.Run("testShouldSplit_CYRILLIC", func(t *testing.T) {
		octetLim := uint(140)
		splitter, ok := CYRILLIC.(Splitter)
		require.True(t, ok, "CYRILLIC must implement Splitter interface")

		// 140 Cyrillic chars
		msg140 := "аааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааа"
		require.False(t, splitter.ShouldSplit(msg140, octetLim), "140 chars should not split")

		// 141 Cyrillic chars
		msg141 := "ааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааа"
		require.True(t, splitter.ShouldSplit(msg141, octetLim), "141 chars should split")
	})

	t.Run("testShouldSplit_HEBREW", func(t *testing.T) {
		octetLim := uint(140)
		splitter, ok := HEBREW.(Splitter)
		require.True(t, ok, "HEBREW must implement Splitter interface")

		// 140 Hebrew chars
		msg140 := "אאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאא"
		require.False(t, splitter.ShouldSplit(msg140, octetLim), "140 chars should not split")

		// 141 Hebrew chars
		msg141 := "אאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאא"
		require.True(t, splitter.ShouldSplit(msg141, octetLim), "141 chars should split")
	})

	t.Run("testShouldSplit_GSM7BIT", func(t *testing.T) {
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

	t.Run("testShouldSplit_UCS2", func(t *testing.T) {
		octetLim := uint(140)
		expect := map[string]bool{
			"":  false,
			"1": false,
			"ởỀÊộẩừỰÉÊỗọễệớỡồỰỬỪựởặỬ̀ỵổẤỨợỶẰỢộứẶHữẹ̃ẾỆằỄéậÃỡẰộ̀ỀỗứẲữỪữộÊỵòALữộòC":  false, /* 70 UCS2 chars */
			"ợÁÊGỷẹííỡỮÂIỆàúễẠỮỊệÂỖÍắẵYẠừẲíộờíẵỠựẤằờởể̃ởỵởềệổồUỡỵầễÁÝởÝNè̉ỚổôỊộợKỨệ́": true,  /* 71 UCS2 chars */
		}

		splitter, _ := UCS2.(Splitter)
		for k, v := range expect {
			ok := splitter.ShouldSplit(k, octetLim)
			require.Equalf(t, ok, v, "Test case %s", k)
		}
	})

	t.Run("testShouldSplit_GSM7BITPACKED", func(t *testing.T) {
		octetLim := uint(140)
		expect := map[string]bool{
			"":  false,
			"1": false,
			"12312312311231231231123123123112312312311231231231123123123112312312311231231231123123123112312312311231231231123123123112312312311234121212":                      false,
			"gjwklgjkwP123+?sasdasdaqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdqwdqwDQWdqwdqwdqwdqwwqwdqwdqwddqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdwqdqwqwdqwdqwqwdqw":  false, /* 160 regular basic alphabet chars */
			"gjwklgjkwP123+?sasdasdaqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdqwdqwDQWdqwdqwdqwdqwwqwdqwdqwddqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdwqdqwqwdqwdqwqwdqwd": true,  /* 161 regular basic alphabet chars */
			"gjwklgjkwP123+?sasdasdaqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdqwdqwDQWdqwdqwdqwdqwwqwdqwdqwddqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdwqdqwqwdqwdqwqwdqw{": true,  /* 159 regular basic alphabet chars + 1 escape char at the end */
			"|}€€|]|€[~€^]€~{~^{|]]|[{|~€^|]^[[{€^]^{€}}^~~]€]~€[€€[]~~[}}]{^}{|}~~]]€^{^|€{^":                                                                                  false, /* 80 escape chars */
			"|}€€|]|€[~€^]€~{~^{|]]|[{|~€^|]^[[{€^]^{€}}^~~]€]~€[€€[]~~[}}]{^}{|}~~]]€^{^|€{^{":                                                                                 true,  /* 81 escape chars */
		}

		splitter, _ := GSM7BITPACKED.(Splitter)
		for k, v := range expect {
			ok := splitter.ShouldSplit(k, octetLim)
			require.Equalf(t, ok, v, "Test case %s", k)
		}
	})
}
func TestSplit(t *testing.T) {
	require.EqualValues(t, 0o0, GSM7BITPACKED.DataCoding())

	t.Run("testSplitASCII", func(t *testing.T) {
		// 212 ASCII chars - should split into 2 segments: 134 + 78
		msg := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		splitter, ok := ASCII.(Splitter)
		require.True(t, ok)

		segments, err := splitter.EncodeSplit(msg, 134)
		require.Nil(t, err)
		require.Equal(t, 2, len(segments))
		require.Equal(t, 134, len(segments[0]))
		require.Equal(t, 78, len(segments[1]))

		// Verify decoded content
		decoded1, err := ASCII.Decode(segments[0])
		require.Nil(t, err)
		decoded2, err := ASCII.Decode(segments[1])
		require.Nil(t, err)
		require.Equal(t, msg, decoded1+decoded2)
	})

	t.Run("testSplitCYRILLIC", func(t *testing.T) {
		// 200 Cyrillic chars - should split into 2 segments
		msg := "аааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааааа"
		splitter, ok := CYRILLIC.(Splitter)
		require.True(t, ok)

		segments, err := splitter.EncodeSplit(msg, 134)
		require.Nil(t, err)
		require.Equal(t, 2, len(segments))
		require.LessOrEqual(t, len(segments[0]), 134)
		require.LessOrEqual(t, len(segments[1]), 134)

		// Verify decoded content
		decoded1, err := CYRILLIC.Decode(segments[0])
		require.Nil(t, err)
		decoded2, err := CYRILLIC.Decode(segments[1])
		require.Nil(t, err)
		require.Equal(t, msg, decoded1+decoded2)
	})

	t.Run("testSplitHEBREW", func(t *testing.T) {
		// 200 Hebrew chars - should split into 2 segments
		msg := "אאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאאא"
		splitter, ok := HEBREW.(Splitter)
		require.True(t, ok)

		segments, err := splitter.EncodeSplit(msg, 134)
		require.Nil(t, err)
		require.Equal(t, 2, len(segments))
		require.LessOrEqual(t, len(segments[0]), 134)
		require.LessOrEqual(t, len(segments[1]), 134)

		// Verify decoded content
		decoded1, err := HEBREW.Decode(segments[0])
		require.Nil(t, err)
		decoded2, err := HEBREW.Decode(segments[1])
		require.Nil(t, err)
		require.Equal(t, msg, decoded1+decoded2)
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

func TestSplit_GSM7BITPACKED(t *testing.T) {
	require.EqualValues(t, 0o0, GSM7BITPACKED.DataCoding())

	t.Run("testSplit_Escape_GSM7BITPACKED", func(t *testing.T) {
		testEncodingSplit(t, GSM7BITPACKED,
			134,
			"gjwklgjkwP123+?sasdasdaqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdqwdqwDQWdqwdqwdqwdqwwqwdqwdqwddqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdwqdqwqwdqwdqwqwdqw{",
			[]string{
				"ceeafb9a7d56afefd0986cb6facdc37372784e0ec7efe4f89d1cbf93e37772fc4e8edfc9f13b397e27c7efe4f89d1cbf93e37772fc4e8edfc9f13b397e27c7efe438397e27c7efc4e8951cbf93e37772fc4e8edfeff13b397e27c7ef6472fc4e8edfc9f13b397e27c7efe4f89d1cbf93e37772fc4e8edfc9f13b394ebec7c9f17bfc4e8edfc9",
				"e2f7f89d1cbf6f50",
			},
			[]string{
				"gjwklgjkwP123+?sasdasdaqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdqwdqwDQWdqwdqwdqwdqwwqwdqwdqwddqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqwdqdwqdqwqwdqwd",
				"qwqwdqw{",
			})
	})

	/*
		Total char count = 160,
		Esc char count = 1,
		Regular char count = 159,
		Seg1 => 153->€
		Expected behaviour: Should not split in the middle of ESC chars
	*/
	t.Run("testSplit_EscEndOfSeg1_GSM7BITPACKED", func(t *testing.T) {
		testEncodingSplit(t, GSM7BITPACKED,
			134,
			"pppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp€ppppppp",
			[]string{
				"e070381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c31b",
				"3665381c0e87c3e1",
			},
			[]string{
				"pppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp\r",
				"€ppppppp",
			})
	})

	/*
		Total char count = 160,
		Esc char count = 2,
		Regular char count = 158,
		Seg1 => 152-> ....{
		Seg2 => 1-> ....{
		Expected behaviour: Should not split in the middle of ESC chars
	*/
	t.Run("testSplit_EscEndOfSeg1AndSeg2_1_GSM7BITPACKED", func(t *testing.T) {
		testEncodingSplit(t, GSM7BITPACKED,
			134,
			"pppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp{{pppppppp",
			[]string{
				"e070381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0edfa01a",
				"3628381c0e87c3e170",
			},
			[]string{
				"pppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp{\r",
				"{pppppppp",
			})
	})

	/*
		Total char count = 160,
		Esc char count = 2,
		Regular char count = 158,
		Seg1 => 152-> ....€
		Seg2 => 1-> ....€
		Expected behaviour: Should not split in the middle of ESC chars
	*/
	t.Run("testSplit_EscEndOfSeg1AndSeg2_2_GSM7BITPACKED", func(t *testing.T) {
		testEncodingSplit(t, GSM7BITPACKED,
			134,
			"pppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp€€pppppppp",
			[]string{
				"e070381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0e87c3e170381c0edf941b",
				"3665381c0e87c3e170",
			},
			[]string{
				"pppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp€\r",
				"€pppppppp",
			})
	})

	/*
		Total char count = 162,
		Esc char count = 0,
		Regular char count = 162,
		Seg1 => 153
		Seg2 => 9
		Scenario: All charcters in the GSM7Bit Basic Character Set table (non-escape chars) https://en.wikipedia.org/wiki/GSM_03.38
	*/
	t.Run("testSplit_AllGSM7BitBasicCharset_GSM7BITPACKED", func(t *testing.T) {
		testEncodingSplit(t, GSM7BITPACKED,
			134,
			"ΩØ;19Ξòå1-¤6aΞΘANanΣ¡>)òΦ3L;aøΛ-o@>I¥1=-ü!N¤&o9Hmda3jΞ@ÅΣlhEE§/:Çù0Θ&:_&Π;KLÅÅ@fÜ-kFH?ΠB5/ÆΓ?55=<Ω¡N2ñ¥*L¤aÖ! ÖΘ+øF£_Ç?øΔΓ-lèòCìnEBmhÉF*<Åi/aΩ¥CDøfGÇ$/=Λ'ÅA3ò#fkù",
			[]string{
				"2a8b5d2ca7413c622d922dacc9049d613706e84b212433e62ecca0b4de005f7210ebb5fc2127c9f4ce21dbe4f04cad0138306c74b1f87de9120658c6a48b982cbb25d3e10098bdadb511f9b3086b2fcee457abf57815a053d61fa898a4303704e266560c632092f8312093169b80181edc45611bfd31aa788ef42b5c190c890cf3312178f528",
				"4e8ee00c3132af0d",
			},
			[]string{
				"ΩØ;19Ξòå1-¤6aΞΘANanΣ¡>)òΦ3L;aøΛ-o@>I¥1=-ü!N¤&o9Hmda3jΞ@ÅΣlhEE§/:Çù0Θ&:_&Π;KLÅÅ@fÜ-kFH?ΠB5/ÆΓ?55=<Ω¡N2ñ¥*L¤aÖ! ÖΘ+øF£_Ç?øΔΓ-lèòCìnEBmhÉF*<Åi/aΩ¥CDøfGÇ$/=Λ",
				"'ÅA3ò#fkù",
			})
	})

	/*
		Total char count = 81,
		Esc char count = 81,
		Regular char count = 0,
		Seg1 => 153
		Seg2 => 9
		Scenario: All charcters in the GSM7Bit Escape Character Set table https://en.wikipedia.org/wiki/GSM_03.38
	*/
	t.Run("testSplit_AllGSM7BitBasicCharset_GSM7BITPACKED", func(t *testing.T) {
		testEncodingSplit(t, GSM7BITPACKED,
			134,
			"|{[€|^€[{|€{[|^{~[}€|}|^|^[^]€{[]~}€]{{^|^][€]|€~€^[~}^]{]~{^^€^[~|^]|€~|^€{]{~|}",
			[]string{
				"36c00d6ac3db9437c00d6553def036a80d7053dea036bc0d7043d9a036bd0d6f93da9437c04d6a03dc5036c00d65c3db5036be4d7983daf036be4d6f93da9437be0d6a83da5036c00d65e3dbf036e58d6f03dc9437bd4d7943d9f036bd4d6a43d9f836a88d6fd3dba036940d6553de5036bc4d6f03dc5036be0d7053def436c00d6553dea01a",
				"36be0d6ad3db003729",
			},
			[]string{
				"|{[€|^€[{|€{[|^{~[}€|}|^|^[^]€{[]~}€]{{^|^][€]|€~€^[~}^]{]~{^^€^[~|^]|€~|^€{\r",
				"]{~|}",
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

type noOpEncDec struct{}

func (*noOpEncDec) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (*noOpEncDec) Decode(data []byte) (string, error) {
	return string(data), nil
}

func TestCustomEncoding(t *testing.T) {
	enc := NewCustomEncoding(GSM7BITCoding, &noOpEncDec{})
	require.EqualValues(t, GSM7BITCoding, enc.DataCoding())

	encoded, err := enc.Encode("abc")
	require.NoError(t, err)
	require.Equal(t, []byte("abc"), encoded)

	decoded, err := enc.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, "abc", decoded)
}
