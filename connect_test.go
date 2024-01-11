package gosmpp

import (
	"github.com/linxGnu/gosmpp/pdu"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

var currentAuth int32

var auths = [][2]string{
	{"689528", "1a97ae"},
}

const (
	smscAddr = "127.0.0.1:2775"
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
	t.Run("RX", func(t *testing.T) {
		addrRange := pdu.AddressRange{}
		err := addrRange.SetAddressRange("31218")
		require.NoError(t, err)
		checker(t, RXConnector(NonTLSDialer, nextAuth(), WithAddressRange(addrRange)))
	})

	t.Run("TRX", func(t *testing.T) {
		checker(t, TRXConnector(NonTLSDialer, nextAuth()))
	})

	t.Run("TRX", func(t *testing.T) {
		addrRange := pdu.AddressRange{}
		err := addrRange.SetAddressRange("31218")
		require.NoError(t, err)
		checker(t, TRXConnector(NonTLSDialer, nextAuth(), WithAddressRange(addrRange)))
	})
}

func TestBindingSMSC_Error(t *testing.T) {
	auth := Auth{SMSC: smscAddr, SystemID: "invalid"}
	checker := func(t *testing.T, c Connector) {
		conn, err := c.Connect()
		require.ErrorContains(t, err, "Invalid System ID")
		_ = conn.Close()
	}

	t.Run("TX", func(t *testing.T) {
		checker(t, TXConnector(NonTLSDialer, auth))
	})

	t.Run("RX", func(t *testing.T) {
		checker(t, RXConnector(NonTLSDialer, auth))
	})

	t.Run("TRX", func(t *testing.T) {
		checker(t, TRXConnector(NonTLSDialer, auth))
	})
}
