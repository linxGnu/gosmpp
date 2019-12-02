package pdu

import "github.com/linxGnu/gosmpp/utils"

// UnsuccessSME indicates submission was unsuccessful and the respective errors.
type UnsuccessSME struct {
	Address
	errorStatusCode int32
}

// NewUnsuccessSME returns new UnsuccessSME
func NewUnsuccessSME() (c UnsuccessSME) {
	c = UnsuccessSME{
		Address: NewAddress(),
	}
	return
}

// NewUnsuccessSMEWithAddr returns new UnsuccessSME with address.
func NewUnsuccessSMEWithAddr(addr string) (c UnsuccessSME, err error) {
	c = NewUnsuccessSME()
	err = c.SetAddress(addr)
	return
}

// NewUnsuccessSMEWithMaxLength returns UnsuccessSME with max length in C of address.
func NewUnsuccessSMEWithMaxLength(len int) (c UnsuccessSME) {
	c = NewUnsuccessSME()
	c.maxAddressLength = len
	return
}

// NewUnsuccessSMEWithTonNpiLen create new address with ton, npi, max length.
func NewUnsuccessSMEWithTonNpiLen(ton, npi byte, len int) UnsuccessSME {
	return UnsuccessSME{
		Address: NewAddressWithTonNpiLen(ton, npi, len),
	}
}

// Unmarshal from buffer.
func (c *UnsuccessSME) Unmarshal(b *utils.ByteBuffer) (err error) {
	if err = c.Address.Unmarshal(b); err == nil {
		c.errorStatusCode, err = b.ReadInt()
	}
	return
}

// Marshal to buffer.
func (c *UnsuccessSME) Marshal(b *utils.ByteBuffer) {
	c.Address.Marshal(b)
	b.WriteInt(c.errorStatusCode)
}

// SetErrorStatusCode sets error status code.
func (c *UnsuccessSME) SetErrorStatusCode(v int32) {
	c.errorStatusCode = v
}

// ErrorStatusCode returns assigned status code.
func (c *UnsuccessSME) ErrorStatusCode() int32 {
	return c.errorStatusCode
}
