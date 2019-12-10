package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestReplaceSM(t *testing.T) {
	v := NewReplaceSM().(*ReplaceSM)
	require.True(t, v.CanResponse())
	require.True(t, v.Message.withoutDataCoding)

	validate(t,
		v.GetResponse(),
		"00000010800000070000000000000001",
		data.REPLACE_SM_RESP,
	)

	v.MessageID = "ID_Her"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)
	v.RegisteredDelivery = 83
	_ = v.Message.SetMessageWithEncoding("nightwish", data.GSM7BIT)
	v.Message.message = ""

	validate(t,
		v,
		"0000002e00000007000000000000000149445f486572001c1d416c696365720000005300096e6967687477697368",
		data.REPLACE_SM,
	)
}
