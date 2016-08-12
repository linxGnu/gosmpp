package gosmpp

import (
	"fmt"
	"os"
	// "runtime/debug"
	"sync"

	"github.com/linxGnu/gosmpp/Data"
	"github.com/linxGnu/gosmpp/Exception"
	"github.com/linxGnu/gosmpp/PDU"
	"github.com/linxGnu/gosmpp/Utils"
)

const (
	RECEIVER_THREAD_NAME = "Receiver"
	RECEIVE_CHAN_SIZE    = 10000
)

type Receiver struct {
	ReceiverBase
	transmitter         *Transmitter
	connection          IConnection
	unprocessed         *Utils.Unprocessed
	pduListener         ServerPDUEventListener
	asynchronousSending bool
	automaticNack       bool // If true then GenericNack messages will be sent automatically if message can't be parsed
	lock                sync.Mutex
}

func NewReceiver() *Receiver {
	a := &Receiver{}
	a.ReceiverBase.Construct()
	a.automaticNack = true
	a.asynchronousSending = false
	a.unprocessed = Utils.NewUnprocessed()
	a.RegisterReceiver(a)

	return a
}

func NewReceiverWithConnection(con IConnection) *Receiver {
	a := NewReceiver()
	a.connection = con

	return a
}

func NewReceiverWithTransmitterCon(trans *Transmitter, con IConnection) *Receiver {
	a := NewReceiverWithConnection(con)
	a.transmitter = trans

	return a
}

func (c *Receiver) SetServerPDUEventListener(listener ServerPDUEventListener) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.pduListener = listener
	c.asynchronousSending = c.pduListener != nil
}

/**
 * Resets unprocessed data and starts receiving on the background.
 *
 * @see ReceiverBase#start()
 */
func (c *Receiver) Start() {
	c.unprocessed.Reset()
	c.StartProcess()
}

func (c *Receiver) Stop() {
	c.StopProcess()
}

func (c *Receiver) StopByException(e *Exception.Exception) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e.Error)
	}

	c.Stop()
}

/**
 * This method receives a PDU or returns PDU received on background,
 * if there is any. It tries to receive a PDU for the specified timeout.
 * If the receiver is asynchronous, then no attempt to receive a PDU
 * and <code>null</code> is returned.
 * The function calls are nested as follows:<br>
 * <ul>
 *   <li>No background receiver thread<br><code>
 *       Receiver.receive(long)<br>
 *       ReceiverBase.tryReceivePDUWithTimeout(Connection,PDU,long)<br>
 *       Receiver.tryReceivePDU(Connection,PDU)<br>
 *       ReceiverBase.receivePDUFromConnection<br>
 *       Connection.receive()</code>
 *   <li>Has background receiver thread<br><code>
 *       Receiver.receive(long)<br>
 *       ReceiverBase.tryReceivePDUWithTimeout(Connection,PDU,long)<br>
 *       Receiver.tryReceivePDU(Connection,PDU)<br>
 *       Queue.dequeue(PDU)</code><br>
 *       and the ReceiverBase.run() function which actually receives the
 *       PDUs and stores them to a queue looks as follows:<br><code>
 *       ReceiverBase.run()<br>
 *       Receiver.receiveAsync()<br>
 *       ReceiverBase.receivePDUFromConnection<br>
 *       Connection.receive()</code>
 *
 * @param timeout for how long is tried to receive a PDU
 * @return the received PDU or null if none received for the spec. timeout
 *
 * @exception IOException exception during communication
 * @exception PDUException incorrect format of PDU
 * @exception TimeoutException rest of PDU not received for too long time
 * @exception UnknownCommandIdException PDU with unknown id was received
 * @see ReceiverBase#tryReceivePDUWithTimeout(Connection,PDU,long)
 */
func (c *Receiver) ReceiveSyncWTimeout(timeout int64) (PDU.IPDU, *Exception.Exception) {
	if c.asynchronousSending {
		return nil, nil
	}

	return c.tryReceivePDUWithCustomTimeout(c.connection, nil, timeout)
}

/**
 * Called from session to receive a response for previously sent request.
 *
 * @param expectedPDU the template for expected PDU; the PDU returned
 *                    must have the same sequence number
 * @return the received PDU or null if none
 * @see ReceiverBase#tryReceivePDUWithTimeout(Connection,PDU,long)
 */
func (c *Receiver) ReceiveSyncWithExpectedPDU(pdu PDU.IPDU) (PDU.IPDU, *Exception.Exception) {
	if c.asynchronousSending {
		return nil, nil
	}

	return c.tryReceivePDUWithTimeout(c.connection, pdu)
}

func (c *Receiver) TryReceivePDU(conn IConnection, expected PDU.IPDU) (pduResult PDU.IPDU, expc *Exception.Exception) {
	defer func() {
		if errs := recover(); errs != nil {
			pduResult = nil
			expc = Exception.NewException(fmt.Errorf("%v", errs))
		}
	}()

	// debug.PrintStack()

	var pdu PDU.IPDU
	pdu, err := c.ReceivePDUFromConnection(c.connection, c.unprocessed)
	if err != nil {
		return pdu, err
	}

	if expected != nil && pdu.IsEquals(expected) {
		return pdu, nil
	}

	return nil, nil
}

/**
 * This method receives a PDU from connection and stores it into
 * <code>pduQueue</code>. It's called from the <code>ReceiverBase</code>'s
 * p<code>process</code> method which is called in loop from
 * <code>ProcessingThread</code>'s <code>run</code> method.
 * <p>
 * If an exception occurs during receiving, depending on type
 * of the exception this method either just reports the exception to
 * debug & event objects or stops processing to indicate
 * that it isn't able to process the exception. The function
 * <code>setTermException</code> is then called with the caught exception.
 *
 * @see ReceiverBase#run()
 */
func (c *Receiver) ReceiveAsync() {
	defer func() {
		if errs := recover(); errs != nil {
			c.Stop()
		}

		// debug.PrintStack()
	}()

	if c.connection != nil && !c.connection.IsOpened() {
		c.Stop()
		return
	}

	pdu, err := c.ReceivePDUFromConnection(c.connection, c.unprocessed)
	if err != nil {
		if err == Exception.InvalidPDUException {
			var seqNr int32
			if pdu != nil {
				seqNr = pdu.GetSequenceNumber()
			}

			if c.automaticNack {
				c.sendGenericNack(Data.ESME_RINVMSGLEN, seqNr)
			} else {
				pdu = PDU.NewGenericNackWithCmStatusSeqNum(Data.ESME_RINVMSGLEN, seqNr)
			}
		} else if err == Exception.UnknownCommandIdException {
			var seqNr int32
			if pdu != nil {
				seqNr = pdu.GetSequenceNumber()
			}

			if c.automaticNack {
				c.sendGenericNack(Data.ESME_RINVCMDID, seqNr)
			} else {
				pdu = PDU.NewGenericNackWithCmStatusSeqNum(Data.ESME_RINVCMDID, seqNr)
			}
		} else if pdu != nil {
			if c.automaticNack {
				c.sendGenericNack(err.ErrorCode, pdu.GetSequenceNumber())
			} else {
				pdu = PDU.NewGenericNackWithCmStatusSeqNum(err.ErrorCode, pdu.GetSequenceNumber())
			}
		}
	}

	if pdu != nil {
		if c.asynchronousSending {
			c.handle(pdu)
		}
	}
}

func (c *Receiver) handle(pdu PDU.IPDU) {
	if pdu == nil {
		return
	}

	if c.pduListener != nil {
		c.pduListener.HandleEvent(NewServerPDUEvent(c, c.connection, pdu))
	} else {
		t, _, _ := pdu.GetData()
		if t != nil {
			fmt.Fprintf(os.Stdout, "async receiver doesn't have ServerPDUEventListener, "+"discarding "+t.GetHexDump()+"\n")
		}
	}
}

func (c *Receiver) sendGenericNack(commandStatus, seqNum int32) {
	if c.transmitter != nil {
		gnack := PDU.NewGenericNackWithCmStatusSeqNum(commandStatus, seqNum)
		c.transmitter.Send(gnack)
	}
}
