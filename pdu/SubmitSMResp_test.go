package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestSubmitSMResp(t *testing.T) {
	req := NewSubmitSM().(*SubmitSM)
	req.SequenceNumber = 13

	v := NewSubmitSMRespFromReq(req).(*SubmitSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	v.MessageID = "football"

	validate(t,
		v,
		"0000001980000004000000000000000d666f6f7462616c6c00",
		data.SUBMIT_SM_RESP,
	)
}
