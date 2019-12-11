package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestUnbindResp(t *testing.T) {
	v := NewUnbindResp().(*UnbindResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"00000010800000060000000000000001",
		data.UNBIND_RESP,
	)
}
