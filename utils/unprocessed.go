package utils

import "time"

// Unprocessed wraps over buffer and its state.
type Unprocessed struct {
	*ByteBuffer
	expected         int32
	lastTimeReceived int64 // in nano second
}

// NewUnprocessed returns new Unprocessed.
func NewUnprocessed() (v *Unprocessed) {
	v = &Unprocessed{ByteBuffer: NewBuffer(nil)}
	return
}

// Reset underlying buffer and state.
func (c *Unprocessed) Reset() {
	c.Reset()
	c.expected = 0
}

// SetExpected sets expected.
func (c *Unprocessed) SetExpected(v int32) {
	c.expected = v
}

// SetLastTimeReceived sets timestamp when received data.
func (c *Unprocessed) SetLastTimeReceived(v int64) {
	c.lastTimeReceived = v
}

// SetLastTimeReceivedCurTime sets timestamp to current.
func (c *Unprocessed) SetLastTimeReceivedCurTime() {
	c.lastTimeReceived = time.Now().UnixNano()
}

// HasUnprocessed checks if there is data to read.
func (c *Unprocessed) HasUnprocessed() bool {
	return c.Len() > 0
}

// Expected returns expected size.
func (c *Unprocessed) Expected() int32 {
	return c.expected
}

// LastTimeReceived returns timestamp when data is received.
func (c *Unprocessed) LastTimeReceived() int64 {
	return c.lastTimeReceived
}
