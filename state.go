package gosmpp

const (
	Alive int32 = iota
	Closed
)

// State represents Transmitter/Receiver/Transceiver state.
type State byte

const (
	// ExplicitClosing indicates that Transmitter/Receiver/Transceiver is closed
	// explicitly (from outside).
	ExplicitClosing State = iota

	// StoppingProcessOnly stops daemons but does not close underlying net conn.
	StoppingProcessOnly

	// InvalidStreaming indicates Transceiver/Receiver data reading state is
	// invalid due to network connection or SMSC responsed with an invalid PDU
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

// String interface.
func (s *State) String() string {
	switch *s {
	case ExplicitClosing:
		return "ExplicitClosing"

	case StoppingProcessOnly:
		return "StoppingProcessOnly"

	case InvalidStreaming:
		return "InvalidStreaming"

	case ConnectionIssue:
		return "ConnectionIssue"

	case UnbindClosing:
		return "UnbindClosing"

	default:
		return ""
	}
}
