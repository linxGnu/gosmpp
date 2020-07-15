package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
)

func TestDestinationAddress(t *testing.T) {
	t.Run("validDESTAddr", func(t *testing.T) {
		addr := NewAddress()
		require.Nil(t, addr.SetAddress("Bob1"))
		d1 := NewDestinationAddress()
		d1.SetAddress(addr)
		require.EqualValues(t, data.SM_DEST_SME_ADDRESS, d1.destFlag)
		require.True(t, d1.IsAddress())
		require.False(t, d1.IsDistributionList())
		require.True(t, d1.HasValue())
		require.Equal(t, "Bob1", d1.Address().Address())
		require.Equal(t, "", d1.DistributionList().Name())

		dl, err := NewDistributionList("List1")
		require.Nil(t, err)
		d2 := NewDestinationAddress()
		d2.SetDistributionList(dl)
		require.EqualValues(t, data.SM_DEST_DL_NAME, d2.destFlag)
		require.False(t, d2.IsAddress())
		require.True(t, d2.IsDistributionList())
		require.True(t, d2.HasValue())
		require.Equal(t, "", d2.Address().Address())
		require.Equal(t, "List1", d2.DistributionList().Name())
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
