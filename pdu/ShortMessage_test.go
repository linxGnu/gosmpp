package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"

	"github.com/stretchr/testify/require"
)

func TestShortMessage(t *testing.T) {
	t.Run("invalidCoding", func(t *testing.T) {
		var s ShortMessage
		require.NotNil(t, s.SetMessageWithEncoding("agjwklgjkwPфngưỡng", data.LATIN1))
	})

	t.Run("invalidSize", func(t *testing.T) {
		var s ShortMessage
		require.Equal(t, errors.ErrShortMessageLengthTooLarge,
			s.SetMessageWithEncoding("agjwklgjkwPфngưỡngasdfasdfasdfasdagjwklgjkwPфngưỡngasdfasdfasdfasdagjwklgjkwPфngưỡngasdfasdfasdfasdagjwklgjkwPфngưỡngasdfasdfasdfasd", data.UCS2))
	})

	t.Run("getMessageWithoutCoding", func(t *testing.T) {
		var s ShortMessage
		s.messageData = []byte{0x61, 0xf1, 0x18}

		m, err := s.GetMessage()
		require.Nil(t, err)
		require.Equal(t, "abc", m)
	})

	t.Run("marshalWithoutCoding", func(t *testing.T) {
		var s ShortMessage
		s.messageData = []byte("abc")
		s.messageData = append(s.messageData, 0)
		s.enc = nil

		buf := NewBuffer(nil)
		s.Marshal(buf)
		require.Equal(t, "00000461626300", toHex(buf.Bytes()))
	})

	t.Run("marshalWithCoding", func(t *testing.T) {
		s, err := NewShortMessageWithEncoding("abc", data.GSM7BIT)
		require.NoError(t, err)

		buf := NewBuffer(nil)
		s.Marshal(buf)
		require.Equal(t, "00000361f118", toHex(buf.Bytes()))
	})

	t.Run("marshalWithCoding160chars", func(t *testing.T) {
		s, err := NewShortMessageWithEncoding("abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcaabcabcabcabcabcabcabcabcacc123121", data.GSM7BIT)
		require.NoError(t, err)

		buf := NewBuffer(nil)
		s.Marshal(buf)

		// should be able to fit 160 char into 140 octets ( + 3 bytes header )
		require.Equal(t, 143, len(buf.Bytes()))
	})

	t.Run("marshalGSM7WithUDHConcat", func(t *testing.T) {
		s, err := NewShortMessageWithEncoding("abc", data.GSM7BIT)
		require.NoError(t, err)
		require.NoError(t, s.SetUDH(UDH{NewIEConcatMessage(2, 1, 12)}))

		buf := NewBuffer(nil)
		s.Marshal(buf)
		require.Equal(t, "0000090500030c020161f118", toHex(buf.Bytes()))
	})
}
