package gosmpp

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	conn, err := net.Dial("tcp", "smscsim.melroselabs.com:2775")
	require.Nil(t, err)

	c := NewConnection(conn)
	defer func() {
		_ = c.Close()
	}()
	t.Log(c.LocalAddr())
	t.Log(c.RemoteAddr())

	require.Nil(t, c.SetDeadline(time.Now().Add(5*time.Second)))
	require.Nil(t, c.SetWriteDeadline(time.Now().Add(5*time.Second)))
	require.Nil(t, c.SetReadDeadline(time.Now().Add(5*time.Second)))
}
