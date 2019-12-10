package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestDataSM(t *testing.T) {
	v := NewDataSM().(*DataSM)
	require.True(t, v.CanResponse())

	validate(t,
		v.GetResponse(),
		"0000001180000103000000000000000100",
		data.DATA_SM_RESP,
	)

	v.ServiceType = "abc"
	_ = v.SourceAddr.SetAddress("Alicer")
	v.SourceAddr.SetTon(28)
	v.SourceAddr.SetNpi(29)
	_ = v.DestAddr.SetAddress("Bobo")
	v.DestAddr.SetTon(30)
	v.DestAddr.SetNpi(31)
	v.EsmClass = 77
	v.RegisteredDelivery = 83
	v.DataCoding = 91
	v.RegisterOptionalParam(Field{Tag: TagDestBearerType, Data: []byte{95}})

	tagged, ok := v.OptionalParameters[TagDestBearerType]
	require.True(t, ok)
	require.Equal(t, TagDestBearerType, tagged.Tag)
	require.Equal(t, []byte{95}, tagged.Data)

	validate(t,
		v,
		"0000002c000001030000000000000001616263001c1d416c69636572001e1f426f626f004d535b000700015f",
		data.DATA_SM,
	)
}
