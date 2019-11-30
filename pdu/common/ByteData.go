package common

import (
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/utils"
)

// IByteData interface.
type IByteData interface {
	SetData(*utils.ByteBuffer) *Exception.Exception
	GetData() (*utils.ByteBuffer, *Exception.Exception)
}
