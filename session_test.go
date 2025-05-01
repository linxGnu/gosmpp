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
			WindowedRequestTracking: &WindowedRequestTracking{
				OnReceivedPduRequest: handleReceivedPduRequest(t),
				MaxWindowSize:        10,
				StoreAccessTimeOut:   100 * time.Millisecond,
			},
		}, 2*time.Second)
	require.Nil(t, err)
	size, err := s.GetWindowSize()
	if err != nil {
		return
	}
	require.Nil(t, err)
	require.Equal(t, 0, size)
	err = s.Close()
	require.Nil(t, err)

	s, err = NewSession(
		RXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
			WindowedRequestTracking: &WindowedRequestTracking{
				OnReceivedPduRequest: handleReceivedPduRequest(t),
				MaxWindowSize:        10,
				StoreAccessTimeOut:   100 * time.Millisecond,
			},
		}, 2*time.Second)
	require.Nil(t, err)
	size, err = s.GetWindowSize()
	if err != nil {
		return
	}
	require.Nil(t, err)
	require.Equal(t, -1, size)
	err = s.Close()
	require.Nil(t, err)

	s, err = NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
			WindowedRequestTracking: &WindowedRequestTracking{
				MaxWindowSize:      10,
				StoreAccessTimeOut: 100 * time.Millisecond,
			},
		}, 2*time.Second)
	require.NoError(t, err)
	size, err = s.GetWindowSize()
	if err != nil {
		return
	}
	require.Nil(t, err)
	require.Equal(t, 0, size)
	err = s.Close()
	require.Nil(t, err)

	s, err = NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
			WindowedRequestTracking: &WindowedRequestTracking{
				ExpireCheckTimer:   5 * time.Second,
				PduExpireTimeOut:   10 * time.Second,
				MaxWindowSize:      10,
				StoreAccessTimeOut: 100 * time.Millisecond,
			},
		}, 2*time.Second)
	require.NoError(t, err)
	size, err = s.GetWindowSize()
	if err != nil {
		return
	}
	require.Nil(t, err)
	require.Equal(t, -1, size)
	err = s.Close()
	require.Nil(t, err)
}
