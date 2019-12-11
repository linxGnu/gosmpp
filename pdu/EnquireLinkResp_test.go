package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestEnquireLinkResp(t *testing.T) {
	v := NewEnquireLinkResp().(*EnquireLinkResp)
	require.False(t, v.CanResponse())
	require.Nil(t, v.GetResponse())

	validate(t,
		v,
		"00000010800000150000000000000001",
		data.ENQUIRE_LINK_RESP,
	)
}
