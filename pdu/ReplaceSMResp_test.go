package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestReplaceSMResp(t *testing.T) {
	v := NewReplaceSMResp().(*ReplaceSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"00000010800000070000000000000001",
		data.REPLACE_SM_RESP,
	)
}
