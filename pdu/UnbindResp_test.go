package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestUnbindResp(t *testing.T) {
	req := NewUnbind().(*Unbind)
	req.SequenceNumber = 13

	v := NewUnbindRespFromReq(req).(*UnbindResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"0000001080000006000000000000000d",
		data.UNBIND_RESP,
	)
}
