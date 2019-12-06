package errors

import (
	"fmt"

	"github.com/linxGnu/gosmpp/data"
)

// SmppErr indicates smpp error(s), compatible with OpenSMPP.
type SmppErr struct {
	err              string
	serialVersionUID int64
}

// Error interface.
func (s *SmppErr) Error() string {
	return fmt.Sprintf("Error happened: [%s]. SerialVersionUID: [%d]", s.err, s.serialVersionUID)
}

var (
	// ErrTCPNotOpened indicates tcp connection not opened.
	ErrTCPNotOpened error = &SmppErr{err: "TCP connection not opened"}

	// ErrUnknown indicates ESME_RUNKNOWNERR.
	ErrUnknown error = &SmppErr{err: "Unknown", serialVersionUID: 5174606627714918071}

	// ErrHeaderIncomplete indicates header incompletion.
	ErrHeaderIncomplete error = &SmppErr{err: "Incomplete header", serialVersionUID: 3813030249230646966}

	// ErrIntegerOutOfRange indicates integer out of range error.
	ErrIntegerOutOfRange error = &SmppErr{err: "Integer out of range", serialVersionUID: 750364511680192335}

	// ErrInvalidPDU indicates invalid pdu payload.
	ErrInvalidPDU error = &SmppErr{err: "PDU payload is invalid", serialVersionUID: -6985061862208729984}

	// ErrMessageIncomplete indicates invalid message payload.
	ErrMessageIncomplete error = &SmppErr{err: "Message payload is incomplete", serialVersionUID: 9081917316954328823}

	// ErrTLV indicates invalid TLV payload.
	ErrTLV error = &SmppErr{err: "TLV payload is invalid", serialVersionUID: -6659626685298184198}

	// ErrTooManyValues indicates too many values.
	ErrTooManyValues error = &SmppErr{err: "Too many values", serialVersionUID: -2777016699062489252}

	// ErrUnexpectedOptionalParameter indicates unexpected optional parameter.
	ErrUnexpectedOptionalParameter error = &SmppErr{err: "Unexpected optional parameter", serialVersionUID: -1284359967986779783}

	// ErrUnknownCommandID indicates unknown command id.
	ErrUnknownCommandID error = &SmppErr{err: "Unknown command id", serialVersionUID: -5091873576710864441}

	// ErrValueNotSet indicates value not set.
	ErrValueNotSet error = &SmppErr{err: "Value not set", serialVersionUID: -4595064103809398438}

	// ErrWrongDateFormat indicates wrong date format.
	ErrWrongDateFormat error = &SmppErr{err: "Wrong date format", serialVersionUID: 5831937612139037591}

	// ErrWrongDestFlag indicates wrong destination flag.
	ErrWrongDestFlag error = &SmppErr{err: "Wrong destination flag", serialVersionUID: 6266749651012701472}

	// ErrWrongLength indicates wrong length.
	ErrWrongLength error = &SmppErr{err: "Wrong length", serialVersionUID: 7935018427341458286}

	// ErrWrongStringLength indicates wrong string length.
	ErrWrongStringLength error = &SmppErr{err: "Wrong string length", serialVersionUID: 8604133584902790266}

	// ErrShortMessageLengthTooLarge indicates short message length is too large.
	ErrShortMessageLengthTooLarge error = &SmppErr{err: fmt.Sprintf("Encoded short message data exceeds size of %d", data.SM_MSG_LEN), serialVersionUID: 78237205927624}

	// ErrBufferFull indicates buffer full.
	ErrBufferFull error = &SmppErr{err: "Buffer is full", serialVersionUID: -3720107899765064964}

	// ErrTerminatingZeroNotFound indicates terminating byte (zero) not found.
	ErrTerminatingZeroNotFound error = &SmppErr{err: "Terminating byte (zero) not found", serialVersionUID: 7028315742573472677}

	// ErrUnsupportedEncoding indicates unsupported encoding.
	ErrUnsupportedEncoding error = &SmppErr{err: "Unsupported encoding", serialVersionUID: 698236878268332}

	// ErrTimeout indicates error timeout.
	ErrTimeout error = &SmppErr{err: "Timeout", serialVersionUID: 4873432724200896611}

	// ErrSessionNotSync indicates session is not synchoronous.
	ErrSessionNotSync error = &SmppErr{err: "Session is not synchoronous", serialVersionUID: -2785891348929001265}

	// ErrWrongSessionState indicates wrong session state.
	ErrWrongSessionState error = &SmppErr{err: "Wrong session state", serialVersionUID: 7296414687928430713}
)
