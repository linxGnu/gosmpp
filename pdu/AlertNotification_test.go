package pdu

import (
	"testing"

	"github.com/linxGnu/gosmpp/data"

	"github.com/stretchr/testify/require"
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

	b := NewBuffer(nil)
	a.Marshal(b)

	expectAfterParse(t, b, a, data.ALERT_NOTIFICATION)
}
