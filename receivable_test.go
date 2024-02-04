package gosmpp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/linxGnu/gosmpp/pdu"

	"github.com/stretchr/testify/require"
)

func TestReceive(t *testing.T) {
	auth := nextAuth()
	receiver, err := NewSession(
		RXConnector(NonTLSDialer, auth),
		Settings{
			ReadTimeout: 2 * time.Second,

			OnReceivingError: func(err error) {
				t.Log(err)
			},

			OnRebindingError: func(err error) {
				t.Log(err)
			},

			OnPDU: func(p pdu.PDU, _ bool) {
				t.Log(p)
			},

			OnClosed: func(state State) {
				t.Log(state)
			},
		}, 5*time.Second)
	require.Nil(t, err)
	require.NotNil(t, receiver)
	defer func() {
		_ = receiver.Close()
	}()

	require.Equal(t, "MelroseLabsSMSC", receiver.Receiver().SystemID())

	time.Sleep(time.Second)
	receiver.rebind()
}

func Test_receivable_handleAllPdu(t1 *testing.T) {
	type fields struct {
		settings Settings
	}
	type args struct {
		p pdu.PDU
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantClosing bool
	}{
		{
			name:        "nil setting",
			fields:      fields{},
			args:        args{},
			wantClosing: false,
		},
		{
			name:        "nil pdu",
			fields:      fields{},
			args:        args{},
			wantClosing: false,
		},
		{
			name: "DeliverSM pdu",
			fields: fields{
				settings: Settings{
					OnAllPDU: receivableHandleAllPDU(t1),
				},
			},
			args: args{
				p: pdu.NewDeliverSM(),
			},
			wantClosing: false,
		},
		{
			name: "EnquireLink pdu",
			fields: fields{
				settings: Settings{
					OnAllPDU: receivableHandleAllPDU(t1),
				},
			},
			args: args{
				p: pdu.NewEnquireLink(),
			},
			wantClosing: false,
		},
		/*{
			name: "Undind pdu", // run this as the last test case
			fields: fields{
				settings: Settings{
					OnAllPDU: receivableHandleAllPDU(t1),
				},
			},
			args: args{
				p: pdu.NewUnbind(),
			},
			wantClosing: true,
		},*/
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &receivable{
				settings: tt.fields.settings,
			}
			assert.Equalf(t1, tt.wantClosing, t.handleAllPdu(tt.args.p), "handleAllPdu(%v)", tt.args.p)
		})
	}
}

func receivableHandleAllPDU(t1 *testing.T) func(pdu.PDU) (pdu.PDU, bool) {
	return func(p pdu.PDU) (pdu.PDU, bool) {
		switch pd := p.(type) {
		case *pdu.Unbind:
			fmt.Println("Unbind Received")
			return pd.GetResponse(), true

		case *pdu.UnbindResp:
			t1.Log("UnbindResp Received")

		case *pdu.SubmitSMResp:
			t1.Log("SubmitSMResp Received")

		case *pdu.GenericNack:
			t1.Log("GenericNack Received")

		case *pdu.EnquireLinkResp:
			t1.Log("EnquireLinkResp Received")

		case *pdu.EnquireLink:
			t1.Log("EnquireLink Received")
			return pd.GetResponse(), false

		case *pdu.DataSM:
			t1.Log("DataSM received")
			return pd.GetResponse(), false

		case *pdu.DeliverSM:
			t1.Log("DeliverSM received")
			return pd.GetResponse(), false
		}
		return nil, false
	}
}

func Test_receivable_handleOrClose(t1 *testing.T) {
	type fields struct {
		settings Settings
	}
	type args struct {
		p pdu.PDU
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantClosing bool
	}{
		{
			name:        "nil setting",
			fields:      fields{},
			args:        args{},
			wantClosing: false,
		},
		{
			name:        "nil pdu",
			fields:      fields{},
			args:        args{},
			wantClosing: false,
		},
		{
			name:   "EnquireLink pdu",
			fields: fields{},
			args: args{
				p: pdu.NewEnquireLink(),
			},
			wantClosing: false,
		},
		/*{
			name:   "Undind pdu", // run this as the last test case
			fields: fields{},
			args: args{
				p: pdu.NewUnbind(),
			},
			wantClosing: true,
		},*/
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &receivable{
				settings: tt.fields.settings,
			}
			assert.Equalf(t1, tt.wantClosing, t.handleOrClose(tt.args.p), "handleOrClose(%v)", tt.args.p)
		})
	}
}

func Test_receivable_handleWindowPdu(t1 *testing.T) {
	type fields struct {
		settings Settings
	}
	type args struct {
		p pdu.PDU
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantClosing bool
	}{
		{
			name:        "nil setting",
			fields:      fields{},
			args:        args{},
			wantClosing: false,
		},
		{
			name:        "nil pdu",
			fields:      fields{},
			args:        args{},
			wantClosing: false,
		},
		{
			name:   "EnquireLink pdu",
			fields: fields{},
			args: args{
				p: pdu.NewEnquireLink(),
			},
			wantClosing: false,
		},
		{
			name:   "EnquireLinkResp pdu",
			fields: fields{},
			args: args{
				p: pdu.NewEnquireLink().GetResponse(),
			},
			wantClosing: false,
		},
		{
			name:   "DeliverSM pdu",
			fields: fields{},
			args: args{
				p: pdu.NewDeliverSM(),
			},
			wantClosing: false,
		},
		{
			name:   "SubmitSMResp pdu",
			fields: fields{},
			args: args{
				p: pdu.NewSubmitSM().GetResponse(),
			},
			wantClosing: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &receivable{
				settings: tt.fields.settings,
			}
			assert.Equalf(t1, tt.wantClosing, t.handleWindowPdu(tt.args.p), "handleWindowPdu(%v)", tt.args.p)
		})
	}
}
