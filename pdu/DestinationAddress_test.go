package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestDestinationAddress(t *testing.T) {
	t.Run("validDESTAddr", func(t *testing.T) {
		d1, err := NewDestinationAddressFromAddress("Bob1")
		require.Nil(t, err)
		require.EqualValues(t, data.SM_DEST_SME_ADDRESS, d1.destFlag)
		require.True(t, d1.IsAddress())
		require.False(t, d1.IsDistributionList())
		require.True(t, d1.HasValue())
		require.Equal(t, "Bob1", d1.Address().Address())
		require.Equal(t, "", d1.DistributionList().Name())

		d2, err := NewDestinationAddressFromDistributionList("List1")
		require.Nil(t, err)
		require.EqualValues(t, data.SM_DEST_DL_NAME, d2.destFlag)
		require.False(t, d2.IsAddress())
		require.True(t, d2.IsDistributionList())
		require.True(t, d2.HasValue())
		require.Equal(t, "", d2.Address().Address())
		require.Equal(t, "List1", d2.DistributionList().Name())

		d3 := NewDestinationAddress()
		require.EqualValues(t, data.DFLT_DEST_FLAG, d3.destFlag)
		require.False(t, d3.IsAddress())
		require.False(t, d3.IsDistributionList())
		require.False(t, d3.HasValue())
		require.Nil(t, d3.SetDistributionList("List2"))
		require.Equal(t, "", d3.Address().Address())
		require.Equal(t, "List2", d3.DistributionList().Name())
		require.False(t, d3.IsAddress())
		require.True(t, d3.IsDistributionList())
		require.EqualValues(t, data.SM_DEST_DL_NAME, d3.destFlag)
	})

	t.Run("invalidDEST", func(t *testing.T) {
		buf := NewBuffer(nil)
		_ = buf.WriteByte(51)
		var d DestinationAddress
		require.NotNil(t, d.Unmarshal(buf))
	})

	t.Run("invalidDESTs", func(t *testing.T) {
		buf := NewBuffer(nil)
		_ = buf.WriteByte(1)
		_ = buf.WriteByte(51)
		var d DestinationAddresses
		require.NotNil(t, d.Unmarshal(buf))
	})
}
