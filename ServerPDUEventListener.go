package gosmpp

import "github.com/tsocial/gosmpp/Exception"

type ServerPDUEventListener interface {
	HandleEvent(event *ServerPDUEvent) *Exception.Exception
}
