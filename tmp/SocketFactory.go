package gosmpp

import (
	"net"
	"strconv"
)

type SocketFactory struct {
}

func (c SocketFactory) CreateTCP(address string, port int) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address+":"+strconv.Itoa(int(port)))
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp", nil, tcpAddr)
}
