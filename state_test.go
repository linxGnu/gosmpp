package gosmpp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		name string
		s    State
		want string
	}{
		{
			name: "ExplicitClosing",
			s:    ExplicitClosing,
			want: "ExplicitClosing",
		},
		{
			name: "StoppingProcessOnly",
			s:    StoppingProcessOnly,
			want: "StoppingProcessOnly",
		},
		{
			name: "InvalidStreaming",
			s:    InvalidStreaming,
			want: "InvalidStreaming",
		},
		{
			name: "ConnectionIssue",
			s:    ConnectionIssue,
			want: "ConnectionIssue",
		},
		{
			name: "UnbindClosing",
			s:    UnbindClosing,
			want: "UnbindClosing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.String(), "String()")
		})
	}
}
