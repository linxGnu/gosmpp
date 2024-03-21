package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressRange(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		a := NewAddressRangeWithAddr("abc")
		require.Equal(t, "abc", a.AddressRange)
	})

	t.Run("newTonNpi", func(t *testing.T) {
		a := NewAddressRangeWithTonNpi(3, 7)
		a.AddressRange = "123456789"
		require.EqualValues(t, 3, a.Ton)
		require.EqualValues(t, 7, a.Npi)
		require.Equal(t, "123456789", a.AddressRange)
	})

	t.Run("newTonNpiAddr", func(t *testing.T) {
		a := NewAddressRangeWithTonNpiAddr(3, 7, "123456789")
		require.EqualValues(t, 3, a.Ton)
		require.EqualValues(t, 7, a.Npi)
		require.Equal(t, "123456789", a.AddressRange)
	})

	t.Run("unmarshal", func(t *testing.T) {
		buf := NewBuffer(fromHex("315b7068616e746f6d537472696b6500"))
		var a AddressRange
		require.Nil(t, a.Unmarshal(buf))
		require.Zero(t, buf.Len())
		require.Equal(t, "phantomStrike", a.AddressRange)
		require.EqualValues(t, 49, a.Ton)
		require.EqualValues(t, 91, a.Npi)
	})

	t.Run("marshal", func(t *testing.T) {
		a := AddressRange{}
		a.AddressRange = "phantomOpera"
		a.Ton = 95
		a.Npi = 13

		buf := NewBuffer(nil)
		a.Marshal(buf)

		require.Equal(t, fromHex("5f0d7068616e746f6d4f7065726100"), buf.Bytes())
	})
}
