package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestEnquireLink(t *testing.T) {
	v := NewEnquireLink().(*EnquireLink)
	require.True(t, v.CanResponse())

	validate(t,
		v.GetResponse(),
		"00000010800000150000000000000001",
		data.ENQUIRE_LINK_RESP,
	)

	validate(t,
		v,
		"00000010000000150000000000000001",
		data.ENQUIRE_LINK,
	)
}
