package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestQuerySM(t *testing.T) {
	v := NewQuerySM().(*QuerySM)
	require.True(t, v.CanResponse())
	v.SequenceNumber = 13

	validate(t,
		v.GetResponse(),
		"0000001480000003000000000000000d00000000",
		data.QUERY_SM_RESP,
	)

	v.MessageID = "away"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)

	validate(t,
		v,
		"0000001e00000003000000000000000d61776179001c1d416c6963657200",
		data.QUERY_SM,
	)
}
