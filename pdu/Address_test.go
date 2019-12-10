package pdu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress(t *testing.T) {
	a, err := NewAddressWithAddr("abc")
	require.Nil(t, err)
	require.Equal(t, "abc", a.Address())

	a, err = NewAddressWithAddr("1234567890123456789012")
	require.NotNil(t, err)

	a = NewAddressWithMaxLength(10)
	require.NotNil(t, a.SetAddress("12345678901"))

	a = NewAddressWithTonNpiLen(3, 7, 9)
	require.Nil(t, a.SetAddress("123456789"))
	require.EqualValues(t, 3, a.Ton())
	require.EqualValues(t, 7, a.Npi())
	require.EqualValues(t, 9, a.maxAddressLength)
	require.Equal(t, "123456789", a.Address())
	a.SetTon(11)
	a.SetNpi(19)
	require.EqualValues(t, 11, a.Ton())
	require.EqualValues(t, 19, a.Npi())
}
