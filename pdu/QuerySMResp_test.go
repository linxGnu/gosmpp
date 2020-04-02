package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestQuerySMResp(t *testing.T) {
	req := NewQuerySM().(*QuerySM)
	req.SequenceNumber = 13

	v := NewQuerySMRespFromReq(req).(*QuerySMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"0000001480000003000000000000000d00000000",
		data.QUERY_SM_RESP,
	)
}
