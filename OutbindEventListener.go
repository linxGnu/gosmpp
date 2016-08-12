package gosmpp

import "github.com/linxGnu/gosmpp/Exception"

type OutbindEventListener interface {
	HandleOutbind(outbind *OutbindEvent) *Exception.Exception
}
