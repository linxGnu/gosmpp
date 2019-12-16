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
	require.Equal(t, data.GSM7BIT, v.Message.Encoding())

	message, err := v.Message.GetMessage()
	require.Nil(t, err)
	require.Equal(t, "nightwish", message)

	validate(t,
		v,
		"0000002d00000007000000000000000149445f486572001c1d416c69636572000000530008eef4194dbfa7e768",
		data.REPLACE_SM,
	)
}
