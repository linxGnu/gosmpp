package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestDeliverSMResp(t *testing.T) {
	req := NewDeliverSM().(*DeliverSM)
	req.SequenceNumber = 13

	v := NewDeliverSMRespFromReq(req).(*DeliverSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	v.MessageID = "testMID"

	validate(t,
		v,
		"0000001880000005000000000000000d746573744d494400",
		data.DELIVER_SM_RESP,
	)
}
