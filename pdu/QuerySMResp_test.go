package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestQuerySMResp(t *testing.T) {
	v := NewQuerySMResp().(*QuerySMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"0000001480000003000000000000000100000000",
		data.QUERY_SM_RESP,
	)
}
