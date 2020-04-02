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
	v.SequenceNumber = 13

	validate(t,
		v,
		"0000001080000000000000000000000d",
		data.GENERIC_NACK,
	)
}
