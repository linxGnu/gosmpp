package Exception

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/linxGnu/gosmpp/Data"
)

type Exception struct {
	Error            error
	ErrorCode        int32
	SerialVersionUID int64
}

func NewExceptionFromStr(str string) *Exception {
	if len(str) <= 0 {
		return nil
	}

	return &Exception{
		errors.New(str),
		0,
		0,
	}
}

func NewException(err error) *Exception {
	if err == nil {
		return nil
	}

	return &Exception{
		fmt.Errorf("%v with stack trace: %v", err.Error(), string(debug.Stack())),
		0,
		0,
	}
}

var EXCEPTION_TCP_NOT_OPEN *Exception = NewExceptionFromStr("TCP connection not opened")

var PDUException *Exception = &Exception{errors.New("Unknown"), Data.ESME_RUNKNOWNERR, 5174606627714918071}

var HeaderIncompleteException *Exception = &Exception{errors.New("HeaderIncompleteException"), Data.ESME_RUNKNOWNERR, 3813030249230646966}

var IntegerOutOfRangeException *Exception = &Exception{errors.New("IntegerOutOfRangeException"), Data.ESME_RUNKNOWNERR, 750364511680192335}

var InvalidPDUException *Exception = &Exception{errors.New("InvalidPDUException"), Data.ESME_RUNKNOWNERR, -6985061862208729984}

var MessageIncompleteException *Exception = &Exception{errors.New("MessageIncompleteException"), Data.ESME_RUNKNOWNERR, 9081917316954328823}

var TLVException *Exception = &Exception{errors.New("TLVException"), Data.ESME_RUNKNOWNERR, -6659626685298184198}

var TooManyValuesException *Exception = &Exception{errors.New("TooManyValuesException"), Data.ESME_RUNKNOWNERR, -2777016699062489252}

var UnexpectedOptionalParameterException *Exception = &Exception{errors.New("UnexpectedOptionalParameterException"), Data.ESME_RUNKNOWNERR, -1284359967986779783}

var UnknownCommandIdException *Exception = &Exception{errors.New("UnknownCommandIdException"), Data.ESME_RUNKNOWNERR, -5091873576710864441}

var ValueNotSetException *Exception = &Exception{errors.New("ValueNotSetException"), Data.ESME_RUNKNOWNERR, -4595064103809398438}

var WrongDateFormatException *Exception = &Exception{errors.New("WrongDateFormatException"), Data.ESME_RUNKNOWNERR, 5831937612139037591}

var WrongDestFlagException *Exception = &Exception{errors.New("WrongDestFlagException"), Data.ESME_RUNKNOWNERR, 6266749651012701472}

var WrongLengthException *Exception = &Exception{errors.New("WrongLengthException"), Data.ESME_RUNKNOWNERR, 7935018427341458286}

var WrongLengthOfStringException *Exception = &Exception{errors.New("WrongLengthOfStringException"), Data.ESME_RUNKNOWNERR, 8604133584902790266}

var NotEnoughDataInByteBufferException *Exception = &Exception{errors.New("NotEnoughDataInByteBufferException"), Data.ESME_RUNKNOWNERR, -3720107899765064964}

var TerminatingZeroNotFoundException *Exception = &Exception{errors.New("TerminatingZeroNotFoundException"), Data.ESME_RUNKNOWNERR, 7028315742573472677}

var UnsupportedEncodingException *Exception = &Exception{errors.New("UnsupportedEncodingException"), Data.ESME_RUNKNOWNERR, 7028315742573472677}

var TimeoutException *Exception = &Exception{errors.New("TimeoutException"), Data.ESME_RUNKNOWNERR, 4873432724200896611}

var NotSynchronousException *Exception = &Exception{errors.New("NotSynchronousException"), Data.ESME_RUNKNOWNERR, -2785891348929001265}

var WrongSessionStateException *Exception = &Exception{errors.New("WrongSessionStateException"), Data.ESME_RUNKNOWNERR, 7296414687928430713}

var ConnectionClosingDueToError *Exception = &Exception{errors.New("Closing connection because of IOException receive via TCPIPConnection"), 0, 0}
