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

	t.Run("marshalBinaryUDHConcatMessage (8 bit)", func(t *testing.T) {
		u := UDH{NewIEConcatMessage(2, 1, 12)}
		b, err := u.MarshalBinary()
		require.NoError(t, err)
		require.Equal(t, "0500030c0201", toHex(b))

		totalParts, sequence, reference, found := u.GetConcatInfo()
		require.True(t, found)
		require.Equal(t, totalParts, byte(2))
		require.Equal(t, sequence, byte(1))
		require.Equal(t, reference, uint(12))
	})

	t.Run("marshalBinaryUDHConcatMessage (16 bit)", func(t *testing.T) {
		u := UDH{NewIEConcatMessage16Bit(2, 1, 1234)}
		b, err := u.MarshalBinary()
		require.NoError(t, err)
		require.Equal(t, "06080404d20201", toHex(b))

		totalParts, sequence, reference, found := u.GetConcatInfo()
		require.True(t, found)
		require.Equal(t, totalParts, byte(2))
		require.Equal(t, sequence, byte(1))
		require.Equal(t, reference, uint(1234))
	})

	t.Run("unmarshalBinaryUDHConcatMessage (8 bit)", func(t *testing.T) {
		u, rd := new(UDH), []byte{0x05, 0x00, 0x03, 0x0c, 0x02, 0x01}
		read, err := u.UnmarshalBinary(rd)
		require.False(t, read <= 0)

		require.NoError(t, err)

		b, err := u.MarshalBinary()
		require.NoError(t, err)
		require.Equal(t, "0500030c0201", toHex(b))
	})

	t.Run("unmarshalBinaryUDHConcatMessage (16 bit)", func(t *testing.T) {
		u, rd := new(UDH), []byte{0x06, 0x08, 0x04, 0x04, 0xd2, 0x02, 0x01}
		read, err := u.UnmarshalBinary(rd)
		require.False(t, read <= 0)

		require.NoError(t, err)

		totalParts, sequence, reference, found := u.GetConcatInfo()
		require.True(t, found)
		require.Equal(t, totalParts, byte(2))
		require.Equal(t, sequence, byte(1))
		require.Equal(t, reference, uint(1234))

		b, err := u.MarshalBinary()
		require.NoError(t, err)
		require.Equal(t, "06080404d20201", toHex(b))
	})

	t.Run("unmarshalBinaryUDHConcatMessage failed", func(t *testing.T) {
		failedList := [][]byte{
			{0x04, 0x00, 0x02, 0x02, 0x01},       // Invalid 8-bit UDH (wrong length)
			{0x04, 0x08, 0x02, 0x02, 0x01},       // Invalid 16-bit UDH (wrong length)
			{0x05, 0x08, 0x03, 0x04, 0xd2, 0x01}, // Invalid 16-bit UDH (missing part number)
		}
		u := new(UDH)
		for _, data := range failedList {
			_, _ = u.UnmarshalBinary(data)
			_, _, _, found := u.GetConcatInfo()
			require.False(t, found, data)
		}
	})
	t.Run("marshalBinaryTruncateLongIE", func(t *testing.T) {
		u := UDH{NewIEConcatMessage(2, 1, 12)}
		for i := 0; i < 255; i++ {
			u = append(u, NewIEConcatMessage(2, 1, 12))
		}

		require.LessOrEqual(t, u.UDHL(), 256) // UDHL must not exceed 256 ( including UDHL byte )

		_, err := u.MarshalBinary()
		require.Error(t, err)
	})
}
