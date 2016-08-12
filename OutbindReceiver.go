package gosmpp

const (
	OUTBIND_RECEIVER_THREAD_NAME string = "OutbindRcv"
)

// Since Outbind is for server, we won't support it at the moment
type OutbindReceiver struct {
	ReceiverBase
}
