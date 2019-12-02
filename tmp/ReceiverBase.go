package gosmpp

import (
	"context"

	"github.com/linxGnu/gosmpp/Data"
)

// base struct for receiver.
type receiverBase struct {
	ctx                         context.Context
	cancel                      context.CancelFunc
	receiveTimeout              int64
	messageIncompleteRetryCount byte
	r                           IReceiver
}

func newReceiverBase(receiver IReceiver) (r *receiverBase) {
	r = &receiverBase{
		r:              receiver,
		receiveTimeout: Data.RECEIVER_TIMEOUT,
	}
	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.init()
	return r
}

func (c *receiverBase) init() {
	c.receiveTimeout = Data.RECEIVER_TIMEOUT
	c.messageIncompleteRetryCount = 0
}

func (c *receiverBase) start() {
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return

			default:
				if c.r != nil {
					c.r.Receive()
				}
			}
		}
	}()
}

func (c *receiverBase) stop() {
	c.cancel()
}

func (c *receiverBase) SetReceiveTimeout(timeout int64) {
	c.receiveTimeout = timeout
}
