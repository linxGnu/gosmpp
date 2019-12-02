package pdu

import (
	"github.com/linxGnu/gosmpp/data"
	"github.com/linxGnu/gosmpp/utils"
)

// UnsuccessSME indicates submission was unsuccessful and the respective errors.
type UnsuccessSME struct {
	Address
	errorStatusCode int32
}

// NewUnsuccessSME returns new UnsuccessSME
func NewUnsuccessSME() (c UnsuccessSME) {
	c = UnsuccessSME{
		Address:         NewAddress(),
		errorStatusCode: data.ESME_ROK,
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
		Address:         NewAddressWithTonNpiLen(ton, npi, len),
		errorStatusCode: data.ESME_ROK,
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

// UnsuccessSMEs represents list of UnsuccessSME.
type UnsuccessSMEs struct {
	l []UnsuccessSME
}

// NewUnsuccessSMEs returns list of UnsuccessSME.
func NewUnsuccessSMEs() (u UnsuccessSMEs) {
	u.l = make([]UnsuccessSME, 0, 8)
	return
}

// Add to list.
func (c *UnsuccessSMEs) Add(u UnsuccessSME) {
	c.l = append(c.l, u)
}

// Get list.
func (c *UnsuccessSMEs) Get() []UnsuccessSME {
	return c.l
}

// Unmarshal from buffer.
func (c *UnsuccessSMEs) Unmarshal(b *utils.ByteBuffer) (err error) {
	var n byte
	if n, err = b.ReadByte(); err == nil {
		c.l = make([]UnsuccessSME, n)

		var i byte
		for ; i < n; i++ {
			if err = c.l[i].Unmarshal(b); err != nil {
				return
			}
		}
	}
	return
}

// Marshal to buffer.
func (c *UnsuccessSMEs) Marshal(b *utils.ByteBuffer) {
	n := byte(len(c.l))
	_ = b.WriteByte(n)

	var i byte
	for ; i < n; i++ {
		c.l[i].Marshal(b)
	}
}
