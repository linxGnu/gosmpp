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
		}, 2*time.Second)
	require.NoError(t, err)
	require.Equal(t, s.GetWindowSize(), 0)
	err = s.Close()
	if err != nil {
		t.Log(err)
	}

	s, err = NewSession(
		RXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
		}, 2*time.Second)
	require.NoError(t, err)
	require.Equal(t, s.GetWindowSize(), -1)
	err = s.Close()
	if err != nil {
		t.Log(err)
	}

	s, err = NewSession(
		TRXConnector(NonTLSDialer, auth),
		Settings{
			EnquireLink: 5 * time.Second,
			ReadTimeout: 10 * time.Second,
		}, 2*time.Second)
	require.NoError(t, err)
	require.Equal(t, s.GetWindowSize(), 0)
	err = s.Close()
	if err != nil {
		t.Log(err)
	}
}
