package gosmpp

// State represents Transmitter/Receiver/Transceiver state.
type State byte

const (
	// ExplicitClosing indicates that Transmitter/Receiver/Transceiver is closed
	// explicitly (from outside).
	ExplicitClosing State = iota

	// InvalidStreaming indicates Transceiver/Receiver data reading state is
	// invalid due to network connection/ or SMSC responsed with an invalid PDU
	// which potentially damages other following PDU(s).
	//
	// In both cases, Transceiver/Receiver is closed implicitly.
	InvalidStreaming

	// ConnectionIssue indicates that Transmitter/Receiver/Transceiver is closed
	// due to network connection issue or SMSC is not available anymore.
	ConnectionIssue

	// UnbindClosing indicates Receiver got unbind request from SMSC and closed due to this request.
	UnbindClosing
)
