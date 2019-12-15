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
}
