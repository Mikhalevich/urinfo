package trace

import (
	"crypto/tls"
	"net/http/httptrace"
	"time"
)

// Trace struct contains http trace timings.
type Trace struct {
	Start                time.Time
	GetConn              time.Time
	GotConn              time.Time
	GotFirstResponseByte time.Time
	DNSStart             time.Time
	DNSDone              time.Time
	ConnectStart         time.Time
	ConnectDone          time.Time
	TLSStart             time.Time
	TLSDone              time.Time
	Done                 time.Time
}

// New constructor for Trace struct, returns Trace itself and ClientTrace for request context.
func New() (*Trace, *httptrace.ClientTrace) {
	trace := Trace{}

	return &trace, &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			trace.GetConn = time.Now()
		},
		GotConn: func(_ httptrace.GotConnInfo) {
			trace.GotConn = time.Now()
		},
		GotFirstResponseByte: func() {
			trace.GotFirstResponseByte = time.Now()
		},
		DNSStart: func(_ httptrace.DNSStartInfo) {
			trace.DNSStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			trace.DNSDone = time.Now()
		},
		ConnectStart: func(_, _ string) {
			trace.ConnectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			trace.ConnectDone = time.Now()
		},
		TLSHandshakeStart: func() {
			trace.TLSStart = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			trace.TLSDone = time.Now()
		},
	}
}
