package gosmpp

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

var currentAuth int32

var auths = [][2]string{
	{"561913", "ff65dc"},
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
	// valid
	connection, err := ConnectAsTransceiver(NonTLSDialer, nextAuth())
	require.Nil(t, err)
	require.NotNil(t, connection)
	_ = connection.Close()

	// valid
	connection, err = ConnectAsReceiver(NonTLSDialer, nextAuth())
	require.Nil(t, err)
	require.NotNil(t, connection)
	_ = connection.Close()

	// valid
	connection, err = ConnectAsTransmitter(NonTLSDialer, nextAuth())
	require.Nil(t, err)
	require.NotNil(t, connection)
	_ = connection.Close()
}
