package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// BindingType indicates type of binding.
type BindingType byte

const (
	// Receiver indicates Receiver binding.
	Receiver BindingType = iota
	// Transceiver indicates Transceiver binding.
	Transceiver
	// Transmitter indicate Transmitter binding.
	Transmitter
)

// BindRequest represents a bind request.
type BindRequest struct {
	base
	SystemID         string
	Password         string
	SystemType       string
	AddressRange     AddressRange
	InterfaceVersion byte
	Type             BindingType
}

// NewBindRequest returns new bind request.
func NewBindRequest(t BindingType) (b *BindRequest) {
	b = &BindRequest{
		base:             newBase(),
		Type:             t,
		SystemID:         data.DFLT_SYSID,
		Password:         data.DFLT_PASS,
		SystemType:       data.DFLT_SYSTYPE,
		AddressRange:     *NewAddressRange(),
		InterfaceVersion: data.SMPP_V34,
	}

	switch t {
	case Transceiver:
		b.CommandID = data.BIND_TRANSCEIVER

	case Receiver:
		b.CommandID = data.BIND_RECEIVER

	case Transmitter:
		b.CommandID = data.BIND_TRANSMITTER
	}

	return
}

// CanResponse implements PDU interface.
func (b *BindRequest) CanResponse() bool {
	return true
}

// GetResponse implements PDU interface.
func (b *BindRequest) GetResponse() PDU {
	return NewBindResp(b)
}

// Marshal implements PDU interface.
func (b *BindRequest) Marshal(w *utils.ByteBuffer) {
	b.base.marshal(w, func(w *utils.ByteBuffer) {
		w.Grow(len(b.SystemID) + len(b.Password) + len(b.SystemType) + 4)
		w.WriteCString(b.SystemID)
		w.WriteCString(b.Password)
		w.WriteCString(b.SystemType)
		_ = w.WriteByte(b.InterfaceVersion)
		b.AddressRange.Marshal(w)
	})
}

// Unmarshal implements PDU interface.
func (b *BindRequest) Unmarshal(w *utils.ByteBuffer) error {
	return b.base.unmarshal(w, func(w *utils.ByteBuffer) (err error) {
		if b.SystemID, err = w.ReadCString(); err == nil {
			if b.Password, err = w.ReadCString(); err == nil {
				if b.SystemType, err = w.ReadCString(); err == nil {
					if b.InterfaceVersion, err = w.ReadByte(); err == nil {
						err = b.AddressRange.Unmarshal(w)
					}
				}
			}
		}
		return
	})
}
