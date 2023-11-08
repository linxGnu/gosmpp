package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/errors"

	"github.com/stretchr/testify/require"
)

func TestParsePDU(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		buf := NewBuffer(fromHex("00000010800000060000000000000001"))
		_, err := Parse(buf)
		require.Nil(t, err)
	})

	t.Run("submit_sm_resp with command_status != 0 ", func(t *testing.T) {
		buf := NewBuffer(fromHex("00000010800000040000005800000001"))
		_, err := Parse(buf)
		require.Nil(t, err)
	})

	t.Run("eof", func(t *testing.T) {
		buf := NewBuffer(nil)
		_, err := Parse(buf)
		require.NotNil(t, err)
	})

	t.Run("invalidCmdLength", func(t *testing.T) {
		buf := NewBuffer(fromHex("0000000f800000060000000000000001"))
		_, err := Parse(buf)
		require.Equal(t, errors.ErrInvalidPDU, err)

		buf = NewBuffer(fromHex("3800000f800000060000000000000001"))
		_, err = Parse(buf)
		require.Equal(t, errors.ErrInvalidPDU, err)
	})

	t.Run("invalidBody", func(t *testing.T) {
		buf := NewBuffer(fromHex("0000001e00000003000000000000000161776179001c1d416c69636572"))
		_, err := Parse(buf)
		require.NotNil(t, err)
	})

	t.Run("invalidPayload", func(t *testing.T) {
		buf := NewBuffer(fromHex("000000118000000400000000000000010012"))
		var b base
		require.NotNil(t, b.unmarshal(buf, func(buf *ByteBuffer) error {
			return nil
		}))

		buf = NewBuffer(fromHex("000000118000000400000000000000010012333333333333333333"))
		require.NotNil(t, b.unmarshal(buf, func(buf *ByteBuffer) error {
			_, _ = buf.ReadN(8)
			return nil
		}))
	})
}
