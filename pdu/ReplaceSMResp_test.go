package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestReplaceSMResp(t *testing.T) {
	req := NewReplaceSM().(*ReplaceSM)
	req.SequenceNumber = 13

	v := NewReplaceSMRespFromReq(req).(*ReplaceSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"0000001080000007000000000000000d",
		data.REPLACE_SM_RESP,
	)
}
