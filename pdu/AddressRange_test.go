package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressRange(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		a, err := NewAddressRangeWithAddr("abc")
		require.Nil(t, err)
		require.Equal(t, "abc", a.AddressRange())
	})

	t.Run("newWithAddr", func(t *testing.T) {
		_, err := NewAddressRangeWithAddr("12345678901234567890121234567890123456789012")
		require.NotNil(t, err)
	})

	t.Run("newTonNpi", func(t *testing.T) {
		a := NewAddressRangeWithTonNpi(3, 7)
		require.Nil(t, a.SetAddressRange("123456789"))
		require.EqualValues(t, 3, a.Ton())
		require.EqualValues(t, 7, a.Npi())
		require.Equal(t, "123456789", a.AddressRange())
		a.SetTon(11)
		a.SetNpi(19)
		require.EqualValues(t, 11, a.Ton())
		require.EqualValues(t, 19, a.Npi())
	})

	t.Run("newTonNpiAddr", func(t *testing.T) {
		a, err := NewAddressRangeWithTonNpiAddr(3, 7, "123456789")
		require.Nil(t, err)
		require.EqualValues(t, 3, a.Ton())
		require.EqualValues(t, 7, a.Npi())
		require.Equal(t, "123456789", a.AddressRange())
		a.SetTon(11)
		a.SetNpi(19)
		require.EqualValues(t, 11, a.Ton())
		require.EqualValues(t, 19, a.Npi())
	})

	t.Run("unmarshal", func(t *testing.T) {
		buf := NewBuffer(fromHex("315b7068616e746f6d537472696b6500"))
		var a AddressRange
		require.Nil(t, a.Unmarshal(buf))
		require.Zero(t, buf.Len())
		require.Equal(t, "phantomStrike", a.AddressRange())
		require.EqualValues(t, 49, a.Ton())
		require.EqualValues(t, 91, a.Npi())
	})

	t.Run("marshal", func(t *testing.T) {
		a, err := NewAddressRangeWithAddr("phantomOpera")
		require.Nil(t, err)
		a.SetTon(95)
		a.SetNpi(13)

		buf := NewBuffer(nil)
		a.Marshal(buf)

		require.Equal(t, fromHex("5f0d7068616e746f6d4f7065726100"), buf.Bytes())
	})
}
