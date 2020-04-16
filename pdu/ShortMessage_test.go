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
		s.messageData = []byte{0x61, 0x62, 0x63}

		m, err := s.GetMessage()
		require.Nil(t, err)
		require.Equal(t, "abc", m)
	})

	t.Run("marshalWithoutCoding", func(t *testing.T) {
		var s ShortMessage
		s.SetMessageData([]byte("abc"))
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
		require.Equal(t, "000003616263", toHex(buf.Bytes()))
	})

	t.Run("marshalWithCoding160chars", func(t *testing.T) {
		s, err := NewShortMessageWithEncoding("abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcabcab", data.GSM7BIT)
		require.NoError(t, err)

		buf := NewBuffer(nil)
		s.Marshal(buf)

		require.Equal(t, 116, len(buf.Bytes()))
	})

	t.Run("marshalGSM7WithUDHConcat", func(t *testing.T) {
		s, err := NewShortMessageWithEncoding("abc", data.GSM7BIT)
		require.NoError(t, err)
		s.SetUDH(UDH{NewIEConcatMessage(2, 1, 12)})

		buf := NewBuffer(nil)
		s.Marshal(buf)
		require.Equal(t, "0000090500030c0201616263", toHex(buf.Bytes()))
	})

	t.Run("unmarshalGSM7WithUDHConcat", func(t *testing.T) {
		s := &ShortMessage{}

		buf := NewBuffer([]byte{0x00, 0x00, 0x09, 0x05, 0x00, 0x03, 0x0c, 0x02, 0x01, 0x61, 0x62, 0x63})

		// check encoding
		require.NoError(t, s.Unmarshal(buf, true))
		require.Equal(t, data.GSM7BIT, s.Encoding())

		// check message
		message, err := s.GetMessageWithEncoding(s.Encoding())
		require.NoError(t, err)
		require.Equal(t, "abc", message)
	})

	t.Run("shortMessageSplitGSM7_169chars", func(t *testing.T) {
		// over gsm7 chars limit ( 169/160 ), split
		sm, err := NewShortMessageWithEncoding("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz1234123456789", data.GSM7BIT)
		require.NoError(t, err)

		multiSM, err := sm.Split()
		require.NoError(t, err)
		require.Equal(t, 2, len(multiSM))
	})

	t.Run("shortMessageSplitGSM7_160chars", func(t *testing.T) {
		// over gsm7 chars limit ( 160/160 ), split
		sm, err := NewShortMessageWithEncoding("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz1234", data.GSM7BIT)
		require.NoError(t, err)

		multiSM, err := sm.Split()
		require.NoError(t, err)
		require.Equal(t, 2, len(multiSM))
	})

	t.Run("shortMessageSplitUCS2_89chars", func(t *testing.T) {
		// over UCS2 chars limit (89/67), split
		sm, err := NewShortMessageWithEncoding("biggest gift của Christmas là có nhiều big/challenging/meaningful problems để sấp mặt làm", data.UCS2)
		require.NoError(t, err)

		multiSM, err := sm.Split()
		require.NoError(t, err)
		require.Equal(t, 2, len(multiSM))
	})

	t.Run("shortMessageSplitUCS2_67chars", func(t *testing.T) {
		// still within UCS2 chars limit (67/67), not split
		sm, err := NewShortMessageWithEncoding("biggest gift của Christmas là có nhiều big/challenging/meaningful p", data.UCS2)
		require.NoError(t, err)

		multiSM, err := sm.Split()
		require.NoError(t, err)
		require.Equal(t, 1, len(multiSM))
	})

	t.Run("shortMessageSplitGSM7_empty", func(t *testing.T) {
		// over UCS2 chars limit (89/67), split
		sm, err := NewShortMessageWithEncoding("", data.GSM7BIT)
		require.NoError(t, err)

		multiSM, err := sm.Split()
		require.NoError(t, err)
		require.Equal(t, 1, len(multiSM))
	})

	t.Run("indempotentMarshal", func(t *testing.T) {
		// over gsm7 chars limit ( 160/160 ), split
		sm, err := NewShortMessageWithEncoding("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz1234", data.GSM7BIT)
		require.NoError(t, err)

		multiSM, err := sm.Split()
		require.NoError(t, err)
		for i := range multiSM {
			b1, b2 := NewBuffer(nil), NewBuffer(nil)
			multiSM[i].Marshal(b1)
			multiSM[i].Marshal(b2)
			require.Equal(t, b1.Bytes(), b2.Bytes())
		}
	})
}
