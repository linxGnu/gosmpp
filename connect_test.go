package gosmpp

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

var currentAuth int32

var auths = [][2]string{
	{"415847", "ea445b"},
}

const (
	smscAddr = "smscsim.melroselabs.com:2775"
	mess     = "Thử nghiệm: chuẩn bị nế mễ"
)

func nextAuth() Auth {
	pair := int(atomic.AddInt32(&currentAuth, 1)) % len(auths)
	return Auth{
		SMSC:       smscAddr,
		SystemID:   auths[pair][0],
		Password:   auths[pair][1],
		SystemType: "",
	}
}

func TestBindingSMSC(t *testing.T) {
	checker := func(t *testing.T, c Connector) {
		conn, err := c.Connect()
		require.Nil(t, err)
		require.NotNil(t, conn)
		_ = conn.Close()
	}

	t.Run("TX", func(t *testing.T) {
		checker(t, TXConnector(NonTLSDialer, nextAuth()))
	})

	t.Run("RX", func(t *testing.T) {
		checker(t, RXConnector(NonTLSDialer, nextAuth()))
	})

	t.Run("TRX", func(t *testing.T) {
		checker(t, TRXConnector(NonTLSDialer, nextAuth()))
	})
}
