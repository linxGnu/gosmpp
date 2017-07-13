package Utils

import "time"

type Unprocessed struct {
	unprocessed      *ByteBuffer
	expected         int32
	hasUnprocessed   bool
	lastTimeReceived int64 // in nano second
}

func NewUnprocessed() *Unprocessed {
	a := &Unprocessed{}
	a.unprocessed = NewBuffer([]byte{})
	a.expected = 0
	a.lastTimeReceived = 0
	a.hasUnprocessed = false

	return a
}

func (c *Unprocessed) Reset() {
	c.hasUnprocessed = false
	c.unprocessed.Reset()
	c.expected = 0
}

func (c *Unprocessed) Check() {
	c.hasUnprocessed = c.unprocessed.Len() > 0
}

func (c *Unprocessed) SetHasUnprocessed(value bool) {
	c.hasUnprocessed = value
}

func (c *Unprocessed) SetExpected(value int32) {
	c.expected = value
}

func (c *Unprocessed) SetLastTimeReceived(value int64) {
	c.lastTimeReceived = value
}

func (c *Unprocessed) SetLastTimeReceivedCurTime() {
	c.lastTimeReceived = time.Now().UnixNano()
}

func (c *Unprocessed) GetUnprocessed() *ByteBuffer {
	return c.unprocessed
}

func (c *Unprocessed) GetHasUnprocessed() bool {
	return c.hasUnprocessed
}

func (c *Unprocessed) GetExpected() int32 {
	return c.expected
}

func (c *Unprocessed) GetLastTimeReceived() int64 {
	return c.lastTimeReceived
}
