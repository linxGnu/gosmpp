package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestUnbind(t *testing.T) {
	v := NewUnbind().(*Unbind)
	require.True(t, v.CanResponse())

	validate(t,
		v.GetResponse(),
		"00000010800000060000000000000001",
		data.UNBIND_RESP,
	)

	validate(t,
		v,
		"00000010000000060000000000000001",
		data.UNBIND,
	)
}
