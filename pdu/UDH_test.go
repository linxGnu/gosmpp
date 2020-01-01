package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserDataHeader(t *testing.T) {
	t.Run("marshalBinaryUDHConcatMessage", func(t *testing.T) {
		u := UDH{NewIEConcatMessage(2, 1, 12)}
		b, err := u.MarshalBinary()
		require.NoError(t, err)
		require.Equal(t, "0500030c0201", toHex(b))
	})

	t.Run("unmarshalBinaryUDHConcatMessage", func(t *testing.T) {
		u, rd := new(UDH), []byte{0x05, 0x00, 0x03, 0x0c, 0x02, 0x01}
		_, err := u.UnmarshalBinary(rd)

		require.NoError(t, err)

		b, err := u.MarshalBinary()
		require.NoError(t, err)

		require.Equal(t, "0500030c0201", toHex(b))
	})
}
