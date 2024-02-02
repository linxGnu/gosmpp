package gosmpp

import (
	"bufio"
	"net"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp/pdu"
	"github.com/orcaman/concurrent-map/v2"
)

var lock = sync.RWMutex{}

// Connection wraps over net.Conn with buffered data reader.
type Connection struct {
	systemID string
	conn     net.Conn
	reader   *bufio.Reader
	window   cmap.ConcurrentMap[string, pdu.Request]
}

// NewConnection returns a Connection.
func NewConnection(conn net.Conn) (c *Connection) {
	c = &Connection{
		conn:   conn,
		reader: bufio.NewReaderSize(conn, 128<<10),
		window: cmap.New[pdu.Request](),
	}
	return
}

// Read reads data from the connection.
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *Connection) Read(b []byte) (n int, err error) {
	n, err = c.reader.Read(b)
	return
}

// Write writes data to the connection.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *Connection) Write(b []byte) (n int, err error) {
	n, err = c.conn.Write(b)
	return
}

// WritePDU data to the connection.
func (c *Connection) WritePDU(p pdu.PDU) (n int, err error) {
	buf := pdu.NewBuffer(make([]byte, 0, 64))
	p.Marshal(buf)
	n, err = c.conn.Write(buf.Bytes())
	return
}

// Close closes the connection.
// Any blocked Read or Write operations will be unblocked and return errors.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// LocalAddr returns the local network address.
func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
//
// Note that if a TCP connection has keep-alive turned on,
// which is the default unless overridden by Dialer.KeepAlive
// or ListenConfig.KeepAlive, then a keep-alive failure may
// also return a timeout error. On Unix systems a keep-alive
// failure on I/O can be detected using
// errors.Is(err, syscall.ETIMEDOUT).
func (c *Connection) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *Connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// SetReadTimeout is equivalent to ReadDeadline(now + timeout)
func (c *Connection) SetReadTimeout(t time.Duration) error {
	return c.conn.SetReadDeadline(time.Now().Add(t))
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *Connection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

// SetWriteTimeout is equivalent to WriteDeadline(now + timeout)
func (c *Connection) SetWriteTimeout(t time.Duration) error {
	return c.conn.SetWriteDeadline(time.Now().Add(t))
}
