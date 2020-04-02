package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestUnbind(t *testing.T) {
	v := NewUnbind().(*Unbind)
	require.True(t, v.CanResponse())
	v.SequenceNumber = 13

	validate(t,
		v.GetResponse(),
		"0000001080000006000000000000000d",
		data.UNBIND_RESP,
	)

	validate(t,
		v,
		"0000001000000006000000000000000d",
		data.UNBIND,
	)
}
