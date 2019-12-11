package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"
	"github.com/stretchr/testify/require"
)

func TestGNack(t *testing.T) {
	v := NewGenerickNack().(*GenerickNack)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())
	require.True(t, v.IsGNack())

	validate(t,
		v,
		"00000010800000000000000000000001",
		data.GENERIC_NACK,
	)
}
