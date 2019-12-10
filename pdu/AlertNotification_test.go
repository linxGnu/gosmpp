package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"

	"github.com/linxGnu/gosmpp/utils"
)

func TestAlertNotification(t *testing.T) {
	a := NewAlertNotification().(*AlertNotification)
	require.False(t, a.CanResponse())
	require.Nil(t, a.GetResponse())

	_ = a.SourceAddr.SetAddress("Alice")
	a.SourceAddr.SetTon(13)
	a.SourceAddr.SetNpi(15)
	_ = a.EsmeAddr.SetAddress("Bob")
	a.EsmeAddr.SetTon(19)
	a.EsmeAddr.SetNpi(7)

	b := utils.NewBuffer(nil)
	a.Marshal(b)
	b.WriteInt(119)

	c, err := Parse(b)
	require.Nil(t, err)
	require.Equal(t, a, c)
	require.EqualValues(t, data.ALERT_NOTIFICATION, c.(*AlertNotification).CommandID)

}
