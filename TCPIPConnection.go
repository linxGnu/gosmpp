package gosmpp

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

const (
	// CONN_SERVER              byte  = 2
	CONN_NONE                byte = 0
	CONN_CLIENT              byte = 1
	DFLT_IO_BUF_SIZE         int  = 2 * 1024
	DFLT_RECEIVE_BUFFER_SIZE int  = 4 * 1024
	DFLT_MAX_RECEIVE_SIZE    int  = 128 * 1024
)

type TCPIPConnection struct {
	Connection
	port              int
	socket            *net.TCPConn
	opened            bool
	connType          byte
	ioBufferSize      int
	receiveBufferSize int
	receiveBuffer     []byte
	maxReceiveSize    int
	socketFactory     *SocketFactory
	lock              sync.RWMutex
}

// NewTCPIPConnection new tcp ip connection
func NewTCPIPConnection() *TCPIPConnection {
	a := &TCPIPConnection{}
	a.SetDefault()
	a.port = 0
	a.socket = nil
	a.opened = false
	a.connType = CONN_NONE
	a.ioBufferSize = DFLT_IO_BUF_SIZE
	a.maxReceiveSize = DFLT_MAX_RECEIVE_SIZE
	a.socketFactory = &SocketFactory{}

	return a
}

// NewTCPIPConnectionWithAddrPort new tcp/ip connection with addr and port
func NewTCPIPConnectionWithAddrPort(addr string, port int) (*TCPIPConnection, error) {
	if port < int(Data.MIN_VALUE_PORT) || port > int(Data.MAX_VALUE_PORT) {
		return nil, errors.New("TCPIPConnection: connection port is invalid")
	}

	addrLength := len(addr)
	if addrLength < int(Data.MIN_LENGTH_ADDRESS) {
		return nil, errors.New("TCPIPConnection: connection address length is invalid")
	}

	a := NewTCPIPConnection()
	a.port = port
	a.address = addr
	a.connType = CONN_CLIENT
	a.SetReceiveBufferSize(DFLT_RECEIVE_BUFFER_SIZE)

	return a, nil
}

// NewTCPIPConnectionWithSocket new tcp/ip connection with preallocated socket. Useful for tls enabled connection.
func NewTCPIPConnectionWithSocket(soc *net.TCPConn) (*TCPIPConnection, error) {
	if soc == nil || soc.RemoteAddr() == nil {
		return nil, errors.New("TCPIPConnection: socket init is nil")
	}

	a := NewTCPIPConnection()
	a.connType = CONN_CLIENT
	a.socket = soc

	addresstmp := soc.RemoteAddr().String()
	add, _port, err := net.SplitHostPort(addresstmp)
	if err != nil {
		return nil, err
	}
	a.address = add

	port, err := strconv.Atoi(_port)
	if err != nil {
		return nil, err
	}
	a.port = port
	a.opened = true
	a.SetReceiveBufferSize(DFLT_RECEIVE_BUFFER_SIZE)

	return a, nil
}

// Open connection
func (c *TCPIPConnection) Open() *Exception.Exception {
	if !c.IsOpened() {
		if c.connType == CONN_CLIENT {
			socket, err := c.socketFactory.CreateTCP(c.address, c.port)
			if err != nil {
				return Exception.NewException(err)
			}
			c.socket = socket
			c.socket.SetKeepAlive(true)
			c.socket.SetKeepAlivePeriod(5 * time.Second)
			c.socket.SetReadBuffer(int(c.receiveBufferSize))
			c.socket.SetWriteBuffer(int(c.receiveBufferSize))
			c.opened = true
		} else {
			return Exception.NewExceptionFromStr("TCPIPConnection: unknown connection type = " + strconv.Itoa(int(c.connType)))
		}
	}

	return nil
}

// Close connection
func (c *TCPIPConnection) Close() *Exception.Exception {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.connType == CONN_CLIENT {
		defer func() {
			c.socket = nil
			c.opened = false
		}()

		if c.socket != nil {
			fmt.Fprintf(os.Stdout, "Closing connection to "+c.GetAddress()+"\n")
			err := c.socket.Close()
			if err != nil {
				fmt.Fprintf(os.Stdout, "Connection closed error: "+err.Error()+"\n")
				return Exception.NewException(err)
			}
			fmt.Fprintf(os.Stdout, "Connection closed successfully")
		}
	}

	return nil
}

// Send buffered data
func (c *TCPIPConnection) Send(data *Utils.ByteBuffer) (err *Exception.Exception) {
	defer func() {
		if e := recover(); e != nil {
			err = Exception.NewException(fmt.Errorf("%v", e))
		}
	}()

	if !c.IsOpened() {
		return Exception.EXCEPTION_TCP_NOT_OPEN
	}

	if data == nil {
		return nil
	}

	if c.connType == CONN_CLIENT {
		sendTimeout := c.GetSendTimeout()
		if sendTimeout >= 0 {
			c.socket.SetWriteDeadline(time.Now().Add(time.Duration(sendTimeout) * time.Millisecond))
		}

		_, err := c.socket.Write(data.Bytes())
		if err != nil {
			c.Close()
			return Exception.NewException(err)
		}

		return nil
	}

	return nil
}

// Receive message from connection in form of buffer
func (c *TCPIPConnection) Receive() (result *Utils.ByteBuffer, err *Exception.Exception) {
	defer func() {
		if e := recover(); e != nil {
			err = Exception.NewException(fmt.Errorf("%v", e))
			result = nil
		}
	}()

	if !c.IsOpened() {
		return nil, Exception.EXCEPTION_TCP_NOT_OPEN
	}

	if c.connType == CONN_CLIENT {
		recBufferLen := len(c.receiveBuffer)
		data := make([]byte, c.maxReceiveSize+recBufferLen)

		recTimeout := c.GetReceiveTimeout()
		deadLineRead := time.Now().Add(time.Duration(recTimeout) * time.Millisecond)

		totalRead := 0

		if recTimeout >= 0 {
			c.socket.SetReadDeadline(deadLineRead)
		}

		for totalRead < c.maxReceiveSize && (recTimeout < 0 || time.Now().Before(deadLineRead)) {
			bytesRead, e := c.socket.Read(c.receiveBuffer)

			if e != nil {
				if nerr, ok := e.(net.Error); ok && nerr.Timeout() {
					fmt.Println(e)
					break
				}

				c.Close()
				return nil, Exception.ConnectionClosingDueToError
			}

			if bytesRead > 0 {
				copy(data[totalRead:], c.receiveBuffer[:bytesRead])
				totalRead += bytesRead
			}

			if bytesRead < recBufferLen {
				break
			}
		}

		return Utils.NewBuffer(data[:totalRead]), nil
	}

	return nil, nil
}

// SetReceiveBufferSize set buffer size for receiving message over socket
func (c *TCPIPConnection) SetReceiveBufferSize(size int) {
	c.receiveBufferSize = size
	c.receiveBuffer = make([]byte, size)
}

// SetIOBufferSize set io buffer size
func (c *TCPIPConnection) SetIOBufferSize(size int) {
	if !c.IsOpened() {
		c.ioBufferSize = size
	}
}

// SetMaxReceiveSize set max length of message allowed to receive over socket
func (c *TCPIPConnection) SetMaxReceiveSize(size int) {
	c.maxReceiveSize = size
}

// IsOpened check if connection is opened
func (c *TCPIPConnection) IsOpened() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.opened
}
