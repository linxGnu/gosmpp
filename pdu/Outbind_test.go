package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestOutbind(t *testing.T) {
	v := NewOutbind().(*Outbind)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())
	require.True(t, v.IsOk())

	v.SystemID = "inventory"
	v.Password = "ipassword"

	validate(t,
		v,
		"000000240000000b0000000000000001696e76656e746f7279006970617373776f726400",
		data.OUTBIND,
	)
}
