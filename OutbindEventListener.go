package gosmpp

import "github.com/tsocial/gosmpp/Exception"

type OutbindEventListener interface {
	HandleOutbind(outbind *OutbindEvent) *Exception.Exception
}
