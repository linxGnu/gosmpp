package gosmpp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInvalidSessionSettings(t *testing.T) {
	auth := nextAuth()

	_, err := NewSession(
		TXConnector(NonTLSDialer, auth),
		Settings{}, 2*time.Second)
	require.Error(t, err)

	_, err = NewSession(
		RXConnector(NonTLSDialer, auth),
		Settings{
			ReadTimeout: 200 * time.Millisecond,
			EnquireLink: 333 * time.Millisecond,
		}, 2*time.Second)
	require.Error(t, err)
}

func TestGetWindowSize(t *testing.T) {

	auth := nextAuth()

	s, err := NewSession(
		TXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
			WindowPDUHandlerConfig: &WindowPDUHandlerConfig{
				OnReceivedPduRequest: handleReceivedPduRequest(t),
				MaxWindowSize:        10,
			},
		}, 2*time.Second)
	require.Nil(t, err)
	require.Equal(t, 0, s.GetWindowSize())
	err = s.Close()
	require.Nil(t, err)

	s, err = NewSession(
		RXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
			WindowPDUHandlerConfig: &WindowPDUHandlerConfig{
				OnReceivedPduRequest: handleReceivedPduRequest(t),
				MaxWindowSize:        10,
			},
		}, 2*time.Second)
	require.Nil(t, err)
	require.Equal(t, -1, s.GetWindowSize())
	err = s.Close()
	require.Nil(t, err)

	s, err = NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
			WindowPDUHandlerConfig: &WindowPDUHandlerConfig{
				MaxWindowSize: 10,
			},
		}, 2*time.Second)
	require.NoError(t, err)
	require.Equal(t, 0, s.GetWindowSize())
	err = s.Close()
	require.Nil(t, err)

	s, err = NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
			WindowPDUHandlerConfig: &WindowPDUHandlerConfig{
				ExpireCheckTimer: 5,
				PduExpireTimeOut: 10,
				MaxWindowSize:    10,
			},
		}, 2*time.Second)
	require.NoError(t, err)
	require.Equal(t, 0, s.GetWindowSize())
	time.Sleep(6)
	err = s.Close()
	require.Nil(t, err)
}
