package pdu

import (
	"github.com/linxGnu/gosmpp/data"
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
func NewUnsuccessSMEWithAddr(addr string, errCode int32) (c UnsuccessSME, err error) {
	c = NewUnsuccessSME()
	if err = c.SetAddress(addr); err == nil {
		c.SetErrorStatusCode(errCode)
	}
	return
}

// NewUnsuccessSMEWithTonNpi create new address with ton, npi and error code.
func NewUnsuccessSMEWithTonNpi(ton, npi byte, errCode int32) UnsuccessSME {
	return UnsuccessSME{
		Address:         NewAddressWithTonNpi(ton, npi),
		errorStatusCode: errCode,
	}
}

// Unmarshal from buffer.
func (c *UnsuccessSME) Unmarshal(b *ByteBuffer) (err error) {
	if err = c.Address.Unmarshal(b); err == nil {
		c.errorStatusCode, err = b.ReadInt()
	}
	return
}

// Marshal to buffer.
func (c *UnsuccessSME) Marshal(b *ByteBuffer) {
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
func (c *UnsuccessSMEs) Add(us ...UnsuccessSME) {
	c.l = append(c.l, us...)
}

// Get list.
func (c *UnsuccessSMEs) Get() []UnsuccessSME {
	return c.l
}

// Unmarshal from buffer.
func (c *UnsuccessSMEs) Unmarshal(b *ByteBuffer) (err error) {
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
func (c *UnsuccessSMEs) Marshal(b *ByteBuffer) {
	n := byte(len(c.l))
	_ = b.WriteByte(n)

	var i byte
	for ; i < n; i++ {
		c.l[i].Marshal(b)
	}
}
