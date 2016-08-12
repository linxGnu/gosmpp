package gosmpp

import "github.com/linxGnu/gosmpp/Exception"

type ServerPDUEventListener interface {
	HandleEvent(event *ServerPDUEvent) *Exception.Exception
}
