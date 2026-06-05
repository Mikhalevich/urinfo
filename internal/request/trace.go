package request

import (
	"crypto/tls"
	"net/http/httptrace"
	"time"
)

type Trace struct {
	GetConn              time.Time
	GotConn              time.Time
	GotFirstResponseByte time.Time
	DNSStart             time.Time
	DNSDone              time.Time
	ConnectStart         time.Time
	ConnectDone          time.Time
	TLSStart             time.Time
	TLSDone              time.Time
}

func NewTrace() (*Trace, *httptrace.ClientTrace) {
	t := Trace{}

	return &t, &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			t.GetConn = time.Now()
		},
		GotConn: func(_ httptrace.GotConnInfo) {
			t.GotConn = time.Now()
		},
		GotFirstResponseByte: func() {
			t.GotFirstResponseByte = time.Now()
		},
		DNSStart: func(_ httptrace.DNSStartInfo) {
			t.DNSStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			t.DNSDone = time.Now()
		},
		ConnectStart: func(_, _ string) {
			t.ConnectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			t.ConnectDone = time.Now()
		},
		TLSHandshakeStart: func() {
			t.TLSStart = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			t.TLSDone = time.Now()
		},
	}
}
