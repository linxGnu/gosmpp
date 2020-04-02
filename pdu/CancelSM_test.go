package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestCancelSM(t *testing.T) {
	v := NewCancelSM().(*CancelSM)
	require.True(t, v.CanResponse())
	v.SequenceNumber = 13

	validate(t,
		v.GetResponse(),
		"0000001080000008000000000000000d",
		data.CANCEL_SM_RESP,
	)

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
		"0000002800000008000000000000000d61626300646566001c1d416c69636572001e1f426f626f00",
		data.CANCEL_SM,
	)
}
