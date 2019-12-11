package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestDeliverSMResp(t *testing.T) {
	v := NewDeliverSMResp().(*DeliverSMResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	v.MessageID = "testMID"

	validate(t,
		v,
		"00000018800000050000000000000001746573744d494400",
		data.DELIVER_SM_RESP,
	)
}
