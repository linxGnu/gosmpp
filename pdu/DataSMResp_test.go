package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestDataSMResp(t *testing.T) {
	req := NewDataSM().(*DataSM)
	req.SequenceNumber = 13

	v := NewDataSMRespFromReq(req).(*DataSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	v.MessageID = "testMID"

	validate(t,
		v,
		"0000001880000103000000000000000d746573744d494400",
		data.DATA_SM_RESP,
	)
}
