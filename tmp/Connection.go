package gosmpp

import (
	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/Utils"
)

// IConnection smpp connection interface
type IConnection interface {
	SetReceiveTimeout(t int64)
	GetReceiveTimeout() int64
	SetSendTimeout(t int64)
	GetSendTimeout() int64
	SetAddress(addr string)
	GetAddress() string
	Open() *Exception.Exception
	IsOpened() bool
	Close() *Exception.Exception
	Send(data *Utils.ByteBuffer) *Exception.Exception
	Receive() (*Utils.ByteBuffer, *Exception.Exception)
}

type Connection struct {
	// receiveTimeout Timeout for receiving data from connection and for accepting new connection.
	receiveTimeout int64

	// sendTimeout Timeout for sending data.
	sendTimeout int64

	// address destination host
	address string
}

func (c *Connection) SetDefault() {
	c.receiveTimeout = Data.CONNECTION_RECEIVE_TIMEOUT
	c.sendTimeout = Data.CONNECTION_SEND_TIMEOUT
}

func (c *Connection) SetReceiveTimeout(t int64) {
	c.receiveTimeout = t
}

func (c *Connection) SetSendTimeout(t int64) {
	c.sendTimeout = t
}

func (c *Connection) GetReceiveTimeout() int64 {
	return c.receiveTimeout
}

func (c *Connection) GetSendTimeout() int64 {
	return c.sendTimeout
}

func (c *Connection) SetAddress(addr string) {
	c.address = addr
}

func (c *Connection) GetAddress() string {
	return c.address
}
