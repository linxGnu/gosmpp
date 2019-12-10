package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestCancelSM(t *testing.T) {
	v := NewCancelSM().(*CancelSM)
	require.True(t, v.CanResponse())

	resp := v.GetResponse()
	require.NotNil(t, resp)
	require.EqualValues(t, data.CANCEL_SM_RESP, resp.GetHeader().CommandID)

	v.ServiceType = "abc"
	v.MessageID = "def"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)
	_ = v.DestAddr.SetAddress("Bobo")
	v.DestAddr.SetTon(30)
	v.DestAddr.SetNpi(31)

	validate(t,
		v,
		"0000002800000008000000000000000161626300646566001c1d416c69636572001e1f426f626f00",
		data.CANCEL_SM,
	)
}
