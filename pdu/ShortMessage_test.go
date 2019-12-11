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
		s.messageData = []byte("abc")

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
}
