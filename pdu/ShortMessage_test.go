package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/errors"

	"github.com/stretchr/testify/require"
)

type customEncoder struct{}

func (*customEncoder) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (*customEncoder) Decode(data []byte) (string, error) {
	return string(data), nil
}

func TestShortMessage(t *testing.T) {
	t.Run("invalidCoding", func(t *testing.T) {
		var s ShortMessage
		require.NotNil(t, s.SetMessageWithEncoding("agjwklgjkwPфngưỡng", data.LATIN1))
	})

	t.Run("customCoding", func(t *testing.T) {
		var s ShortMessage

		customCoding := data.NewCustomEncoding(246, &customEncoder{})
		err := s.SetMessageDataWithEncoding([]byte{0x61, 0x62, 0x63}, customCoding) // "abc"
		require.NoError(t, err)
		require.EqualValues(t, 246, s.Encoding().DataCoding())

		m, err := s.GetMessage()
		require.Nil(t, err)
		require.Equal(t, "abc", m)

		// try to get message string with other encoding
		m, err = s.GetMessageWithEncoding(data.FromDataCoding(data.UCS2Coding))
		require.Nil(t, err)
		require.NotEqual(t, "abc", m)

		// get message string with custom encoding
		m, err = s.GetMessageWithEncoding(customCoding)
		require.Nil(t, err)
		require.Equal(t, "abc", m)
	})

	t.Run("customCodingFromPeer", func(t *testing.T) {
		var senderSM ShortMessage

		// set custom data coding for test
		customCoding := data.NewCustomEncoding(0x19, &customEncoder{})

		err := senderSM.SetMessageDataWithEncoding([]byte{0x61, 0x62, 0x63}, customCoding) // "abc"
		require.NoError(t, err)

		b := NewBuffer(nil)
		senderSM.Marshal(b)

		// From here the message is not know anymore to the receiver in terms of encoding methods, but he wants to know the encoding code while once received the packet
		var receivedSM ShortMessage
		err = receivedSM.Unmarshal(b, false)
		require.NoError(t, err)

		require.NotNil(t, receivedSM.Encoding())
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

	t.Run("getMessageData", func(t *testing.T) {
		s, err := NewBinaryShortMessage([]byte{0x00, 0x01, 0x02, 0x03})
		require.NoError(t, err)

		messageData, err := s.GetMessageData()
		require.NoError(t, err)
		require.Equal(t, "00010203", toHex(messageData))
	})

	t.Run("marshalBinaryMessage", func(t *testing.T) {
		s, err := NewBinaryShortMessage([]byte{0x00, 0x01, 0x02, 0x03, 0x04})
		require.NoError(t, err)

		buf := NewBuffer(nil)
		s.Marshal(buf)

		require.Equal(t, "0400050001020304", toHex(buf.Bytes()))
	})

	t.Run("marshalWithoutCoding", func(t *testing.T) {
		var s ShortMessage
		err := s.SetMessageDataWithEncoding([]byte("abc"), nil)
		require.NoError(t, err)
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

	t.Run("unmarshalBinaryWithUDHConcat", func(t *testing.T) {
		s := &ShortMessage{}

		buf := NewBuffer([]byte{0x04, 0x00, 0x09, 0x05, 0x00, 0x03, 0x0c, 0x02, 0x01, 0x01, 0x02, 0x03})

		// check encoding
		require.NoError(t, s.Unmarshal(buf, true))
		require.Equal(t, data.BINARY8BIT2, s.Encoding())

		// check message
		messageData, err := s.GetMessageData()
		require.NoError(t, err)
		require.Equal(t, []byte{0x01, 0x02, 0x03}, messageData)
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
		sm, err := NewLongMessageWithEncoding("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz1234123456789", data.GSM7BIT)
		require.NoError(t, err)

		require.Equal(t, 2, len(sm))
	})

	t.Run("shortMessageSplitGSM7_160chars", func(t *testing.T) {
		// over gsm7 chars limit ( 160/160 ), split
		sm, err := NewLongMessageWithEncoding("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz1234", data.GSM7BIT)
		require.NoError(t, err)

		require.Equal(t, 2, len(sm))
	})

	t.Run("shortMessageSplitUCS2_89chars", func(t *testing.T) {
		// over UCS2 chars limit (89/67), split
		sm, err := NewLongMessageWithEncoding("biggest gift của Christmas là có nhiều big/challenging/meaningful problems để sấp mặt làm", data.UCS2)
		require.NoError(t, err)

		require.Equal(t, 2, len(sm))
	})

	t.Run("shortMessageSplitUCS2_67chars", func(t *testing.T) {
		// still within UCS2 chars limit (67/67), not split
		sm, err := NewLongMessageWithEncoding("biggest gift của Christmas là có nhiều big/challenging/meaningful p", data.UCS2)
		require.NoError(t, err)

		require.Equal(t, 1, len(sm))
	})

	t.Run("shortMessageSplitGSM7_empty", func(t *testing.T) {
		// over UCS2 chars limit (89/67), split
		sm, err := NewLongMessageWithEncoding("", data.GSM7BIT)
		require.NoError(t, err)

		require.Equal(t, 1, len(sm))
	})

	t.Run("indempotentMarshal", func(t *testing.T) {
		// over gsm7 chars limit ( 160/160 ), split
		multiSM, err := NewLongMessageWithEncoding("abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz1234", data.GSM7BIT)
		require.NoError(t, err)

		for i := range multiSM {
			b1, b2 := NewBuffer(nil), NewBuffer(nil)
			multiSM[i].Marshal(b1)
			multiSM[i].Marshal(b2)
			require.Equal(t, b1.Bytes(), b2.Bytes())
		}
	})
}
