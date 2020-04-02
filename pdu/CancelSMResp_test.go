package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestCancelSMResp(t *testing.T) {
	req := NewCancelSM().(*CancelSM)
	req.SequenceNumber = 11

	v := NewCancelSMRespFromReq(req).(*CancelSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"0000001080000008000000000000000b",
		data.CANCEL_SM_RESP,
	)
}
